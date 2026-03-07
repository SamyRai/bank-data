# Project Roadmap & TODOs

Tracking work for the `bank-data` core and expansion.

## Active Tasks (Prioritized)

- [x] **Registry Automation**: Implement `go:generate` for country metadata.
- [x] **Checksum Integrity**: Migrate to streaming MOD-97 algorithm.
- [x] **API Unification**: Single `iban.Service` facade with DI.
- [x] **Typed Errors**: Exported, informative error constants.
- [x] **Public Accessibility**: Refactor `internal/iban` to `pkg/iban` for external users.
- [x] **Security Hardening**
  - [x] Enforce 34-char input limit.
  - [x] Implement constant-time comparison for checksum components.
- [x] **Bank Directory Integration**
  - [x] Implement DE bank code (BLZ) to BIC mapping.
  - [x] Support lazy-loading for large datasets.
- [x] **Data Type Expansion**
  - [x] Detailed validation rules for BIC/SWIFT.
  - [x] SEPA Creditor ID support.

## Next Priorities

1. **Additional Countries**: Expand bank code mapping for other EU countries (e.g., FR, AT).
2. **Performance Benchmarking**: Add more detailed benchmarks for lazy loading.
3. **Documentation**: Complete API documentation and usage examples.
