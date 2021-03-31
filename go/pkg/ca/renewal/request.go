// Copyright 2020 Anapaya Systems
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package renewal

import (
	"bytes"
	"context"
	"crypto/x509"
	"time"

	"google.golang.org/protobuf/proto"

	"github.com/scionproto/scion/go/lib/scrypto"
	"github.com/scionproto/scion/go/lib/scrypto/cms/protocol"
	"github.com/scionproto/scion/go/lib/scrypto/cppki"
	"github.com/scionproto/scion/go/lib/scrypto/signed"
	"github.com/scionproto/scion/go/lib/serrors"
	cppb "github.com/scionproto/scion/go/pkg/proto/control_plane"
	cryptopb "github.com/scionproto/scion/go/pkg/proto/crypto"
	"github.com/scionproto/scion/go/pkg/trust"
)

// NewChainRenewalRequest builds a ChainRenewalRequest given a serialized CSR
// and a signer enveloped in a CMS SignedData.
func NewChainRenewalRequest(ctx context.Context, csr []byte,
	signer trust.Signer) (*cppb.ChainRenewalRequest, error) {

	signedCMS, err := signer.SignCMS(ctx, csr)
	if err != nil {
		return nil, err
	}
	return &cppb.ChainRenewalRequest{
		CmsSignedRequest: signedCMS,
	}, nil
}

// NewLegacyChainRenewalRequest builds a ChainRenewalRequest given a serialized CSR
// and a signer enveloped in a protobuf SignedMessage.
func NewLegacyChainRenewalRequest(ctx context.Context, csr []byte,
	signer trust.Signer) (*cppb.ChainRenewalRequest, error) {

	body := &cppb.ChainRenewalRequestBody{
		Csr: csr,
	}
	rawBody, err := proto.Marshal(body)
	if err != nil {
		return nil, err
	}
	signedMsg, err := signer.Sign(ctx, rawBody)
	if err != nil {
		return nil, err
	}
	return &cppb.ChainRenewalRequest{
		SignedRequest: signedMsg,
	}, nil
}

type TRCFetcher interface {
	SignedTRC(ctx context.Context, id cppki.TRCID) (cppki.SignedTRC, error)
}

type RequestVerifier struct {
	TRCFetcher TRCFetcher
}

// VerifyPbSignedRenewalRequest verifies a renewal request that is encapsulated in a protobuf
// SignedMessage envelop. It checks that the contained CSR is valid and correctly self-signed, and
// that the signature is valid and can be verified by a chain in the given chains.
func (r RequestVerifier) VerifyPbSignedRenewalRequest(ctx context.Context,
	req *cryptopb.SignedMessage, chains [][]*x509.Certificate) (*x509.CertificateRequest, error) {

	var authCert *x509.Certificate
	var msg *signed.Message
	for _, chain := range chains {
		m, err := signed.Verify(req, chain[0].PublicKey)
		if err == nil {
			msg, authCert = m, chain[0]
			break
		}
	}
	if msg == nil {
		return nil, serrors.New("no provided chain can verify the signature")
	}
	var body cppb.ChainRenewalRequestBody
	if err := proto.Unmarshal(msg.Body, &body); err != nil {
		return nil, serrors.WrapStr("parsing request body", err)
	}
	csr, err := x509.ParseCertificateRequest(body.Csr)
	if err != nil {
		return nil, serrors.WrapStr("parsing CSR", err)
	}

	return r.processCSR(csr, authCert)
}

// VerifyCMSSignedRenewalRequest verifies a renewal request that is encapsulated in a CMS
// envelop. It checks that the contained CSR is valid and correctly self-signed, and
// that the signature is valid and can be verified by the chain included in the CMS envelop.
func (r RequestVerifier) VerifyCMSSignedRenewalRequest(ctx context.Context,
	req []byte) (*x509.CertificateRequest, error) {

	ci, err := protocol.ParseContentInfo(req)
	if err != nil {
		return nil, serrors.WrapStr("parsing ContentInfo", err)
	}
	sd, err := ci.SignedDataContent()
	if err != nil {
		return nil, serrors.WrapStr("parsing SignedData", err)
	}
	if sd.Version != 1 {
		return nil, serrors.New("unsupported SignedData version", "actual", sd.Version,
			"supported", 1)
	}

	chain, err := ExtractChain(sd)
	if err != nil {
		return nil, serrors.WrapStr("extracting signing certificate chain", err)
	}

	if len(sd.SignerInfos) != 1 {
		return nil, serrors.New("unexpected number of signers", "expected", 1,
			"actual", len(sd.SignerInfos))
	}
	si := sd.SignerInfos[0]
	if _, err := si.FindCertificate(chain); err != nil {
		return nil, serrors.WrapStr("selecting client certificate", err)
	}
	if err := r.verifyClientChain(ctx, chain); err != nil {
		return nil, serrors.WrapStr("verifying client chain", err)
	}

	if !sd.EncapContentInfo.IsTypeData() {
		return nil, serrors.New("unsupported EncapContentInfo type",
			"type", sd.EncapContentInfo.EContentType)
	}
	pld, err := sd.EncapContentInfo.EContentValue()
	if err != nil {
		return nil, serrors.WrapStr("reading payload", err)
	}

	if err := verifySignerInfo(pld, chain[0], si); err != nil {
		return nil, serrors.WrapStr("verifying signer info", err)
	}

	csr, err := x509.ParseCertificateRequest(pld)
	if err != nil {
		return nil, serrors.WrapStr("parsing CSR", err)
	}

	return r.processCSR(csr, chain[0])
}

