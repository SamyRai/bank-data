package validation

import (
	"errors"
	"testing"
)

type testStruct struct {
	A       int
	B       int
	Country string
}

func TestCrossFieldValidatorImpl_ValidateFields(t *testing.T) {
	rule1 := func(input testStruct) (bool, error) {
		if input.A < input.B {
			return true, nil
		}
		return false, errors.New("A must be less than B")
	}
	rule2 := func(input testStruct) (bool, error) {
		if input.A > 0 && input.Country != "DE" {
			return false, errors.New("Country must be 'DE' if A > 0")
		}
		return true, nil
	}
	validator := &CrossFieldValidatorImpl[testStruct]{Rules: []CrossFieldRule[testStruct]{rule1, rule2}}

	res := validator.ValidateFields(testStruct{A: 3, B: 2, Country: "FR"})
	if res.Valid {
		t.Fatalf("expected invalid result")
	}
	if res.Code == "" || res.Message == "" {
		t.Fatalf("expected code and message to be set, got %+v", res)
	}
}
