# MW-JSON Standard (v1.0)

MW-JSON is the official transaction data model for the Malawi Pay Standard. It ensures that transaction data is consistent, signed, and verifiable across all Malawian nodes.

## Core Principles
- **Idempotency**: Every request must be idempotent to prevent double-charging on Malawian mobile networks.
- **Normalization**: Phone numbers (MSISDNs) are automatically normalized to the `265...` format.
- **Security**: Mandatory Ed25519 cryptographic signing for every transaction.

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
    "sender": { "id": "265881234567", "alias": "@john" },
    "receiver": { "id": "265991122334", "alias": "@mubas_cafe" }
  },
  "trust_layer": {
    "signature": "...",
    "pub_key": "..."
  }
}
```

## Error Codes
| Code | Meaning | Context |
|------|---------|---------|
| `MW001` | Insufficient Funds | Sender has less than `Amount` |
| `MW002` | MSISDN Invalid | Phone number format error |
| `MW003` | Signature Mismatch | Trust Layer verification failed |
| `MW004` | TTL Expired | Transaction sent too long ago |
