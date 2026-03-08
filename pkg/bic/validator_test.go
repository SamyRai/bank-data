package bic

import "testing"

func TestValidator_Validate(t *testing.T) {
	v := NewValidator()
	cases := []struct {
		bic  string
		ok   bool
		name string
	}{
		{"DEUTDEFF", true, "valid 8-char BIC"},
		{"DEUTDEFF500", true, "valid 11-char BIC"},
		{"DEUTDEFF5", false, "invalid length"},
		{"DEUTDEFFF0", false, "invalid length (10)"},
		{"DEUTUS33", true, "valid country token format"},
		{"DEUTDEFF@#", false, "invalid chars"},
	}
	for _, c := range cases {
		res := v.Validate(c.bic)
		if res.Valid != c.ok {
			t.Errorf("%s: Validate(%q) = valid:%v code:%s msg:%s, want ok=%v", c.name, c.bic, res.Valid, res.Code, res.Message, c.ok)
		}
	}
}

func TestParse(t *testing.T) {
	info, err := Parse("DEUTDEFF500")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	if info.Institution != "DEUT" || info.Country != "DE" || info.Branch != "500" {
		t.Fatalf("unexpected parse result: %+v", info)
	}
	info8, _ := Parse("DEUTDEFF")
	if info8.Branch != "XXX" {
		t.Errorf("expected XXX for 8-char BIC, got %s", info8.Branch)
	}
}
