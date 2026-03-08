# Bank Code Registry Research

> Research document for **Milestone 1 — Country Coverage & Data Expansion**.
> Covers official data sources, formats, licensing, update cadence, and Go integration notes for each target country.
> Last updated: 2026-03-07

---

## Overview

Each country entry follows the same structure:

| Field | Meaning |
|---|---|
| **National code** | The domestic bank identifier embedded in the IBAN |
| **IBAN structure** | Positions of bank code within the country's IBAN |
| **Official source** | The authoritative body publishing the dataset |
| **Data URL / API** | Direct link to download or API endpoint |
| **Format** | File format(s) available |
| **Update frequency** | How often the source publishes updates |
| **License / access** | Free / registration required / commercial |
| **Go integration notes** | Practical implementation guidance |

---

## 🇩🇪 DE (Germany) — BLZ *(already implemented, reference)*

| Field | Value |
|---|---|
| National code | **BLZ** (Bankleitzahl) — 8 digits |
| IBAN structure | `DE` + 2 check + **8-digit BLZ** + 10-digit account |
| Official source | Deutsche Bundesbank |
| Data URL | `https://www.bundesbank.de/de/aufgaben/unbarer-zahlungsverkehr/serviceangebot/bankleitzahlen/download-bankleitzahlen-602592` |
| Format | TXT (fixed-width), CSV |
| Update frequency | Quarterly |
| License | Free, no registration |
| Go integration | Existing implementation in `pkg/bank/` — reference model for new countries |

---

## 🇫🇷 FR (France)

| Field | Value |
|---|---|
| National code | **Code banque** (5 digits) + **Code guichet** (5 digits) |
| IBAN structure | `FR` + 2 check + **5-digit bank code** + 5-digit branch + 11-digit account + 2 national check |
| Official source | Banque de France (BdF) Webstat data portal |
| Data URL | `https://webstat.banque-france.fr/` (API key required) |
| Format | JSON, XLSX, CSV (via Webstat API) |
| Update frequency | Monthly |
| License | Free with registration (API key) |
| Supplemental | **SWIFT SwiftRef** BIC directory — `https://www.swiftref.com/` (licensed) |
| Supplemental | **OpenSanctions BIC dataset** — daily, free for non-commercial: `https://www.opensanctions.org/datasets/bic/` |

### Notes
- The Banque de France does **not** publish a standalone BIC-to-bank-code mapping for download. The Webstat API covers monetary/statistical data.
- **Recommended approach**: Combine Banque de France's open statistical register (for French bank entity names/IDs) with the **SWIFT ISO 9362 BIC directory** or the **OpenSanctions BIC dataset** (free, daily-updated JSON/CSV, community-friendly license).
- French IBAN bank code is a 5-digit "code banque" assigned by BdF — a lookup table of ~1,000 active institutions suffices.
- The community dataset at `github.com/jschauma/bics` compiles SWIFT data into machine-readable JSON and may serve as a bootstrap.

### Go Integration
```go
// pkg/bank/fr/loader.go
// Source: OpenSanctions BIC CSV + BdF Banque dataset
// Fields: BankCode (5 chars), Name, BIC, City
type FRBankRecord struct {
    CodeBanque string // 5-digit
    Name       string
    BIC        string
    City       string
}
```

---

## 🇦🇹 AT (Austria)

| Field | Value |
|---|---|
| National code | **BLZ** (Bankleitzahl) — 5 digits |
| IBAN structure | `AT` + 2 check + **5-digit BLZ** + 11-digit account |
| Official source | Oesterreichische Nationalbank (OeNB) |
| Data URL | `https://www.oenb.at/Zahlungsverkehr/Zahlungsverkehrsverzeichnis.html` → "SEPA-Zahlungsverkehrs-Verzeichnis (gesamt)" |
| Format | **CSV** |
| Update frequency | Quarterly (with interim updates for mergers/closures) |
| License | Free, no registration — official public dataset |
| Fields included | BLZ, Bank name, BIC, Address, SEPA participation flag |

### Notes
- **Best-in-class source among all target countries**: official CSV published directly by the central bank, free, complete BLZ→BIC mapping.
- Schema is stable; column order has not changed since 2016.
- Approx. ~700 active entries.

