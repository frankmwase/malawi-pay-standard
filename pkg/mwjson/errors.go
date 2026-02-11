package mwjson

import "fmt"

type MWErrorCode string

const (
	// Transaction Errors
	ErrInsufficientFunds MWErrorCode = "MW001"
	ErrGhostTransaction  MWErrorCode = "MW408" // Timeout/TTL expired
	ErrSchemaValidation  MWErrorCode = "MW400"
	ErrDuplicateTx       MWErrorCode = "MW409" // Idempotency conflict

	// Authentication/Authorization Errors
	ErrInvalidSignature MWErrorCode = "MW401"
	ErrUnauthorized     MWErrorCode = "MW403"

	// Resource Errors
	ErrAliasNotFound MWErrorCode = "MW404"
	ErrProviderDown  MWErrorCode = "MW503"

	// System Errors
	ErrInternalError MWErrorCode = "MW500"
)

// MWError represents a standardized error in the Malawian Digital Exchange ecosystem.
type MWError struct {
	Code    MWErrorCode `json:"code"`
	Message string      `json:"message"`
	Details string      `json:"details,omitempty"`
}

// Error implements the error interface.
func (e *MWError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("[%s] %s: %s", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// NewMWError creates a new standardized error.
func NewMWError(code MWErrorCode, msg string, details string) *MWError {
	return &MWError{
		Code:    code,
		Message: msg,
		Details: details,
	}
}
