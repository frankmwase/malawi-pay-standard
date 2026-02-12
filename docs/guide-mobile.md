# Mobile USSD Router Guide

The Malawi Pay Standard includes a reference implementation for a **Mobile USSD Router**. This allows Android devices to act as automated gateways for Airtel Money and TNM Mpamba.

## How it Works
1. **Request**: The app receives an MW-JSON payment request (via QR or Deep Link).
2. **Routing**: The app identifies the provider (e.g., Airtel) and selects the correct SIM card.
3. **Execution**: The app uses the **Accessibility Service** to dial USSD codes (e.g., `*211#`) and navigate through the menus automatically.

## Native Integration
The core logic is implemented in Java/Kotlin to provide high-performance telephony access:

```java
// Example of dialing a carrier-specific USSD code
mTelephonyManager.sendUssdRequest("*211#", callback, handler);
```

## Security
- **Sandbox**: The USSD router only executes codes authorized by the MW-Standard middleware.
- **Privacy**: No PINs are stored locally; they are entered by the user via a secure native overlay.
