# Project Roadmap & Long-Term TODO

> **bank-data** — Go library for IBAN validation, parsing, financial data, and banking utilities.
> All previously tracked tasks have been completed. This document tracks the next phase of development.

---

## 🏁 Milestone 1 — Country Coverage & Data Expansion
**Target:** Q2 2026 — _Broaden the library's utility across the full EU/SEPA geography_

### Bank Code Registry Expansion
- [ ] **FR (France)** — BIC mapping from Banque de France directory
- [ ] **AT (Austria)** — BLZ-to-BIC lookup from Österreichische Nationalbank dataset
- [ ] **NL (Netherlands)** — Dutch bank code (Bankrekening) dataset integration
- [ ] **ES (Spain)** — CCC (Código Cuenta Cliente) bank code mapping
- [ ] **IT (Italy)** — ABI code-to-BIC mapping from Banca d'Italia
- [ ] **PL (Poland)** — NRB / NBP bank code registry
- [ ] **SE (Sweden)** — Clearing number to BIC mapping (Bankgirot)
- [ ] **CH (Switzerland)** — BC-Nr to BIC mapping from SIX Group
- [ ] **UK (United Kingdom)** — Sort Code to BIC via Open Banking directory (post-Brexit SEPA gap)
- [ ] Add automated data refresh pipeline for official banking registries (`gen/` tooling)
- [ ] Add `datasets/` versioning strategy (tag source + date on each downloaded dataset)
- [ ] Write per-country integration tests with real-world IBAN samples

### IBAN Format Extensions
- [ ] Validate IBAN structure for all 80 registered SWIFT country definitions
- [ ] Cross-verify registry against the latest EPC SEPA Country Rulebook (update cycle: annually)
- [ ] Add fuzzy-match correction suggestions (e.g., `GB29 NWBK...` → detect transposed digits)

---

## 🏁 Milestone 2 — Financial Identifier Suite
**Target:** Q3 2026 — _Expand from IBAN-only to a complete financial identifier toolkit_

### New Data Types (from `FUTURE_FEATURES.md`)
- [x] **VAT Number (Phase A/B initial)** — Country routing + checksums for `DE`, `FR`, `NL`, `IT`, `ES`
- [x] **LEI (Legal Entity Identifier)** — MOD-97-10 checksum per ISO 17442
- [x] **ISIN** — Luhn checksum + ISO 6166 country prefix validation
- [x] **Card PAN** — Luhn algorithm with IIN/BIN prefix detection (Visa / MC / Amex / Maestro)
- [ ] **National Account Numbers** — UK Sort Code + Account Number pair validation
- [x] **SEPA Creditor ID** — MOD-97 check digits

### Unified `pkg/financial` Validator API
- [x] Implement canonical `financial.Service` facade (`Detect`, `Validate`, `Parse`)
- [x] Implement `financial.Detect(input string) (Type, error)` — auto-detect identifier type
- [x] Expose typed `financial.Validate(input string, hint Type) (ValidationReport, error)`
- [ ] Release as `v0.2.0` with semver-stable API

---

## 🏁 Milestone 3 — SEPA & Regulatory Compliance
**Target:** Q3–Q4 2026 — _Stay ahead of mandatory EU payment regulation changes_

> **⚠️ REGULATORY DEADLINE: October 9, 2025** — SEPA Verification of Payee (VoP) became mandatory for
> Eurozone PSPs. Non-Eurozone deadline: **July 9, 2027**. Integrate support now for downstream users.

### Verification of Payee (VoP) Support
- [ ] Design `pkg/vop` package with `MatchRequest` / `MatchResponse` types  
  - Response categories: `Match`, `CloseMatch`, `NoMatch`, `Unavailable`
- [ ] Implement name matching heuristics (levenshtein, diacritic normalization, legal suffix stripping)
- [ ] Provide an interface `vop.Verifier` for integrating with bank API backends
- [ ] Document integration pattern for PSPs using the library
- [ ] Add fuzzy name matching benchmarks (target < 1ms per check)

### ISO 20022 Alignment
- [ ] Parse and validate ISO 20022 `pain.001` payment initiation messages (SCT schema)
- [ ] Support ISO 20022 structured address fields (mandatory from **November 2026**)
  - Fields: `TwnNm`, `Ctry`, `PstCd` + max 2 × 70-char address lines
- [ ] Provide `iban.ToISO20022Creditor()` / `iban.ToISO20022Debtor()` conversion helpers
- [ ] Validate SEPA Instant Credit Transfer (SCT Inst) specific rules (amount ≤ 100,000 EUR, BIC mandatory)

### SEPA Rulebook Updates
- [ ] Track and implement new EPC SEPA Rulebook version (updated annually each November)
- [ ] Support `Alias/Proxy` IBAN derivation field (email / phone → IBAN lookup interface)
- [ ] Add `IsSEPAInstant(countryCode string) bool` (not all SEPA countries support SCT Inst)

---

## 🏁 Milestone 4 — Performance & Scalability
**Target:** Q4 2026 — _Production-grade throughput for high-volume financial systems_

### Batch Processing Engine (`pkg/batch`)
- [ ] Implement `batch.Processor` with configurable worker pool size
- [ ] Use Go channels + `sync.WaitGroup` for concurrent validation pipeline
- [ ] Support streaming I/O: `io.Reader` CSV input → `io.Writer` JSON/CSV output
- [ ] Target: > 100,000 IBAN validations/sec on a single 4-core machine
- [ ] Add `benchmarks/batch_test.go` suite and publish results in `docs/BENCHMARKING.md`
- [ ] Lazy-load datasets in batch mode to cap memory footprint under 50 MB

