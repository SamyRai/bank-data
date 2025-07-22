# Future Data Types & Validations (Europe Focus)

| Data Type         | Validation Methods                |
|-------------------|----------------------------------|
| IBAN              | Format, length, MOD-97 checksum, country rules |
| BIC/SWIFT         | Regex, length, country/bank lookup            |
| SEPA Creditor ID  | Regex, length, MOD-97, country rules         |
| National Accounts | Country-specific, format, checksum           |
| Card Number (PAN) | Luhn, length, prefix                         |
| VAT Number        | Regex, country rules, checksum               |
| LEI               | Length, MOD-97-10                            |
| ISIN              | Luhn, country code                           |
| Card Expiry       | Format, date logic                           |
| CVC/CVV           | Length, numeric                              |
| Sort Code         | Format, modulus (some)                       |
| Tax IDs           | Regex, length, checksum                      |
| EORI              | Regex, length                                |

### Details and Examples

- **IBAN**: International Bank Account Number. Format, length, and MOD-97 checksum validation. Country-specific rules.
- **BIC (SWIFT)**: 8 or 11 alphanumeric. Regex, length, and country/bank code lookup.
- **SEPA Creditor Identifier**: Country code, check digits, business code, national ID. Regex, length, MOD-97.
- **National Bank Account Numbers**: e.g., German Kontonummer/BLZ, French RIB, Spanish CCC. Country-specific format and checksum.
- **Card Number (PAN)**: 13–19 digits. Luhn algorithm, length, prefix (BIN/IIN).
- **VAT Number (VATIN)**: Country code + alphanumeric. Regex, country rules, some with checksums.
- **LEI**: Legal Entity Identifier. 20 alphanumeric, ISO 17442, MOD-97-10.
- **ISIN**: International Securities Identification Number. 12 alphanumeric, Luhn checksum.
- **Card Expiry Date**: MM/YY or MM/YYYY. Format, valid month/year, not in the past.
- **CVC/CVV**: 3 or 4 digits. Length, numeric.
- **Sort Code**: UK/Ireland, 6 digits. Format, sometimes modulus checks.
- **Tax IDs**: e.g., NIF/NIE (Spain), Codice Fiscale (Italy), TINs. Regex, length, checksum.
- **EORI**: Economic Operators Registration and Identification. Country code + up to 15 alphanumeric. Regex, length.
