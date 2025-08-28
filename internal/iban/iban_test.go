package iban_test

import (
	"testing"

	"github.com/SamyRai/bank-data/internal/bic/map"
	ibanimpl "github.com/SamyRai/bank-data/internal/iban"
	"github.com/SamyRai/bank-data/pkg/iban"
)

func TestValidator_Validate(t *testing.T) {
	validator := ibanimpl.NewValidator()
	tests := []struct {
		name    string
		iban    string
		wantErr bool
	}{
		{"valid DE IBAN", "DE89370400440532013000", false},
		{"invalid checksum", "DE89370400440532013001", true},
		{"too short", "DE89", true},
		{"invalid chars", "DE89!@#$%^", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Validate(tt.iban)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestParser_Parse(t *testing.T) {
	parser := ibanimpl.NewParser()
	ibanStr := "DE89370400440532013000"
	info, err := parser.Parse(ibanStr)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if info.CountryCode != "DE" || info.BankCode != "37040044" || info.AccountNumber != "0532013000" {
		t.Errorf("Parse() got = %+v", info)
	}
}

func TestDetector_Detect(t *testing.T) {
	detector := ibanimpl.NewDetector()
	ibanStr := "DE89370400440532013000"
	structure, err := detector.Detect(ibanStr)
	if err != nil {
		t.Fatalf("Detect() error = %v", err)
	}
	if structure.CountryCode != "DE" || structure.Length != 22 {
		t.Errorf("Detect() got = %+v", structure)
	}
}

func TestValidator_ValidateAndBankInfo(t *testing.T) {
	validator := ibanimpl.NewValidator()
	mockBicMap := bicmap.BankBICMap{
		"DE": {
			"37040044": bicmap.BankBICEntry{
				BIC:      "COBADEFFXXX",
				BankName: "COMMERZBANK",
				Country:  "DE",
			},
		},
	}

	tests := []struct {
		name      string
		iban      string
		bicMap    bicmap.BankBICMap
		wantBic   string
		wantErr   bool
		errType   error
	}{
		{
			name:      "valid IBAN and bank info found",
			iban:      "DE89370400440532013000",
			bicMap:    mockBicMap,
			wantBic:   "COBADEFFXXX",
			wantErr:   false,
			errType:   nil,
		},
		{
			name:      "valid IBAN but bank info not found",
			iban:      "DE89370400440532013000",
			bicMap:    bicmap.BankBICMap{},
			wantErr:   true,
			errType:   iban.ErrBankInfoNotFound,
		},
		{
			name:      "invalid IBAN",
			iban:      "invalid-iban",
			bicMap:    mockBicMap,
			wantErr:   true,
			errType:   iban.ErrInvalidChars,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := validator.ValidateAndBankInfo(tt.iban, tt.bicMap)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateAndBankInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				if err.Error() != tt.errType.Error() {
					t.Errorf("ValidateAndBankInfo() error = %v, want %v", err, tt.errType)
				}
			} else {
				if got.BIC != tt.wantBic {
					t.Errorf("ValidateAndBankInfo() BIC = %v, want %v", got.BIC, tt.wantBic)
				}
			}
		})
	}
}
