package mwjson

import (
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"time"
)

// SignTransaction generates a signature for the transaction using the sender's private key.
// It populates the TrustLayer.Signature field.
func (t *Transaction) SignTransaction(privateKey ed25519.PrivateKey) error {
	
	// 1. Create the canonical string to sign
	// We need to sign the immutable parts: Header and Payload.
	// We exclude TrustLayer itself to avoid recursion, though IntegrityHash is part of it.
	// A common pattern is to sign the hash of the payload + header.

	// Concatenate Header + Payload
	// In a real spec, we'd need a very strict canonicalization (e.g., JCS).
	// For this standard, we'll assume the byte buffer matches what is sent,
	// OR we sign a constructed string like "msg_id|timestamp|amount|sender|receiver" to be safe against JSON formatting issues.
	// Let's go with the constructed string approach for robustness in this MVP.

	canonicalString := fmt.Sprintf("%s|%s|%.2f|%s|%s",
		t.Header.MsgID,
		t.Header.Timestamp.UTC().Format(time.RFC3339),
		t.Payload.Amount,
		t.Payload.Sender.ID,
		t.Payload.Receiver.ID,
	)

	// 2. Sign
	signature := ed25519.Sign(privateKey, []byte(canonicalString))

	// 3. Encode to Hex and set
	t.TrustLayer.Signature = hex.EncodeToString(signature)

	// Can also set integrity hash (SHA256 of the canonical string or payload) .... good for "TrustLayer"
	// t.TrustLayer.IntegrityHash = ... (Skipping for now as Signature covers it)

	return nil
}

// VerifySignature checks if the transaction signature is valid for the given public key.
func (t *Transaction) VerifySignature(publicKey ed25519.PublicKey) error {
	if t.TrustLayer.Signature == "" {
		return NewMWError(ErrInvalidSignature, "Missing Signature", "")
	}

	// 1. Reconstruct Canonical String
	canonicalString := fmt.Sprintf("%s|%s|%.2f|%s|%s",
		t.Header.MsgID,
		t.Header.Timestamp.UTC().Format(time.RFC3339),
		t.Payload.Amount,
		t.Payload.Sender.ID,
		t.Payload.Receiver.ID,
	)

	// 2. Decode Signature
	sigBytes, err := hex.DecodeString(t.TrustLayer.Signature)
	if err != nil {
		return NewMWError(ErrInvalidSignature, "Invalid Signature Format", "Not Hex")
	}

	if len(sigBytes) != ed25519.SignatureSize {
		return NewMWError(ErrInvalidSignature, "Invalid Signature Length", "")
	}

	// 3. Verify
	valid := ed25519.Verify(publicKey, []byte(canonicalString), sigBytes)
	if !valid {
		return NewMWError(ErrInvalidSignature, "Signature Verification Failed", "")
	}

	return nil
}
