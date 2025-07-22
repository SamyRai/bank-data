// Package iban provides interfaces for IBAN validation, parsing, and information extraction.
package iban

// Validator defines the interface for IBAN validation.
type Validator interface {
	Validate(iban string) error
}

// Parser defines the interface for IBAN parsing and information extraction.
type Parser interface {
	Parse(iban string) (*IBANInfo, error)
}

// Detector defines the interface for IBAN country and structure detection.
type Detector interface {
	Detect(iban string) (*IBANStructure, error)
}
