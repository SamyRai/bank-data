package iban_test

import (
	"testing"

	internal "github.com/SamyRai/bank-data/internal/iban"
	"github.com/SamyRai/bank-data/pkg/iban"
)

func FuzzIBANValidate(f *testing.F) {
	// Seed with valid and invalid IBANs
	f.Add("DE89370400440532013000")      // valid
	f.Add("GB82WEST12345698765432")      // valid
	f.Add("FR1420041010050500013M02606") // valid
	f.Add("INVALIDIBAN123")              // invalid
	f.Add("")                            // invalid

	service := iban.NewService(
		internal.NewValidator(),
		internal.NewParser(),
		internal.NewDetector(),
	)
	f.Fuzz(func(_ *testing.T, s string) {
		_ = service.Validate(s) // Should not panic or crash
	})
}
