// Package iban provides the public API for IBAN validation, parsing, and detection.
package iban

import (
	"context"
	"iter"

	bicmap "github.com/SamyRai/bank-data/internal/bic/map"
	"github.com/SamyRai/bank-data/pkg/bank"
	"github.com/SamyRai/bank-data/pkg/bic"
	"github.com/SamyRai/bank-data/pkg/sepa"
	"github.com/SamyRai/bank-data/pkg/validation"
)

// IBANInfo holds parsed IBAN details such as country code, bank code, account number, and check digits.
type IBANInfo struct {
	CountryCode   string // ISO country code (2 letters)
	BankCode      string // Bank identifier (country-specific)
	BranchCode    string // Branch identifier (if applicable)
	AccountNumber string // Account number (country-specific)
	CheckDigits   string // IBAN check digits (2 digits)
	BIC           string // BIC/SWIFT code (if found)
	BankName      string // Bank name (if found)
	Raw           string // Original normalized IBAN string
}

// IBANStructure holds metadata about IBAN format for a country.
type IBANStructure struct {
	CountryCode string // ISO country code
	Length      int    // IBAN length for the country
	Structure   string // Structure string (e.g. "CCKKBBBBBBBBCCCCCCCCCC")
}

// IBANErrorCode defines error codes for IBAN operations.
type IBANErrorCode string

const (
	// ErrCodeInvalidChars indicates invalid characters in the IBAN.
	ErrCodeInvalidChars IBANErrorCode = "invalid_chars"
	// ErrCodeWrongLength indicates the IBAN has the wrong length.
	ErrCodeWrongLength IBANErrorCode = "wrong_length"
	// ErrCodeChecksum indicates a failed checksum validation.
	ErrCodeChecksum IBANErrorCode = "checksum_failed"
	// ErrCodeUnsupportedCountry indicates the country is not supported.
	ErrCodeUnsupportedCountry IBANErrorCode = "unsupported_country"
)

// IBANError provides structured error information for IBAN operations.
type IBANError struct {
	Code    IBANErrorCode // Machine-readable error code
	Field   string        // Field or aspect (e.g., "country", "length", "checksum")
	Message string        // Human-readable error message
	Value   string        // The value that caused the error
}

// Error returns the error message for IBANError.
func (e *IBANError) Error() string {
	if e.Field != "" {
		return e.Field + ": " + e.Message
	}
	return e.Message
}

// NewIBANError constructs a new IBANError with the given code, field, message, and value.
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
	// ErrInvalidChars is returned when the IBAN contains invalid characters.
	ErrInvalidChars = NewIBANError(ErrCodeInvalidChars, "characters", "IBAN contains invalid characters (only A-Z, 0-9 allowed)", "")
	// ErrWrongLength is returned when the IBAN has the wrong length for the country or is too short.
	ErrWrongLength = NewIBANError(ErrCodeWrongLength, "length", "IBAN has wrong length for country or is too short", "")
	// ErrChecksum is returned when the IBAN checksum validation fails.
	ErrChecksum = NewIBANError(ErrCodeChecksum, "checksum", "IBAN checksum validation failed", "")
	// ErrUnsupportedCountry is returned when the IBAN country code is not supported.
	ErrUnsupportedCountry = NewIBANError(ErrCodeUnsupportedCountry, "country", "IBAN country code is not supported", "")
)

// Service is the main entry point for IBAN operations. It provides validation, parsing, and structure detection.
type Service struct {
	Validator Validator
	Parser    Parser
	Detector  Detector
	BICMap    *bicmap.BankBICMap             // For bank metadata lookups
	Registry  *validation.ValidationRegistry // For extensible, tag-based validation
}

// NewService constructs a Service with the provided Validator, Parser, Detector, and optional ValidationRegistry.
// If any component is nil, the default implementation is used.
func NewService(v Validator, p Parser, d Detector, reg *validation.ValidationRegistry) *Service {
	if v == nil {
		v = NewValidator()
	}
	if p == nil {
		p = NewParser()
	}
	if d == nil {
		d = NewDetector()
	}
	if reg == nil {
		reg = validation.NewValidationRegistry()
	}

	// Register default validators
	reg.Register("iban", v)
	reg.Register("bic", bic.NewValidator())
	reg.Register("sci", sepa.NewCreditorIDValidator())

	bm, _ := bicmap.DefaultLoader() // best effort

	return &Service{
		Validator: v,
		Parser:    p,
		Detector:  d,
		BICMap:    bm,
		Registry:  reg,
	}
}

// Validate checks if the IBAN is valid (format and checksum). Returns an error if invalid, or nil if valid.
func (s *Service) Validate(ibanStr string) error {
	return s.Validator.Validate(ibanStr)
}

