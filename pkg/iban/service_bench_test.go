package iban

import (
	"context"
	"fmt"
	"testing"
)

func BenchmarkService_Validate_Matrix(b *testing.B) {
	svc := NewService(nil, nil, nil, nil)

	cases := []struct {
		name    string
		iban    string
		country string
		mode    string
	}{
		// Valid cases
		{name: "Valid/country=DE", iban: "DE89370400440532013000", country: "DE", mode: "valid"},
		{name: "Valid/country=FR", iban: "FR1420041010050500013M02606", country: "FR", mode: "valid"},
		{name: "Valid/country=CH", iban: "CH9300762011623852957", country: "CH", mode: "valid"},
		{name: "Valid/country=GB", iban: "GB82WEST12345698765432", country: "GB", mode: "valid"},

		// Failure modes
		{name: "Error/type=Length", iban: "DE89", country: "DE", mode: "invalid_length"},
		{name: "Error/type=Chars", iban: "DE8937040044053201300!", country: "DE", mode: "invalid_chars"},
		{name: "Error/type=Checksum", iban: "DE89370400440532013001", country: "DE", mode: "invalid_checksum"},
	}

	for _, tc := range cases {
		b.Run(tc.name, func(b *testing.B) {
			b.SetBytes(int64(len(tc.iban)))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = svc.Validate(tc.iban)
			}
		})
	}
}

func BenchmarkService_ValidateBatch_Scaling(b *testing.B) {
	svc := NewService(nil, nil, nil, nil)
	ctx := context.Background()

	sizes := []int{100, 1000, 10000}
	for _, size := range sizes {
		inputs := make([]string, size)
		for i := 0; i < size; i++ {
			inputs[i] = "DE89370400440532013000"
		}

		b.Run(fmt.Sprintf("size=%d", size), func(b *testing.B) {
			b.SetBytes(int64(len(inputs[0]) * size))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = svc.ValidateBatch(ctx, inputs)
			}
			b.ReportMetric(float64(b.N*size)/b.Elapsed().Seconds(), "IBANs/s")
		})
	}
}

func BenchmarkService_Parse_Matrix(b *testing.B) {
	svc := NewService(nil, nil, nil, nil)
	input := "DE89370400440532013000"

	b.SetBytes(int64(len(input)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = svc.Parse(input)
	}
}
