// Package iban implements IBAN validation, parsing, and detection logic.
package iban

import (
	"crypto/subtle"
	"regexp"
	"strings"

	bicmap "github.com/SamyRai/bank-data/internal/bic/map"
	"github.com/SamyRai/bank-data/internal/countrymeta"
	"github.com/SamyRai/bank-data/internal/log"
	"github.com/SamyRai/bank-data/pkg/bank"
	"github.com/SamyRai/bank-data/pkg/iban"
)

// validator implements the iban.Validator interface for IBAN validation.
type validator struct{}

// NewValidator returns a new IBAN Validator implementing the Validator interface.
func NewValidator() iban.Validator {
	return &validator{}
}

// Validate performs a comprehensive validation of an IBAN.
// It checks for correct length, valid characters, supported country code,
// country-specific format (using a regex), and checksum validity.
// Returns a typed error if the IBAN is invalid, otherwise nil.
func (v *validator) Validate(ibanStr string) error {
	ibanStrNorm := strings.ToUpper(strings.ReplaceAll(ibanStr, " ", ""))
	log.Debug("Validating IBAN", log.Fields{"iban": ibanStrNorm, "operation": "validate"})

	// Basic checks for length and characters
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

	// Country-specific validation
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

	// BBAN format validation using a pre-compiled regex
	if meta.BBANRegex != nil {
		bban := ibanStrNorm[4:]
		if !meta.BBANRegex.MatchString(bban) {
			err := *iban.ErrInvalidFormat
			err.Value = ibanStr
			log.Warn("IBAN validation failed: regex", log.Fields{"iban": ibanStrNorm, "code": err.Code, "error": err.Message})
			return &err
		}
	}

	// Checksum validation
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
func (v *validator) ValidateAndBankInfo(ibanStr string, bicMap bicmap.BankBICMap) (*bank.BankInfo, error) {
	ibanStrNorm := strings.ToUpper(strings.ReplaceAll(ibanStr, " ", ""))
	if len(ibanStrNorm) < 4 {
		return nil, iban.ErrWrongLength
	}
	if !regexp.MustCompile(`^[A-Z0-9]+$`).MatchString(ibanStrNorm) {
		return nil, iban.ErrInvalidChars
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

// validateIBANChecksum implements the IBAN checksum validation algorithm using streaming MOD-97.
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
	// Constant-time comparison to prevent timing attacks
	return subtle.ConstantTimeEq(int32(rem), 1) == 1
}
