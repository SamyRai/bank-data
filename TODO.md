# TODO List for IBAN Core Refactor and Productionization (Prioritized 80/20)

- [x] **Expand country registry: Add top EU/SEPA and broad set of IBAN countries to countrymeta.Registry**
- [ ] **Automate registry generation from SWIFT CSV (go:generate)**
  - Priority: Highest | Effort: 1d
- [x] **Upgrade and expand tests: Table-driven, fuzzing, edge cases**
- [x] **Switch to streaming MOD-97 algorithm**
- [x] **Unify public API surface (Service type, DI)**
- [x] **Introduce typed, exported errors**
- [x] **Security & compliance hardening**
  - Cap input at 34 chars
  - Use constant-time comparison for check digits
  - Add comprehensive edge case tests for length, format, characters, country, and checksum
  - Priority: High | Effort: 0.5d
- [x] **Add bank directory lookups (e.g. DE BLZ→BIC)**
  - Implemented and production ready: IBAN parsing extracts bank code, lookup performed using loaded datasets, public API supports enrichment, and tested in bic_test.go
  - Priority: High | Effort: 1d
- [ ] **Developer ergonomics**
  - [x] One public facade (iban.Service) — done
  - [ ] Chainable API (iban.Must(s).BankCode()) — TODO: Implement fluent, chainable API for ergonomic access to parsed IBAN fields and validation results. Recommended approach:
    - Add `MustParse(iban string) IBANInfo` helper to panic on error, for one-liners in tests and scripts.
    - Add fluent accessors (e.g., `BankCode()`, `CountryCode()`, `AccountNumber()`) to `IBANInfo`.
    - Optionally, add a wrapper type (e.g., `IBANResult`) with chainable methods for validation, parsing, and enrichment.
    - Document usage patterns and provide examples in README.
    - Priority: Medium | Effort: 1d
- [ ] **Documentation and CI**
  - Add README, LICENSE, and CI config for multi-platform, multi-version testing
  - Priority: Medium | Effort: 0.5d
- [x] **Remove redundant IBANRegistrySource abstraction in favor of general encoder/decoder**
  - Priority: Low | Effort: 5m

---

## Suggested Next Priorities

1. **Automate country registry updates**: Implement and test go:generate-based automation for countrymeta.Registry using official SWIFT/EC CSV. This will ensure future-proof, error-free expansion and maintenance. (Highest priority)
2. **Developer ergonomics**: Implement chainable API for IBAN parsing and validation, as described above.
3. **Documentation and CI**: Add README, LICENSE, and set up CI for robust, multi-platform testing.

*Update this list after each milestone. Use checkboxes to track progress. Add TODO comments in code for context, priority, and estimated effort.*
