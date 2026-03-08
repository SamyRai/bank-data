# Objective 06: Performance SLOs and Regression Control

## Why This Matters

Benchmarks and gates exist, but they should be tied to explicit SLOs and package-level budgets to guide optimization and avoid premature complexity.

## Current State

- Bench baseline exists in `benchmarks/financial_baseline.txt`.
- Regression check exists in `scripts/check_bench_regression.sh`.
- Benchmarks exist for `pkg/financial`, `pkg/iban`, and `pkg/vop` matcher.

## Plan

### Phase 1: SLO Definition

- Define target SLOs per public API path:
  - `Detect`, `Validate`, `Parse`, `Suggest`
  - batch throughput and max p95 duration
- Define max allocation budgets for hot paths.

### Phase 2: Benchmark Matrix Expansion

- Add benchmarks for:
  - all identifier hints through `pkg/financial`
  - enriched IBAN parse path with registry lookups
  - representative invalid/adversarial input sets
- Store baseline files per benchmark family.

### Phase 3: Controlled Optimization

- Use `pprof` before optimization changes.
- Apply optimizations only when benchmark evidence shows value.
- Keep optimizations local and reversible; avoid cross-cutting complexity.

## Primary Touchpoints

- `pkg/financial/service_bench_test.go`
- `pkg/iban/service_bench_test.go`
- `pkg/vop/matcher_bench_test.go`
- `scripts/check_bench_regression.sh`
- `docs/BENCHMARKING.md`

## Risks

- Benchstat noise on small sample counts can create false positives.
- Optimizing for microbenchmarks may hurt real batch behavior.

## Exit Criteria

- SLO table is documented and mapped to benchmark names.
- CI checks benchmark regressions against objective-specific baselines.
- Any optimization merge includes before/after benchmark evidence.
