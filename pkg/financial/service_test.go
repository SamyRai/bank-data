package financial

import (
	"reflect"
	"testing"

)

// ensure CanonicalService implements Service at compile time
var _ Service = (*canonicalService)(nil)

// ensure types match API specification snapshot
func TestAPISnapshot(t *testing.T) {
	// IdentifierType
	if TypeUnknown != "unknown" {
		t.Errorf("TypeUnknown should be 'unknown'")
	}
	if TypeIBAN != "iban" {
		t.Errorf("TypeIBAN should be 'iban'")
	}
	if TypeBIC != "bic" {
		t.Errorf("TypeBIC should be 'bic'")
	}
	if TypeSEPACreditorID != "sepa_creditor_id" {
		t.Errorf("TypeSEPACreditorID should be 'sepa_creditor_id'")
	}

	// Structural checks for ValidationReport
	rtReport := reflect.TypeOf(ValidationReport{})
	if _, ok := rtReport.FieldByName("Valid"); !ok {
		t.Errorf("ValidationReport missing Valid field")
	}
	if _, ok := rtReport.FieldByName("IdentifierType"); !ok {
		t.Errorf("ValidationReport missing IdentifierType field")
	}
	if _, ok := rtReport.FieldByName("Normalized"); !ok {
		t.Errorf("ValidationReport missing Normalized field")
	}
	if _, ok := rtReport.FieldByName("Error"); !ok {
		t.Errorf("ValidationReport missing Error field")
	}

	// Structural checks for ParsedIdentifier
	rtParsed := reflect.TypeOf(ParsedIdentifier{})
	if _, ok := rtParsed.FieldByName("IdentifierType"); !ok {
		t.Errorf("ParsedIdentifier missing IdentifierType field")
	}
	if _, ok := rtParsed.FieldByName("Normalized"); !ok {
		t.Errorf("ParsedIdentifier missing Normalized field")
	}
	if _, ok := rtParsed.FieldByName("Components"); !ok {
		t.Errorf("ParsedIdentifier missing Components field")
	}
}

func TestDeterministicDetectPrecedence(t *testing.T) {
	svc := NewService()

	tests := []struct {
		name     string
		input    string
		expected IdentifierType
	}{
		{
			name:     "Valid IBAN",
			input:    "DE89370400440532013000",
			expected: TypeIBAN,
		},
		{
			name:     "Valid BIC 8",
			input:    "DEUTDEFF",
			expected: TypeBIC,
		},
		{
			name:     "Valid BIC 11",
			input:    "DEUTDEFFXXX",
			expected: TypeBIC,
		},
		{
			name:     "Valid SEPA Creditor ID",
			input:    "DE98ZZZ09999999999",
			expected: TypeSEPACreditorID,
		},
		{
			name:     "Invalid/Unknown string",
			input:    "INVALIDSTRING123",
			expected: TypeUnknown,
		},
		{
			name:     "Empty string",
			input:    "",
			expected: TypeUnknown,
		},
		{
			name:     "Whitespace only",
			input:    "   ",
			expected: TypeUnknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := svc.Detect(tt.input)
			if got != tt.expected {
				t.Errorf("Detect(%q) = %v; want %v", tt.input, got, tt.expected)
			}
		})
	}
}

func TestParseValidateNormalizedValueParity(t *testing.T) {
	svc := NewService()

	tests := []struct {
		name       string
		input      string
		normalized string
		valid      bool
		idType     IdentifierType
	}{
		{
			name:       "IBAN with spaces",
			input:      "DE89 3704 0044 0532 0130 00",
			normalized: "DE89370400440532013000",
			valid:      true,
			idType:     TypeIBAN,
		},
		{
			name:       "BIC with lowercase and spaces",
			input:      " deut de ff xxx ",
			normalized: "DEUTDEFFXXX",
			valid:      true,
			idType:     TypeBIC,
		},
		{
			name:       "SEPA Creditor ID with spaces",
			input:      " DE 98 ZZZ 09999999999 ",
			normalized: "DE98ZZZ09999999999",
			valid:      true,
			idType:     TypeSEPACreditorID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test Validate
			report := svc.Validate(tt.input)
			if report.Normalized != tt.normalized {
				t.Errorf("Validate().Normalized = %q; want %q", report.Normalized, tt.normalized)
			}
			if report.Valid != tt.valid {
				t.Errorf("Validate().Valid = %v; want %v", report.Valid, tt.valid)
			}
			if report.IdentifierType != tt.idType {
				t.Errorf("Validate().IdentifierType = %v; want %v", report.IdentifierType, tt.idType)
			}

			// Test Parse
			parsed, err := svc.Parse(tt.input)
			if err != nil {
				t.Fatalf("Parse() unexpected error: %v", err)
			}
			if parsed.Normalized != tt.normalized {
				t.Errorf("Parse().Normalized = %q; want %q", parsed.Normalized, tt.normalized)
			}
			if parsed.IdentifierType != tt.idType {
				t.Errorf("Parse().IdentifierType = %v; want %v", parsed.IdentifierType, tt.idType)
			}
		})
	}
}

func TestErrorCodeConsistency(t *testing.T) {
	svc := NewService()

	// Test IBAN specific error
	invalidIBAN := "DE89370400440532013000X" // wrong length
	report := svc.Validate(invalidIBAN)
	if report.Valid {
		t.Errorf("Expected invalid IBAN to be invalid")
	}

	if report.Error == nil {
		t.Fatalf("Expected error for invalid IBAN")
	}

	if report.Error == nil {
		t.Fatalf("Expected error to be returned for invalid IBAN")
	}

	// Test BIC specific error
	invalidBIC := "DEUTDEFFX" // wrong length 9 (8 or 11 allowed)
	reportBIC := svc.Validate(invalidBIC)
	if reportBIC.Valid {
		t.Errorf("Expected invalid BIC to be invalid")
	}

	if reportBIC.Error == nil {
		t.Fatalf("Expected error for invalid BIC")
	}
}
