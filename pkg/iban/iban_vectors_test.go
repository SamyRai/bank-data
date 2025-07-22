package iban_test

import (
	"bufio"
	"os"
	"strings"
	"testing"

	internal "github.com/SamyRai/bank-data/internal/iban"
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
