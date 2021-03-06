// Package api provides primitives to interact the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen DO NOT EDIT.
package api

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/pkg/errors"
)

// CA defines model for CA.
type CA struct {
	CertValidity Validity     `json:"cert_validity"`
	Policy       Policy       `json:"policy"`
	Subject      Subject      `json:"subject"`
	SubjectKeyId SubjectKeyID `json:"subject_key_id"`
}

// Certificate defines model for Certificate.
type Certificate struct {
	DistinguishedName string       `json:"distinguished_name"`
	IsdAs             IsdAs        `json:"isd_as"`
	SubjectKeyAlgo    string       `json:"subject_key_algo"`
	SubjectKeyId      SubjectKeyID `json:"subject_key_id"`
	Validity          Validity     `json:"validity"`
}

// Chain defines model for Chain.
type Chain struct {
	Issuer  Certificate `json:"issuer"`
	Subject Certificate `json:"subject"`
}

// ChainBrief defines model for ChainBrief.
type ChainBrief struct {
	Id       ChainID  `json:"id"`
	Issuer   IsdAs    `json:"issuer"`
	Subject  IsdAs    `json:"subject"`
	Validity Validity `json:"validity"`
}

// ChainID defines model for ChainID.
type ChainID string

// Hop defines model for Hop.
type Hop struct {
	Interface int   `json:"interface"`
	IsdAs     IsdAs `json:"isd_as"`
}

// IsdAs defines model for IsdAs.
type IsdAs string

// LogLevel defines model for LogLevel.
type LogLevel struct {

	// Logging level
	Level string `json:"level"`
}

// Policy defines model for Policy.
type Policy struct {
	ChainLifetime string `json:"chain_lifetime"`
}

// Problem defines model for Problem.
type Problem struct {

	// A human readable explanation specific to this occurrence of the problem that is helpful to locate the problem and give advice on how to proceed. Written in English and readable for engineers, usually not suited for non technical stakeholders and not localized.
	Detail *string `json:"detail,omitempty"`

	// A URI reference that identifies the specific occurrence of the problem, e.g. by adding a fragment identifier or sub-path to the problem type. May be used to locate the root of this problem in the source code.
	Instance *string `json:"instance,omitempty"`

	// The HTTP status code generated by the origin server for this occurrence of the problem.
	Status int `json:"status"`

	// A short summary of the problem type. Written in English and readable for engineers, usually not suited for non technical stakeholders and not localized.
	Title string `json:"title"`

	// A URI reference that uniquely identifies the problem type only in the context of the provided API. Opposed to the specification in RFC-7807, it is neither recommended to be dereferencable and point to a human-readable documentation nor globally unique for the problem type.
	Type *string `json:"type,omitempty"`
}

// Segment defines model for Segment.
type Segment struct {
	Expiration  time.Time `json:"expiration"`
	Hops        []Hop     `json:"hops"`
	Id          SegmentID `json:"id"`
	LastUpdated time.Time `json:"last_updated"`
	Timestamp   time.Time `json:"timestamp"`
}

// SegmentBrief defines model for SegmentBrief.
type SegmentBrief struct {
	EndIsdAs IsdAs     `json:"end_isd_as"`
	Id       SegmentID `json:"id"`

	// Length of the segment.
	Length     int   `json:"length"`
	StartIsdAs IsdAs `json:"start_isd_as"`
}

// SegmentID defines model for SegmentID.
type SegmentID string

// SegmentIDs defines model for SegmentIDs.
type SegmentIDs []SegmentID

// Signer defines model for Signer.
type Signer struct {
	AsCertificate Certificate `json:"as_certificate"`

	// Signer expiration imposed by chain and TRC validity.
	Expiration time.Time `json:"expiration"`
	TrcId      TRCID     `json:"trc_id"`

	// TRC used as trust root is in grace period, and the latest TRC cannot
	// be used as trust root.
	TrcInGracePeriod bool `json:"trc_in_grace_period"`
}

