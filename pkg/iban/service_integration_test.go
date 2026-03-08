package iban_test

import (
	"testing"

	"github.com/SamyRai/bank-data/pkg/bank"
	"github.com/SamyRai/bank-data/pkg/iban"
)

type lookup struct{}

func (lookup) LookupBank(country, bankCode string) (*bank.BankInfo, bool) {
	if country == "DE" && bankCode == "37040044" {
		return &bank.BankInfo{CountryCode: country, BankCode: bankCode, BIC: "COBADEFFXXX", BankName: "Commerzbank"}, true
	}
	return nil, false
}

func TestService_Parse_BICEnrichment(t *testing.T) {
	svc := iban.NewService(nil, nil, nil, lookup{})
	info, err := svc.Parse("DE89370400440532013000")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	if info.BIC != "COBADEFFXXX" {
		t.Fatalf("expected enriched BIC, got %q", info.BIC)
	}
	if info.BankName != "Commerzbank" {
		t.Fatalf("expected enriched BankName, got %q", info.BankName)
	}
}

func TestService_LookupBank(t *testing.T) {
	svc := iban.NewService(nil, nil, nil, lookup{})
	bi, err := svc.LookupBank("DE", "37040044")
	if err != nil {
		t.Fatalf("LookupBank failed: %v", err)
	}
	if bi.BIC == "" || bi.BankName == "" {
		t.Fatalf("expected populated bank info, got %+v", bi)
	}
}
