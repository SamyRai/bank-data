# Objective 07: Integration Surfaces (CLI/HTTP/gRPC)

## Why This Matters

Many consumers cannot embed Go libraries directly. Thin service/CLI surfaces can broaden adoption while preserving `pkg/financial` as the single source of truth.

## Current State

- Core API is library-first (`pkg/financial`).
- No public CLI or service wrapper package exists in-repo for runtime integration.
- Generator tooling exists under `gen/`, but not user-facing validation tooling.

## Plan

### Phase 1: CLI Reference Tool

- Add `cmd/bank-data` with subcommands:
  - `detect`
  - `validate`
  - `parse`
  - `suggest`
- Keep output deterministic (JSON by default, optional CSV for batch).

### Phase 2: HTTP Reference Server

- Add minimal REST server package (internal or `cmd/` scoped):
  - `POST /detect`
  - `POST /validate`
  - `POST /parse`
  - `POST /suggest`
- Reuse existing `pkg/financial` types and error codes.

### Phase 3: gRPC Feasibility

- Define protobuf contracts that mirror `pkg/financial` semantics.
- Add compatibility tests that compare gRPC/REST outputs to direct library calls.

## Primary Touchpoints

- `pkg/financial/*`
- `pkg/financial/types.go`
- `docs/api_reference.md`
- `TODO.md` SDK/client section

## Risks

- API drift between wrappers and library if contracts are copied manually.
- Error semantics divergence across transport formats.

## Exit Criteria

- CLI tool can process single and batch inputs with stable output schema.
- HTTP wrapper passes contract tests against direct library calls.
- gRPC plan is documented with clear decision to proceed or defer.
