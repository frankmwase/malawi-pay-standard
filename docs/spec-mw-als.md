# Malawi Alias Lookup Service (MW-ALS)

The MW-ALS prototype resolves human-readable aliases from a single local JSON store. It is not certified, federated or connected to a bank-operated blockchain. `HybridService` and the Besu contract are separate experimental components, not wired into the running HTTP server.

## Resolving an Alias
A loopback development node supports `GET /resolve/@chifundo`. Its Ed25519 response signature covers all returned fields, but clients still need a trusted public key. Endpoints registered by an operator are **not verified as owned by the named account holder**.

### Legacy prototype blind-token illustration
Manually seeded private records can return opaque tokens; there is no redemption API. New private registrations are disabled, so these tokens must not be used for payment routing:

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
`POST /register` is disabled unless `MW_ALS_REGISTRATION_TOKEN` is configured. Then send an operator bearer token; the request is limited to 16 KiB and aliases must match 3–32 lowercase letters, numbers or underscores (an optional `@` is accepted). This is **operator authentication only**, not identity certification or proof of destination ownership. Never permit untrusted callers to register payment destinations. See [pilot readiness](pilot-readiness.md).
