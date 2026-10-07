# Universal Malawi QR (UMQR)

**Universal Malawi QR** is an experimental EMV-style TLV/CRC encoding proposal. It has not passed EMVCo certification or interoperability testing; do not accept its output as an authenticated payment request.

## Encoding Format
UMQR uses a Tag-Length-Value (TLV) format with a checksum at the end.

### Mandatory Tags
- **00**: Payload Format Indicator (Fixed `01`)
- **01**: Point of Initiation Method (`11` for Static, `12` for Dynamic)
- **26**: Merchant Account Information (The Malawi Interop Sub-tags)
- **53**: Transaction Currency (`454` for MWK)
- **58**: Country Code (`MW`)
- **59**, **60**: Merchant Name and City
- **63**: CRC16-CCITT checksum (error detection, **not** authentication)

The current Go encoder uses tag **26** with nested global ID, account type and alias. Tag **54** (amount) is optional for static QRs and emitted only for positive values. The encoder does not reject values over the 99-byte TLV length limit, so validate fields before use. A corresponding decoder does not yet exist.

## Examples

### Static Merchant QR
Generate a test payload with the Go example below. No fixed sample is published here until a decoder and independent conformance vectors are available.

### QR Generation (Go)
```go
import (
    "github.com/frankmwase/malawi-pay-standard/pkg/umqr"
    "github.com/shopspring/decimal"
)

qr := umqr.GenerateMerchantQR("MUBAS Cafe", "Blantyre", "@mubas_cafe", "AIRTEL_MONEY", decimal.NewFromInt(2500), "LUNCH-45")
_ = qr // demonstration only; do not use for real payment acceptance
```
