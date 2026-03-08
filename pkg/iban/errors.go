package iban

import "fmt"

// IBANErrorCode defines machine-readable categories for IBAN failures.
type IBANErrorCode string

const (
	ErrCodeInvalidChars       IBANErrorCode = "invalid_chars"
	ErrCodeWrongLength        IBANErrorCode = "wrong_length"
	ErrCodeChecksum           IBANErrorCode = "checksum_failed"
	ErrCodeUnsupportedCountry IBANErrorCode = "unsupported_country"
	ErrCodeInvalidFormat      IBANErrorCode = "invalid_format"
)

// IBANError is the canonical typed error for the package.
type IBANError struct {
	Code    IBANErrorCode
	Field   string
	Message string
	Value   string
}

func (e *IBANError) Error() string {
	if e.Field == "" {
		return e.Message
	}
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

func newIBANError(code IBANErrorCode, field, message, value string) *IBANError {
	return &IBANError{Code: code, Field: field, Message: message, Value: value}
}

var (
	ErrInvalidChars       = newIBANError(ErrCodeInvalidChars, "characters", "IBAN contains invalid characters", "")
	ErrWrongLength        = newIBANError(ErrCodeWrongLength, "length", "IBAN length is invalid", "")
	ErrChecksum           = newIBANError(ErrCodeChecksum, "checksum", "IBAN checksum validation failed", "")
	ErrUnsupportedCountry = newIBANError(ErrCodeUnsupportedCountry, "country", "IBAN country code is not supported", "")
	ErrInvalidFormat      = newIBANError(ErrCodeInvalidFormat, "format", "IBAN does not match country format", "")
	ErrNilIBANInfo        = &IBANError{Code: "nil_iban_info", Message: "IBANInfo is nil"}
	ErrBankInfoNotFound   = &IBANError{Code: "bank_info_not_found", Message: "No bank info found for IBAN"}
)
