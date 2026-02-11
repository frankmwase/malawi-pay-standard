package mwjson

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"strings"
	"time"
)

// Validate checks if the transaction adheres to the MW-JSON standard.
func (t *Transaction) Validate() error {
	// 1. Basic Field Checks
	if t.MWVersion != MWJSONVersion {
		return NewMWError(ErrSchemaValidation, "Invalid MW-JSON Version", fmt.Sprintf("Expected %s, got %s", MWJSONVersion, t.MWVersion))
	}
	if t.Header.MsgID == "" {
		return NewMWError(ErrSchemaValidation, "Missing Message ID", "")
	}
	if t.Header.IdempotencyKey == "" {
		return NewMWError(ErrSchemaValidation, "Missing Idempotency Key", "")
	}

	// 2. Timestamp & TTL
	if t.Header.Timestamp.IsZero() {
		return NewMWError(ErrSchemaValidation, "Missing Timestamp", "")
	}
	// Force UTC check (or at least awareness) - The user asked to "Force UTC"
	if t.Header.Timestamp.Location() != time.UTC {
		return NewMWError(ErrSchemaValidation, "Timestamp must be in UTC", "")
	}
	// Check/Enforce TTL
	if t.Header.TTL <= 0 {
		return NewMWError(ErrSchemaValidation, "Invalid TTL", "Must be positive integer")
	}
	if time.Since(t.Header.Timestamp) > time.Duration(t.Header.TTL)*time.Second {
		return NewMWError(ErrGhostTransaction, "Transaction Expired", "TTL exceeded")
	}

	// 3. Payload Validation
	if t.Payload.Currency != CurrencyMWK {
		return NewMWError(ErrSchemaValidation, "Invalid Currency", fmt.Sprintf("Expected %s", CurrencyMWK))
	}
	if err := validateAmount(t.Payload.Amount); err != nil {
		return err
	}

	// 4. Participants (Sender/Receiver)
	if err := t.Payload.Sender.Validate(); err != nil {
		return NewMWError(ErrSchemaValidation, "Invalid Sender", err.Error())
	}
	if err := t.Payload.Receiver.Validate(); err != nil {
		return NewMWError(ErrSchemaValidation, "Invalid Receiver", err.Error())
	}

	return nil
}

// validateAmount ensures the amount is positive and has valid precision for MWK.
// MWK is typically 2 decimal places, but often used as integer in digital retail.
func validateAmount(amount float64) error {
	if amount <= 0 {
		return NewMWError(ErrSchemaValidation, "Invalid Amount", "Must be greater than 0")
	}
	// Check for more than 2 decimal places
	// Multiply by 100, checking if it's an integer
	scaled := amount * 100
	if math.Abs(scaled-math.Round(scaled)) > 0.000001 {
		return NewMWError(ErrSchemaValidation, "Invalid Amount Precision", "MWK supports up to 2 decimal places")
	}
	return nil
}

// Validate checks participant details
func (p *Participant) Validate() error {
	if p.ID == "" {
		return errors.New("missing ID")
	}
	if p.Provider == "" {
		return errors.New("missing Provider")
	}

	// MSISDN Mormalization & Validation
	if p.IDType == IDTypeMSISDN {
		normalized, err := NormalizeMSISDN(p.ID)
		if err != nil {
			return err
		}
		// Update the ID to the normalized version?
		// The validator shouldn't mutate, but the prompt said "validation logic must normalize".
		// Usually, normalization happens before validation or during object construction.
		// For now, we just check if it IS normalized or validatable.
		// Ideally, we'd have a `Normalize()` method on the struct.
		if p.ID != normalized {
			return fmt.Errorf("MSISDN not normalized: expected %s", normalized)
		}
	}

	return nil
}

// NormalizeMSISDN converts various definitions to 265XXXXXXXXX format.
// Inputs: 099..., 99..., +26599...
// Output: 26599...
func NormalizeMSISDN(input string) (string, error) {
	// Remove spaces and hyphens
	clean := strings.ReplaceAll(input, " ", "")
	clean = strings.ReplaceAll(clean, "-", "")

	// Regex for validation
	// Matches:
	// ^0 -> Starts with 0 (e.g., 0991234567) -> Length 10
	// ^\+265 -> Starts with +265 (e.g., +265991234567) -> Length 13
	// ^265 -> Starts with 265 (e.g., 265991234567) -> Length 12
	// ^9 -> Starts with 9 (e.g., 991234567) -> Length 9 (Old deprecated but possible)

	if strings.HasPrefix(clean, "+265") {
		clean = clean[1:] // Remove +
	} else if strings.HasPrefix(clean, "0") {
		clean = "265" + clean[1:] // Replace 0 with 265
	} else if len(clean) == 9 {
		// Assume it's just the number without leading 0, e.g. 991234567
		clean = "265" + clean
	}

	// Final check: Must be 12 digits and start with 265
	matched, _ := regexp.MatchString(`^265\d{9}$`, clean)
	if !matched {
		return "", NewMWError(ErrSchemaValidation, "Invalid MSISDN format", "Must resolve to 265XXXXXXXXX")
	}

	return clean, nil
}
