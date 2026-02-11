package mwjson

import (
	"context"
	"sync"
)

// PaymentProvider defines the standard interface for all Malawian payment rails.
// Mobile Money (Airtel/TNM) and Banks (NBM/Standard) must implement this.
type PaymentProvider interface {
	// Authorize checks if a user has funds and approves the transaction.
	// Returns a transaction ID or error.
	Authorize(ctx context.Context, entry *Transaction) (string, error)

	// Transfer executes the actual movement of funds.
	// Often called after Authorize.
	Transfer(ctx context.Context, entry *Transaction) (string, error)

	// QueryStatus checks the state of a transaction.
	QueryStatus(ctx context.Context, msgID string) (*TransactionStatus, error)
}

// TransactionStatus holds the result of a status query
type TransactionStatus struct {
	MsgID   string                 `json:"msg_id"`
	Status  string                 `json:"status"` // PENDING, SUCCESS, FAILED
	RawData map[string]interface{} `json:"raw_data,omitempty"`
}

// ProviderRegistry manages the available payment providers.
type ProviderRegistry struct {
	providers map[Provider]PaymentProvider
	mu        sync.RWMutex
}

var (
	// globalRegistry is the default registry.
	globalRegistry = &ProviderRegistry{
		providers: make(map[Provider]PaymentProvider),
	}
)

// RegisterProvider allows a new integration to register itself at runtime.
// e.g., mwjson.RegisterProvider(mwjson.ProviderAirtelMoney, &AirtelAdapter{})
func RegisterProvider(name Provider, p PaymentProvider) {
	globalRegistry.mu.Lock()
	defer globalRegistry.mu.Unlock()
	globalRegistry.providers[name] = p
}

// GetProvider retrieves a registered provider.
func GetProvider(name Provider) (PaymentProvider, error) {
	globalRegistry.mu.RLock()
	defer globalRegistry.mu.RUnlock()
	p, ok := globalRegistry.providers[name]
	if !ok {
		return nil, NewMWError(ErrProviderDown, "Provider Not Registered", string(name))
	}
	return p, nil
}

// ListProviders returns all registered providers.
func ListProviders() []Provider {
	globalRegistry.mu.RLock()
	defer globalRegistry.mu.RUnlock()
	keys := make([]Provider, 0, len(globalRegistry.providers))
	for k := range globalRegistry.providers {
		keys = append(keys, k)
	}
	return keys
}
