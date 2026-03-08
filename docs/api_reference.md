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

## v1 API Stability

The `pkg/financial` package serves as the canonical v1 public API for `bank-data`.

### Stability Scope

**Stable:**
*   Types: `IdentifierType`, `ValidationReport`, `ParsedIdentifier`, `Suggestion`
*   Methods: All exported methods on `financial.Service`
*   Constants: All defined `IdentifierType` constants

**Unstable / Internal:**
*   Any package under `internal/*` (including generated metadata formats)
*   These are not covered by semantic versioning guarantees and may change at any time.

### Versioning Rules

*   **Additive changes:** New fields on structs (e.g., `ValidationReport`, `ParsedIdentifier.Fields`) or new exported methods may be added in minor version releases. Consumers should not rely on exact struct field counts or exhaustive matching without a default case.
*   **Breaking changes:** There will be no breaking type changes, method signature changes, or removals in the v1.x series for the stable scope.

### Compatibility & Deprecation Policy

The legacy packages (`pkg/iban`, `pkg/bic`, `pkg/sepa`) are currently maintained for compatibility.
*   **Deprecation Window:** These packages are considered deprecated and are scheduled for removal in `v2.0` (or after a minimum 6-month deprecation window).
*   **Migration:** All new integrations must use `pkg/financial`. Existing integrations should migrate to `pkg/financial` according to the migration matrix below.

#### Migration Matrix

| Legacy Package / Method | Recommended `pkg/financial` Equivalent |
| :--- | :--- |
| `pkg/iban.NewService().Validate(iban)` | `financial.NewService().Validate(input, financial.IdentifierIBAN)` |
| `pkg/iban.NewService().Parse(iban)` | `financial.NewService().Parse(input, financial.IdentifierIBAN)` |
| `pkg/bic.NewValidator().Validate(bic)` | `financial.NewService().Validate(input, financial.IdentifierBIC)` |
| `pkg/bic.Parse(bic)` | `financial.NewService().Parse(input, financial.IdentifierBIC)` |
| `pkg/sepa.NewCreditorIDValidator().Validate(id)` | `financial.NewService().Validate(input, financial.IdentifierSEPACreditor)` |

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
