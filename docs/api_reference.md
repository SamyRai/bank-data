# API Reference: Financial Package (v1 Public API Contract)

The `financial` package is the canonical facade for validating, parsing, and detecting all supported financial identifiers (IBAN, BIC, and SEPA Creditor ID).

## v1 API Stability

To ensure stability and predictability for downstream systems, `pkg/financial` adheres to the following versioning rules:

*   **Stable Scope**: The `IdentifierType`, `ValidationReport`, `ParsedIdentifier`, `Suggestion` types, and all exported methods on the `Service` interface are guaranteed to be stable.
*   **Internal Scope**: Anything under `internal/*` or generated code formats are considered unstable and subject to change.
*   **Additive Changes**: New fields to structs (like `ValidationReport` or `ParsedIdentifier`) will only be added in minor versions (e.g., v1.1.0, v1.2.0).
*   **No Breaking Changes**: There will be no breaking type changes or method signature modifications within the `v1.x` series.
*   **Deprecation Policy**: Compatibility packages (`pkg/iban`, `pkg/bic`, `pkg/sepa`) are now deprecated. They will receive critical bug fixes but no new features. They will be removed entirely in `v2.0.0` (estimated mid-2027). Users must migrate to `pkg/financial`.

## Migration Matrix

To migrate from legacy packages to the unified `pkg/financial` package, update your imports and types as follows:

| Legacy Package / Method | v1 Canonical Equivalent (`pkg/financial`) | Notes |
| :--- | :--- | :--- |
| `iban.NewService()` | `financial.NewService()` | Instantiates the combined service for all types. |
| `iban.Service.Validate()` | `financial.Service.Validate()` | Returns a comprehensive `ValidationReport` instead of just an `error`. |
| `iban.Service.Parse()` | `financial.Service.Parse()` | Returns a `ParsedIdentifier` containing a `Components` map. |
| `iban.Service.Detect()` | `financial.Service.Detect()` | Now returns a structured `IdentifierType` constant. |
| `bic.Validator.Validate()` | `financial.Service.Validate()` | Same method covers BIC validation. |
| `sepa.CreditorIDValidator.Validate()` | `financial.Service.Validate()` | Same method covers SEPA validation. |

---

## The Financial Service

The `Service` is the central entry point. For optimal performance, initialize it once and reuse it across your application.

```go
import "github.com/SamyRai/bank-data/pkg/financial"

// Initialize the universal service
svc := financial.NewService()
```

### Methods

#### `Validate(input string) ValidationReport`
Validates a string against all supported financial types. It returns a `ValidationReport` containing the determined type, normalization, and any error encountered.

#### `Parse(input string) (*ParsedIdentifier, error)`
Deconstructs a normalized identifier into its canonical components (Country Code, Bank Code, Account Number, etc., depending on the detected type).

#### `Detect(input string) IdentifierType`
Retrieves the specific `IdentifierType` (e.g., `TypeIBAN`, `TypeBIC`, `TypeSEPACreditorID`) of the given string based on format and structure.

---

## Errors and Status Codes

When a validation fails, the `ValidationReport.Error` field will be populated. The underlying error types depend on the detected identifier type (e.g., `*iban.IBANError` for IBANs).

### Idiomatic Error Handling
```go
report := svc.Validate(input)
if !report.Valid {
    // Check if it's specifically an IBAN error
    var ibanErr *iban.IBANError
    if errors.As(report.Error, &ibanErr) {
        fmt.Printf("IBAN Validation failed: %s (Code: %s)\n", ibanErr.Message, ibanErr.Code)
    } else {
        fmt.Printf("Validation failed: %v\n", report.Error)
    }
}
```
