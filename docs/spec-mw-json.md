# MW-JSON Standard (v1.0)

MW-JSON is an experimental Go transaction data model. No institution has certified this schema or a network protocol around it.

## Core Principles
- **Idempotency**: A key is required by `Validate()`, but deduplication and replay protection must be implemented by the payment processor.
- **Normalization**: Call `NormalizeMSISDN()` explicitly; `Validate()` rejects non-normalized MSISDNs and does not mutate the transaction.
- **Security**: Signing and verification are separate Go SDK operations; `Validate()` alone does not authenticate a transaction. Signatures cover version, header, and payload using Go `encoding/json` struct serialization. Cross-language canonicalization and key discovery are unspecified.

## Data Structure

```json
{
  "mw_version": "1.0",
  "header": {
    "msg_id": "9988776655",
    "timestamp": "2026-02-12T20:00:00Z",
    "ttl": 300,
    "idempotency_key": "unique-uuid-here"
  },
  "payload": {
    "amount": 5000.00,
    "currency": "MWK",
    "type": "C2B",
    "sender": { "id": "265881234567", "id_type": "MSISDN", "provider": "TNM_MPAMBA", "alias": "@john" },
    "receiver": { "id": "265991122334", "id_type": "MSISDN", "provider": "AIRTEL_MONEY", "alias": "@mubas_cafe" }
  },
  "trust_layer": {
    "integrity_hash": "",
    "kyc_verified": false,
    "extension_signature": "ed25519-signature-in-hex"
  }
}
```

## Error Codes
| Code | Meaning | Context |
|------|---------|---------|
| `MW400` | Schema validation | Invalid field, including MSISDN format |
| `MW401` | Invalid signature | Verification failed |
| `MW408` | Expired transaction | TTL exceeded |
| `MW409` | Duplicate transaction | Reserved for processor-side idempotency handling |

The SDK does not verify KYC flags, integrity hashes or ownership of payment endpoints. Use the named key and a durable replay store in any processor.
