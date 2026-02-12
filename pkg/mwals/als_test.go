package mwals_test

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"strings"
	"testing"

	"github.com/frankmwase/malawi-pay-standard/pkg/mwals"
)

func TestNormalizer(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"@Koda_Dev", "koda_dev"},
		{"  @Alice  ", "alice"},
		{"BOB", "bob"},
	}

	for _, tt := range tests {
		got := mwals.Normalizer(tt.input)
		if got != tt.expected {
			t.Errorf("Normalizer(%s) = %s; want %s", tt.input, got, tt.expected)
		}
	}
}

func TestResolve(t *testing.T) {
	_, privKey, _ := ed25519.GenerateKey(rand.Reader)
	svc, _ := mwals.NewService(privKey, "")
	svc.Seed(&mwals.AliasRecord{
		Alias:             "koda_dev",
		Status:            mwals.AliasStatusActive,
		IdentityMask:      "K**** M*******",
		VerificationProof: "654321",
		Endpoints: []mwals.Endpoint{
			{Priority: 1, Provider: "AIRTEL_MONEY", Type: mwals.EndpointTypeWallet, Destination: "26599..."},
		},
	})

	ctx := context.Background()
	resp, err := svc.Resolve(ctx, "@Koda_Dev")
	if err != nil {
		t.Fatalf("Resolve failed: %v", err)
	}

	if resp.Alias != "@koda_dev" {
		t.Errorf("Expected leading @ in response alias, got %s", resp.Alias)
	}

	if resp.SecuritySig == "" || resp.SecuritySig == "unsigned" {
		t.Errorf("Expected valid SecuritySig, got %s", resp.SecuritySig)
	}

	if len(resp.Endpoints) != 1 {
		t.Errorf("Expected 1 endpoint, got %d", len(resp.Endpoints))
	}
}

func TestAttestAlias(t *testing.T) {
	_, privKey, _ := ed25519.GenerateKey(rand.Reader)
	svc, _ := mwals.NewService(privKey, "")
	alias := "alice"
	proof := "998877"
	svc.Seed(&mwals.AliasRecord{
		Alias:             alias,
		Status:            mwals.AliasStatusActive,
		IdentityMask:      "A****",
		VerificationProof: proof,
		Attestation:       mwals.AttestationUnverified,
	})

	// Valid Attestation
	if err := svc.AttestAlias(alias, mwals.AttestationVerified, proof); err != nil {
		t.Errorf("Expected valid attestation to succeed, got %v", err)
	}

	// Invalid Attestation
	if err := svc.AttestAlias(alias, mwals.AttestationVerified, "wrong"); err == nil {
		t.Error("Expected invalid attestation to fail, got nil")
	}
}

func TestReserved(t *testing.T) {
	if !mwals.IsReserved("@Airtel") {
		t.Errorf("Expected @Airtel to be reserved")
	}
	if mwals.IsReserved("lusekel") {
		t.Errorf("Expected lusekel NOT to be reserved")
	}
}

func TestRegister(t *testing.T) {
	_, key, _ := ed25519.GenerateKey(rand.Reader)
	svc, _ := mwals.NewService(key, "")

	record := &mwals.AliasRecord{
		Alias:        "newuser",
		IdentityMask: "N******",
		Endpoints: []mwals.Endpoint{
			{Provider: "TNM", Destination: "0888111222"},
		},
	}

	err := svc.Register(context.Background(), record)
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	// Try to register again
	err = svc.Register(context.Background(), record)
	if err == nil {
		t.Error("expected error for duplicate registration, got nil")
	}
}

func TestBlindResolution(t *testing.T) {
	_, key, _ := ed25519.GenerateKey(rand.Reader)
	svc, _ := mwals.NewService(key, "")
	svc.Seed(&mwals.AliasRecord{
		Alias:        "private_user",
		Status:       mwals.AliasStatusActive,
		IdentityMask: "P***********",
		IsPrivate:    true,
		Endpoints: []mwals.Endpoint{
			{Priority: 1, Provider: "AIRTEL", Destination: "0999000111"},
		},
	})

	resp, err := svc.Resolve(context.Background(), "@private_user")
	if err != nil {
		t.Fatalf("Resolve failed: %v", err)
	}

	if len(resp.Endpoints) == 0 {
		t.Fatal("expected at least one endpoint")
	}

	if !strings.HasPrefix(resp.Endpoints[0].Destination, "TOKEN:") {
		t.Errorf("expected destination to be a TOKEN, got %s", resp.Endpoints[0].Destination)
	}
}
