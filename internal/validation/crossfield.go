// Package validation provides cross-field validation helpers.
package validation

import "github.com/SamyRai/bank-data/pkg/validation"

// CrossFieldRule defines a cross-field rule.
type CrossFieldRule[T any] func(input T) (bool, error)

// CrossFieldValidatorImpl applies ordered rules to one input.
type CrossFieldValidatorImpl[T any] struct {
	Rules []CrossFieldRule[T]
}

func (v *CrossFieldValidatorImpl[T]) ValidateFields(input T) validation.ValidationResult {
	result := validation.ValidationResult{Valid: true}
	for _, rule := range v.Rules {
		ok, err := rule(input)
		if !ok {
			result.Valid = false
			if err != nil {
				result.Code = "cross_field_failed"
				result.Message = err.Error()
			}
			break
		}
	}
	return result
}
