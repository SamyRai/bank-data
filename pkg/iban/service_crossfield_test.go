package iban

import (
	"testing"

	"github.com/SamyRai/bank-data/internal/countrymeta"
)

func TestBuildIBANStructureString_BBANIndexing(t *testing.T) {
	meta := countrymeta.Registry["DE"]
	got := buildIBANStructureString(meta)
	want := "CCKKBBBBBBBBAAAAAAAAAA"
	if got != want {
		t.Fatalf("buildIBANStructureString(DE) = %s, want %s", got, want)
	}
}
