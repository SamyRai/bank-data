# Objective 09: VoP Scoring Engine Maturity

## Why This Matters

`pkg/vop` now provides deterministic local matching. To support production verification workflows, it needs stronger scoring semantics, explainability, and policy controls.

## Current State

- `pkg/vop/matcher.go` includes normalization, suffix stripping, and Levenshtein score.
- Response includes categorical match and numeric score.
- Thresholds are fixed in `NewMatcher()` with no policy profile abstraction.

## Plan

### Phase 1: Scoring Explainability

- Extend `MatchResponse` to include explainability fields:
  - matched tokens
  - penalties applied
  - normalization transforms used
- Keep backward compatibility by making new fields additive.

### Phase 2: Policy Profiles

- Introduce configurable profiles (strict, balanced, lenient).
- Support region-aware normalization plugins without network dependency.
- Add deterministic fixtures for profile behavior.

### Phase 3: Fraud/Risk Hooks

- Add optional hooks for downstream risk systems:
  - threshold override callback
  - custom rule scoring extension
- Keep default flow pure and deterministic.

## Primary Touchpoints

- `pkg/vop/matcher.go`
- `pkg/vop/matcher_test.go`
- `pkg/vop/matcher_bench_test.go`
- `pkg/vop/vop_fuzz_test.go`

## Risks

- More scoring dimensions can reduce predictability if not carefully documented.
- Locale-specific rules can create bias and false positives.

## Exit Criteria

- Profile-driven matching behavior is tested with deterministic vectors.
- Match responses include enough explanation for audit/debug use.
- p95 latency remains within documented target under benchmark load.