### Go Integration
```go
// pkg/bank/at/loader.go
// Direct CSV download from OeNB
// Columns: BLZ | Name | BIC | PLZ | Ort | ...
type ATBankRecord struct {
    BLZ  string // 5-digit
    Name string
    BIC  string
    ZIP  string
    City string
}
```

**Automation**: Use `go:generate` + `curl` in `gen/at/` to fetch and embed via `go:embed`.

---

## 🇳🇱 NL (Netherlands)

| Field | Value |
|---|---|
| National code | No standalone national bank code system; BIC is the primary bank identifier |
| IBAN structure | `NL` + 2 check + **4-char BIC prefix (bank code)** + 10-digit account |
| Official source | De Nederlandsche Bank (DNB) — no public bulk download |
| Data URL | `https://www.dnb.nl/` (no direct CSV/API for BIC directory) |
| Supplemental 1 | **OpenSanctions BIC** — `https://www.opensanctions.org/datasets/bic/` (JSON/CSV, daily, free) |
| Supplemental 2 | **SWIFT SwiftRef BIC Directory** — licensed bulk download |
| Supplemental 3 | De Nederlandsche Bank supervised entities register (Excel) — institutions list only |
| Format | JSON, CSV (OpenSanctions); Excel (DNB) |
| Update frequency | Daily (OpenSanctions); monthly (SWIFT) |
| License | OpenSanctions: ODbL (open); SWIFT: commercial |

### Notes
- The Dutch IBAN embeds a **4-char bank code** (position 5–8), which directly maps to the first 4 chars of the bank's BIC.
- No equivalent of the German BLZ or Austrian BLZ exists in the Netherlands; BIC is used natively.
- **Practical approach**: Parse 4-char IBAN bank codes and resolve via OpenSanctions SWIFT BIC dataset.
- DNB publishes a supervised-entities register but it does not include the bank-code-to-BIC mapping needed.

### Go Integration
```go
// pkg/bank/nl/loader.go
// 4-char bank code extracted from IBAN positions 4–8
// Resolved via OpenSanctions BIC JSON dataset
type NLBankRecord struct {
    BankCode string // 4-char (IBAN chars 5-8)
    BIC      string
    Name     string
}
```

---

## 🇪🇸 ES (Spain)

| Field | Value |
|---|---|
| National code | **CCC bank code** — first 4 digits of the 20-digit CCC (Código Cuenta Cliente) |
| IBAN structure | `ES` + 2 check + **4-digit bank code** + 4-digit branch + 2 CCC check + 10-digit account |
| Official source | Banco de España (BdE) |
| Data URL | `https://www.bde.es/webbde/es/estadis/infoest/si_1_1.pdf` (institution list PDF); statistical API at `https://www.bde.es/webbde/en/estadis/webbde_es_estadis_descarga.html` |
| Format | JSON, CSV, XLS (BdE statistical API); PDF for institution list |
| Update frequency | Monthly |
| License | Free, no registration for API; registration for full datasets |
| Coverage | ~200 active bank codes |

### Notes
- CCC was replaced by IBAN for all operations since January 2016 but **still embedded** within ES IBAN.
- Banco de España publishes the **Registro de Entidades** (institution register), which includes 4-digit bank codes with institution names — downloadable as Excel/PDF.
- There is **no public BdE CSV mapping bank-code → BIC**; augment with SWIFT SwiftRef or OpenSanctions.
- The `ES` IBAN bank code (4 digits) is assigned by BdE and stable for established institutions.

### Go Integration
```go
// pkg/bank/es/loader.go
// 4-digit BdE entity code from IBAN
// Banco de España Registro de Entidades → supplemented with OpenSanctions BIC
type ESBankRecord struct {
    EntityCode string // 4-digit BdE code
    Name       string
    BIC        string
}
```

---

## 🇮🇹 IT (Italy)

| Field | Value |
|---|---|
| National code | **ABI code** — 5 digits; **CAB code** (branch) — 5 digits |
| IBAN structure | `IT` + 2 check + 1 CIN char + **5-digit ABI** + 5-digit CAB + 12-digit account |
| Official source | Banca d'Italia — Supervisory Registers |
| Data URL | `https://www.bancaditalia.it/compiti/vigilanza/albi-elenchi/index.html` |
| Format | Excel (.xls/.xlsx), PDF |
| Update frequency | Monthly (supervisory registers) |
| License | Free, no registration |
| Supplemental | CodiciBancari.it — ABI + CAB + BIC + branch lookup (~93,000 branches, last full update 2016, partial updates continue) |
| Supplemental | SWIFT SwiftRef for BIC mapping |

