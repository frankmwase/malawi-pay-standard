package mwals

import (
	"context"
	"fmt"
)

// BlockchainResolver defines the interface for chain-based lookups.
type BlockchainResolver interface {
	Resolve(ctx context.Context, alias string) (*AliasRecord, error)
	// RegisterOnChain pushes a record to the blockchain
	RegisterOnChain(ctx context.Context, record *AliasRecord) error
}

// HybridService orchestration: JSON Cache + Blockchain Truth
type HybridService struct {
	local  *Service
	remote BlockchainResolver
}

func NewHybridService(local *Service, remote BlockchainResolver) *HybridService {
	return &HybridService{
		local:  local,
		remote: remote,
	}
}

// Resolve attempts to use local cache first, then falls back to blockchain.
func (h *HybridService) Resolve(ctx context.Context, alias string) (*ResolutionResponse, error) {
	// 1. Check local cache (JSON)
	resp, err := h.local.Resolve(ctx, alias)
	if err == nil {
		return resp, nil
	}

	// 2. Fallback to Blockchain if local lookup fails
	fmt.Printf("Local cache miss for %s. Querying blockchain...\n", alias)
	record, err := h.remote.Resolve(ctx, alias)
	if err != nil {
		return nil, fmt.Errorf("blockchain resolution failed: %w", err)
	}

	// 3. Update local cache with blockchain data (Sync)
	if err := h.local.Register(ctx, record); err != nil {
		fmt.Printf("Warning: Failed to cache blockchain result: %v\n", err)
	}

	return h.local.Resolve(ctx, alias)
}

// Register ensures data is saved to both blockchain and local store.
func (h *HybridService) Register(ctx context.Context, record *AliasRecord) error {
	// 1. Submit to Blockchain (Source of Truth)
	if err := h.remote.RegisterOnChain(ctx, record); err != nil {
		return fmt.Errorf("failed to register on blockchain: %w", err)
	}

	// 2. Update local cache
	return h.local.Register(ctx, record)
}
