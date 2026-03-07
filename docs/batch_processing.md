# High-Performance Batch Processing

This document outlines the architecture, tutorials, and development plan for the `bank-data` high-performance batching engine, optimized for Go 1.26.

## Overview

Validating financial records at scale (e.g., millions of IBANs during a database migration or daily settlement) requires more than simple loops. The batch engine leverages Go 1.26's runtime improvements to achieve maximum throughput with minimal resource footprint.

## Tutorials & Ideas

### Basic Batch Usage
```go
package main

import (
    "context"
    "fmt"
    "github.com/SamyRai/bank-data/pkg/iban"
)

func main() {
    svc := iban.NewService(nil, nil, nil, nil)
    inputs := []string{"DE89...", "GB33...", "IT..."} // 10,000+ records

    // High-performance batch validation
    results := svc.ValidateBatch(context.Background(), inputs)

    for _, res := range results {
        if !res.Valid {
            fmt.Printf("Error in %s: %v\n", res.Input, res.Error)
        }
    }
}
```

### Streaming with Iterators (Go 1.26)
For processing datasets that don't fit in memory, use the streaming iterator API:

```go
// Stream validated results as they arrive
iter := svc.StreamValidate(context.Background(), inputSourceIter)
for res := range iter {
    // Process results as soon as they are ready
}
```

## Architectural Approach

### 1. Hybrid Pipeline Architecture
Instead of a single "worker pool," the engine uses a 3-stage pipeline:
1. **Pre-scan**: Fast length and character checks (highly parallel).
2. **Checksum Layer**: SIMD-accelerated MOD-97 (Go 1.26 experimental).
3. **Registry Enrichment**: Lookups for country metadata (cached).

### 2. Concurrency Primitives
- **Dynamic Workers**: The engine adjusts worker count automatically based on detected CPU limits (Go 1.25+ cgroup awareness), avoiding context-switching overhead in containerized environments (Kubernetes).
- **Lock-Free Result Channels**: Minimizes mutex contention during the fan-in phase.

### 3. Go 1.26 Optimizations
- **Green Tea GC**: We use pointer-free structs for `ValidationResult` to maximize the "Green Tea" collector's efficiency in scanning small objects.
- **SIMD MOD-97**: Utilizing `simd/archsimd` to process multiple IBAN segments in a single CPU cycle.
- **Fast Small-Object Allocation**: Leveraging Go 1.26's 30% speedup for objects < 512 bytes during result generation.

## Development Plan

### Phase 1: Core Engine (Current)
- [ ] Implement `pkg/validation/engine.go` with dynamic worker scaling.
- [ ] Add `iter.Seq` support for input/output streaming.

### Phase 2: Implementation Optimization
- [ ] Integrate experimental `simd/archsimd` into the IBAN checksum logic.
- [ ] Profile memory allocations and apply `sync.Pool` where the Green Tea GC benefits most.

### Phase 3: Monitoring & Polish
- [ ] Add OpenTelemetry metrics for batch throughput and latency.
- [ ] Create a comprehensive benchmark suite for typical batch sizes (1k, 10k, 100k).

## Future Ideas
- **WASM Support**: Offload validation to client-side browsers using the same high-performance engine.
- **GPU Acceleration**: Researching Vulkan/Metal compute kernels for extreme-scale validation (millions of records per second).
