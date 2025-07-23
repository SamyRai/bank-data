# TODO List for IBAN Core Refactor and Productionization (Prioritized 80/20)

- [x] **Expand country registry: Add top EU/SEPA and broad set of IBAN countries to countrymeta.Registry**
- [ ] **Automate registry generation from SWIFT CSV (go:generate)**
  - Priority: Highest | Effort: 1d
- [x] **Upgrade and expand tests: Table-driven, fuzzing, edge cases**
- [x] **Switch to streaming MOD-97 algorithm**
- [x] **Unify public API surface (Service type, DI)**
- [x] **Introduce typed, exported errors**
- [ ] **Security & compliance hardening**
  - Cap input at 34 chars
  - Use constant-time comparison for check digits
  - Priority: High | Effort: 0.5d
- [ ] **Add bank directory lookups (e.g. DE BLZ→BIC)**
  - Priority: High | Effort: 1d
- [ ] **Developer ergonomics**
  - One public facade (iban.Service)
  - Chainable API (iban.Must(s).BankCode())
  - Priority: Medium | Effort: 1d
- [ ] **Documentation and CI**
  - Add README, LICENSE, and CI config for multi-platform, multi-version testing
  - Priority: Medium | Effort: 0.5d
- [x] **Remove redundant IBANRegistrySource abstraction in favor of general encoder/decoder**
  - Priority: Low | Effort: 5m

---

## Suggested Next Priorities

1. **Automate country registry updates**: Implement and test go:generate-based automation for countrymeta.Registry using official SWIFT/EC CSV. This will ensure future-proof, error-free expansion and maintenance. (Highest priority)
2. **Security & compliance hardening**: Enforce input length, add constant-time check digit comparison, and review for edge-case handling.
3. **Bank directory lookup**: For DE and other countries, map bank codes (e.g. BLZ) to BIC, using static or lazy-loaded data.
4. **Developer ergonomics**: Refine the public API for ease of use and chaining, and document usage patterns.
5. **Documentation and CI**: Add README, LICENSE, and set up CI for robust, multi-platform testing.

*Update this list after each milestone. Use checkboxes to track progress. Add TODO comments in code for context, priority, and estimated effort.*
