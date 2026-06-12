package iban_test

import (
	"bufio"
	"os"
	"strings"
	"testing"

	internal "github.com/SamyRai/bank-data/internal/validationformats/iban"
	"github.com/SamyRai/bank-data/pkg/iban"
)

func TestIBANVectors(t *testing.T) {
	f, err := os.Open("../../testdata/iban_test_vectors.txt")
	if err != nil {
		t.Fatalf("failed to open test vectors: %v", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			t.Fatalf("failed to close test vectors: %v", err)
		}
	}()
	scanner := bufio.NewScanner(f)
	service := iban.NewService(
		internal.NewValidator(),
		internal.NewParser(),
		internal.NewDetector(),
		nil, // pass nil for registry
	)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") || len(strings.TrimSpace(line)) == 0 {
			continue
		}
		parts := strings.Split(line, ",")
		if len(parts) < 2 {
			continue
		}
		ibanStr := parts[0]
		valid := parts[1] == "1"
		err := service.Validate(ibanStr)
		if valid && err != nil {
			t.Errorf("expected valid IBAN %q, got error: %v", ibanStr, err)
		}
		if !valid && err == nil {
			t.Errorf("expected invalid IBAN %q, got no error", ibanStr)
		}
	}
	if err := scanner.Err(); err != nil {
		t.Fatalf("scanner error: %v", err)
	}
}

func TestIBANInputLengthCapping(t *testing.T) {
	service := iban.NewService(
		internal.NewValidator(),
		internal.NewParser(),
		internal.NewDetector(),
		nil,
	)
	// Overlength IBAN (35 chars, valid country code)
	ibanStr := "DE" + strings.Repeat("1", 33)
	err := service.Validate(ibanStr)
	if err == nil {
		t.Errorf("expected error for overlength IBAN, got nil")
	}
	if !strings.Contains(err.Error(), "length") {
		t.Errorf("expected length error, got: %v", err)
	}
	// Max length IBAN (34 chars, valid country code, but invalid format)
	ibanStr = "DE" + strings.Repeat("1", 32)
	err = service.Validate(ibanStr)
	if err == nil {
		t.Errorf("expected error for invalid format, got nil")
	}
	if !strings.Contains(err.Error(), "wrong length") && !strings.Contains(err.Error(), "format") {
		t.Errorf("expected wrong length or format error for 34-char IBAN, got: %v", err)
	}
}

func TestIBANConstantTimeCheckDigit(t *testing.T) {
	// Valid IBAN with correct check digit
	ibanStr := "DE89370400440532013000"
	if !internal.TestValidateIBANChecksum(ibanStr) {
		t.Errorf("expected valid checksum for %s", ibanStr)
	}
	// Invalid IBAN with wrong check digit
	ibanStrWrong := "DE89370400440532013001"
	if internal.TestValidateIBANChecksum(ibanStrWrong) {
		t.Errorf("expected invalid checksum for %s", ibanStrWrong)
	}
}

func TestIBANInvalidChars(t *testing.T) {
	service := iban.NewService(
		internal.NewValidator(),
		internal.NewParser(),
		internal.NewDetector(),
		nil,
	)
	ibanStr := "DE89!@#400440532013000"
	err := service.Validate(ibanStr)
	if err == nil || !strings.Contains(err.Error(), "characters") {
		t.Errorf("expected invalid character error, got: %v", err)
	}
}

func TestIBANUnsupportedCountry(t *testing.T) {
	service := iban.NewService(
		internal.NewValidator(),
		internal.NewParser(),
		internal.NewDetector(),
		nil,
	)
	ibanStr := "ZZ89370400440532013000"
	err := service.Validate(ibanStr)
	if err == nil || !strings.Contains(err.Error(), "country") {
		t.Errorf("expected unsupported country error, got: %v", err)
	}
}

func TestIBANChecksumEdgeCases(t *testing.T) {
	service := iban.NewService(
		internal.NewValidator(),
		internal.NewParser(),
		internal.NewDetector(),
		nil,
	)
	// Valid IBAN
	ibanStr := "DE89370400440532013000"
	err := service.Validate(ibanStr)
	if err != nil {
		t.Errorf("expected valid IBAN, got error: %v", err)
	}
	// IBAN with valid format but wrong checksum
	ibanStrWrong := "DE89370400440532013001"
	err = service.Validate(ibanStrWrong)
	if err == nil || !strings.Contains(err.Error(), "checksum") {
		t.Errorf("expected checksum error, got: %v", err)
	}
}
