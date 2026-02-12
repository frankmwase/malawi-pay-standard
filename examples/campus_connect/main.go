package main

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"fmt"
	"log"
	"time"

	"github.com/frankmwase/malawi-pay-standard/pkg/mwals"
	"github.com/frankmwase/malawi-pay-standard/pkg/mwjson"
	"github.com/frankmwase/malawi-pay-standard/pkg/umqr"
)

func main() {
	fmt.Println("=== Campus-Connect: University Value Transfer Example ===")

	// 1. Setup the Alias Lookup Service (ALS)
	_, alsKey, _ := ed25519.GenerateKey(rand.Reader)
	als, _ := mwals.NewService(alsKey, "")
	als.Seed(&mwals.AliasRecord{
		Alias:             "mubas_cafe",
		Status:            mwals.AliasStatusActive,
		IdentityMask:      "M**** C********",
		VerificationProof: "123456",
		Endpoints: []mwals.Endpoint{
			{Priority: 1, Provider: "AIRTEL_MONEY", Type: mwals.EndpointTypeWallet, Destination: "265991112223"},
		},
	})

	// 2. Merchant (Cafe) Generates a UMQR for a student to scan
	fmt.Println("\n[Cafe] Generating Dynamic UMQR for Lunch...")
	lunchAmount := 2500.00
	qr := umqr.GenerateMerchantQR("MUBAS Cafeteria", "Blantyre", "@mubas_cafe", "AIRTEL_MONEY", lunchAmount, "LUNCH-45")
	fmt.Printf("Produced UMQR: %s\n", qr)

	// 3. Student scans QR and resolves the alias to find where to send money
	fmt.Println("\n[Student App] Resolving @mubas_cafe via ALS...")
	res, err := als.Resolve(context.Background(), "@mubas_cafe")
	if err != nil {
		log.Fatalf("ALS Error: %v", err)
	}
	fmt.Printf("Resolved to: %s (Provider: %s)\n", res.IdentityMask, res.Endpoints[0].Provider)

	// 4. Student App constructs a MW-JSON Transaction
	fmt.Println("\n[Student App] Constructing MW-JSON Transaction...")
	pubKey, privKey, _ := ed25519.GenerateKey(rand.Reader) // In reality, keys are stored on device

	tx := &mwjson.Transaction{
		MWVersion: mwjson.MWJSONVersion,
		Header: mwjson.Header{
			MsgID:          "TXN-STU-99",
			Timestamp:      time.Now().UTC(),
			TTL:            300,
			IdempotencyKey: "unique-stu-12345",
		},
		Payload: mwjson.Payload{
			Amount:   lunchAmount,
			Currency: mwjson.CurrencyMWK,
			Type:     mwjson.TxTypeC2B,
			Sender: mwjson.Participant{
				ID:       "265991234567",
				IDType:   mwjson.IDTypeMSISDN,
				Provider: mwjson.ProviderAirtelMoney,
				Alias:    "@student_john",
			},
			Receiver: mwjson.Participant{
				ID:       res.Endpoints[0].Destination,
				IDType:   mwjson.IDTypeMSISDN,
				Provider: mwjson.Provider(res.Endpoints[0].Provider),
				Alias:    res.Alias,
			},
		},
	}

	// 5. Student signs the transaction
	err = tx.SignTransaction(privKey)
	if err != nil {
		log.Fatalf("Signing Error: %v", err)
	}
	fmt.Printf("Transaction Signed. Signature: %s...\n", tx.TrustLayer.Signature[:16])

	// 6. Verification by the Gateway/Provider
	fmt.Println("\n[Gateway] Verifying Standard Compliance...")
	if err := tx.Validate(); err != nil {
		log.Fatalf("Validation Failed: %v", err)
	}
	fmt.Println("Validation SUCCESS.")

	if err := tx.VerifySignature(pubKey); err != nil {
		log.Fatalf("Signature Invalid: %v", err)
	}
	fmt.Println("Signature SUCCESS.")

	fmt.Println("\n=== Transaction Complete: Lunch Paid! ===")
}
