package bank_test

import (
	"testing"

	"github.com/SamyRai/bank-data/pkg/bank"
)

func TestNewBankInfo(t *testing.T) {
	info := bank.NewBankInfo("DE", "12345678", "DEUTDEFF", "Deutsche Bank")
	if info.CountryCode != "DE" {
		t.Errorf("expected CountryCode DE, got %s", info.CountryCode)
	}
	if info.BankCode != "12345678" {
		t.Errorf("expected BankCode 12345678, got %s", info.BankCode)
	}
	if info.BIC != "DEUTDEFF" {
		t.Errorf("expected BIC DEUTDEFF, got %s", info.BIC)
	}
	if info.BankName != "Deutsche Bank" {
		t.Errorf("expected BankName Deutsche Bank, got %s", info.BankName)
	}
}

func TestSetAddress(t *testing.T) {
	info := bank.NewBankInfo("DE", "12345678", "DEUTDEFF", "Deutsche Bank")
	info.SetAddress("Taunusanlage 12", "Frankfurt am Main")
	if info.Address != "Taunusanlage 12" {
		t.Errorf("expected Address Taunusanlage 12, got %s", info.Address)
	}
	if info.City != "Frankfurt am Main" {
		t.Errorf("expected City Frankfurt am Main, got %s", info.City)
	}
}

func TestAddCode(t *testing.T) {
	info := bank.NewBankInfo("DE", "12345678", "DEUTDEFF", "Deutsche Bank")
	info.AddCode("ABA", "123456789")
	if info.OtherCodes["ABA"] != "123456789" {
		t.Errorf("expected ABA 123456789, got %s", info.OtherCodes["ABA"])
	}
}

func TestIsSEPA(t *testing.T) {
	info := bank.NewBankInfo("DE", "12345678", "DEUTDEFF", "Deutsche Bank")
	if info.IsSEPA() {
		t.Error("expected IsSEPA to be false for now")
	}
}

func TestValidate(t *testing.T) {
	info := bank.NewBankInfo("DE", "12345678", "DEUTDEFF", "Deutsche Bank")
	if err := info.Validate(); err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	info.BankCode = ""
	if err := info.Validate(); err == nil {
		t.Error("expected error for missing BankCode, got nil")
	}
}

func TestBankInfoError_Error(t *testing.T) {
	err := &bank.BankInfoError{Msg: "test error"}
	if err.Error() != "test error" {
		t.Errorf("expected 'test error', got '%s'", err.Error())
	}
}
