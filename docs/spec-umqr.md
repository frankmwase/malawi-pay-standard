# Universal Malawi QR (UMQR)

The **Universal Malawi QR** is an EMVCo-compatible standard designed to represent payment requests in a way that any Malawian banking app can understand.

## Encoding Format
UMQR uses a Tag-Length-Value (TLV) format with a checksum at the end.

### Mandatory Tags
- **00**: Payload Format Indicator (Fixed `01`)
- **01**: Point of Initiation Method (`11` for Static, `12` for Dynamic)
- **26**: Merchant Account Information (The Malawi Interop Sub-tags)
- **54**: Transaction Amount
- **58**: Country Code (Fixed `MW`)
- **63**: CRC16 Checksum

## Examples

### Static Merchant QR
```text
00020101021126440012@mubas_cafe0112AIRTEL_MONEY52040000530345454072500.005802MW63045E2A
```

### QR Generation (Go)
```go
import "github.com/frankmwase/malawi-pay-standard/pkg/umqr"

qr := umqr.GenerateMerchantQR(
    "MUBAS Cafe", 
    "Blantyre", 
    "@mubas_cafe", 
    "AIRTEL_MONEY", 
    2500.00, 
    "LUNCH-45",
)
```
