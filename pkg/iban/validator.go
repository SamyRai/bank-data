package iban

import (
	"crypto/subtle"
	"strings"

	bicmap "github.com/SamyRai/bank-data/internal/bic/map"
	"github.com/SamyRai/bank-data/internal/countrymeta"
	"github.com/SamyRai/bank-data/internal/log"
	"github.com/SamyRai/bank-data/pkg/bank"
)

// validator implements the Validator interface for IBAN validation.
type validator struct{}

// NewValidator returns a new IBAN Validator implementing the Validator interface.
func NewValidator() Validator {
	return &validator{}
}

// Validate checks if the IBAN is valid (format and checksum). Returns an error if invalid, or nil if valid.
func (v *validator) Validate(ibanStr string) error {
	ibanStrNorm := strings.ToUpper(strings.ReplaceAll(ibanStr, " ", ""))
	log.Debug("Validating IBAN", log.Fields{"iban": ibanStrNorm, "operation": "validate"})
	if len(ibanStrNorm) < 4 || len(ibanStrNorm) > 34 {
		err := *ErrWrongLength
		err.Value = ibanStr
		log.Warn("IBAN validation failed: wrong length", log.Fields{"iban": ibanStrNorm, "code": err.Code, "error": err.Message})
		return &err
	}
	// Fast alphanumeric check avoiding heavy regex
	for i := 0; i < len(ibanStrNorm); i++ {
		r := ibanStrNorm[i]
		if !((r >= '0' && r <= '9') || (r >= 'A' && r <= 'Z')) {
			err := *ErrInvalidChars
			err.Value = ibanStr
			log.Warn("IBAN validation failed: invalid characters", log.Fields{"iban": ibanStrNorm, "code": err.Code, "error": err.Message})
			return &err
		}
	}
	country := ibanStrNorm[:2]
	meta, ok := countrymeta.Registry[country]
	if !ok {
		err := *ErrUnsupportedCountry
		err.Value = country
		log.Warn("IBAN validation failed: unsupported country", log.Fields{"iban": ibanStrNorm, "code": err.Code, "error": err.Message})
		return &err
	}
	if len(ibanStrNorm) != meta.Length {
		err := *ErrWrongLength
		err.Value = ibanStr
		log.Warn("IBAN validation failed: wrong length for country", log.Fields{"iban": ibanStrNorm, "code": err.Code, "error": err.Message})
		return &err
	}
	// Regex validation using pre-compiled regex in meta.Regex
	if meta.Regex != nil {
		if !meta.Regex.MatchString(ibanStrNorm) {
			err := *ErrInvalidFormat
			err.Value = ibanStr
			log.Warn("IBAN validation failed: regex", log.Fields{"iban": ibanStrNorm, "code": err.Code, "error": err.Message})
			return &err
		}
	}
	if !validateIBANChecksum(ibanStrNorm) {
		err := *ErrChecksum
		err.Value = ibanStr
		log.Warn("IBAN validation failed: checksum", log.Fields{"iban": ibanStrNorm, "code": err.Code, "error": err.Message})
		return &err
	}
	log.Info("IBAN validated successfully", log.Fields{"iban": ibanStrNorm, "operation": "validate"})
	return nil
}

// ValidateAndBankInfo validates the IBAN and returns BankInfo if valid.
func (v *validator) ValidateAndBankInfo(ibanStr string, bicMap *bicmap.BankBICMap) (*bank.BankInfo, error) {
	ibanStrNorm := strings.ToUpper(strings.ReplaceAll(ibanStr, " ", ""))
	if len(ibanStrNorm) < 4 {
		return nil, ErrWrongLength
	}
	country := ibanStrNorm[:2]
	meta, ok := countrymeta.Registry[country]
	if !ok {
		return nil, ErrUnsupportedCountry
	}
	if len(ibanStrNorm) != meta.Length {
		return nil, ErrWrongLength
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
		return nil, ErrBankInfoNotFound
	}
	return bankInfo, nil
}

// validateIBANChecksum implements the IBAN checksum validation algorithm using a highly optimized loop.
// It avoids allocations by using a stack-based buffer for the rearranged string.
func validateIBANChecksum(ibanStr string) bool {
	// IBAN is max 34 chars. Stack allocation is safe and fast.
	var buf [34]byte
	copy(buf[:], ibanStr[4:])
	copy(buf[len(ibanStr)-4:], ibanStr[:4])

	rem := 0
	for i := 0; i < len(ibanStr); i++ {
		r := buf[i]
		switch {
		case r >= '0' && r <= '9':
			rem = (rem*10 + int(r-'0')) % 97
		case r >= 'A' && r <= 'Z':
			v := int(r - 'A' + 10)
			// MOD-97 is $R = (R * 100 + V) \pmod{97}$
			rem = (rem*100 + v) % 97
		default:
			return false
		}
	}
	// Constant-time check for result (security hardening as per TODO)
	return subtle.ConstantTimeEq(int32(rem), 1) == 1
}