### Notes
- Banca d'Italia provides the **Albo delle Banche** (banking register) in Excel — contains ABI codes, institution names, and status.
- ABI→BIC mapping is **not directly published** by Banca d'Italia; augment with SWIFT data or `bancheitalia.it`.
- CodiciBancari.it is the most comprehensive community source for ABI+CAB+BIC in machine-readable form.
- Italy has ~700 active ABI codes (banks); ~93,000 branches (CAB codes) — store only head-office ABI→BIC, branch CAB is optional for basic validation.

### Go Integration
```go
// pkg/bank/it/loader.go
// 5-digit ABI from IBAN, resolved to BIC
// Source: Banca d'Italia Excel register + SWIFT BIC augmentation
type ITBankRecord struct {
    ABI  string // 5-digit
    Name string
    BIC  string
}
```

---

## 🇵🇱 PL (Poland)

| Field | Value |
|---|---|
| National code | **NRB bank ID** — first 8 digits of the 26-digit NRB (includes 3-digit bank group + 5-digit branch) |
| IBAN structure | `PL` + 2 IBAN check + **8-digit NRB bank/branch code** + 16-digit account |
| Official source | Krajowa Izba Rozliczeniowa (KIR) — ELIXIR system participant list |
| Data URL | `https://www.kir.pl/` (no public bulk CSV, participant info only) |
| Supplemental | National Bank of Poland (NBP) statistical publications: `https://www.nbp.pl/` |
| Supplemental | OpenSanctions BIC dataset for Polish banks |
| Format | No official open-data CSV; NBP publishes supervised entities as PDF/XLS |
| Update frequency | KIR: real-time (internal); Public datasets: quarterly |
| License | KIR data: not publicly downloadable; NBP: free |

### Notes
- **Most difficult country** among the 9: KIR does not publish a downloadable NRB→BIC registry for public use.
- Polish NRB bank code = 8 digits (first 8 of 26-digit account), where digits 1–3 = bank group, 4–8 = branch.
- **Practical approach**:
  1. Curate a static mapping of major Polish banks (PKO, Pekao, mBank, Santander PL, etc.) from public sources
  2. Use OpenSanctions BIC dataset to cross-reference BIC codes
  3. Implement NRB checksum validation (MOD-97) independently of bank code lookup
- Consider contributing to / forking `github.com/kbkk/ibantools` Polish data.

### Go Integration
```go
// pkg/bank/pl/loader.go
// Static curated map + OpenSanctions BIC augmentation
// NRB validation is separable from bank lookup
type PLBankRecord struct {
    BankGroupCode string // 3-digit prefix
    Name          string
    BIC           string
}

// Standalone NRB checksum (MOD-97) validation available without lookup
func ValidateNRB(nrb string) error { ... }
```

---

## 🇸🇪 SE (Sweden)

