package mwals

import "time"

// AliasStatus represents the state of an alias in the registry.
type AliasStatus string

const (
	AliasStatusActive    AliasStatus = "ACTIVE"
	AliasStatusSuspended AliasStatus = "SUSPENDED"
	AliasStatusPending   AliasStatus = "PENDING"
)

// EndpointType defines the type of financial destination.
type EndpointType string

const (
	EndpointTypeWallet      EndpointType = "WALLET"
	EndpointTypeBankAccount EndpointType = "BANK_ACCOUNT"
)

// Endpoint represents a specific financial destination for an alias.
type Endpoint struct {
	Priority         int          `json:"priority"`
	Provider         string       `json:"provider"`
	Type             EndpointType `json:"type"`
	Destination      string       `json:"destination"` // Encrypted or masked
	SupportedMethods []string     `json:"supported_methods"`
}

// ResolutionResponse is the object returned when an alias is resolved.
type ResolutionResponse struct {
	Alias               string      `json:"alias"`
	Status              AliasStatus `json:"status"`
	IdentityMask        string      `json:"identity_mask"`
	ResolutionTimestamp time.Time   `json:"resolution_timestamp"`
	Endpoints           []Endpoint  `json:"endpoints"`
	SecuritySig         string      `json:"security_sig"`
}

// AttestationLevel defines the trust level of an alias.
type AttestationLevel int

const (
	AttestationUnverified AttestationLevel = 1 // Basic registration
	AttestationVerified   AttestationLevel = 2 // OTP verified
	AttestationCertified  AttestationLevel = 3 // National ID verified
)
