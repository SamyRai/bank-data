package financial

import "testing"

func BenchmarkFinancialValidate_Matrix(b *testing.B) {
	svc := NewService()
	cases := []struct {
		name  string
		input string
		hint  IdentifierType
	}{
		{name: "IBAN/valid", input: "DE89370400440532013000", hint: IdentifierIBAN},
		{name: "IBAN/invalid", input: "DE89370400440532013001", hint: IdentifierIBAN},
		{name: "LEI/valid", input: "529900T8BM49AURSDO78", hint: IdentifierLEI},
		{name: "ISIN/valid", input: "US0378331005", hint: IdentifierISIN},
		{name: "PAN/valid", input: "4111111111111111", hint: IdentifierPAN},
		{name: "VAT/valid", input: "DE136695976", hint: IdentifierVAT},
		{name: "AUTO/detect", input: "DEUTDEFF", hint: ""},
	}

	for _, tc := range cases {
		b.Run(tc.name, func(b *testing.B) {
			b.SetBytes(int64(len(tc.input)))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _ = svc.Validate(tc.input, tc.hint)
			}
		})
	}
}
