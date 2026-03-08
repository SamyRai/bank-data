package iban

import (
	"errors"
	"strings"

	ibanid "github.com/SamyRai/bank-data/internal/identifiers/iban"
	"github.com/SamyRai/bank-data/pkg/bank"
)

// validator implements Validator.
type validator struct {
	module *ibanid.Module
}

func NewValidator() Validator {
	return &validator{module: ibanid.New()}
}

func (v *validator) Validate(ibanStr string) error {
	n := v.module.Normalize(ibanStr)
	if err := v.module.Validate(n); err != nil {
		return toIBANError(n, err)
	}
	return nil
}

// ValidateAndBankInfo validates and enriches with bank metadata using a public lookup seam.
func (v *validator) ValidateAndBankInfo(ibanStr string, lookup BankLookup) (*bank.BankInfo, error) {
	if lookup == nil {
		return nil, ErrBankInfoNotFound
	}
	n := v.module.Normalize(ibanStr)
	fields, err := v.module.Parse(n)
	if err != nil {
		return nil, toIBANError(n, err)
	}
	bi, ok := lookup.LookupBank(fields["country_code"], fields["bank_code"])
	if !ok {
		return nil, ErrBankInfoNotFound
	}
	return bi, nil
}

func toIBANError(value string, err error) error {
	if err == nil {
		return nil
	}
	msg := err.Error()
	code := ErrCodeInvalidFormat
	field := "iban"

	switch {
	case containsCode(err, "invalid_chars"):
		code = ErrCodeInvalidChars
		field = "characters"
	case containsCode(err, "wrong_length"):
		code = ErrCodeWrongLength
		field = "length"
	case containsCode(err, "checksum_failed"):
		code = ErrCodeChecksum
		field = "checksum"
	case containsCode(err, "unsupported_country"):
		code = ErrCodeUnsupportedCountry
		field = "country"
	case containsCode(err, "invalid_format"):
		code = ErrCodeInvalidFormat
		field = "format"
	}

	return &IBANError{Code: code, Field: field, Message: msg, Value: value}
}

func containsCode(err error, code string) bool {
	var ie interface{ Error() string }
	if errors.As(err, &ie) {
		return strings.Contains(ie.Error(), code)
	}
	return false
}
