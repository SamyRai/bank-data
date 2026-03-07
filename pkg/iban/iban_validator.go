package iban

import (
	"errors"
	"time"

	"github.com/SamyRai/bank-data/internal/log"
	"github.com/SamyRai/bank-data/pkg/validation"
)

// IBANBatchValidator implements validation.Validator for IBAN strings.
type IBANBatchValidator struct {
	Validator Validator // dependency injection for testability
}

// NewIBANBatchValidator constructs an IBANBatchValidator with a custom validator.
func NewIBANBatchValidator(v Validator) *IBANBatchValidator {
	return &IBANBatchValidator{Validator: v}
}

// Validate returns a validation.ValidationResult for the given IBAN string, with structured logging and rich details.
func (v *IBANBatchValidator) Validate(input string) validation.ValidationResult {
	start := time.Now()
	log.Debug("Batch validating IBAN", log.Fields{"iban": input, "operation": "batch_validate"})
	var err error
	if v.Validator != nil {
		err = v.Validator.Validate(input)
	} else {
		err = NewValidator().Validate(input)
	}
	vr := validation.ValidationResult{
		Input:   input,
		Valid:   err == nil,
		Error:   err,
		Details: map[string]any{"duration_ms": time.Since(start).Milliseconds()},
	}
	if err != nil {
		var ibanErr *IBANError
		if errors.As(err, &ibanErr) {
			vr.Details["code"] = ibanErr.Code
			vr.Details["field"] = ibanErr.Field
			vr.Details["message"] = ibanErr.Message
			vr.Details["value"] = ibanErr.Value
			log.Warn("IBAN batch validation failed", log.Fields{"iban": input, "code": ibanErr.Code, "field": ibanErr.Field, "msg": ibanErr.Message})
		} else {
			log.Warn("IBAN batch validation failed (unknown error)", log.Fields{"iban": input, "err": err.Error()})
		}
	} else {
		log.Info("IBAN batch validation succeeded", log.Fields{"iban": input})
	}
	return vr
}
