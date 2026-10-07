# MW-Standard: The Open Foundation for Malawian Digital Exchange 🇲🇼

> "Building the plumbing so Malawi can build the future."

In Malawi, our digital economy is trapped in "Walled Gardens." Airtel Money, TNM Mpamba, and our commercial banks often don't speak the same language. This forces developers to write redundant code, merchants to display five different QR codes, and users to pay high fees for "offnet" transfers.

**MW-Standard** is an open source initiative to build the missing foundations of the Malawian digital ecosystem. Inspired by global standards like India’s UPI and Singapore’s PayNow, we are building the primitives for identity, payments, and discovery.

##  The Three Pillars

### 1. MW-JSON (The Language)
A standardized, lightweight JSON schema for transactions. It abstracts the complexity of different providers into a single object.
- **Provider-neutral data model**: Describes wallet/bank transactions; provider adapters are not implemented.
- **Replay metadata**: Includes a bounded TTL and idempotency key; consumers must implement durable deduplication themselves.
- **Signed SDK messages**: Go Ed25519 helpers sign the version, header and full payload; key distribution and cross-language canonicalization are not specified.

### 2. UMQR (The Interface)
An experimental EMV-style TLV/CRC encoder for merchant QR payloads. It is **not certified EMVCo interoperability**: there is no validated decoder or end-to-end payment integration, and generated fields are not yet length-checked.

### 3. MW-ALS (The Discovery)
A local JSON-backed alias lookup prototype, not a decentralized production directory.
- Resolves aliases to stored endpoints and signs complete responses.
- Public registrations require operator authorization; **endpoint ownership is not verified**. Private token redemption, KYC and blockchain-backed verification are not implemented.

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
package main

import (
    "log"
    "time"

    "github.com/frankmwase/malawi-pay-standard/pkg/mwjson"
    "github.com/shopspring/decimal"
)

func main() {
    txn := &mwjson.Transaction{
        MWVersion: mwjson.MWJSONVersion,
        Header: mwjson.Header{
            MsgID: "TXN-123", Timestamp: time.Now().UTC(),
            TTL: 300, IdempotencyKey: "unique-key",
        },
        Payload: mwjson.Payload{
            Amount: decimal.NewFromInt(15000), Currency: mwjson.CurrencyMWK,
            Type: mwjson.TxTypeP2P,
            Sender: mwjson.Participant{ID: "265991234567", IDType: mwjson.IDTypeMSISDN, Provider: mwjson.ProviderAirtelMoney},
            Receiver: mwjson.Participant{ID: "265881234567", IDType: mwjson.IDTypeMSISDN, Provider: mwjson.ProviderTNMPamba},
        },
    }
    if err := txn.Validate(); err != nil {
        log.Fatal(err)
    }
    // SignTransaction requires a securely provisioned Ed25519 private key.
}
```

## Pilot status
This is a **prototype for synthetic-data demonstrations only**, not a usable payment network or a deployable university payment app. A campus intranet does not replace provider authorization, TLS or payment controls. See the [pilot readiness checklist](docs/pilot-readiness.md) for release blockers and operator guidance.

##  How to Contribute
We aren't just looking for code; we are looking for Founders.
1. **Review the Specs**: Check the code for MW-JSON 1.0 and UMQR 1.0 specifications.
2. **Build a Driver**: Help us write the "Adapter" for different Malawian banks.
3. **Optimize**: Help us make the binary encoding smaller for USSD systems.
4. **Report**: Open an issue if you find an edge case in how we handle Kwacha denominations or local IDs.

##  Roadmap
- [x] **Alpha**: MW-JSON Schema & Go SDK Core.
- [ ] **Beta**: Validated UMQR encoder/decoder and independently tested interoperability.
- [x] **Prototype**: Local alias lookup service (ALS).
- [ ] **Pilot readiness**: Provider authorization, QR camera and payment integration, security review, and operational runbooks.
- [ ] **V1.0**: National Interoperability Framework Proposal.

> "If you want to go fast, go alone. If you want to go far, go together." 

This project is open-source and always will be. It belongs to the Malawian developer community.
