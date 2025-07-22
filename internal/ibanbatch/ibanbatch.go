// Package ibanbatch provides a batch validator for IBAN using the generic batch validation engine.
package ibanbatch

import (
	"errors"

	internaliban "github.com/SamyRai/bank-data/internal/iban"
	"github.com/SamyRai/bank-data/pkg/iban"
	"github.com/SamyRai/bank-data/pkg/validation"
)

// IBANBatchValidator implements validation.Validator for IBAN strings.
type IBANBatchValidator struct{}

func (v *IBANBatchValidator) Validate(input string) validation.ValidationResult {
	err := internaliban.NewValidator().Validate(input)
	vr := validation.ValidationResult{
		Input:   input,
		Valid:   err == nil,
		Error:   err,
		Details: map[string]any{},
	}
	if err != nil {
		var ibanErr *iban.IBANError
		if errors.As(err, &ibanErr) {
			vr.Details["code"] = ibanErr.Code
			vr.Details["field"] = ibanErr.Field
			vr.Details["message"] = ibanErr.Message
			vr.Details["value"] = ibanErr.Value
		}
	}
	return vr
}
