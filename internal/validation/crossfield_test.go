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
	// Rule: A must be less than B
	rule1 := func(input testStruct) (bool, error) {
		if input.A < input.B {
			return true, nil
		}
		return false, errors.New("A must be less than B")
	}
	// Rule: Country must be 'DE' if A > 0
	rule2 := func(input testStruct) (bool, error) {
		if input.A > 0 && input.Country != "DE" {
			return false, errors.New("Country must be 'DE' if A > 0")
		}
		return true, nil
	}
	validator := &CrossFieldValidatorImpl[testStruct]{
		Rules: []CrossFieldRule[testStruct]{rule1, rule2},
	}

	tests := []struct {
		name   string
		input  testStruct
		valid  bool
		errMsg string
	}{
		{"valid", testStruct{A: 1, B: 2, Country: "DE"}, true, ""},
		{"A not less than B", testStruct{A: 3, B: 2, Country: "DE"}, false, "A must be less than B"},
		{"Country not DE", testStruct{A: 1, B: 2, Country: "FR"}, false, "Country must be 'DE' if A > 0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := validator.ValidateFields(tt.input)
			if res.Valid != tt.valid {
				t.Errorf("expected valid=%v, got %v", tt.valid, res.Valid)
			}
			if !tt.valid && (res.Error == nil || res.Error.Error() != tt.errMsg) {
				t.Errorf("expected error '%s', got '%v'", tt.errMsg, res.Error)
			}
		})
	}
}
