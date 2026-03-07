#!/bin/bash
set -e

# Configuration
# Default to storing results in the repo for regression checking
OUTPUT_FILE=${1:-"benchmarks/latest.txt"}
PROFILE_DIR="benchmarks/profiles"
CPU_PROFILE="$PROFILE_DIR/cpu.prof"
MEM_PROFILE="$PROFILE_DIR/mem.prof"

# Ensure tools are installed
if ! command -v benchstat &> /dev/null; then
    echo "Installing benchstat..."
    go install golang.org/x/perf/cmd/benchstat@latest
fi

# Help message
show_help() {
    echo "Usage: $0 [output_file] [--profile] [--compare base_file]"
    echo ""
    echo "Default output: benchmarks/latest.txt"
    echo ""
    echo "Options:"
    echo "  --profile         Generate CPU and Memory profiles in $PROFILE_DIR"
    echo "  --compare FILE    Compare results against a baseline using benchstat"
}

DO_PROFILE=false
COMPARE_FILE=""

# Parse flags
# We handle the first argument as OUTPUT_FILE if it doesn't start with --
if [[ "$1" != --* ]] && [[ -n "$1" ]]; then
    OUTPUT_FILE=$1
    shift
fi

while [[ "$#" -gt 0 ]]; do
    case $1 in
        --profile) DO_PROFILE=true ;;
        --compare) COMPARE_FILE="$2"; shift ;;
        -h|--help) show_help; exit 0 ;;
    esac
    shift
done

mkdir -p benchmarks

echo "Running benchmarks and saving to $OUTPUT_FILE..."

# Standard benchmark run
# -benchmem for allocation stats
# -count 5 for statistical significance in benchstat
go test -bench=. -benchmem -count 5 ./... > "$OUTPUT_FILE"

# Profiling if requested
if [ "$DO_PROFILE" = true ]; then
    echo "Generating profiles in $PROFILE_DIR..."
    mkdir -p "$PROFILE_DIR"
    # Run a specific representative benchmark for profiling to keep it focused
    go test -bench=BenchmarkService_ValidateBatch_Scaling -cpuprofile="$CPU_PROFILE" -memprofile="$MEM_PROFILE" ./pkg/iban
    echo "Profiles generated: $CPU_PROFILE, $MEM_PROFILE"
    echo "To view CPU profile: go tool pprof $CPU_PROFILE"
fi

# Compare if requested
if [ -n "$COMPARE_FILE" ]; then
    echo "--- Comparison vs $COMPARE_FILE ---"
    benchstat "$COMPARE_FILE" "$OUTPUT_FILE"
else
    echo "--- Benchmark Summary ---"
    benchstat "$OUTPUT_FILE"
fi
