#!/usr/bin/env bash
set -euo pipefail

baseline="benchmarks/financial_baseline.txt"
current="benchmarks/current_financial.txt"
report="benchmarks/benchstat_financial.txt"

if [[ ! -f "$baseline" ]]; then
  echo "missing benchmark baseline: $baseline" >&2
  exit 1
fi

if ! command -v benchstat >/dev/null 2>&1; then
  echo "installing benchstat..."
  go install golang.org/x/perf/cmd/benchstat@latest
fi

mkdir -p benchmarks

go test -run=^$ -bench=BenchmarkFinancialValidate_Matrix -benchmem -count 3 ./pkg/financial > "$current"
benchstat "$baseline" "$current" > "$report"

cat "$report"

if ! awk '
  /sec\/op/ {in_sec = 1}
  in_sec && /B\/s/ {in_sec = 0}
  in_sec {
    line = $0
    while (match(line, /\+[0-9]+(\.[0-9]+)?%/)) {
      pct = substr(line, RSTART + 1, RLENGTH - 2) + 0
      if (pct >= 10) {
        bad = 1
        printf("benchmark regression >=10%% detected in sec/op: %s\n", $0) > "/dev/stderr"
      }
      line = substr(line, RSTART + RLENGTH)
    }
  }
  END { exit bad ? 1 : 0 }
' "$report"; then
  echo "benchmark regression gate failed: detected >=10% sec/op regression" >&2
  exit 1
fi

echo "benchmark regression gate passed"
