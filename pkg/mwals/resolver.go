package mwals

import (
	"context"
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
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
	VerificationProof string `json:"-"` // Ephemeral challenge; never persist or return it.
	// IsPrivate indicates that endpoints should be returned as signed tokens (Blind Resolution)
	IsPrivate bool
}

func NewService(key ed25519.PrivateKey, dataPath string) (*Service, error) {
	if len(key) != ed25519.PrivateKeySize {
		return nil, fmt.Errorf("an Ed25519 signing key is required")
	}
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

	if err := json.Unmarshal(data, &s.store); err != nil {
		return err
	}
	if s.store == nil {
		return fmt.Errorf("registry must be a JSON object")
	}
	for alias, record := range s.store {
		if record == nil || !aliasPattern.MatchString(alias) {
			return fmt.Errorf("invalid stored alias record")
		}
	}
	return nil
}

// saveLocked writes a private, atomic snapshot while the caller holds s.mu.
func (s *Service) saveLocked() error {
	if s.persistencePath == "" {
		return nil
	}
	data, err := json.MarshalIndent(s.store, "", "  ")
	if err != nil {
		return err
	}
	path := filepath.Clean(s.persistencePath)
	file, err := os.CreateTemp(filepath.Dir(path), ".als-*")
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())
	defer file.Close()
	if err := file.Chmod(0600); err != nil {
		return err
	}
	if _, err := file.Write(data); err != nil {
		return err
	}
	if err := file.Sync(); err != nil {
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}
	return os.Rename(file.Name(), path)
}

func copyRecord(record *AliasRecord) *AliasRecord {
	copy := *record
	copy.Endpoints = append([]Endpoint(nil), record.Endpoints...)
	for i := range copy.Endpoints {
		copy.Endpoints[i].SupportedMethods = append([]string(nil), record.Endpoints[i].SupportedMethods...)
	}
	return &copy
}

var aliasPattern = regexp.MustCompile(`^[a-z0-9_]{3,32}$`)

// Resolve implements the Resolver interface.
func (s *Service) Resolve(ctx context.Context, alias string) (*ResolutionResponse, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	clean := Normalizer(alias)
	s.mu.RLock()
	stored, ok := s.store[clean]
	if !ok {
		s.mu.RUnlock()
		return nil, fmt.Errorf("alias not found: %s", alias)
	}
	record := copyRecord(stored)
	s.mu.RUnlock()

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
			// Prototype blind token only: no redemption API exists yet.
			expiry := time.Now().UTC().Add(10 * time.Minute).Format(time.RFC3339)
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
	// SecuritySig is empty until after serialization. All returned fields are covered.
	data, err := json.Marshal(resp)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(ed25519.Sign(s.signingKey, data)), nil
}

// Seed adds a record to the in-memory store for test/demo use only.
func (s *Service) Seed(record *AliasRecord) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.store[Normalizer(record.Alias)] = copyRecord(record)
}

// Register adds a new alias and persists it before reporting success.
func (s *Service) Register(ctx context.Context, record *AliasRecord) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if record == nil {
		return fmt.Errorf("record is required")
	}
	clean := Normalizer(record.Alias)
	if !aliasPattern.MatchString(clean) || IsReserved(clean) {
		return fmt.Errorf("invalid or reserved alias")
	}
	if record.IsPrivate {
		return fmt.Errorf("private registrations require a token redemption service")
	}
	if len(record.Endpoints) == 0 || len(record.Endpoints) > 10 {
		return fmt.Errorf("between 1 and 10 endpoints required")
	}
	for _, ep := range record.Endpoints {
		if ep.Provider == "" || ep.Destination == "" || len(ep.Destination) > 128 {
			return fmt.Errorf("endpoint provider and destination required (max 128 bytes)")
		}
	}
	clone := copyRecord(record)
	clone.Alias = clean
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.store[clean]; exists {
		return fmt.Errorf("alias already registered: %s", clean)
	}
	s.store[clean] = clone
	if err := s.saveLocked(); err != nil {
		delete(s.store, clean)
		return fmt.Errorf("persist registration: %w", err)
	}
	return nil
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
	s.mu.Lock()
	defer s.mu.Unlock()
	record, ok := s.store[clean]
	if !ok {
		return fmt.Errorf("alias not found")
	}
	if level != AttestationVerified {
		return fmt.Errorf("only verification with an issued challenge is supported")
	}
	if record.VerificationProof == "" || proof != record.VerificationProof {
		return fmt.Errorf("invalid or missing verification proof")
	}
	previous := copyRecord(record)
	record.Attestation = level
	record.VerificationProof = ""
	if err := s.saveLocked(); err != nil {
		s.store[clean] = previous
		return fmt.Errorf("persist attestation: %w", err)
	}
	return nil
}
