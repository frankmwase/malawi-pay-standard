# Malawi Alias Lookup Service (MW-ALS)

The MW-ALS acts as the "Discovery" layer. It translates human-readable aliases into actionable financial endpoints.

## The Hybrid Model
MW-ALS uses a **Hybrid Trust Model**:
1. **Source of Truth**: A Hyperledger Besu blockchain (Ethereum-compatible) managed by Malawian banks.
2. **Local Cache**: JSON-based storage on local server nodes for performance and offline reliability.

## Resolving an Alias
Send a `GET` request to any certified ALS node:
`GET /resolve/@chifundo`

### Response (Blind Token Mode)
To protect privacy, the ALS returns a **Blind Token** for private endpoints instead of a phone number:

```json
{
  "alias": "@chifundo",
  "status": "ACTIVE",
  "identity_mask": "C*** F***",
  "endpoints": [
    {
      "priority": 1,
      "provider": "AIRTEL",
      "destination": "TOKEN:998877...:2026-02-12T20:30:00Z"
    }
  ],
  "security_sig": "..."
}
```

## Registering an Alias
A `POST /register` request with an identity certificate (e.g., NRIS hash) is required to claim an alias.
