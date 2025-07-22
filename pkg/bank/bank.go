// Package bank provides a future-proof, extensible domain entity for banking institutions.
package bank

// BankInfo represents a banking institution with global identifiers and metadata.
type BankInfo struct {
	CountryCode string // ISO 3166-1 alpha-2
	BankCode    string // National bank code (e.g., BLZ, Sort Code, ABI, etc.)
	BIC         string // SWIFT BIC
	BankName    string
	Address     string            // Optional: physical address
	City        string            // Optional
	OtherCodes  map[string]string // e.g., {"SWIFT": "...", "ABA": "..."}
}

// NewBankInfo creates a new BankInfo instance.
func NewBankInfo(country, bankCode, bic, name string) *BankInfo {
	return &BankInfo{
		CountryCode: country,
		BankCode:    bankCode,
		BIC:         bic,
		BankName:    name,
		OtherCodes:  make(map[string]string),
	}
}

// SetAddress sets the address and city for the bank.
func (b *BankInfo) SetAddress(address, city string) {
	b.Address = address
	b.City = city
}

// AddCode adds or updates an additional code (e.g., ABA, SWIFT, etc.).
func (b *BankInfo) AddCode(codeType, code string) {
	if b.OtherCodes == nil {
		b.OtherCodes = make(map[string]string)
	}
	b.OtherCodes[codeType] = code
}

// IsSEPA returns true if the bank is in a SEPA country (future-proof, stub).
func (b *BankInfo) IsSEPA() bool {
	// TODO: Implement SEPA country check [priority:medium, effort:low]
	return false
}

// Validate checks if the essential fields are present and valid.
func (b *BankInfo) Validate() error {
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

// TODO: Add enrichment methods (e.g., from IBAN, from registry) [priority:high, effort:medium]
// TODO: Add support for multiple branches per bank [priority:medium, effort:high]
// TODO: Add audit/logging hooks for compliance [priority:low, effort:high]