// StandardError defines model for StandardError.
type StandardError struct {

	// Error message
	Error string `json:"error"`
}

// Subject defines model for Subject.
type Subject struct {
	IsdAs IsdAs `json:"isd_as"`
}

// SubjectKeyID defines model for SubjectKeyID.
type SubjectKeyID string

// TRC defines model for TRC.
type TRC struct {
	AuthoritativeAses []IsdAs  `json:"authoritative_ases"`
	CoreAses          []IsdAs  `json:"core_ases"`
	Description       string   `json:"description"`
	Id                TRCID    `json:"id"`
	Validity          Validity `json:"validity"`
}

// TRCBrief defines model for TRCBrief.
type TRCBrief struct {
	Id TRCID `json:"id"`
}

// TRCID defines model for TRCID.
type TRCID struct {
	BaseNumber   int `json:"base_number"`
	Isd          int `json:"isd"`
	SerialNumber int `json:"serial_number"`
}

// Topology defines model for Topology.
type Topology struct {
	AdditionalProperties map[string]interface{} `json:"-"`
}

// Validity defines model for Validity.
type Validity struct {
	NotAfter  time.Time `json:"not_after"`
	NotBefore time.Time `json:"not_before"`
}

// BadRequest defines model for BadRequest.
type BadRequest StandardError

// GetCertificatesParams defines parameters for GetCertificates.
type GetCertificatesParams struct {
	IsdAs   *IsdAs     `json:"isd_as,omitempty"`
	ValidAt *time.Time `json:"valid_at,omitempty"`
	All     *bool      `json:"all,omitempty"`
}

// SetLogLevelJSONBody defines parameters for SetLogLevel.
type SetLogLevelJSONBody LogLevel

// GetSegmentsParams defines parameters for GetSegments.
type GetSegmentsParams struct {

	// Start ISD-AS of segment.
	StartIsdAs *IsdAs `json:"start_isd_as,omitempty"`

	// Terminal AS of segment.
	EndIsdAs *IsdAs `json:"end_isd_as,omitempty"`
}

// GetTrcsParams defines parameters for GetTrcs.
type GetTrcsParams struct {
	Isd *[]int `json:"isd,omitempty"`
	All *bool  `json:"all,omitempty"`
}

// SetLogLevelJSONRequestBody defines body for SetLogLevel for application/json ContentType.
type SetLogLevelJSONRequestBody SetLogLevelJSONBody

// Getter for additional properties for Topology. Returns the specified
// element and whether it was found
func (a Topology) Get(fieldName string) (value interface{}, found bool) {
	if a.AdditionalProperties != nil {
		value, found = a.AdditionalProperties[fieldName]
	}
	return
}

// Setter for additional properties for Topology
func (a *Topology) Set(fieldName string, value interface{}) {
	if a.AdditionalProperties == nil {
		a.AdditionalProperties = make(map[string]interface{})
	}
	a.AdditionalProperties[fieldName] = value
}

// Override default JSON handling for Topology to handle AdditionalProperties
func (a *Topology) UnmarshalJSON(b []byte) error {
	object := make(map[string]json.RawMessage)
	err := json.Unmarshal(b, &object)
	if err != nil {
		return err
	}

	if len(object) != 0 {
		a.AdditionalProperties = make(map[string]interface{})
		for fieldName, fieldBuf := range object {
			var fieldVal interface{}
			err := json.Unmarshal(fieldBuf, &fieldVal)
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("error unmarshaling field %s", fieldName))
			}
			a.AdditionalProperties[fieldName] = fieldVal
		}
	}
	return nil
}

// Override default JSON handling for Topology to handle AdditionalProperties
func (a Topology) MarshalJSON() ([]byte, error) {
	var err error
	object := make(map[string]json.RawMessage)

	for fieldName, field := range a.AdditionalProperties {
		object[fieldName], err = json.Marshal(field)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("error marshaling '%s'", fieldName))
		}
	}
	return json.Marshal(object)
}
