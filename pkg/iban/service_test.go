package iban

import "testing"

func TestService_Validate_Parse_Detect(t *testing.T) {
	svc := NewService(NewValidator(), NewParser(), NewDetector(), nil)

	if err := svc.Validate("DE89370400440532013000"); err != nil {
		t.Fatalf("Validate() unexpected error: %v", err)
	}
	if err := svc.Validate("DE89370400440532013001"); err == nil {
		t.Fatalf("Validate() expected checksum error")
	}

	info, err := svc.Parse("DE89370400440532013000")
	if err != nil {
		t.Fatalf("Parse() unexpected error: %v", err)
	}
	if info.BankCode != "37040044" {
		t.Fatalf("Parse() bank code = %s, want 37040044", info.BankCode)
	}

	st, err := svc.Detect("DE89370400440532013000")
	if err != nil {
		t.Fatalf("Detect() unexpected error: %v", err)
	}
	if st.Structure != "CCKKBBBBBBBBAAAAAAAAAA" {
		t.Fatalf("Detect() structure = %s, want CCKKBBBBBBBBAAAAAAAAAA", st.Structure)
	}
}
