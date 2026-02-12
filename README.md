# MW-Standard: The Open Foundation for Malawian Digital Exchange 🇲🇼

> "Building the plumbing so Malawi can build the future."

In Malawi, our digital economy is trapped in "Walled Gardens." Airtel Money, TNM Mpamba, and our commercial banks often don't speak the same language. This forces developers to write redundant code, merchants to display five different QR codes, and users to pay high fees for "offnet" transfers.

**MW-Standard** is an open source initiative to build the missing foundations of the Malawian digital ecosystem. Inspired by global standards like India’s UPI and Singapore’s PayNow, we are building the primitives for identity, payments, and discovery.

##  The Three Pillars

### 1. MW-JSON (The Language)
A standardized, lightweight JSON schema for transactions. It abstracts the complexity of different providers into a single object.
- **Interoperable**: Works across all wallets and banks.
- **Resilient**: Includes TTL  and Idempotency keys for Malawi’s spotty network conditions.
- **Secure**: Features a built in TrustLayer for cryptographic signatures using Ed25519.

### 2. UMQR (The Interface)
The Universal Malawian QR Code. Based on EMVCo standards, UMQR allows a single sticker to accept payments from any compliant app. No more "QR clutter" on merchant desks.

### 3. MW-ALS (The Discovery)
The Alias Lookup Service. A decentralized "DNS for Money."
- **Resolve @alias** (e.g., `@chifundo`) to a bank account or phone number.
- **Privacy-first**: Resolves to an endpoint without exposing full personal details to the sender.

## Tech Stack
- **Language**: Go (Golang)  chosen for its performance, concurrency, and tiny binary size.
- **Encoding**: JSON (Standard) & Protobuf (for low bandwidth USSD/GPRS).
- **Security**: Ed25519 for transaction signing.

##  Installation (For Developers)
To start using the standard in your Go project:

```bash
go get github.com/frankmwase/malawi-pay-standard
```

### Quick Example: Create a Standard Transaction
```go
import (
    "time"
    "github.com/frankmwase/malawi-pay-standard/pkg/mwjson"
)

func main() {
    // Create a new transaction (using the schema)
    txn := &mwjson.Transaction{
        MWVersion: "1.0",
        Header: mwjson.Header{
            MsgID: "TXN-123",
            Timestamp: time.Now().UTC(),
            TTL: 300,
            IdempotencyKey: "unique-key",
        },
        Payload: mwjson.Payload{
            Amount: 15000.00,
            Currency: "MWK",
            Type: mwjson.TxTypeP2P,
            Sender: mwjson.Participant{
                ID: "26599...",
                IDType: mwjson.IDTypeMSISDN,
                Provider: mwjson.ProviderAirtelMoney,
            },
            Receiver: mwjson.Participant{
                ID: "26588...",
                IDType: mwjson.IDTypeMSISDN,
                Provider: mwjson.ProviderTNMPamba,
            },
        },
    }

    // Validate the transaction against Malawian regulations
    if err := txn.Validate(); err != nil {
        log.Fatal("Invalid Transaction: ", err)
    }
}
```

##  The University Pilot
We are focusing initial adoption on Malawian Universities (MUBAS, MUST, UNIMA, MZUNI). By deploying these standards on campus intranets, we create a high-trust laboratory where students can build apps that interact with local campus economies without needing expensive internet data.

##  How to Contribute
We aren't just looking for code; we are looking for Founders.
1. **Review the Specs**: Check the code for MW-JSON 1.0 and UMQR 1.0 specifications.
2. **Build a Driver**: Help us write the "Adapter" for different Malawian banks.
3. **Optimize**: Help us make the binary encoding smaller for USSD systems.
4. **Report**: Open an issue if you find an edge case in how we handle Kwacha denominations or local IDs.

##  Roadmap
- [x] **Alpha**: MW-JSON Schema & Go SDK Core.
- [x] **Beta**: UMQR Encoder/Decoder & Example Walkthrough.
- [x] **Gamma**: Prototype Alias Lookup Service (ALS).
- [ ] **V1.0**: National Interoperability Framework Proposal.

> "If you want to go fast, go alone. If you want to go far, go together." 

This project is open-source and always will be. It belongs to the Malawian developer community.
