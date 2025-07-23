// Package iban implements IBAN validation, parsing, and detection logic.
package iban

import (
	"regexp"
	"strings"

	bicmap "github.com/SamyRai/bank-data/internal/bic/map"
	"github.com/SamyRai/bank-data/internal/countrymeta"
	"github.com/SamyRai/bank-data/internal/log"
	"github.com/SamyRai/bank-data/pkg/bank"
	"github.com/SamyRai/bank-data/pkg/iban"
)

// IBANMaxLength is the maximum allowed IBAN length per ISO 13616.
const IBANMaxLength = 34

// validator implements the iban.Validator interface for IBAN validation.
type validator struct{}

// NewValidator returns a new IBAN Validator implementing the Validator interface.
func NewValidator() iban.Validator {
	return &validator{}
}

// Validate checks if the IBAN is valid (format and checksum). Returns an error if invalid, or nil if valid.
func (v *validator) Validate(ibanStr string) error {
	ibanStrNorm := strings.ToUpper(strings.ReplaceAll(ibanStr, " ", ""))
	if len(ibanStrNorm) > IBANMaxLength {
		err := *iban.ErrWrongLength
		err.Value = ibanStr
		log.Warn("IBAN validation failed: exceeds max length", log.Fields{"iban": ibanStrNorm, "code": err.Code, "error": err.Message})
		return &err
	}
	if len(ibanStrNorm) < 4 {
		err := *iban.ErrWrongLength
		err.Value = ibanStr
		log.Warn("IBAN validation failed: wrong length", log.Fields{"iban": ibanStrNorm, "code": err.Code, "error": err.Message})
		return &err
	}
	if !regexp.MustCompile(`^[A-Z0-9]+$`).MatchString(ibanStrNorm) {
		err := *iban.ErrInvalidChars
		err.Value = ibanStr
		log.Warn("IBAN validation failed: invalid characters", log.Fields{"iban": ibanStrNorm, "code": err.Code, "error": err.Message})
		return &err
	}
	country := ibanStrNorm[:2]
	meta, ok := countrymeta.Registry[country]
	if !ok {
		err := *iban.ErrUnsupportedCountry
		err.Value = country
		log.Warn("IBAN validation failed: unsupported country", log.Fields{"iban": ibanStrNorm, "code": err.Code, "error": err.Message})
		return &err
	}
	if len(ibanStrNorm) != meta.Length {
		err := *iban.ErrWrongLength
		err.Value = ibanStr
		log.Warn("IBAN validation failed: wrong length for country", log.Fields{"iban": ibanStrNorm, "code": err.Code, "error": err.Message})
		return &err
	}
	// Regex validation using pre-compiled regex in meta.Regex
	if meta.Regex != nil {
		if !meta.Regex.MatchString(ibanStrNorm) {
			err := *iban.ErrInvalidFormat
			err.Value = ibanStr
			log.Warn("IBAN validation failed: regex", log.Fields{"iban": ibanStrNorm, "code": err.Code, "error": err.Message})
			return &err
		}
	}
	if !validateIBANChecksum(ibanStrNorm) {
		err := *iban.ErrChecksum
		err.Value = ibanStr
		log.Warn("IBAN validation failed: checksum", log.Fields{"iban": ibanStrNorm, "code": err.Code, "error": err.Message})
		return &err
	}
	log.Info("IBAN validated successfully", log.Fields{"iban": ibanStrNorm, "operation": "validate"})
	return nil
}

// ValidateAndBankInfo validates the IBAN and returns BankInfo if valid.
func (v *validator) ValidateAndBankInfo(ibanStr string, bicMap bicmap.BankBICMap) (*bank.Info, error) {
	ibanStrNorm := strings.ToUpper(strings.ReplaceAll(ibanStr, " ", ""))
	if len(ibanStrNorm) > IBANMaxLength {
		return nil, iban.ErrWrongLength
	}
	if len(ibanStrNorm) < 4 {
		return nil, iban.ErrWrongLength
	}
	country := ibanStrNorm[:2]
	meta, ok := countrymeta.Registry[country]
	if !ok {
		return nil, iban.ErrUnsupportedCountry
	}
	if len(ibanStrNorm) != meta.Length {
		return nil, iban.ErrWrongLength
	}
	// Extract bank code using meta
	bankCode := ""
	bbanOffset := 4
	if meta.BankStart > 0 && meta.BankEnd > 0 && meta.BankEnd > meta.BankStart {
		start := bbanOffset + (meta.BankStart - 1)
		end := bbanOffset + meta.BankEnd
		if end <= len(ibanStrNorm) && start < end {
			bankCode = ibanStrNorm[start:end]
		}
	}
	bankInfo, ok := bicMap.LookupBankInfo(country, bankCode)
	if !ok {
		return nil, iban.ErrBankInfoNotFound
	}
	return bankInfo, nil
}

// validateIBANChecksum implements the IBAN checksum validation algorithm using streaming MOD-97 and constant-time comparison.
func validateIBANChecksum(ibanStr string) bool {
	rearranged := ibanStr[4:] + ibanStr[:4]
	rem := 0
	for _, r := range rearranged {
		switch {
		case r >= '0' && r <= '9':
			rem = (rem*10 + int(r-'0')) % 97
		case r >= 'A' && r <= 'Z':
			v := int(r - 'A' + 10)
			rem = (rem*10 + v/10) % 97
			rem = (rem*10 + v%10) % 97
		default:
			return false
		}
	}
	// Constant-time comparison for security
	return constantTimeEq(rem, 1)
}

// constantTimeEq returns true if a == b, in constant time.
func constantTimeEq(a, b int) bool {
	// TODO: Consider using crypto/subtle for even stricter constant-time guarantees. Priority: medium, Effort: 1h
	return (a ^ b) == 0
}

// TestValidateIBANChecksum is an exported version for unit testing.
func TestValidateIBANChecksum(ibanStr string) bool {
	return validateIBANChecksum(ibanStr)
}

// TODO: Add more comprehensive unit tests for input length capping and constant-time check digit validation. Priority: high, Effort: 1h
// TODO: Refactor error handling to use domain-specific error types everywhere. Priority: medium, Effort: 2h
