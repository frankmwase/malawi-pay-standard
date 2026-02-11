package mwjson

import (
	"encoding/json"
	"time"
)

// Global constants
const (
	MWJSONVersion = "1.0"
	CurrencyMWK   = "MWK"
)

// Transaction represents the root MW-JSON object
type Transaction struct {
	MWVersion  string     `json:"mw_version"`
	Header     Header     `json:"header"`
	Payload    Payload    `json:"payload"`
	TrustLayer TrustLayer `json:"trust_layer"`
}

// Header contains metadata about the transaction
type Header struct {
	MsgID          string    `json:"msg_id"`
	Timestamp      time.Time `json:"timestamp"` // Will be enforced to UTC
	TTL            int       `json:"ttl"`       // Seconds to live
	IdempotencyKey string    `json:"idempotency_key"`
}

// Payload contains the business logic of the transaction
type Payload struct {
	Amount   float64     `json:"amount"` // Validated for precision later
	Currency string      `json:"currency"`
	Type     TxType      `json:"type"`
	Sender   Participant `json:"sender"`
	Receiver Participant `json:"receiver"`
}

// Participant represents a sender or receiver within the transaction
type Participant struct {
	ID       string   `json:"id"`
	IDType   IDType   `json:"id_type"`
	Provider Provider `json:"provider"`
	Alias    string   `json:"alias,omitempty"`
}

// TrustLayer contains security and verification data
type TrustLayer struct {
	IntegrityHash string `json:"integrity_hash"`
	KYCVerified   bool   `json:"kyc_verified"`
	Signature     string `json:"extension_signature"` // Ed25519 signature
}

// Enums

type Provider string

const (
	ProviderAirtelMoney  Provider = "AIRTEL_MONEY"
	ProviderTNMPamba     Provider = "TNM_MPAMBA"
	ProviderNationalBank Provider = "NBM"
	ProviderStandardBank Provider = "STANDARD_BANK"
	ProviderFDH          Provider = "FDH"
	// Extend as needed via Registry later
)

type IDType string

const (
	IDTypeNRIS   IDType = "NRIS"   // National ID
	IDTypeMSISDN IDType = "MSISDN" // Phone Number
	IDTypeIBAN   IDType = "IBAN"   // Bank Account
)

type TxType string

const (
	TxTypeP2P TxType = "P2P" // Person to Person
	TxTypeC2B TxType = "C2B" // Customer to Business
	TxTypeB2C TxType = "B2C" // Business to Customer
)

// Helper for strict JSON marshaling if needed
func (t *Transaction) ToJSON() ([]byte, error) {
	return json.Marshal(t)
}

func FromJSON(data []byte) (*Transaction, error) {
	var t Transaction
	err := json.Unmarshal(data, &t)
	return &t, err
}
