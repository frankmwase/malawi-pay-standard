package mwjson_test

import (
	"crypto/ed25519"
	"crypto/rand"
	"testing"
	"time"

	"github.com/frankmwase/malawi-pay-standard/pkg/mwjson"
)

func TestTransactionValidation(t *testing.T) {
	// Create a valid base transaction
	tx := &mwjson.Transaction{
		MWVersion: mwjson.MWJSONVersion,
		Header: mwjson.Header{
			MsgID:          "TXN-TEST-001",
			Timestamp:      time.Now().UTC(),
			TTL:            300,
			IdempotencyKey: "unique-key-123",
		},
		Payload: mwjson.Payload{
			Amount:   15000.00,
			Currency: mwjson.CurrencyMWK,
			Type:     mwjson.TxTypeP2P,
			Sender: mwjson.Participant{
				ID:       "265991234567",
				IDType:   mwjson.IDTypeMSISDN,
				Provider: mwjson.ProviderAirtelMoney,
				Alias:    "@sender",
			},
			Receiver: mwjson.Participant{
				ID:       "265881234567",
				IDType:   mwjson.IDTypeMSISDN,
				Provider: mwjson.ProviderTNMPamba,
				Alias:    "@receiver",
			},
		},
		TrustLayer: mwjson.TrustLayer{
			KYCVerified: true,
		},
	}

	// Test 1: Valid Transaction
	if err := tx.Validate(); err != nil {
		t.Errorf("Expected valid transaction, got error: %v", err)
	}

	// Test 2: Invalid Amount Precision
	tx.Payload.Amount = 15000.123
	if err := tx.Validate(); err == nil {
		t.Error("Expected error for invalid amount precision, got nil")
	}
	tx.Payload.Amount = 15000.00 // Reset

	// Test 3: Invalid MSISDN
	tx.Payload.Sender.ID = "123"
	if err := tx.Validate(); err == nil {
		t.Error("Expected error for invalid MSISDN, got nil")
	}
	tx.Payload.Sender.ID = "265991234567" // Reset

	// Test 4: Expired TTL
	tx.Header.Timestamp = time.Now().UTC().Add(-10 * time.Minute)
	if err := tx.Validate(); err == nil {
		t.Error("Expected error for expired TTL, got nil")
	}
	tx.Header.Timestamp = time.Now().UTC() // Reset
}

func TestSigning(t *testing.T) {
	// Generate keys
	pubKey, privKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("Failed to generate keys: %v", err)
	}

	tx := &mwjson.Transaction{
		MWVersion: mwjson.MWJSONVersion,
		Header: mwjson.Header{
			MsgID:          "TXN-SIG-001",
			Timestamp:      time.Now().UTC(),
			TTL:            300,
			IdempotencyKey: "sig-key-123",
		},
		Payload: mwjson.Payload{
			Amount:   5000.00,
			Currency: mwjson.CurrencyMWK,
			Type:     mwjson.TxTypeP2P,
			Sender: mwjson.Participant{
				ID:       "265991234567",
				IDType:   mwjson.IDTypeMSISDN,
				Provider: mwjson.ProviderAirtelMoney,
			},
			Receiver: mwjson.Participant{
				ID:       "265881234567",
				IDType:   mwjson.IDTypeMSISDN,
				Provider: mwjson.ProviderTNMPamba,
			},
		},
	}

	// Sign
	if err := tx.SignTransaction(privKey); err != nil {
		t.Fatalf("Failed to sign transaction: %v", err)
	}

	if tx.TrustLayer.Signature == "" {
		t.Fatal("Signature is empty after signing")
	}

	// Verify
	if err := tx.VerifySignature(pubKey); err != nil {
		t.Errorf("Signature verification failed: %v", err)
	}

	// Tamper
	tx.Payload.Amount = 6000.00
	if err := tx.VerifySignature(pubKey); err == nil {
		t.Error("Expected verification failure after tampering, got nil")
	}
}
