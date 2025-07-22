// Package iban implements IBAN validation, parsing, and detection logic.
package iban

import (
	"regexp"
	"strings"

	"github.com/SamyRai/bank-data/internal/log"
	"github.com/SamyRai/bank-data/pkg/iban"
)

// validator implements the iban.Validator interface.
type validator struct{}

// NewValidator returns a new IBAN Validator.
func NewValidator() iban.Validator {
	return &validator{}
}

// Validate checks if the IBAN is valid (format and checksum).
func (v *validator) Validate(ibanStr string) error {
	ibanStrNorm := strings.ToUpper(strings.ReplaceAll(ibanStr, " ", ""))
	log.Debug("Validating IBAN", log.Fields{"iban": ibanStrNorm, "operation": "validate"})
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
	if !validateIBANChecksum(ibanStrNorm) {
		err := *iban.ErrChecksum
		err.Value = ibanStr
		log.Warn("IBAN validation failed: checksum", log.Fields{"iban": ibanStrNorm, "code": err.Code, "error": err.Message})
		return &err
	}
	log.Info("IBAN validated successfully", log.Fields{"iban": ibanStrNorm, "operation": "validate"})
	return nil
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
	return rem == 1
}
