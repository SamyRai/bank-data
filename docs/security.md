# Security Policy

The `bank-data` project is committed to providing a secure and reliable library for handling sensitive financial data formats.

## Implemented Security Features

### 1. Robust Input Validation
- **Length Capping**: All inputs are capped at a maximum of **34 characters** (the maximum theoretical IBAN length) to prevent memory allocation attacks and DoS.
- **Strict Character Sets**: Only `A-Z` and `0-9` are permitted in normalized IBANs.

### 2. Algorithmic Integrity
- **MOD-97 Streaming**: The IBAN checksum uses a streaming MOD-97 algorithm that avoids large integer overflows and ensures constant-time performance for the validation loop.

### 3. Dependency Management
- **Minimal Dependencies**: The core business logic is implemented using the Go standard library only, reducing the attack surface from third-party supply chain vulnerabilities.

---

## Best Practices for Consumers

1. **Normalize Early**: Always pass raw input through `financial.Service.Validate` before storing or processing it further.
2. **Prefer Explicit Hints**: When identifier type is known, pass a `IdentifierType` hint to reduce ambiguity and avoid accidental auto-detection drift.
3. **Use `Must()` with Caution**: The `Must()` wrappers panic on error and should be limited to scripts/tooling where panic is acceptable.

---

## Reporting Vulnerabilities

If you discover a security vulnerability, do **not** open a public issue.
Please use GitHub private vulnerability reporting:
- [Report a vulnerability](https://github.com/SamyRai/bank-data/security/advisories/new)

We aim to acknowledge all security reports within 48 hours and provide a fix or mitigation path within 1 week for critical issues.
