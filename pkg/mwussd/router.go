package mwussd

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/frankmwase/malawi-pay-standard/pkg/mwjson"
)

// UssdAction defines the type of interaction for a step.
type UssdAction string

const (
	ActionDial  UssdAction = "DIAL"
	ActionReply UssdAction = "REPLY"
)

// UssdStep represents a single interaction in a USSD session.
type UssdStep struct {
	Action  UssdAction // DIAL or REPLY
	Content string     // The USSD code (e.g. *121#) or the input value (e.g. "1")
	Expect  string     // Regex to match the screen content before proceeding (for Accessibility Svc)
	Masked  bool       // If true, the UI should mask this input (e.g. PIN)
}

// Router handles the translation of a transaction into USSD steps.
type Router struct{}

func NewRouter() *Router {
	return &Router{}
}

// GenerateSession creates a sequence of USSD steps for the given transaction.
// The 'pin' argument is the user's PIN for the provider.
func (r *Router) GenerateSession(txn *mwjson.Transaction, pin string) ([]UssdStep, error) {
	if txn == nil {
		return nil, errors.New("transaction cannot be nil")
	}

	// Route based on sender's provider (who is initiating the USSD?)
	senderProvider := txn.Payload.Sender.Provider

	switch senderProvider {
	case mwjson.ProviderAirtelMoney:
		return r.generateAirtelSession(txn, pin)
	case mwjson.ProviderTNMPamba:
		return r.generateTNMSession(txn, pin)
	default:
		return nil, fmt.Errorf("unsupported provider for USSD: %s", senderProvider)
	}
}

// generateAirtelSession creates the steps for Airtel Money (Example Flow)
// Note: This is a hypothetical flow for "Send Money". Real flow may vary.
func (r *Router) generateAirtelSession(txn *mwjson.Transaction, pin string) ([]UssdStep, error) {
	steps := []UssdStep{
		// 1. Dial the root menu
		{
			Action:  ActionDial,
			Content: "*121#",
			Expect:  "(?i)Menu.*Send Money", // Match "Menu" and "Send Money" case-insensitive
		},
		// 2. Select "Send Money" (Assuming option 1)
		{
			Action:  ActionReply,
			Content: "1",
			Expect:  "(?i)Enter.*Number", // Match "Enter Number"
		},
		// 3. Enter Recipient Number
		{
			Action:  ActionReply,
			Content: txn.Payload.Receiver.ID, // Already normalized
			Expect:  "(?i)Enter.*Amount",     // Match "Enter Amount"
		},
		// 4. Enter Amount
		{
			Action:  ActionReply,
			Content: strconv.FormatFloat(txn.Payload.Amount, 'f', 0, 64), // No decimals for USSD usually
			Expect:  "(?i)Enter.*PIN",                                    // Match "Enter PIN"
		},
		// 5. Enter PIN
		{
			Action:  ActionReply,
			Content: pin,
			Expect:  "(?i)Success|Confirmed", // Match success message
			Masked:  true,
		},
	}
	return steps, nil
}

// generateTNMSession creates the steps for TNM Mpamba (Example Flow)
func (r *Router) generateTNMSession(txn *mwjson.Transaction, pin string) ([]UssdStep, error) {
	steps := []UssdStep{
		// 1. Dial the root menu
		{
			Action:  ActionDial,
			Content: "*444#",
			Expect:  "(?i)Mpamba.*Send Money",
		},
		// 2. Select "Send Money" (Assuming option 1)
		{
			Action:  ActionReply,
			Content: "1",
			Expect:  "(?i)Enter.*Number",
		},
		// 3. Enter Recipient Number
		{
			Action:  ActionReply,
			Content: txn.Payload.Receiver.ID,
			Expect:  "(?i)Enter.*Amount",
		},
		// 4. Enter Amount
		{
			Action:  ActionReply,
			Content: strconv.FormatFloat(txn.Payload.Amount, 'f', 0, 64),
			Expect:  "(?i)Enter.*PIN",
		},
		// 5. Enter PIN
		{
			Action:  ActionReply,
			Content: pin,
			Expect:  "(?i)Success|Transferred",
			Masked:  true,
		},
	}
	return steps, nil
}
