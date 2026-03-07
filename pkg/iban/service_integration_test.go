package iban_test

import (
	"testing"

	"github.com/SamyRai/bank-data/pkg/iban"
)

func TestService_Parse_BICEnrichment(t *testing.T) {
	// Initialize service (will load datasets/blz_bic.csv from repo root if found)
	svc := iban.NewService(nil, nil, nil, nil)

	// Example German IBAN (Deutsche Bank BLZ 10070000)
	// BLZ: 10070000
	// Account: 0123456789
	// Check: 89
	input := "DE89100700000123456789"

	info, err := svc.Parse(input)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if info.CountryCode != "DE" {
		t.Errorf("expected country DE, got %s", info.CountryCode)
	}

	if info.BankCode != "10070000" {
		t.Errorf("expected bank code 10070000, got %s", info.BankCode)
	}

	// Verify BIC enrichment
	if info.BIC == "" {
		t.Error("expected BIC to be enriched, but got empty")
	}

	if info.BankName == "" {
		t.Error("expected BankName to be enriched, but got empty")
	}

	t.Logf("Enriched BIC: %s, BankName: %s", info.BIC, info.BankName)
}

func TestService_LookupBank(t *testing.T) {
	svc := iban.NewService(nil, nil, nil, nil)

	bank, err := svc.LookupBank("DE", "10070000")
	if err != nil {
		t.Fatalf("LookupBank failed: %v", err)
	}

	if bank.BIC == "" {
		t.Error("expected BIC in bank info")
	}

	if bank.BankName == "" {
		t.Error("expected BankName in bank info")
	}
}
