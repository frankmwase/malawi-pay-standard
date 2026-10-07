package mwjson

import (
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
)

// signingBytes serializes every payment-relevant field in a fixed Go struct order.
// This is a Go SDK signing format, not a cross-language JSON canonicalization standard.
func (t *Transaction) signingBytes() ([]byte, error) {
	return json.Marshal(struct {
		Version string  `json:"mw_version"`
		Header  Header  `json:"header"`
		Payload Payload `json:"payload"`
	}{t.MWVersion, t.Header, t.Payload})
}

// SignTransaction signs the version, header and complete payload (not the mutable trust layer).
func (t *Transaction) SignTransaction(privateKey ed25519.PrivateKey) error {
	if len(privateKey) != ed25519.PrivateKeySize {
		return NewMWError(ErrInvalidSignature, "Invalid private key", "")
	}
	data, err := t.signingBytes()
	if err != nil {
		return err
	}
	t.TrustLayer.Signature = hex.EncodeToString(ed25519.Sign(privateKey, data))
	return nil
}

// VerifySignature checks the complete signed transaction against the supplied public key.
func (t *Transaction) VerifySignature(publicKey ed25519.PublicKey) error {
	if len(publicKey) != ed25519.PublicKeySize {
		return NewMWError(ErrInvalidSignature, "Invalid public key", "")
	}
	sig, err := hex.DecodeString(t.TrustLayer.Signature)
	if err != nil || len(sig) != ed25519.SignatureSize {
		return NewMWError(ErrInvalidSignature, "Invalid signature format", "")
	}
	data, err := t.signingBytes()
	if err != nil {
		return err
	}
	if !ed25519.Verify(publicKey, data, sig) {
		return NewMWError(ErrInvalidSignature, "Signature verification failed", "")
	}
	return nil
}
