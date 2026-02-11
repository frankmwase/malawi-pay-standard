# Contributing to MW-Standard 🇲🇼

Welcome, Builder.

Thank you for your interest in contributing to the **Open Standard for Malawian Digital Exchange**. By contributing to this project, you aren't just writing Go code; you are building the digital "plumbing" that will power the next generation of Malawian fintech, e-commerce, and digital identity.

We are looking for contributions from everyone: students, professional engineers, fintech veterans, and enthusiasts.

##  The Golden Rules
1. **It's a Standard, not just an App**: Code must be backward compatible. We cannot break the API once banks or apps start using it.
2. **Local Context Matters**: Always consider the Malawian spotty internet, USSD limitations, and low priority GPRS data.
3. **Trust is Code**: Security (Ed25519) and strict validation are not optional.

## 🛠 Getting Started

### Prerequisites
- **Go**: Version 1.21 or later.
- **Git**: Basic knowledge of branching and PRs.
- **Knowledge**: Familiarity with EMVCo QR standards and JSON/Protobuf.

### Setup
1. Fork the repository on GitHub.
2. Clone your fork locally:
   ```bash
   git clone https://github.com/frankmwase/malawi-pay-standard.git
   cd malawi-pay-standard
   ```
3. Install Dependencies:
   ```bash
   go mod download
   ```

##  Development Workflow

### 1. Branching Strategy
Create a feature branch with a descriptive name:
```bash
git checkout -b feat/add-tnm-validator
```

### 2. Coding Standards (The "Go Way")
We strictly follow idiomatic Go patterns.
- **Formatting**: Run `go fmt ./...` before committing.
- **Error Handling**: Use our custom `MWError` types (defined in `pkg/mwjson/errors.go`) so API consumers know exactly why a transaction failed.
- **Comments**: Exported functions must have comments explaining their purpose.

```go
// NormalizeMSISDN converts various definitions to 265XXXXXXXXX format.
// Inputs: 099..., 99..., +26599...
func NormalizeMSISDN(input string) (string, error) { ... }
```

### 3. Malawian Context Guidelines
- **Precision**: While we use `float64` for the API, always validate precision to 2 decimal places (Tambala) as seen in `pkg/mwjson/validation.go`.
- **Time**: **Force UTC**. All timestamps must be in UTC to avoid server-sync issues.
- **Connectivity**: Assume the network will fail. Use the `Idempotency-Key` logic in `Header` to prevent duplicate processing.

## 🧪 Testing
"If it isn't tested, it doesn't exist."

We use **Table-Driven Tests** for all validation and encoding logic.
Example from `pkg/mwals/als_test.go`:

```go
func TestNormalizer(t *testing.T) {
    tests := []struct {
        input    string
        expected string
    }{
        {"@Koda_Dev", "koda_dev"},
        {"  @Alice  ", "alice"},
    }

    for _, tt := range tests {
        got := mwals.Normalizer(tt.input)
        if got != tt.expected {
            t.Errorf("Normalizer(%s) = %s; want %s", tt.input, got, tt.expected)
        }
    }
}
```

Run all tests before submitting:
```bash
go test ./...
```

##  Submitting a Pull Request (PR)
- **Commit Messages**: Use Conventional Commits (e.g., `feat: allow 12-digit TNM numbers`, `fix: resolve panic on nil header`).
- **PR Template**: Explain *what* changed, *why* it's needed, and *how* you tested it.

##  Community Code of Conduct
We are building for Malawi. We treat each other with respect.
- **Be Inclusive**: We welcome code from a 1st-year student at MUST/MUBAS just as warmly as a Senior Engineer at a bank.
- **Constructive Feedback**: Critique the code, not the person.

## Need Help?
- **Join the Discussion**: princefranklinemwase@gmail.com
- **Read the Implementation**: Explore `examples/campus_connect` for the reference flow.

**Let's build. 🇲🇼**
