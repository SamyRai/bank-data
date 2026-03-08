package financial

import (
	"bufio"
	"os"
	"strings"
	"testing"
)

func TestFinancialVectors(t *testing.T) {
	f, err := os.Open("../../testdata/financial_test_vectors.csv")
	if err != nil {
		t.Fatalf("open vectors: %v", err)
	}
	defer func() { _ = f.Close() }()

	svc := NewService()
	sc := bufio.NewScanner(f)
	lineNo := 0
	for sc.Scan() {
		lineNo++
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if lineNo == 1 {
			continue // header
		}

		parts := strings.Split(line, ",")
		if len(parts) != 3 {
			t.Fatalf("bad vector line %d: %q", lineNo, line)
		}
		typ := parseType(parts[0])
		input := parts[1]
		wantValid := parts[2] == "1"

		report, err := svc.Validate(input, typ)
		if wantValid && err != nil {
			t.Fatalf("line %d: expected valid %s %q, got err %v", lineNo, typ, input, err)
		}
		if !wantValid && err == nil {
			t.Fatalf("line %d: expected invalid %s %q, got valid report %+v", lineNo, typ, input, report)
		}

		if wantValid {
			if _, err := svc.Parse(input, typ); err != nil {
				t.Fatalf("line %d: parse failed for valid vector %s %q: %v", lineNo, typ, input, err)
			}
		}
	}
	if err := sc.Err(); err != nil {
		t.Fatalf("scan vectors: %v", err)
	}
}

func parseType(s string) IdentifierType {
	switch strings.TrimSpace(s) {
	case "IBAN":
		return IdentifierIBAN
	case "BIC":
		return IdentifierBIC
	case "SEPA_CREDITOR_ID":
		return IdentifierSEPACreditor
	case "LEI":
		return IdentifierLEI
	case "ISIN":
		return IdentifierISIN
	case "PAN":
		return IdentifierPAN
	case "VAT":
		return IdentifierVAT
	case "NATIONAL_ACCOUNT_UK":
		return IdentifierNationalAccountUK
	default:
		return ""
	}
}
