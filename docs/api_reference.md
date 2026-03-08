# API Reference

This project exposes one canonical API and several compatibility packages.

## Canonical Public API: `pkg/financial`

### Construction

```go
import "github.com/SamyRai/bank-data/pkg/financial"

svc := financial.NewService()
```

### Core Methods

- `Detect(input string) (IdentifierType, error)`
- `Validate(input string, hint IdentifierType) (ValidationReport, error)`
- `Parse(input string, hint IdentifierType) (ParsedIdentifier, error)`
- `Suggest(input string, hint IdentifierType) ([]Suggestion, error)`
- `ValidateBatch(ctx context.Context, inputs []string, hint IdentifierType) []ValidationReport`
- `StreamValidate(ctx context.Context, inputs iter.Seq[string], hint IdentifierType) iter.Seq[ValidationReport]`

### Identifier Types

- `IdentifierIBAN`
- `IdentifierBIC`
- `IdentifierSEPACreditor`
- `IdentifierLEI`
- `IdentifierISIN`
- `IdentifierPAN`
- `IdentifierVAT`
- `IdentifierNationalAccountUK`

### Result Types

`ValidationReport`:
- `Type IdentifierType`
- `Input string`
- `Normalized string`
- `Valid bool`
- `Error *ValidationError`

`ParsedIdentifier`:
- `Type IdentifierType`
- `Normalized string`
- `Fields map[string]string`

## Compatibility Packages

Use these only for legacy integrations that are already coupled to old package names.

### `pkg/iban`
- `NewService(v Validator, p Parser, d Detector, bankLookup BankLookup) *Service`
- `Service.Validate / Parse / Detect / ValidateBatch / StreamValidate`

### `pkg/bic`
- `NewValidator() *Validator`
- `(*Validator).Validate(input string) ValidationResult`
- `Parse(bic string) (*BICInfo, error)`

### `pkg/sepa`
- `NewCreditorIDValidator() *CreditorIDValidator`
- `(*CreditorIDValidator).Validate(input string) ValidationResult`

### `pkg/validation`
Typed helper interfaces and shared validation result primitives.

## Additional Public Packages

### `pkg/nationalaccount`
- UK sort-code/account-number validation and parsing.

### `pkg/vop`
- Offline Verification of Payee matcher (`Match`, `CloseMatch`, `NoMatch`, `Unavailable`).
- `type Verifier interface { Verify(ctx context.Context, req MatchRequest) (MatchResponse, error) }`
- `Matcher` provides a deterministic local implementation of `Verifier`.

### `pkg/iso20022`
- Minimal `pain.001` parser and validation helpers for SCT/SCT Inst subset rules.

## Non-Public Internal Boundaries

The following are internal implementation details and are **not** API contracts:

- `internal/identifiers/*` (per-identifier normalization/validation/parsing)
- `internal/core/validation` (concurrency engine)
- `internal/countrymeta` (generated country metadata)
- `internal/bankregistry` (offline registry loaders and metadata)

Do not import `internal/*` from external modules.
