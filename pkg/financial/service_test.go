package financial

import (
	"context"
	"iter"
	"slices"
	"testing"
)

func TestService_Validate_AllTypes(t *testing.T) {
	svc := NewService()
	cases := []struct {
		name  string
		input string
		hint  IdentifierType
	}{
		{name: "iban", input: "DE89370400440532013000", hint: IdentifierIBAN},
		{name: "bic", input: "DEUTDEFF", hint: IdentifierBIC},
		{name: "sepa", input: "DE98ZZZ09999999999", hint: IdentifierSEPACreditor},
		{name: "lei", input: "529900T8BM49AURSDO78", hint: IdentifierLEI},
		{name: "isin", input: "US0378331005", hint: IdentifierISIN},
		{name: "pan", input: "4111111111111111", hint: IdentifierPAN},
		{name: "vat", input: "DE136695976", hint: IdentifierVAT},
		{name: "national account", input: "20-00-00 55779911", hint: IdentifierNationalAccountUK},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			report, err := svc.Validate(tc.input, tc.hint)
			if err != nil {
				t.Fatalf("Validate() unexpected error: %v", err)
			}
			if !report.Valid {
				t.Fatalf("Validate() report invalid: %+v", report)
			}
			if report.Type != tc.hint {
				t.Fatalf("Validate() type = %s, want %s", report.Type, tc.hint)
			}
		})
	}
}

func TestService_Detect_AndParse(t *testing.T) {
	svc := NewService()
	type tests struct {
		input string
		want  IdentifierType
		key   string
	}
	cases := []tests{
		{input: "DE89370400440532013000", want: IdentifierIBAN, key: "bank_code"},
		{input: "US0378331005", want: IdentifierISIN, key: "nsin"},
		{input: "4111111111111111", want: IdentifierPAN, key: "network"},
		{input: "20-00-00 55779911", want: IdentifierNationalAccountUK, key: "sort_code"},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			typ, err := svc.Detect(tc.input)
			if err != nil {
				t.Fatalf("Detect() unexpected error: %v", err)
			}
			if typ != tc.want {
				t.Fatalf("Detect() = %s, want %s", typ, tc.want)
			}
			parsed, err := svc.Parse(tc.input, "")
			if err != nil {
				t.Fatalf("Parse() unexpected error: %v", err)
			}
			if parsed.Fields[tc.key] == "" {
				t.Fatalf("Parse() missing field %q in %+v", tc.key, parsed.Fields)
			}
		})
	}
}

func TestService_ValidateBatchAndStream(t *testing.T) {
	svc := NewService()
	inputs := []string{"DE89370400440532013000", "INVALID", "US0378331005", "20-00-00 55779911"}
	batch := svc.ValidateBatch(context.Background(), inputs, "")
	if len(batch) != len(inputs) {
		t.Fatalf("ValidateBatch() len = %d, want %d", len(batch), len(inputs))
	}

	seq := iter.Seq[string](func(yield func(string) bool) {
		for _, in := range inputs {
			if !yield(in) {
				return
			}
		}
	})

	stream := []ValidationReport{}
	for r := range svc.StreamValidate(context.Background(), seq, "") {
		stream = append(stream, r)
	}
	if len(stream) != len(inputs) {
		t.Fatalf("StreamValidate() len = %d, want %d", len(stream), len(inputs))
	}

	types := []IdentifierType{}
	for _, r := range stream {
		types = append(types, r.Type)
	}
	slices.Sort(types)
	if len(types) != 4 {
		t.Fatalf("unexpected stream types: %+v", types)
	}
}

func TestService_SuggestIBAN(t *testing.T) {
	svc := NewService()
	suggestions, err := svc.Suggest("DE89 3704 0044 0532 0130 10", IdentifierIBAN)
	if err != nil {
		t.Fatalf("suggest failed: %v", err)
	}
	if len(suggestions) == 0 {
		t.Fatalf("expected suggestions for close invalid IBAN")
	}
}

func TestService_SuggestUnsupportedType(t *testing.T) {
	svc := NewService()
	if _, err := svc.Suggest("4111111111111111", IdentifierPAN); err == nil {
		t.Fatalf("expected unsupported suggest error for PAN")
	}
}
