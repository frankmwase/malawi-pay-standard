package mwjson_test

import (
	"crypto/ed25519"
	"crypto/rand"
	"testing"
	"time"

	"github.com/frankmwase/malawi-pay-standard/pkg/mwjson"
	"github.com/shopspring/decimal"
)

func TestSignedFieldsAndTimestamp(t *testing.T) {
	pub, key, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	tx := &mwjson.Transaction{
		MWVersion: mwjson.MWJSONVersion,
		Header: mwjson.Header{
			MsgID: "id", Timestamp: time.Now().UTC(), TTL: 300, IdempotencyKey: "unique",
		},
		Payload: mwjson.Payload{
			Amount: decimal.NewFromInt(100), Currency: mwjson.CurrencyMWK, Type: mwjson.TxTypeP2P,
			Sender:   mwjson.Participant{ID: "265991234567", IDType: mwjson.IDTypeMSISDN, Provider: mwjson.ProviderAirtelMoney},
			Receiver: mwjson.Participant{ID: "265881234567", IDType: mwjson.IDTypeMSISDN, Provider: mwjson.ProviderTNMPamba},
		},
	}
	if err := tx.SignTransaction(key); err != nil {
		t.Fatal(err)
	}
	if err := tx.VerifySignature(pub); err != nil {
		t.Fatal(err)
	}
	for name, tamper := range map[string]func(*mwjson.Transaction){
		"currency":    func(tx *mwjson.Transaction) { tx.Payload.Currency = "USD" },
		"idempotency": func(tx *mwjson.Transaction) { tx.Header.IdempotencyKey = "other" },
		"provider":    func(tx *mwjson.Transaction) { tx.Payload.Receiver.Provider = mwjson.ProviderFDH },
	} {
		t.Run(name, func(t *testing.T) {
			changed := *tx
			tamper(&changed)
			if err := changed.VerifySignature(pub); err == nil {
				t.Fatal("tampered transaction verified")
			}
		})
	}
	if err := tx.SignTransaction(nil); err == nil {
		t.Fatal("nil key accepted")
	}
	tx.Header.Timestamp = time.Now().UTC().Add(5 * time.Minute)
	if err := tx.Validate(); err == nil {
		t.Fatal("future timestamp accepted")
	}
	tx.Header.Timestamp = time.Now().UTC()
	tx.Header.TTL = int(^uint(0) >> 1)
	if err := tx.Validate(); err == nil {
		t.Fatal("excessive TTL accepted")
	}
}
