# Benchmarking Guide

This guide explains how to use the benchmarking suite to measure and improve the performance of the `bank-data` library.

## Running Benchmarks

Use the provided automation script to run all benchmarks with statistical analysis:

```bash
./scripts/run_bench.sh
```

For CI regression gating against the committed baseline:

```bash
./scripts/check_bench_regression.sh
```

For VoP matcher latency checks:

```bash
go test -run=^$ -bench=BenchmarkMatcher_Match -benchmem ./pkg/vop
```

### Advanced Options

- **Profiling**: Generate CPU and memory profiles to find bottlenecks.
  ```bash
  ./scripts/run_bench.sh --profile
  ```
  Then view with: `go tool pprof ./benchmarks/profiles/cpu.prof`

- **Regression Testing (PRs)**:
  Compare your changes against the currently checked-in baseline:
  ```bash
  ./scripts/run_bench.sh benchmarks/current_branch.txt --compare benchmarks/latest.txt
  ```

## Quality Gates

- Coverage gate uses `.ci/coverage_threshold.txt` and `./scripts/check_coverage.sh`.
- Benchmark regression gate compares `benchmarks/financial_baseline.txt` against current runs and fails on >=10% regressions.
- VoP matcher benchmark target is sub-millisecond per `MatchRequest` on developer hardware.

## Interpreting Results

- **ns/op**: Time taken per operation. Lower is better.
- **B/op**: Memory allocated per operation. We target **0 B/op** for core validation logic.
- **allocs/op**: Number of heap allocations. We target **0 allocs/op** for core logic.
- **IBANs/s**: Custom metric showing throughput in items per second.

## Go 1.26 Specifics

### Green Tea GC
To measure the impact of the new GC, you can run benchmarks with different GC settings:
```bash
GODEBUG=gcstoptheworld=1 go test -bench=. # Measure STW pauses
```

### SIMD Optimizations
If you are working on SIMD optimizations, ensure your benchmarks use `b.SetBytes` to report `MiB/s`, which is the standard for vector operations.
