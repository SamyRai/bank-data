// Package iban provides interfaces for IBAN validation, parsing, and information extraction.
package iban

import (
	"github.com/SamyRai/bank-data/internal/bic/map"
	"github.com/SamyRai/bank-data/pkg/bank"
)

// Validator defines the interface for IBAN validation.
// Validate returns an error if the IBAN is invalid, or nil if valid.
type Validator interface {
	Validate(iban string) error
	ValidateAndBankInfo(ibanStr string, bicMap bicmap.BankBICMap) (*bank.BankInfo, error)
}

// Parser defines the interface for IBAN parsing and information extraction.
// Parse returns IBANInfo and error if parsing fails.
type Parser interface {
	Parse(iban string) (*IBANInfo, error)
}

// Detector defines the interface for IBAN country and structure detection.
// Detect returns IBANStructure and error if detection fails.
type Detector interface {
	Detect(iban string) (*IBANStructure, error)
}
