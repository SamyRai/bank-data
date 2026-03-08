package iban

import "github.com/SamyRai/bank-data/pkg/bank"

// Validator validates an IBAN string.
type Validator interface {
	Validate(iban string) error
}

// Parser parses an IBAN string into structured components.
type Parser interface {
	Parse(iban string) (*IBANInfo, error)
}

// Detector provides country-level IBAN structure metadata.
type Detector interface {
	Detect(iban string) (*IBANStructure, error)
}

// BankLookup is an optional enrichment seam for bank metadata.
type BankLookup interface {
	LookupBank(country, bankCode string) (*bank.BankInfo, bool)
}