### Core Performance Improvements
- [ ] Profile `pkg/iban` hot paths with `pprof` and `go tool trace`
- [ ] Investigate SIMD-accelerated MOD-97 for bulk calculations (Go assembly or CGo)
- [ ] Add benchmark regression gate to CI (fail if >10% regression vs. baseline)
- [ ] Implement `sync.Pool` for `IBANInfo` struct reuse in hot paths
- [ ] Publish detailed benchmark table in `docs/BENCHMARKING.md` (per country, cold/warm)

---

## 🏁 Milestone 5 — Developer Experience & Ecosystem
**Target:** Q1 2027 — _Make bank-data the go-to library for Go financial developers_

### API & Documentation
- [ ] Complete and publish full `docs/api_reference.md` for all packages
- [ ] Add runnable Go Playground examples to README for each identifier type
- [ ] Generate `pkg.go.dev` documentation with rich package-level examples
- [ ] Write a "Getting Started" guide for fintech engineers (target: 5-minute integration)
- [ ] Add a `CHANGELOG.md` with semantic version history

### SDK / Client Libraries
- [ ] Evaluate gRPC service wrapper for `bank-data` (for polyglot environments)
- [ ] Provide a REST microservice Docker image (`bank-data-api`) for non-Go users
- [ ] Publish pre-built binaries for the `gen/` registry tool via GitHub Releases

### Testing & Quality
- [ ] Add property-based tests using `pgregory.net/rapid` for IBAN generation + validation round-trips
- [ ] Achieve ≥ 90% test coverage across all `pkg/` packages (gate in CI)
- [ ] Add corpus-based fuzz testing for all public `Validate()` and `Parse()` entrypoints
- [ ] Integration test suite with real-world SWIFT IBAN registry test vectors
- [ ] Test against official EPC SEPA test IBANs (published annually by EPC)

### Open Source & Community
- [ ] Add `CONTRIBUTING.md` with contribution guidelines and ADR process
- [ ] Create GitHub Issue templates (bug report, feature request, new country support)
- [ ] Set up GitHub Discussions for RFC-style proposals
- [ ] Publish library to `pkg.go.dev` under a stable `v1.x` tag
- [ ] Reach out to `moov-io` / `gocardless` ecosystems for potential collaboration

---

## 🏁 Milestone 6 — Enterprise & Advanced Features
**Target:** Q2–Q3 2027 — _Capabilities needed by fintechs and banks in production_

### Open Banking Integration
- [ ] Design `pkg/openbanking` with PSD2/PSD3-compatible account lookup interface
- [ ] Implement GoCardless (Nordigen) Open Banking API adapter for account verification
- [ ] Support pluggable bank verification backends (Plaid, TrueLayer, etc.)
- [ ] Provide rate-limiter-aware HTTP client wrappers for external API calls

### Bank Information Enrichment
- [ ] Extend IBAN parser to resolve bank metadata: name, address, BIC, phone
- [ ] Integrate with public SWIFT BIC directory (via licensed dataset or scrape-free API)
- [ ] Optional cache layer for bank metadata (TTL-based, Redis-injectable interface)
- [ ] Support reverse lookup: BIC → IBAN country + bank code

### Fraud & Risk Signals
- [ ] `risk.Score(iban string) (RiskSignal, error)` — heuristic-based risk indicator
  - Signals: newly-registered bank code, inactive country, structured as known fraud pattern
- [ ] VoP name-match score as a numeric confidence value (0–1 range)
- [ ] Provide integration hooks for downstream fraud detection pipelines

### Observability & Reliability
- [ ] Add structured `slog`-compatible logging throughout all packages
- [ ] Expose `prometheus`-compatible metrics hook for validation throughput and error rates
- [ ] Implement per-country circuit-breaker for external bank registry fetches
- [ ] Add OpenTelemetry span instrumentation for batch operations

---

## 📌 Ongoing / Evergreen Tasks
> These tasks repeat on a regular cadence and are never "done"

- [ ] **Annual** — Update country IBAN registry against latest SWIFT/EPC data (every November)
- [ ] **Annual** — Audit SEPA Rulebook changes and implement required updates
- [ ] **Quarterly** — Refresh bank code datasets for active countries (DE, FR, AT, NL, UK…)
- [ ] **Per-release** — Run fuzz tests for 24 hours before any `v1.x` tag
- [ ] **Per-release** — Update `CHANGELOG.md` and tag with semantic version
- [ ] **Monthly** — Dependency audit (`go list -m -u all`) and security scan

---

## 🔭 Research Queue
> Ideas requiring investigation before commitment

- [ ] **R1** — Evaluate `wasm` compilation target for browser-side IBAN validation
- [ ] **R2** — Investigate MICA (EU crypto asset regulation) address validation overlap
- [ ] **R3** — Assess Go generics for a unified `Validator[T]` interface across all financial types
- [ ] **R4** — Research EBA (European Banking Authority) open data APIs for real-time BIC resolution
- [ ] **R5** — Explore formal verification of MOD-97 implementation via TLA+ or Coq
- [ ] **R6** — Evaluate `bank-data` as a Wasm plugin for payment gateways (Envoy/WASM filter)

---

_Last updated: 2026-03-08 | See also [FUTURE_FEATURES.md](FUTURE_FEATURES.md) and [docs/](docs/)_
