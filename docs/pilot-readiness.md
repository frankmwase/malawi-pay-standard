# Pilot readiness (prototype status)

**Not ready for live payments or public Internet exposure.** The current repository is suitable for local development and synthetic-data demonstrations only. Do not put real phone numbers, account numbers, PINs or signing seeds in the repository, logs, screenshots or example payloads.

## What works

- Go MW-JSON structure, decimal amount validation, bounded TTL and Ed25519 signing/verification of version, header and payload.
- Go UMQR experimental TLV encoder and CRC test vector; CRC detects errors, not fraud.
- Single-node JSON ALS lookup with signed responses, authenticated registration when enabled, bounded request body, atomic private-file persistence and default loopback listening.
- Go package tests, race detector and vet in CI.

## Before any controlled pilot

1. Agree on a canonical wire format, version negotiation and signed test vectors for every supported language. Specify identity-to-key binding, rotation, revocation, signature requirements and replay prevention with a durable idempotency store.
2. Implement QR parsing with strict nested TLV validation, 2-digit byte-length bounds, CRC checks, field limits, amount validation and authenticated merchant/provider lookups; compare against independent EMVCo/partner conformance vectors. The current Go generator can produce invalid payloads for oversized fields.
3. Build and test the actual mobile camera/native integration. Remove or replace the unused payment screen and demo-only components. Do not collect PINs in JavaScript; obtain a security review and provider authorization before attempting telephony/USSD execution.
4. Verify ownership and consent before registering a destination. A bearer token identifies an operator, **not** the owner of an alias. Implement authenticated ownership changes, audit logging, rate limiting, abuse handling and anti-enumeration protections. The blind-token prototype has no redemption API; private registration is disabled.
5. Replace the JSON file with a backed-up transactional store if multiple writers/nodes are needed; document backup/restore, key rotation, health checks, monitoring, data retention and incident response. The sample blockchain resolver and smart contract are not integrated into the running HTTP server.
6. Deploy behind authenticated TLS termination and a firewall. Keep Besu JSON-RPC on loopback/private network, and do not turn on ADMIN/DEBUG for untrusted clients. Configure ingress rate limits, body limits and production timeouts. Provision `MW_ALS_SIGNING_SEED` and `MW_ALS_REGISTRATION_TOKEN` from a secrets manager; never pass them via CLI flags. Avoid recording sensitive account destinations in application logs.
7. Complete integration tests with authorized partners using test credentials and synthetic data, including reconciliation, ambiguous USSD outcomes, duplicate payments, lost connectivity, timeout retries and manual recovery. Obtain legal/regulatory and accessibility review before accepting funds.

Go signing changed to cover all payment fields: old prototype signatures are intentionally **incompatible**; re-sign existing test messages. The Go SDK is not a payment processor and a valid signature alone does not establish KYC or authorization.
