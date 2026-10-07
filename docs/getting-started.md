# Getting Started

This guide will help you begin integrating the Malawi Pay Standard into your Malawian financial product.

## Prerequisites
- **Go** 1.24.1 or later (per the module directive)
- **Git**

## 1. Install the SDK
Add the standard as a dependency to your Go project:

```bash
go get github.com/frankmwase/malawi-pay-standard
```

## 2. Generate a Transaction
See the compilable transaction example in the project README. Validate the message with `Validate()` before signing it using `SignTransaction(privateKey)`. The current Go SDK signature format is not a cross-language JSON canonicalization protocol. It does **not** supply key management, recipient verification, or idempotent execution.

## 3. Deployment (Institutional Nodes)
The ALS server is a **local development prototype**, not an institutional node. Provision a persistent Ed25519 seed in the `MW_ALS_SIGNING_SEED` environment variable (64 hex characters) using a secret manager; do not put secrets in shell history or repository files. Start with `go run ./cmd/mwals --listen 127.0.0.1:8080 --data als_data.json`. It binds to loopback by default. Registration is disabled unless `MW_ALS_REGISTRATION_TOKEN` is set; authorized clients must send `Authorization: Bearer <token>`. Never expose this HTTP endpoint or the Besu RPC to an untrusted network without TLS, authentication at the proxy, rate limits, and operational controls. The JSON file is a single-node store, not a replicated database.

---
> [!TIP]
> Always enforce UTC timestamps to avoid drift across Malawian mobile network towers.
