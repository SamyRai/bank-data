# API Reference: Financial Package

`pkg/financial` is the canonical API for identifier detection, validation, and parsing.

## Service

```go
import "github.com/SamyRai/bank-data/pkg/financial"

svc := financial.NewService()
```

### `Detect(input string) (IdentifierType, error)`
Auto-detects identifier type by running deterministic candidate checks and checksum validation.

### `Validate(input string, hint IdentifierType) (ValidationReport, error)`
Validates `input` as `hint`. If `hint == ""`, type is auto-detected first.

### `Parse(input string, hint IdentifierType) (ParsedIdentifier, error)`
Validates and returns parsed fields (`map[string]string`) for the chosen identifier type.

### `ValidateBatch(ctx, inputs, hint)` and `StreamValidate(ctx, seq, hint)`
Concurrent validation helpers built on the shared core engine.

## Types

### `IdentifierType`
- `IdentifierIBAN`
- `IdentifierBIC`
- `IdentifierSEPACreditor`
- `IdentifierLEI`
- `IdentifierISIN`
- `IdentifierPAN`
- `IdentifierVAT`

### `ValidationReport`
- `Type IdentifierType`
- `Input string`
- `Normalized string`
- `Valid bool`
- `Error *ValidationError`

### `ParsedIdentifier`
- `Type IdentifierType`
- `Normalized string`
- `Fields map[string]string`

## Compatibility Packages

`pkg/iban`, `pkg/bic`, and `pkg/sepa` are retained as compatibility layers. New development should target `pkg/financial`.
