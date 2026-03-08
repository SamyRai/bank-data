// Package nationalaccount validates UK sort code + account number pairs.
package nationalaccount

import (
	internalna "github.com/SamyRai/bank-data/internal/identifiers/nationalaccount"
)

// ValidationResult captures validation outcome for UK national accounts.
type ValidationResult struct {
	Input      string
	Normalized string
	Valid      bool
	Code       string
	Message    string
}

// AccountInfo is parsed UK national account data.
type AccountInfo struct {
	SortCode      string
	AccountNumber string
	Raw           string
}

type Validator struct{ module *internalna.Module }

func NewValidator() *Validator { return &Validator{module: internalna.New()} }

func (v *Validator) Validate(input string) ValidationResult {
	n := v.module.Normalize(input)
	if err := v.module.Validate(n); err != nil {
		return ValidationResult{Input: input, Normalized: n, Valid: false, Code: "validation_failed", Message: err.Error()}
	}
	return ValidationResult{Input: input, Normalized: n, Valid: true}
}

func (v *Validator) Parse(input string) (*AccountInfo, error) {
	n := v.module.Normalize(input)
	fields, err := v.module.Parse(n)
	if err != nil {
		return nil, err
	}
	return &AccountInfo{SortCode: fields["sort_code"], AccountNumber: fields["account_number"], Raw: fields["raw"]}, nil
}
