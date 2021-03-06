// Package api provides primitives to interact the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen DO NOT EDIT.
package api

const (
	BearerAuthScopes = "BearerAuth.Scopes"
)

// AS defines model for AS.
type AS string

// AccessCredentials defines model for AccessCredentials.
type AccessCredentials struct {

	// ID of the control service requesting authentication.
	ClientId string `json:"client_id"`

	// Secret that authenticates the control service.
	ClientSecret string `json:"client_secret"`
}

// AccessToken defines model for AccessToken.
type AccessToken struct {

	// The encoded JWT token
	AccessToken string `json:"access_token"`

	// Validity duration of this token in seconds.
	ExpiresIn int `json:"expires_in"`

	// Type of returned access token. Currently always Bearer.
	TokenType string `json:"token_type"`
}

// CertificateChain defines model for CertificateChain.
type CertificateChain []byte

// HealthCheckStatus defines model for HealthCheckStatus.
type HealthCheckStatus struct {
	Status string `json:"status"`
}

// Problem defines model for Problem.
type Problem struct {

	// Identifier to correlate multiple error messages to the same case.
	CorrelationId *string `json:"correlation_id,omitempty"`

	// A human readable explanation specific to this occurrence of the problem that is helpful to locate the problem and give advice on how to proceed. Written in English and readable for engineers, usually not suited for non technical stakeholders and not localized.
	Detail *string `json:"detail,omitempty"`

	// A URI reference that identifies the specific occurrence of the problem, e.g. by adding a fragment identifier or sub-path to the problem type.
	Instance *string `json:"instance,omitempty"`

	// The HTTP status code generated by the server for this occurrence of the problem.
	Status int `json:"status"`

	// A short summary of the problem type. Written in English and readable for engineers, usually not suited for non technical stakeholders and not localized.
	Title string `json:"title"`

	// A URI reference that uniquely identifies the problem type in the context of the provided API.
	Type string `json:"type"`
}

// RenewalRequest defines model for RenewalRequest.
type RenewalRequest struct {

	// Base64 encoded renewal request as described below.
	//
	// The renewal requests consists of a CMS SignedData structure that
	// contains a PKCS#10 defining the parameters of the requested
	// certificate.
	//
	// The following must hold for the CMS structure:
	//
	// - The `certificates` field in `SignedData` MUST contain an existing
	//   and verifiable certificate chain that authenticates the private
	//   key that was used to sign the CMS structure. It MUST NOT contain
	//   any other certificates.
	//
	// - The `eContentType` is set to `id-data`. The contents of `eContent`
	//   is the ASN.1 DER encoded PKCS#10. This ensures backwards
	//   compatibility with PKCS#7, as described in
	//   [RFC5652](https://tools.ietf.org/html/rfc5652#section-5.2.1)
	//
	// - The `SignerIdentifier` MUST be the choice `IssuerAndSerialNumber`,
	//   thus, `version` in `SignerInfo` must be 1, as required by
	//   [RFC5652](https://tools.ietf.org/html/rfc5652#section-5.3)
	Csr []byte `json:"csr"`
}

// RenewalResponse defines model for RenewalResponse.
type RenewalResponse struct {

	// Base64 encoded certificate chain in PKCS#7 format.
	//
	// The certificate chain is part of the `certificates` field of a
	// `SignedData` in the degenerate case, where `signerInfos` is empty.
	CertificateChain CertificateChain `json:"certificate_chain"`
}

// PostAuthTokenJSONBody defines parameters for PostAuthToken.
type PostAuthTokenJSONBody AccessCredentials

// PostCertificateRenewalJSONBody defines parameters for PostCertificateRenewal.
type PostCertificateRenewalJSONBody RenewalRequest

// PostAuthTokenJSONRequestBody defines body for PostAuthToken for application/json ContentType.
type PostAuthTokenJSONRequestBody PostAuthTokenJSONBody

// PostCertificateRenewalJSONRequestBody defines body for PostCertificateRenewal for application/json ContentType.
type PostCertificateRenewalJSONRequestBody PostCertificateRenewalJSONBody
