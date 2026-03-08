package iban

import (
	"errors"
	"testing"

	"github.com/SamyRai/bank-data/pkg/validation"
)

func TestService_ValidateByTag(t *testing.T) {
	reg := validation.NewValidationRegistry()
	service := NewService(
		NewValidator(),
		NewParser(),
		NewDetector(),
		reg,
	)
	reg.Register("iban", &IBANBatchValidator{Validator: NewValidator()})
	tests := []struct {
		name    string
		tag     string
		input   string
		wantOk  bool
		wantErr bool
	}{
		{"valid IBAN by tag", "iban", "DE89370400440532013000", true, false},
		{"invalid IBAN by tag", "iban", "INVALIDIBAN123", false, true},
		{"unsupported tag", "bic", "SOMEVALUE", false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.ValidateByTag(tt.tag, tt.input)
			if tt.tag == "iban" {
				if result.Valid != tt.wantOk {
					t.Errorf("ValidateByTag() Valid = %v, want %v", result.Valid, tt.wantOk)
				}
				if (result.Error != nil) != tt.wantErr {
					t.Errorf("ValidateByTag() error presence = %v, wantErr %v", result.Error != nil, tt.wantErr)
				}
				if result.Error != nil {
					var ibanErr *IBANError
					if !errors.As(result.Error, &ibanErr) {
						t.Errorf("ValidateByTag() error type = %T, want *IBANError", result.Error)
					}
				}
			} else {
				if err == nil {
					t.Errorf("ValidateByTag() error = nil, want error for unsupported tag")
				}
			}
		})
	}
}
