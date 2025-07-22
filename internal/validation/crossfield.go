// Package validation provides cross-field and conditional validation logic.
package validation

import (
	"github.com/SamyRai/bank-data/pkg/validation"
)

// CrossFieldRule defines a function for cross-field or conditional validation.
type CrossFieldRule[T any] func(input T) (bool, error)

// CrossFieldValidatorImpl implements CrossFieldValidator for any struct.
type CrossFieldValidatorImpl[T any] struct {
	Rules []CrossFieldRule[T]
}

// ValidateFields applies all cross-field rules to the input and returns a ValidationResult.
func (v *CrossFieldValidatorImpl[T]) ValidateFields(input T) validation.ValidationResult {
	result := validation.ValidationResult{
		Input:   input,
		Valid:   true,
		Details: map[string]any{},
	}
	var failedRules []int
	var errorsList []error
	for i, rule := range v.Rules {
		ok, err := rule(input)
		if !ok {
			result.Valid = false
			failedRules = append(failedRules, i)
			errorsList = append(errorsList, err)
		}
	}
	if !result.Valid {
		if len(errorsList) > 0 {
			result.Error = errorsList[0]
		} else {
			result.Error = nil
		}
		result.Details["failed_rule_indices"] = failedRules
		result.Details["rule_errors"] = errorsList
	}
	return result
}
