package mwals

import (
	"context"
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

// Resolver defines the core logic for translating an alias to endpoints.
type Resolver interface {
	Resolve(ctx context.Context, alias string) (*ResolutionResponse, error)
}

// Normalizer handles string cleaning for aliases.
func Normalizer(alias string) string {
	// 1. Convert to lowercase
	clean := strings.ToLower(alias)
	// 2. Trim whitespace
	clean = strings.TrimSpace(clean)
	// 3. Remove '@' prefix if present for consistent internal storage
	clean = strings.TrimPrefix(clean, "@")
	// 4. (Optional) Strip other special characters if needed
	return clean
}

// Service is the reference implementation of the Resolver.
type Service struct {
	mu sync.RWMutex
	// store maps clean alias names to records
	store map[string]*AliasRecord
	// Private key for signing resolution responses
	signingKey ed25519.PrivateKey
	// persistencePath is the location of the JSON data store
	persistencePath string
}

// AliasRecord represents the internal database state for an alias.
type AliasRecord struct {
	Alias             string
	Status            AliasStatus
	IdentityMask      string
	Attestation       AttestationLevel
	Endpoints         []Endpoint
	VerificationProof string
	// IsPrivate indicates that endpoints should be returned as signed tokens (Blind Resolution)
	IsPrivate bool
}

func NewService(key ed25519.PrivateKey, dataPath string) (*Service, error) {
	s := &Service{
		store:           make(map[string]*AliasRecord),
		signingKey:      key,
		persistencePath: dataPath,
	}

	if dataPath != "" {
		if err := s.load(); err != nil && !os.IsNotExist(err) {
			return nil, fmt.Errorf("failed to load data store: %w", err)
		}
	}
	return s, nil
}

func (s *Service) load() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := os.ReadFile(s.persistencePath)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &s.store)
}

func (s *Service) save() error {
	if s.persistencePath == "" {
		return nil
	}

	s.mu.RLock()
	data, err := json.MarshalIndent(s.store, "", "  ")
	s.mu.RUnlock()

	if err != nil {
		return err
	}

	return os.WriteFile(s.persistencePath, data, 0644)
}

// Resolve implements the Resolver interface.
func (s *Service) Resolve(ctx context.Context, alias string) (*ResolutionResponse, error) {
	clean := Normalizer(alias)

	record, ok := s.store[clean]
	if !ok {
		return nil, fmt.Errorf("alias not found: %s", alias)
	}

	if record.Status == AliasStatusSuspended {
		return nil, fmt.Errorf("alias is suspended")
	}

	resp := &ResolutionResponse{
		Alias:               "@" + record.Alias,
		Status:              record.Status,
		IdentityMask:        record.IdentityMask,
		ResolutionTimestamp: time.Now().UTC(),
		Endpoints:           make([]Endpoint, len(record.Endpoints)),
	}

	for i, ep := range record.Endpoints {
		resp.Endpoints[i] = ep
		if record.IsPrivate {
			// Blind the destination with a signed resolution token
			// token = provider|dest|expiry signed by ALS
			expiry := time.Now().Add(10 * time.Minute).Format(time.RFC3339)
			payload := fmt.Sprintf("%s|%s|%s", ep.Provider, ep.Destination, expiry)
			sig := ed25519.Sign(s.signingKey, []byte(payload))

			resp.Endpoints[i].Destination = fmt.Sprintf("TOKEN:%s:%s", hex.EncodeToString(sig), expiry)
		}
	}

	// Sign the response
	sig, err := s.signResponse(resp)
	if err != nil {
		return nil, fmt.Errorf("failed to sign response: %w", err)
	}
	resp.SecuritySig = sig

	return resp, nil
}

func (s *Service) signResponse(resp *ResolutionResponse) (string, error) {
	if s.signingKey == nil {
		return "unsigned", nil
	}

	// Create canonical string for signing
	// format: alias|status|timestamp|endpoint_count
	canonical := fmt.Sprintf("%s|%s|%s|%d",
		resp.Alias,
		resp.Status,
		resp.ResolutionTimestamp.Format(time.RFC3339),
		len(resp.Endpoints),
	)

	sig := ed25519.Sign(s.signingKey, []byte(canonical))
	return hex.EncodeToString(sig), nil
}

// Seed adds a record to the mock store.
func (s *Service) Seed(record *AliasRecord) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.store[Normalizer(record.Alias)] = record
}

// Register adds a new alias to the service.
func (s *Service) Register(ctx context.Context, record *AliasRecord) error {
	if IsReserved(record.Alias) {
		return fmt.Errorf("alias is reserved: %s", record.Alias)
	}

	clean := Normalizer(record.Alias)

	s.mu.Lock()
	if _, exists := s.store[clean]; exists {
		s.mu.Unlock()
		return fmt.Errorf("alias already registered: %s", record.Alias)
	}
	s.store[clean] = record
	s.mu.Unlock()

	return s.save()
}

// IsReserved checks for sensitive aliases.
func IsReserved(alias string) bool {
	reserved := []string{"president", "government", "airtel", "tnm", "natswitch", "mw-als", "admin"}
	clean := Normalizer(alias)
	for _, r := range reserved {
		if clean == r {
			return true
		}
	}
	return false
}

// AttestAlias simulates the trust level upgrade process.
func (s *Service) AttestAlias(alias string, level AttestationLevel, proof string) error {
	clean := Normalizer(alias)
	record, ok := s.store[clean]
	if !ok {
		return fmt.Errorf("alias not found")
	}

	// Validation logic without hardcoded values
	switch level {
	case AttestationVerified:
		// Check against the record's specific expected proof (e.g. dynamic OTP)
		if record.VerificationProof == "" {
			return fmt.Errorf("no verification challenge pending for this alias")
		}
		if proof != record.VerificationProof {
			return fmt.Errorf("invalid verification proof")
		}
	case AttestationCertified:
		// Suppose proof is a signed NRIS blob
		if !strings.HasPrefix(proof, "NRIS-") {
			return fmt.Errorf("invalid NRIS attestation proof: missing identity certificate prefix")
		}
	}

	record.Attestation = level
	return nil
}
