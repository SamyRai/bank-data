// Package iban provides the public API for IBAN validation, parsing, and detection.
package iban

// IBANInfo holds parsed IBAN details.
type IBANInfo struct {
	CountryCode   string
	BankCode      string
	BranchCode    string
	AccountNumber string
	CheckDigits   string
	Raw           string
}

// IBANStructure holds metadata about IBAN format for a country.
type IBANStructure struct {
	CountryCode string
	Length      int
	Structure   string // e.g. "CCKKBBBBBBBBCCCCCCCCCC"
}

// IBANErrorCode defines error codes for IBAN operations.
type IBANErrorCode string

const (
	ErrCodeInvalidChars       IBANErrorCode = "invalid_chars"
	ErrCodeWrongLength        IBANErrorCode = "wrong_length"
	ErrCodeChecksum           IBANErrorCode = "checksum_failed"
	ErrCodeUnsupportedCountry IBANErrorCode = "unsupported_country"
)

// IBANError provides structured error information for IBAN operations.
type IBANError struct {
	Code    IBANErrorCode // Machine-readable error code
	Field   string        // Optional: field or aspect (e.g., "country", "length", "checksum")
	Message string        // Human-readable error message
	Value   string        // Optional: the value that caused the error
}

func (e *IBANError) Error() string {
	if e.Field != "" {
		return e.Field + ": " + e.Message
	}
	return e.Message
}

// NewIBANError constructs a new IBANError.
func NewIBANError(code IBANErrorCode, field, message, value string) *IBANError {
	return &IBANError{
		Code:    code,
		Field:   field,
		Message: message,
		Value:   value,
	}
}

// Public, typed errors for all major failure modes.
var (
	ErrInvalidChars       = NewIBANError(ErrCodeInvalidChars, "characters", "IBAN contains invalid characters (only A-Z, 0-9 allowed)", "")
	ErrWrongLength        = NewIBANError(ErrCodeWrongLength, "length", "IBAN has wrong length for country or is too short", "")
	ErrChecksum           = NewIBANError(ErrCodeChecksum, "checksum", "IBAN checksum validation failed", "")
	ErrUnsupportedCountry = NewIBANError(ErrCodeUnsupportedCountry, "country", "IBAN country code is not supported", "")
)

// Service is the main entry point for IBAN operations.
type Service struct {
	Validator Validator
	Parser    Parser
	Detector  Detector
}

// NewService constructs a Service with the provided implementations.
func NewService(v Validator, p Parser, d Detector) *Service {
	return &Service{
		Validator: v,
		Parser:    p,
		Detector:  d,
	}
}

// Validate checks if the IBAN is valid (format and checksum).
func (s *Service) Validate(ibanStr string) error {
	return s.Validator.Validate(ibanStr)
}

// Parse extracts IBANInfo from the given IBAN string.
func (s *Service) Parse(ibanStr string) (*IBANInfo, error) {
	return s.Parser.Parse(ibanStr)
}

// Detect returns IBAN structure metadata for a given IBAN string.
func (s *Service) Detect(ibanStr string) (*IBANStructure, error) {
	return s.Detector.Detect(ibanStr)
}

// (Batch validation removed from public Service API. Use internal/validation for batch operations.)
