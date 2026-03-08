# Objective 04: Official Vector and Conformance Program

## Why This Matters

Property and fuzz tests are present, but long-term correctness requires curated conformance vectors tied to standards and rulebooks.

## Current State

- Existing vectors are in `testdata/financial_test_vectors.csv`.
- Conformance coverage for IBAN structure exists in `internal/identifiers/iban/conformance_test.go`.
- Fuzz tests cover multiple public packages.

## Plan

### Phase 1: Vector Taxonomy

- Split vectors by source and confidence level:
  - standard-derived vectors
  - regression vectors (historical bugs)
  - adversarial edge vectors
- Add metadata columns (`source`, `edition`, `confidence`).

### Phase 2: Package-Level Conformance Suites

- Add standard-focused conformance tests per package:
  - `pkg/iban`, `pkg/sepa`, `pkg/iso20022`, `pkg/vop`, `pkg/nationalaccount`
- Ensure detect/validate/parse behavior is tested consistently from `pkg/financial`.

### Phase 3: Drift Detection

- Add CI job that fails if generated metadata changes without vector refresh notes.
- Add summary output listing newly failing vectors by identifier type.

## Primary Touchpoints

- `testdata/financial_test_vectors.csv`
- `pkg/financial/vectors_test.go`
- `internal/identifiers/iban/conformance_test.go`
- `pkg/*/*_test.go` and `pkg/*/*_fuzz_test.go`

## Risks

- Some official standards have licensing/distribution constraints.
- Mixed vector quality can produce noisy failures unless clearly tagged.

## Exit Criteria

- Vector taxonomy documented and implemented in testdata.
- Each public identifier package has explicit conformance tests.
- CI emits clear, actionable vector failure reports.