| Field | Value |
|---|---|
| National code | **Clearing number** — 4 or 5 digits (varies by bank; some use 4+1 split) |
| IBAN structure | `SE` + 2 check + **3-digit bank code** + 17-digit account |
| Official source | Bankgirot — manages clearing numbers |
| Data URL | `https://www.bankgirot.se/` — search only, no download |
| Supplemental 1 | `bankinfrastruktur.se` — industry reference |
| Supplemental 2 | Community dataset: [github.com/jop-io/kontonummer](https://jop-io.github.io/kontonummer/) — JSON, clearing numbers + BICs |
| Supplemental 3 | Riksbank open data API: `https://api.riksbank.se/swea/v1/` — monetary data only, not BIC directory |
| Format | JSON (community), no official bulk CSV available |
| Update frequency | Community: semi-annual; Bankgirot: real-time internal |
| License | Community: open (MIT/CC); Bankgirot: not publicly licensed |

### Notes
- Bankgirot explicitly prohibits scraping and does **not** offer a downloadable bank code registry.
- Swedish clearing numbers are maintained by Bankgirot under bank confidentiality constraints.
- The **community GitHub project** (bankinfrastruktur.se data, compiled by `jop-io`) is the best available machine-readable source — ~30 clearing ranges covering all major Swedish banks.
- Swedish IBAN bank code (3 chars) ≠ clearing number (4–5 digits). The IBAN uses a condensed bank code; clearing numbers include branch info.

### Go Integration
```go
// pkg/bank/se/loader.go
// Embedded JSON from community dataset (bankinfrastruktur.se) 
// Clearing range: [{from: 1100, to: 1199, bic: "NDEASESS", name: "Nordea"}]
type SEClearingRange struct {
    From int
    To   int
    BIC  string
    Name string
}
```

**Automation**: Embed the community JSON at build time; add GitHub Action to check for community repo updates semi-annually.

---

## 🇨🇭 CH (Switzerland)

| Field | Value |
|---|---|
| National code | **BC-Nr** (Bankenclearingnummer / IID) — 3 to 5 digits |
| IBAN structure | `CH` + 2 check + **5-digit IID (BC-Nr)** + 12-digit account |
| Official source | SIX Group — "Download Bank Master" |
| Data URL | `https://www.six-group.com/en/products-services/banking-services/interbank-clearing/online-services/download-bank-master.html` |
| API URL | `https://www.six-group.com/en/products-services/banking-services/interbank-clearing/online-services/bank-master-rest.html` |
| Format | **CSV** (Bank Master V 3.0) + **REST API** (JSON, Bank Master REST API V 3.0) |
| Update frequency | Daily (REST API); monthly on 20th (mutations bulletin) |
| License | Free, no registration — SIX Group public service |
| Fields included | IID (BC-Nr), Branch ID, BIC, Bank name, SIC participation, euroSIC participation, address |

### Notes
- **Second best-in-class** after AT: SIX Group provides both a free CSV download and a REST JSON API, both publicly accessible.
- Mutations (additions, mergers, deletions) are published on the 20th of each month.
- BC-Nr is the key: position 4–8 of the Swiss IBAN (padded to 5 digits).
- ~1,000 active entries in the Swiss bank master.

### Go Integration
```go
// pkg/bank/ch/loader.go OR use REST API at runtime
// REST: GET https://www.six-group.com/dam/download/banking-services/interbank-clearing/en/bank-master.json
type CHBankRecord struct {
    IID      string // 3-5 digit BC-Nr (zero-padded to 5)
    BranchID string
    BIC      string
    Name     string
    SIC      bool   // participates in Swiss Interbank Clearing
    EuroSIC  bool   // participates in euroSIC
}
```

**Automation**: REST API call from `gen/ch/gen_registry.go` with `go:generate`.

---

## 🇬🇧 UK (United Kingdom)

| Field | Value |
|---|---|
| National code | **Sort Code** — 6 digits (XX-XX-XX), embedded verbatim in UK IBAN |
| IBAN structure | `GB` + 2 check + **4-char BIC prefix** + **6-digit sort code** + 8-digit account |
| Official source | Pay.UK — Extended Industry Sort Code Directory (EISCD) |
| Data URL | `https://www.wearepay.uk/what-we-do/payment-systems/sort-code-directory/` |
| Format | **TXT** (fixed-width), **XML**, **CSV**, **XLS** |
| API | Automated SFTP / Web API via licensed distributor (e.g., `sortcodes.co.uk`) |
| Update frequency | Weekly (EISCD v2); real-time via API distributors |
| License | EISCD: free download with Pay.UK registration; commercial SFTP/API access via distributor |
| Fields included | Sort code, BIC (Bank Office Identifier), Branch name, Bacs/FPS/CHAPS participation, address |

### Notes
- **Post-Brexit note**: UK left SEPA in January 2021. UK IBANs are not accepted for SEPA transactions, but are valid internationally. Validate separately from SEPA logic.
- The EISCD is the authoritative source — produced by Pay.UK and used by all UK banks.
- The "BIC" in EISCD is a **Bank Office Identifier (BOI)** — it may differ from the institution's primary SWIFT BIC.
- For sort code → SWIFT BIC mapping (international wire), use SWIFT SwiftRef or Open Banking API.
- The UK Open Banking directory (`openbanking.org.uk`) provides ASPSP (bank) BIC lookups via OAuth-protected API.
- ~20,000 active sort codes (but only ~50 distinct BIC institutions).

### Go Integration
```go
// pkg/bank/gb/loader.go
// Source: EISCD TXT/CSV (Pay.UK) — registration required
// Sort code is the primary key; BIC is the BOI from EISCD
type GBBankRecord struct {
    SortCode string // 6-digit, no separators
    BIC      string // Bank Office Identifier (8/11 char)
    Name     string
    Bacs     bool
    FPS      bool  // Faster Payments
    CHAPS    bool
}

// IBAN Sort Code is at positions 8–14 of GB IBAN
func ExtractSortCode(iban string) (string, error) {
    // "GB29NWBK60161331926819" → "601613"
}
```

**Note**: Store EISCD in `datasets/gb/eiscd_YYYY-MM.csv`; update weekly via CI fetch.

---

## Cross-Country Strategy

### Unified Registry Interface

All country loaders should satisfy the same interface in `internal/countrymeta/`:

```go
type BankRecord struct {
    CountryCode string   // ISO 3166-1 alpha-2
    NationalCode string  // Country's domestic bank code (BLZ, BIC prefix, ABI, etc.)
    BIC         string   // SWIFT BIC (8 or 11 chars)
    Name        string   // Bank display name
    City        string
    IsActive    bool
}

type BankRegistry interface {
    Lookup(countryCode, nationalCode string) (*BankRecord, error)
    All(countryCode string) ([]BankRecord, error)
}
```

---

## Data Source Priority Matrix

| Country | Source Quality | Free | Bulk Download | BIC Included | Update Freq |
|---------|---------------|------|---------------|--------------|-------------|
| **DE** | ⭐⭐⭐⭐⭐ Official | ✅ | ✅ CSV | ✅ | Quarterly |
| **AT** | ⭐⭐⭐⭐⭐ Official | ✅ | ✅ CSV | ✅ | Quarterly |
| **CH** | ⭐⭐⭐⭐⭐ Official | ✅ | ✅ CSV + REST | ✅ | Daily |
| **UK** | ⭐⭐⭐⭐ Official | ⚠️ Registration | ✅ TXT/CSV | ✅ BOI | Weekly |
| **IT** | ⭐⭐⭐ Official (partial) | ✅ | ✅ Excel | ❌ (augment) | Monthly |
| **ES** | ⭐⭐⭐ Official (partial) | ✅ | ✅ JSON/XLS | ❌ (augment) | Monthly |
| **FR** | ⭐⭐ Official + 3rd party | ✅/⚠️ | ⚠️ API key | ❌ (augment) | Monthly |
| **NL** | ⭐⭐ 3rd party | ✅ | ✅ JSON | ✅ | Daily |
| **SE** | ⭐⭐ Community | ✅ | ✅ JSON | ✅ | Semi-annual |
| **PL** | ⭐ Curated/Community | ✅ | ❌ (static) | ⚠️ Partial | Manual |

---

## Automated Data Refresh Pipeline (`gen/` tooling)

### Proposed `gen/` structure

```
gen/
├── at/
│   └── gen_registry.go   // curl OeNB CSV → parse → embed
├── ch/
│   └── gen_registry.go   // SIX Group REST API → parse → embed
├── de/
│   └── gen_registry.go   // Bundesbank TXT → parse → embed (existing)
├── gb/
│   └── gen_registry.go   // EISCD CSV (pre-downloaded) → parse → embed
├── it/
│   └── gen_registry.go   // Banca d'Italia Excel → parse → embed
├── es/
│   └── gen_registry.go   // BdE JSON API → parse → embed
├── fr/
│   └── gen_registry.go   // OpenSanctions BIC JSON → filter FR → embed
├── nl/
│   └── gen_registry.go   // OpenSanctions BIC JSON → filter NL → embed
├── pl/
│   └── gen_registry.go   // Static JSON curation → embed
└── se/
    └── gen_registry.go   // Community JSON → embed
```

Each generator:
1. Fetches data from the official source (or a stable mirror)
2. Validates schema / record count (CI gate: fail if < expected minimum)
3. Parses into `[]BankRecord`
4. Writes `internal/countrymeta/<cc>/registry_gen.go` with `//go:embed` or a compiled-in map

**Run**: `go generate ./gen/...`

### `base.go` directive (per country)
```go
// in internal/countrymeta/at/base.go
//go:generate go run ../../../gen/at/gen_registry.go
```

---

## Dataset Versioning Strategy (`datasets/`)

### Directory layout
```
datasets/
├── de/
│   ├── blz_2026-01-01.txt          ← source file, dated
│   └── blz_current -> blz_2026-01.txt  ← symlink to latest
├── at/
│   ├── sepa_verzeichnis_2026-01-01.csv
│   └── sepa_verzeichnis_current -> ...
├── ch/
│   ├── bank_master_2026-01-20.csv
│   └── bank_master_current -> ...
...
```

### Metadata sidecar (`.meta.json`)
Each dataset file ships with a sidecar:
```json
{
  "country": "AT",
  "source_url": "https://www.oenb.at/...",
  "source_name": "OeNB SEPA-Zahlungsverkehrs-Verzeichnis",
  "fetched_at": "2026-01-02T08:00:00Z",
  "record_count": 712,
  "sha256": "abc123...",
  "license": "public-domain",
  "next_refresh": "2026-04-01"
}
```

### CI Refresh Workflow (`.github/workflows/refresh-registries.yml`)
```yaml
on:
  schedule:
    - cron: '0 6 1 */3 *'   # Quarterly on the 1st
  workflow_dispatch:

jobs:
  refresh:
    steps:
      - run: go generate ./gen/...
      - run: go test ./...  # Gate: fail if record count drops > 10%
      - uses: peter-evans/create-pull-request@v6  # Auto-PR with diff
```

---

## Per-Country Integration Test Plan

Each country must have a test file at `pkg/bank/<cc>/<cc>_test.go`:

### Test structure
```go
func TestLookup_RealIBANSamples(t *testing.T) {
    testCases := []struct {
        iban        string
        wantBIC     string
        wantName    string
    }{
        // Use real, publicly known IBANs (central bank, example IBANs from EPC)
        {"AT61 1904 3002 3457 3201", "OPSKATWW", "Österreichische Postsparkasse"},
        // ...
    }
    svc := iban.NewService(nil, nil, nil, nil)
    for _, tc := range testCases {
        info, err := svc.Parse(tc.iban)
        require.NoError(t, err)
        assert.Equal(t, tc.wantBIC, info.BIC)
        assert.Equal(t, tc.wantName, info.BankName)
    }
}
```

### Test IBANs sources per country
| Country | Source for test IBANs |
|---|---|
| **DE** | Bundesbank published test IBANs |
| **AT** | OeNB validation examples |
| **CH** | SIX Group published test IBANs |
| **UK** | Pay.UK `EISCD` documentation examples |
| **IT** | Banca d'Italia IBAN examples + EPC IBAN Checker |
| **ES** | EPC official IBAN test vectors |
| **FR** | BdF official IBAN test vectors |
| **NL** | EPC official IBAN test vectors |
| **SE** | Bankgirot documentation |
| **PL** | KIR ELIXIR documentation examples |

### Minimum coverage gates (per `pkg/bank/<cc>/`)
- ≥ 5 real-world IBAN parse tests (known BIC outcome)
- ≥ 2 invalid IBAN rejection tests
- ≥ 1 edge-case test (closed bank, historic BLZ, etc.)
- Benchmark: `BenchmarkLookup` — gate at < 1µs/op after warm cache

---

## Implementation Order (recommended)

Given data availability and ease of integration:

1. 🟢 **CH** — SIX Group REST API, free, JSON, daily → implement first as template
2. 🟢 **AT** — OeNB CSV, free, direct download → trivial
3. 🟡 **UK** — Pay.UK EISCD, registration + CSV → high value, slightly more setup
4. 🟡 **IT** — Banca d'Italia Excel + augment BIC → medium effort
5. 🟡 **ES** — BdE API + augment BIC → medium effort  
6. 🟠 **FR** — Banque de France API + OpenSanctions → multiple sources needed
7. 🟠 **NL** — OpenSanctions BIC filter → simple but no native bank code
8. 🟠 **SE** — Community JSON embed → stable but community-maintained
9. 🔴 **PL** — Static curated map + NRB checksum → manual maintenance, revisit if KIR opens API

---

*See also: [TODO.md](../TODO.md) · [architecture.md](architecture.md) · [FUTURE_FEATURES.md](../FUTURE_FEATURES.md)*
