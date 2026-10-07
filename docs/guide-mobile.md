# Mobile USSD Router Guide

The mobile app is a **non-payment UI prototype**. Its camera and payment paths are disabled; there is no Android/iOS native telephony module, secure PIN overlay, SIM detection, production QR parser, or payment middleware. Never enter a real PIN or attempt a live transaction with this app.

The Go package `mwussd` generates **hypothetical** Airtel/TNM menu steps. Carrier menus and dialing codes are not verified, and generated steps must not be executed automatically. It rejects fractional MWK amounts rather than silently rounding them. Real deployment needs agreements with providers, device/security review, user consent, tested menu flows, failure handling and authorization. See the pilot readiness checklist in the README.
