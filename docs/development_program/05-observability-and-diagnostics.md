# Objective 05: Observability and Diagnostics

## Why This Matters

As the library grows into high-throughput and service-wrapper use cases, users need consistent telemetry hooks without forcing a runtime dependency by default.

## Current State

- Logging utility exists in `internal/log`, but telemetry is not exposed as stable public hooks.
- No metrics/tracing interfaces in `pkg/financial` or `internal/core/validation`.

## Plan

### Phase 1: Telemetry Interfaces

- Add optional interfaces for:
  - structured logging events
  - counters/histograms (validation count, duration, error codes)
  - tracing spans around batch and parse/validate operations
- Keep default behavior zero-overhead when no observer is configured.

### Phase 2: Instrument Core Paths

- Instrument:
  - `pkg/financial` detect/validate/parse/suggest
  - `internal/core/validation` batch stream execution
  - `internal/bankregistry` load/lookup paths

### Phase 3: Documentation and Examples

- Add examples for integrating:
  - `slog`
  - Prometheus-style metrics adapter
  - OpenTelemetry tracer adapter

## Primary Touchpoints

- `pkg/financial/service.go`
- `internal/core/validation/engine.go`
- `internal/bankregistry/registry.go`
- `internal/log/log.go`

## Risks

- API bloat if telemetry hooks expose too many internals.
- Allocation overhead on hot paths if callbacks are unguarded.

## Exit Criteria

- Telemetry hooks are optional, documented, and benchmarked.
- Metrics and traces can be attached without changing validation semantics.
- Regression tests verify telemetry does not alter functional outputs.
