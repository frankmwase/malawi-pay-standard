# Getting Started

This guide will help you begin integrating the Malawi Pay Standard into your Malawian financial product.

## Prerequisites
- **Go** (1.24+ recommended)
- **Git**

## 1. Install the SDK
Add the standard as a dependency to your Go project:

```bash
go get github.com/frankmwase/malawi-pay-standard
```

## 2. Generate a Transaction
Construct a standard compliant transaction signed with your institution's private key.

```go
tx := &mwjson.Transaction{
    // ... setup transaction fields
}
err := tx.SignTransaction(privKey)
```

## 3. Deployment (Institutional Nodes)
If you are a Malawian bank or MNO, you should run your own ALS node to provide resolution data to your customers.

```bash
go run cmd/mwals/main.go --port 8080 --data my_data.json
```

---
> [!TIP]
> Always enforce UTC timestamps to avoid drift across Malawian mobile network towers.
