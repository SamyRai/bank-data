package nationalaccount

import "testing"

func TestValidator_ValidateAndParse(t *testing.T) {
	v := NewValidator()
	ok := v.Validate("20-00-00 55779911")
	if !ok.Valid {
		t.Fatalf("expected valid account, got %+v", ok)
	}
	bad := v.Validate("00-00-00 00000000")
	if bad.Valid {
		t.Fatalf("expected invalid account for all-zero sort/account")
	}
	info, err := v.Parse("20-00-00 55779911")
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	if info.SortCode != "200000" || info.AccountNumber != "55779911" {
		t.Fatalf("unexpected parsed info: %+v", info)
	}
}
