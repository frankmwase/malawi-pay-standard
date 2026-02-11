package umqr

import (
	"fmt"
	"sort"
)

// TLV represents a Tag-Length-Value entry.
type TLV struct {
	Tag   string
	Value string
}

// Length returns the value's length as a 2-digit string.
func (t TLV) Length() string {
	return fmt.Sprintf("%02d", len(t.Value))
}

// String returns the formatted TLV string (TTLLVV...).
func (t TLV) String() string {
	return t.Tag + t.Length() + t.Value
}

// Encoder handles building the final QR payload string.
type Encoder struct {
	tags map[string]string
}

func NewEncoder() *Encoder {
	return &Encoder{tags: make(map[string]string)}
}

func (e *Encoder) Set(tag string, value string) {
	e.tags[tag] = value
}

// Encode assembles the tags into a single string, appends the CRC tag (63),
// calculates the CRC, and returns the full payload.
func (e *Encoder) Encode() string {
	var keys []string
	for k := range e.tags {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	payload := ""
	for _, k := range keys {
		payload += TLV{Tag: k, Value: e.tags[k]}.String()
	}

	// Append Tag 63 (CRC) placeholder "6304"
	payload += "6304"

	// Calculate CRC
	crcValue := CalculateCRC16CCITT([]byte(payload))
	crcStr := fmt.Sprintf("%04X", crcValue)

	return payload + crcStr
}

// Tag definitions based on EMVCo & Malawi Standard
const (
	TagPayloadFormatIndicator  = "00" // Always "01"
	TagPointOfInitiationMethod = "01" // "11" for Static, "12" for Dynamic
	TagMalawiMerchantAccount   = "26" // Malawi Interop ID
	TagMerchantCategoryCode    = "52"
	TagTransactionCurrency     = "53" // "454" for MWK
	TagTransactionAmount       = "54"
	TagCountryCode             = "58" // "MW"
	TagMerchantName            = "59"
	TagMerchantCity            = "60"
	TagAdditionalData          = "62" // Ref strings, Invoice #
	TagCRC                     = "63"
)

// Tag 26 Sub-tags for Malawi
const (
	SubTagGlobalID    = "00" // "MW.GOV.NATSWITCH"
	SubTagAccountType = "01" // e.g., "AIRTEL_MONEY", "BANK"
	SubTagAliasName   = "02" // e.g., "@koda_dev"
)

// EncodeTag26 handles the complex Tag 26 mapping.
func EncodeTag26(accountType, alias string) string {
	gID := TLV{Tag: SubTagGlobalID, Value: "MW.GOV.NATSWITCH"}.String()
	aType := TLV{Tag: SubTagAccountType, Value: accountType}.String()
	aName := TLV{Tag: SubTagAliasName, Value: alias}.String()
	return gID + aType + aName
}
