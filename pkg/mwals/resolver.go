package mwals

import (
	"context"
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"strings"
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
	// Mock DB for demonstration
	store map[string]*AliasRecord
	// Private key for signing resolution responses
	signingKey ed25519.PrivateKey
}

// AliasRecord represents the internal database state for an alias.
type AliasRecord struct {
	Alias        string
	Status       AliasStatus
	IdentityMask string
	Attestation  AttestationLevel
	Endpoints    []Endpoint
	// In a real system, we might store a specialized challenge/proof for verification
	VerificationProof string
}

func NewService(key ed25519.PrivateKey) *Service {
	return &Service{
		store:      make(map[string]*AliasRecord),
		signingKey: key,
	}
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
		Endpoints:           record.Endpoints,
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
	s.store[Normalizer(record.Alias)] = record
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
