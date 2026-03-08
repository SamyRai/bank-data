package iban

import (
	"errors"

	"github.com/SamyRai/bank-data/pkg/validation"
)

// IBANBatchValidator adapts IBAN validation to pkg/validation result shape.
type IBANBatchValidator struct {
	Validator Validator
}

func NewIBANBatchValidator(v Validator) *IBANBatchValidator {
	return &IBANBatchValidator{Validator: v}
}

func (v *IBANBatchValidator) Validate(input string) validation.ValidationResult {
	validator := v.Validator
	if validator == nil {
		validator = NewValidator()
	}
	err := validator.Validate(input)
	if err == nil {
		return validation.ValidationResult{Input: input, Valid: true}
	}

	res := validation.ValidationResult{Input: input, Valid: false, Code: "validation_failed", Message: err.Error()}
	var ibanErr *IBANError
	if errors.As(err, &ibanErr) {
		res.Code = string(ibanErr.Code)
		res.Message = ibanErr.Message
	}
	return res
}