func (r RequestVerifier) verifyClientChain(ctx context.Context, chain []*x509.Certificate) error {
	ia, err := cppki.ExtractIA(chain[0].Subject)
	if err != nil {
		return err
	}
	tid := cppki.TRCID{
		ISD:    ia.I,
		Serial: scrypto.LatestVer,
		Base:   scrypto.LatestVer,
	}
	trc, err := r.TRCFetcher.SignedTRC(ctx, tid)
	if err != nil {
		return serrors.WrapStr("loading TRC to verify client chain", err)
	}
	if trc.IsZero() {
		return serrors.New("TRC not found", "isd", ia.I)
	}
	opts := cppki.VerifyOptions{TRC: &trc.TRC}
	if err := cppki.VerifyChain(chain, opts); err != nil {
		// If the the previous TRC is in grace period the CA certificate of the chain might
		// have been issued with a previous Root. Try verifying with the TRC in grace period.
		if time.Now().After(trc.TRC.GracePeriodEnd()) {
			return serrors.WrapStr("verifying client chain", err)
		}
		graceID := trc.TRC.ID
		graceID.Serial--
		prevTRC, err := r.TRCFetcher.SignedTRC(ctx, graceID)
		if err != nil {
			return serrors.WrapStr("loading TRC in grace period to verify client chain", err,
				"trc_id", graceID)
		}
		if prevTRC.IsZero() {
			return serrors.New("TRC in grace period not found", "trc_id", graceID)
		}
		if err := cppki.VerifyChain(chain, cppki.VerifyOptions{TRC: &prevTRC.TRC}); err != nil {
			return serrors.WrapStr("verifying client chain", err)
		}
	}
	return nil
}

func verifySignerInfo(pld []byte, cert *x509.Certificate, si protocol.SignerInfo) error {
	hash, err := si.Hash()
	if err != nil {
		return err
	}
	attrDigest, err := si.GetMessageDigestAttribute()
	if err != nil {
		return err
	}
	actualDigest := hash.New()
	actualDigest.Write(pld)
	if !bytes.Equal(attrDigest, actualDigest.Sum(nil)) {
		return serrors.New("message digest does not match")
	}
	sigInput, err := si.SignedAttrs.MarshaledForVerifying()
	if err != nil {
		return err
	}
	algo := si.X509SignatureAlgorithm()
	return cert.CheckSignature(algo, sigInput, si.Signature)
}

func (r RequestVerifier) processCSR(csr *x509.CertificateRequest,
	cert *x509.Certificate) (*x509.CertificateRequest, error) {

	csrIA, err := cppki.ExtractIA(csr.Subject)
	if err != nil {
		return nil, serrors.WrapStr("extracting ISD-AS from CSR", err)
	}
	chainIA, err := cppki.ExtractIA(cert.Subject)
	if err != nil {
		return nil, serrors.WrapStr("extracting ISD-AS from certificate chain", err)
	}
	if !csrIA.Equal(chainIA) {
		return nil, serrors.New("signing subject is different from CSR subject",
			"csr_isd_as", csrIA, "chain_isd_as", chainIA)
	}
	if err := csr.CheckSignature(); err != nil {
		return nil, serrors.WrapStr("invalid CSR signature", err)
	}
	return csr, nil
}

func ExtractChain(sd *protocol.SignedData) ([]*x509.Certificate, error) {
	certs, err := sd.X509Certificates()
	if err == nil {
		if len(certs) == 0 {
			err = protocol.ErrNoCertificate
		} else if len(certs) != 2 {
			err = serrors.New("unexpected number of certificates", "count", len(certs))
		}
	}
	if err != nil {
		return nil, serrors.WrapStr("parsing certificate chain", err)
	}

	certType, err := cppki.ValidateCert(certs[0])
	if err != nil {
		return nil, serrors.WrapStr("checking certificate type", err)
	}
	if certType == cppki.CA {
		certs[0], certs[1] = certs[1], certs[0]
	}
	if err := cppki.ValidateChain(certs); err != nil {
		return nil, serrors.WrapStr("validating chain", err)
	}
	return certs, nil
}
