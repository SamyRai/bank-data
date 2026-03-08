# Objective 02: Dataset Governance and Provenance

## Why This Matters

Validation quality now depends on offline datasets (`datasets/`, generated country metadata, and bank registry fixtures). Governance is required to keep updates reproducible and auditable.

## Current State

- Dataset manifest exists: `datasets/manifest.json`.
- Generated IBAN metadata is produced into `internal/countrymeta/registry.go`.
- Registry loaders exist in `internal/bankregistry/registry.go` with per-country metadata fields.
- No strict schema/version enforcement exists for manifests and fixtures.

## Plan

### Phase 1: Manifest Contract

- Define strict manifest schema:
  - dataset id, source URL, license, version date, checksum, generation timestamp
- Add schema validation command integrated into CI.
- Fail CI when manifest entries are missing required fields.

### Phase 2: Reproducible Generation

- Add deterministic generation docs and scripts for:
  - `datasets/iban-registry.txt` refresh
  - `internal/countrymeta/registry.go` regeneration
- Ensure generator output is stable across runs (ordering + formatting).

### Phase 3: Provenance and Drift Reports

- Produce update reports for data changes:
  - row counts by country
  - added/removed bank codes
  - checksum change summary
- Store a machine-readable diff artifact for review.

## Primary Touchpoints

- `datasets/manifest.json`
- `datasets/iban-registry.txt`
- `gen/cmd/gen_registry/main.go`
- `internal/countrymeta/base.go`
- `internal/bankregistry/registry.go`

## Risks

- Data source licensing changes may force schema changes.
- Non-deterministic generator behavior can create noisy commits.

## Exit Criteria

- Manifest validation gate runs in CI.
- Dataset updates produce deterministic output and reproducible diffs.
- Every loaded country dataset has traceable metadata (`source`, `version_date`, `checksum`, `license`).
