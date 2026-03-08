# Development Program (Strategic)

This directory breaks the next major development cycle into 10 independent but connected objectives.

Each objective has:
- current codebase context
- implementation plan
- dependencies and risks
- concrete acceptance criteria

## Execution Order

1. [Objective 01: v1 Public API Contract](01-v1-public-api-contract.md)
2. [Objective 02: Dataset Governance and Provenance](02-dataset-governance-and-provenance.md)
3. [Objective 03: Country Coverage Expansion](03-country-coverage-expansion.md)
4. [Objective 04: Official Vector and Conformance Program](04-official-vector-and-conformance-program.md)
5. [Objective 05: Observability and Diagnostics](05-observability-and-diagnostics.md)
6. [Objective 06: Performance SLOs and Regression Control](06-performance-slos-and-regression-control.md)
7. [Objective 07: Integration Surfaces (CLI/HTTP/gRPC)](07-integration-surfaces-cli-http-grpc.md)
8. [Objective 08: Security and Supply-Chain Hardening](08-security-and-supply-chain-hardening.md)
9. [Objective 09: VoP Scoring Engine Maturity](09-vop-scoring-engine-maturity.md)
10. [Objective 10: Ecosystem and Adoption](10-ecosystem-and-adoption.md)

## Status Snapshot (2026-03-08)

- Objective 01: complete (v1 contract + conformance suite merged in PR #9).
- Objective 02: next critical objective; governs data provenance for Objectives 03/04/08.
- Objectives 03-10: planned; sequencing below reflects current blockers/unblocks.

## Program Dependencies

- Objective 01 is complete; Objectives 07 and 10 are now unblocked from contract stability concerns.
- Objective 02 remains required before large parts of Objectives 03 and 04.
- Objective 04 should feed tests into Objectives 06 and 08.
- Objective 05 observability should be in place before high-throughput benchmarking in Objective 06.
- Objective 09 depends on Objective 04 test vectors and Objective 06 SLO measurement.

## Dependency Matrix (Objectives 02-10)

| Objective | Depends On | Status | Notes |
| --- | --- | --- | --- |
| 02 Dataset Governance | 01 (completed) | Ready | Highest priority next step; enables reproducible data foundation. |
| 03 Country Coverage | 02 | Blocked | Avoid broad country rollout before manifest/schema/provenance gates are in place. |
| 04 Vector + Conformance | 02 | Partially blocked | Can add local vectors now; official/source-governed expansion should follow 02. |
| 05 Observability | 01 (completed) | Ready | Can proceed independently of dataset expansion. |
| 06 Performance SLOs | 05, 04 | Sequenced | Baselines exist; final SLO enforcement should include observability + richer vectors. |
| 07 CLI/HTTP/gRPC | 01 (completed) | Ready | Contract lock is in place; wrappers can align to stable API surface. |
| 08 Security + Supply Chain | 04, 02 | Partially ready | CI baseline is stable; deeper hardening should consume governed datasets + vectors. |
| 09 VoP Maturity | 04, 06 | Blocked | Needs vector depth and measurable performance budgets first. |
| 10 Ecosystem + Adoption | 01, 07 | Sequenced | Docs/process can advance now; broader adoption push should follow integration surfaces. |

## Current Baseline Anchors

- Public facade: `pkg/financial`
- Identifier modules: `internal/identifiers/*`
- Offline registry: `internal/bankregistry`
- Validation engine: `internal/core/validation`
- CI gates: `.github/workflows/test.yml`, `.github/workflows/lint.yml`, `.github/workflows/gosec.yml`
- Quality scripts: `scripts/check_repo_docs.sh`, `scripts/check_coverage.sh`, `scripts/check_bench_regression.sh`
