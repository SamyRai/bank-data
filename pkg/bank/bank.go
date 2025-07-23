// Package bank provides a future-proof, extensible domain entity for banking institutions.
package bank

import "github.com/SamyRai/bank-data/internal/bankdata"

// BranchInfo represents a branch of a bank.
type BranchInfo struct {
	BranchCode string
	Address    string
	City       string
	Contact    string // Optional: phone/email
}

// Info represents a banking institution with global identifiers and metadata.
type Info struct {
	CountryCode string // ISO 3166-1 alpha-2
	BankCode    string // National bank code (e.g., BLZ, Sort Code, ABI, etc.)
	BIC         string // SWIFT BIC
	BankName    string
	Address     string            // Optional: physical address
	City        string            // Optional
	OtherCodes  map[string]string // e.g., {"SWIFT": "...", "ABA": "..."}
	Branches    []BranchInfo      // List of branches for this bank
}

// NewBankInfo creates a new BankInfo instance.
func NewBankInfo(country, bankCode, bic, name string) *Info {
	return &Info{
		CountryCode: country,
		BankCode:    bankCode,
		BIC:         bic,
		BankName:    name,
		OtherCodes:  make(map[string]string),
	}
}

// SetAddress sets the address and city for the bank.
func (b *Info) SetAddress(address, city string) {
	b.Address = address
	b.City = city
}

// AddCode adds or updates an additional code (e.g., ABA, SWIFT, etc.).
func (b *Info) AddCode(codeType, code string) {
	if b.OtherCodes == nil {
		b.OtherCodes = make(map[string]string)
	}
	b.OtherCodes[codeType] = code
}

// AddBranch adds a branch to the bank.
func (b *Info) AddBranch(branch BranchInfo) {
	b.Branches = append(b.Branches, branch)
}

// GetBranches returns all branches for the bank.
func (b *Info) GetBranches() []BranchInfo {
	return b.Branches
}

// sepaCountries is a set of ISO country codes in SEPA.
var sepaCountries = map[string]struct{}{
	"AT": {}, "BE": {}, "BG": {}, "CH": {}, "CY": {}, "CZ": {}, "DE": {}, "DK": {}, "EE": {}, "ES": {}, "FI": {}, "FR": {}, "GB": {}, "GR": {}, "HR": {}, "HU": {}, "IE": {}, "IS": {}, "IT": {}, "LI": {}, "LT": {}, "LU": {}, "LV": {}, "MC": {}, "MT": {}, "NL": {}, "NO": {}, "PL": {}, "PT": {}, "RO": {}, "SE": {}, "SI": {}, "SK": {}, "SM": {}, "VA": {}, "GI": {}, "AX": {}, "PM": {}, "GF": {}, "GP": {}, "MQ": {}, "RE": {}, "BL": {}, "MF": {}, "YT": {}, "NC": {}, "TF": {}, "WF": {}, "EH": {}, "JE": {}, "IM": {}, "GG": {},
}

// IsSEPA returns true if the bank is in a SEPA country.
func (b *Info) IsSEPA() bool {
	_, ok := sepaCountries[b.CountryCode]
	return ok
}

// EnrichFromIBAN parses IBAN and fills CountryCode and BankCode if possible.
func (b *Info) EnrichFromIBAN(iban string) {
	iban = removeSpaces(iban)
	if len(iban) < 8 {
		return
	}
	b.CountryCode = iban[:2]
	// Example: for DE, bank code is positions 5-12 (0-based 4:12)
	switch b.CountryCode {
	case "DE":
		if len(iban) >= 12 {
			b.BankCode = iban[4:12]
		}
		// Add more country-specific logic as needed
	}
}

// EnrichFromMFI enriches Info with address, city, LEI, and category from the MFI dataset.
func (b *Info) EnrichFromMFI(records []bankdata.MFIRecord) {
	// Try to match by BIC, BankCode, or Name (case-insensitive)
	for _, rec := range records {
		if (b.BIC != "" && rec.LEI == b.BIC) ||
			(b.BankCode != "" && rec.RIADCode == b.BankCode) ||
			(b.BankName != "" && equalFold(rec.Name, b.BankName)) {
			b.Address = rec.Address
			b.City = rec.City
			b.AddCode("LEI", rec.LEI)
			b.AddCode("RIAD", rec.RIADCode)
			b.AddCode("Category", rec.Category)
			return
		}
	}
}

func equalFold(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		ca, cb := a[i], b[i]
		if ca|0x20 != cb|0x20 && ca != cb {
			return false
		}
	}
	return true
}

func removeSpaces(s string) string {
	res := make([]rune, 0, len(s))
	for _, r := range s {
		if r != ' ' {
			res = append(res, r)
		}
	}
	return string(res)
}

// Validate checks if the essential fields are present and valid.
func (b *Info) Validate() error {
	// TODO: Implement comprehensive validation [priority:high, effort:medium]
	if b.CountryCode == "" || b.BankCode == "" || b.BIC == "" {
		return ErrInvalidBankInfo
	}
	return nil
}

// ErrInvalidBankInfo is returned when BankInfo is incomplete or invalid.
var ErrInvalidBankInfo = &BankInfoError{"missing required bank info fields"}

// BankInfoError represents an error related to BankInfo.
type BankInfoError struct {
	Msg string
}

func (e *BankInfoError) Error() string {
	return e.Msg
}

// TODO: Add audit/logging hooks for compliance [priority:low, effort:high]
