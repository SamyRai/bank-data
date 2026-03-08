package financial_test

import (
	"context"
	"iter"
	"strings"
	"testing"

	"github.com/SamyRai/bank-data/pkg/financial"
)

func TestErrorCodeConsistency(t *testing.T) {
	svc := financial.NewService()

	cases := []struct {
		name  string
		input string
		hint  financial.IdentifierType
	}{
		{name: "empty_iban", input: "", hint: financial.IdentifierIBAN},
		{name: "empty_bic", input: " ", hint: financial.IdentifierBIC},
		{name: "invalid_iban", input: "DE12345", hint: financial.IdentifierIBAN},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			report, err := svc.Validate(tc.input, tc.hint)
			if report.Valid {
				t.Fatalf("expected invalid input for %q", tc.name)
			}
			if err == nil {
				t.Fatalf("expected error for invalid input %q", tc.name)
			}
			if report.Error == nil {
				t.Fatalf("expected non-nil report.Error")
			}

			// Empty inputs should always return "empty_input"
			if strings.TrimSpace(tc.input) == "" {
				if report.Error.Code != "empty_input" {
					t.Errorf("expected code empty_input, got %q", report.Error.Code)
				}
			}

			// Verify error wrapper contains valid structure
			if report.Error.Type != tc.hint {
				t.Errorf("expected error type %q, got %q", tc.hint, report.Error.Type)
			}
			if report.Error.Code == "" {
				t.Errorf("expected non-empty error code")
			}
			if report.Error.Message == "" {
				t.Errorf("expected non-empty error message")
			}
		})
	}
}

func TestParseValidateParity(t *testing.T) {
	svc := financial.NewService()

	cases := []struct {
		input string
		hint  financial.IdentifierType
	}{
		{input: "de 89 3704 0044 0532 0130 00", hint: financial.IdentifierIBAN},
		{input: " US0378331005 ", hint: financial.IdentifierISIN},
	}

	for _, tc := range cases {
		t.Run(string(tc.hint), func(t *testing.T) {
			report, err := svc.Validate(tc.input, tc.hint)
			if err != nil {
				t.Fatalf("Validate(%q) unexpected error: %v", tc.input, err)
			}

			parsed, err := svc.Parse(tc.input, tc.hint)
			if err != nil {
				t.Fatalf("Parse(%q) unexpected error: %v", tc.input, err)
			}

			if report.Normalized != parsed.Normalized {
				t.Errorf("Normalized parity mismatch: Validate=%q, Parse=%q", report.Normalized, parsed.Normalized)
			}
			if report.Type != parsed.Type {
				t.Errorf("Type parity mismatch: Validate=%q, Parse=%q", report.Type, parsed.Type)
			}
		})
	}
}

func TestPublicAPISnapshot(_ *testing.T) {
	// This test asserts the public API methods of the Service via compile-time
	// checks of interface signatures. This ensures we don't accidentally break
	// method signatures in a minor or patch release.

	// Assert the core public methods exist on *financial.Service
	type ServiceAPI interface {
		Detect(input string) (financial.IdentifierType, error)
		Validate(input string, hint financial.IdentifierType) (financial.ValidationReport, error)
		Parse(input string, hint financial.IdentifierType) (financial.ParsedIdentifier, error)
		Suggest(input string, hint financial.IdentifierType) ([]financial.Suggestion, error)
		ValidateBatch(ctx context.Context, inputs []string, hint financial.IdentifierType) []financial.ValidationReport
		StreamValidate(ctx context.Context, inputs iter.Seq[string], hint financial.IdentifierType) iter.Seq[financial.ValidationReport]
	}

	// This assignment acts as a compile-time check that *financial.Service
	// implements the ServiceAPI interface.
	var _ ServiceAPI = financial.NewService()
}

func TestDetectPrecedence(t *testing.T) {
	svc := financial.NewService()

	// Given an input that could theoretically be matched by multiple modules
	// (though in practice this is rare due to strict validation), we check
	// the ranked order. Here we rely on known valid inputs and ensure detection
	// successfully finds the *right* type instead of randomly picking one.

	cases := []struct {
		input string
		want  financial.IdentifierType
	}{
		// IBAN has rank 100
		{input: "DE89370400440532013000", want: financial.IdentifierIBAN},
		// LEI has rank 95
		{input: "529900T8BM49AURSDO78", want: financial.IdentifierLEI},
		// ISIN has rank 90
		{input: "US0378331005", want: financial.IdentifierISIN},
		// SEPA Creditor ID has rank 85
		{input: "DE98ZZZ09999999999", want: financial.IdentifierSEPACreditor},
		// BIC has rank 80
		{input: "DEUTDEFF", want: financial.IdentifierBIC},
	}

	for _, tc := range cases {
		t.Run(string(tc.want), func(t *testing.T) {
			typ, err := svc.Detect(tc.input)
			if err != nil {
				t.Fatalf("Detect(%q) unexpected error: %v", tc.input, err)
			}
			if typ != tc.want {
				t.Errorf("Detect(%q) = %v, want %v", tc.input, typ, tc.want)
			}
		})
	}
}
