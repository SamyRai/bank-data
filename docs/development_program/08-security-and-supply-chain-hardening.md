# Objective 08: Security and Supply-Chain Hardening

## Why This Matters

This library validates untrusted financial input and may be embedded in regulated systems. Security confidence must include input hardening, dependency posture, and release integrity.

## Current State

- Security reporting policy exists in root `SECURITY.md`.
- CI includes `gosec` workflow and dependency review workflow.
- Fuzzing exists but primarily as short smoke runs.

## Plan

### Phase 1: Threat Model and Security Baseline

- Add a concise threat model doc for:
  - malformed payload denial-of-service paths
  - parser abuse (especially XML in `pkg/iso20022`)
  - data poisoning in offline datasets
- Classify trust boundaries (`pkg/*` public input vs internal generated data).

### Phase 2: Hardening Controls

- Add input-size guards for expensive parse paths (especially XML).
- Add explicit time/space constraints to batch flows.
- Expand negative tests for parser/resource exhaustion scenarios.

### Phase 3: Supply-Chain and Release Integrity

- Add SBOM generation step for release builds.
- Sign release artifacts and publish checksum files.
- Add dependency policy checks (licenses + critical CVEs).

## Primary Touchpoints

- `pkg/iso20022/pain001.go`
- `pkg/financial/service.go`
- `.github/workflows/gosec.yml`
- `.github/workflows/dependency-review.yml`
- `SECURITY.md`

## Risks

- Overly strict input limits can block legitimate large messages.
- Additional scanning steps may increase CI runtime.

## Exit Criteria

- Threat model published in `docs/` and linked from `SECURITY.md`.
- Security regression tests added for key parser and batch entrypoints.
- Release process emits signed artifacts and SBOM metadata.
