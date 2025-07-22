package iban_test

import (
	"testing"

	internalvalidation "github.com/SamyRai/bank-data/internal/validation"
	pkgvalidation "github.com/SamyRai/bank-data/pkg/validation"
)

type crossFieldIBANInput struct {
	Country  string
	BankCode string
	Account  string
}

func TestService_CrossFieldValidatorByTag(t *testing.T) {
	reg := pkgvalidation.NewValidationRegistry()
	// Example: If country is "DE", bank code must start with "3704"
	rule := func(input crossFieldIBANInput) (bool, error) {
		if input.Country == "DE" && len(input.BankCode) > 0 && input.BankCode[:4] != "3704" {
			return false, pkgvalidation.NewValidationError("bank_code", "must start with 3704 for DE", input.BankCode)
		}
		return true, nil
	}
	crossValidator := &internalvalidation.CrossFieldValidatorImpl[crossFieldIBANInput]{
		Rules: []internalvalidation.CrossFieldRule[crossFieldIBANInput]{rule},
	}
	reg.Register("iban-crossfield", crossValidator)

	// Simulate service layer usage
	v := reg.Get("iban-crossfield").(*internalvalidation.CrossFieldValidatorImpl[crossFieldIBANInput])
	input := crossFieldIBANInput{Country: "DE", BankCode: "12345678", Account: "1234567890"}
	result := v.ValidateFields(input)
	if result.Valid {
		t.Errorf("expected invalid, got valid")
	}
	if len(result.Details["failed_rule_indices"].([]int)) == 0 {
		t.Errorf("expected failed rule indices, got none")
	}
}
