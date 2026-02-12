package mwals

import (
	"context"
	"fmt"
)

// BesuMock provides a development/test implementation of the BlockchainResolver.
type BesuMock struct {
	ChainStore map[string]*AliasRecord
}

func NewBesuMock() *BesuMock {
	return &BesuMock{
		ChainStore: make(map[string]*AliasRecord),
	}
}

func (m *BesuMock) Resolve(ctx context.Context, alias string) (*AliasRecord, error) {
	record, ok := m.ChainStore[Normalizer(alias)]
	if !ok {
		return nil, fmt.Errorf("alias not found on blockchain: %s", alias)
	}
	return record, nil
}

func (m *BesuMock) RegisterOnChain(ctx context.Context, record *AliasRecord) error {
	m.ChainStore[Normalizer(record.Alias)] = record
	fmt.Printf("Blockchain: Registered @%s\n", record.Alias)
	return nil
}

// BesuClient is the placeholder for the production implementation.
// In a full implementation, this would use github.com/ethereum/go-ethereum/ethclient.
type BesuClient struct {
	RPCURL          string
	ContractAddress string
}

func NewBesuClient(url string, address string) *BesuClient {
	return &BesuClient{
		RPCURL:          url,
		ContractAddress: address,
	}
}

func (c *BesuClient) Resolve(ctx context.Context, alias string) (*AliasRecord, error) {
	// This is where you would use abigen bindings or direct JSON-RPC calls.
	// Example: eth_call to the MWAliasRegistry contract at ContractAddress.
	return nil, fmt.Errorf("BesuClient integration requires ethclient dependency. Use BesuMock for now.")
}

func (c *BesuClient) RegisterOnChain(ctx context.Context, record *AliasRecord) error {
	// This is where you would send an eth_sendRawTransaction to call registerAlias.
	return fmt.Errorf("BesuClient integration requires ethclient dependency.")
}
