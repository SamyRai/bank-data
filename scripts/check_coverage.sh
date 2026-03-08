#!/usr/bin/env bash
set -euo pipefail

threshold_file=".ci/coverage_threshold.txt"
cover_file="coverage.out"

if [[ ! -f "$threshold_file" ]]; then
  echo "missing coverage threshold file: $threshold_file" >&2
  exit 1
fi
if [[ ! -f "$cover_file" ]]; then
  echo "missing coverage profile: $cover_file" >&2
  exit 1
fi

threshold=$(tr -d '[:space:]' < "$threshold_file")
actual=$(go tool cover -func="$cover_file" | awk '/^total:/{gsub("%","",$3); print $3}')

if [[ -z "$actual" ]]; then
  echo "failed to parse total coverage from $cover_file" >&2
  exit 1
fi

awk -v actual="$actual" -v threshold="$threshold" 'BEGIN {
  if (actual + 0 < threshold + 0) {
    printf("coverage gate failed: actual %.2f%% < threshold %.2f%%\n", actual, threshold)
    exit 1
  }
  printf("coverage gate passed: actual %.2f%% >= threshold %.2f%%\n", actual, threshold)
}'

awk -v actual="$actual" -v threshold="$threshold" 'BEGIN {
  if (actual - threshold >= 0.5) {
    printf("note: coverage improved materially (%.2f%%). Consider ratcheting .ci/coverage_threshold.txt\n", actual)
  }
}'