// ValidateByTag validates input using a registered validator by tag (e.g., "iban", "bic").
// Returns a ValidationResult and error if the tag is not found or validation fails.
func (s *Service) ValidateByTag(tag string, input string) (validation.ValidationResult, error) {
	v := s.Registry.Get(tag)
	if v == nil {
		return validation.ValidationResult{Input: input, Valid: false, Error: ErrUnsupportedCountry}, ErrUnsupportedCountry
	}
	// Try to use the generic Validator interface
	if validator, ok := v.(validation.Validator[string, validation.ValidationResult]); ok {
		res := validator.Validate(input)
		return res, res.Error
	}
	return validation.ValidationResult{Input: input, Valid: false, Error: ErrUnsupportedCountry}, ErrUnsupportedCountry
}

// ValidateBatch validates a slice of IBANs in parallel.
func (s *Service) ValidateBatch(ctx context.Context, inputs []string) []error {
	engine := validation.NewEngine[string, error](s.Validator, 0)
	return engine.ValidateBatch(ctx, inputs)
}

// StreamValidate provides a streaming iterator for IBAN validation.
func (s *Service) StreamValidate(ctx context.Context, inputs iter.Seq[string]) iter.Seq[error] {
	engine := validation.NewEngine[string, error](s.Validator, 0)
	return engine.StreamValidate(ctx, inputs)
}

// Parse extracts IBANInfo from the given IBAN string. Returns IBANInfo and error if parsing fails.
// It automatically enriches the result with bank information if available in the Service's BICMap.
func (s *Service) Parse(ibanStr string) (*IBANInfo, error) {
	info, err := s.Parser.Parse(ibanStr)
	if err != nil {
		return nil, err
	}
	// Enrich with bank info if mapping is available
	if s.BICMap != nil {
		if bankInfo, ok := s.BICMap.LookupBankInfo(info.CountryCode, info.BankCode); ok {
			info.BIC = bankInfo.BIC
			info.BankName = bankInfo.BankName
		}
	}
	return info, nil
}

// Detect returns IBANStructure metadata for a given IBAN string. Returns IBANStructure and error if detection fails.
func (s *Service) Detect(ibanStr string) (*IBANStructure, error) {
	return s.Detector.Detect(ibanStr)
}

// LookupBank returns bank metadata for a given country and bank code.
func (s *Service) LookupBank(country, bankCode string) (*bank.BankInfo, error) {
	if s.BICMap == nil {
		return nil, ErrBankInfoNotFound
	}
	info, ok := s.BICMap.LookupBankInfo(country, bankCode)
	if !ok {
		return nil, ErrBankInfoNotFound
	}
	return info, nil
}

// EnrichWithBankInfo looks up and attaches BankInfo to a parsed IBANInfo using the provided mapping.
func EnrichWithBankInfo(info *IBANInfo, bicMap *bicmap.BankBICMap) (*bank.BankInfo, error) {
	if info == nil {
		return nil, ErrNilIBANInfo
	}
	if bicMap == nil {
		return nil, ErrBankInfoNotFound
	}
	bankInfo, ok := bicMap.LookupBankInfo(info.CountryCode, info.BankCode)
	if !ok {
		return nil, ErrBankInfoNotFound
	}
	return bankInfo, nil
}

// Must returns a MustService wrapper that panics on error.
func (s *Service) Must() *MustService {
	return &MustService{s}
}

// MustService is a wrapper around Service that panics on error.
type MustService struct {
	*Service
}

// Validate checks if the IBAN is valid and panics if not.
func (s *MustService) Validate(ibanStr string) {
	if err := s.Service.Validate(ibanStr); err != nil {
		panic(err)
	}
}

// Parse extracts IBANInfo and panics if parsing fails.
func (s *MustService) Parse(ibanStr string) *IBANInfo {
	info, err := s.Service.Parse(ibanStr)
	if err != nil {
		panic(err)
	}
	return info
}

// Detect returns IBANStructure and panics if detection fails.
func (s *MustService) Detect(ibanStr string) *IBANStructure {
	structure, err := s.Service.Detect(ibanStr)
	if err != nil {
		panic(err)
	}
	return structure
}

// ErrNilIBANInfo is returned if IBANInfo is nil.
var ErrNilIBANInfo = &IBANError{Code: "nil_iban_info", Message: "IBANInfo is nil"}

// ErrBankInfoNotFound is returned if no bank info is found for the IBAN.
var ErrBankInfoNotFound = &IBANError{Code: "bank_info_not_found", Message: "No bank info found for IBAN"}
