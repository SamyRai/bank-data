# Objective 03: Country Coverage Expansion

## Why This Matters

The library is strongest when bank metadata enrichment and local account semantics are available for more countries, not only core examples.

## Current State

- Offline loader framework supports deterministic CSV fixtures (`internal/bankregistry/LoadCSV`).
- Initial fixture coverage exists for `FR`, `AT`, and `NL` in `testdata/bankregistry/`.
- `TODO.md` tracks additional countries (`ES`, `IT`, `PL`, `SE`, `CH`, `UK`).

## Plan

### Phase 1: Loader Standardization

- Define a per-country loader contract using `CountryLoader`.
- Add canonical schema maps for each supported country source format.
- Keep all integration tests fixture-driven with no network fetches.

### Phase 2: Priority Country Rollout

- Implement deterministic loaders and tests for:
  - `ES`, `IT`, `PL`, `SE`, `CH`, `UK`
- Add enrichment verification tests via `pkg/financial` + `WithBankEnricher`.

### Phase 3: Coverage Reporting

- Add a generated coverage matrix doc showing:
  - countries with bank-code support
  - countries with local account support
  - last dataset version date

## Primary Touchpoints

- `internal/bankregistry/registry.go`
- `internal/bankregistry/registry_test.go`
- `testdata/bankregistry/*`
- `pkg/financial/service.go` (IBAN enrichment path)
- `TODO.md`

## Risks

- Country source formats vary significantly and may require custom normalization.
- BIC mapping assumptions can differ by country and must be explicitly tested.

## Exit Criteria

- At least six additional countries added with deterministic fixtures and tests.
- Enrichment behavior is documented and tested through public API paths.
- Country support matrix is published in docs.
