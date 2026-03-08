# Objective 01: v1 Public API Contract

## Why This Matters

The codebase now has a clear canonical facade (`pkg/financial`) plus compatibility packages (`pkg/iban`, `pkg/bic`, `pkg/sepa`). The next critical step is to freeze a stable contract and remove ambiguity before adding more integration surfaces.

## Current State

- Canonical API is implemented in `pkg/financial/service.go`.
- Public types are in `pkg/financial/types.go`.
- Compatibility wrappers still expose legacy behavior in `pkg/iban`, `pkg/bic`, and `pkg/sepa`.
- There is no explicit compatibility policy document for `v1` guarantees.

## Plan

### Phase 1: Contract Definition

- Define stability scope for `pkg/financial`:
  - stable: `IdentifierType`, `ValidationReport`, `ParsedIdentifier`, `Suggestion`, and all exported `Service` methods
  - unstable/internal: `internal/*` packages and generated metadata format
- Publish versioning rules in docs:
  - additive fields only in minors
  - no breaking type changes in `v1.x`
  - explicit deprecation window for compatibility packages

### Phase 2: API Conformance Lock

- Add API conformance tests focused on:
  - error code consistency per identifier
  - deterministic detect precedence
  - parse/validate normalized value parity
- Add “public API snapshot” checks (compile-time + golden signatures)

### Phase 3: Compatibility Strategy

- Document migration matrix from legacy packages to `pkg/financial`.
- Decide end-of-life schedule for compatibility wrappers (`pkg/iban`, `pkg/bic`, `pkg/sepa`).

## Primary Touchpoints

- `pkg/financial/service.go`
- `pkg/financial/types.go`
- `pkg/financial/service_test.go`
- `docs/api_reference.md`
- `docs/getting_started.md`

## Risks

- Hidden behavior differences between compatibility wrappers and `pkg/financial`.
- Error code drift while adding new identifier modules.

## Exit Criteria

- `v1 API Stability` section documented in `docs/api_reference.md`.
- Dedicated API conformance suite exists and passes in CI.
- Compatibility/deprecation policy published and linked from README.
