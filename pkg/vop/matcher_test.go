package vop

import "testing"

func TestMatcher_Match(t *testing.T) {
	m := NewMatcher()

	exact, err := m.Match(MatchRequest{SuppliedName: "Acme GmbH", ExpectedName: "ACME"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if exact.Category != CategoryMatch {
		t.Fatalf("expected Match, got %+v", exact)
	}

	closeRes, err := m.Match(MatchRequest{SuppliedName: "Acme Trading", ExpectedName: "Acme Tradng"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if closeRes.Category != CategoryCloseMatch && closeRes.Category != CategoryMatch {
		t.Fatalf("expected CloseMatch/Match, got %+v", closeRes)
	}

	noMatch, err := m.Match(MatchRequest{SuppliedName: "Alpha", ExpectedName: "Zeta"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if noMatch.Category != CategoryNoMatch {
		t.Fatalf("expected NoMatch, got %+v", noMatch)
	}

	_, err = m.Match(MatchRequest{SuppliedName: "", ExpectedName: "Acme"})
	if err == nil {
		t.Fatalf("expected error for empty names")
	}
}
