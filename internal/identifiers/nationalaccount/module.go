package nationalaccount

import (
	"strings"

	"github.com/SamyRai/bank-data/internal/identifiers"
)

// Module validates and parses UK sort-code/account-number pairs.
type Module struct{}

func New() *Module { return &Module{} }

func (m *Module) Normalize(input string) string {
	replacer := strings.NewReplacer(" ", "", "-", "")
	return replacer.Replace(strings.TrimSpace(input))
}

func (m *Module) DetectCandidate(normalized string) bool {
	if len(normalized) != 14 {
		return false
	}
	for i := 0; i < len(normalized); i++ {
		if normalized[i] < '0' || normalized[i] > '9' {
			return false
		}
	}
	return true
}

func (m *Module) Validate(normalized string) error {
	if !m.DetectCandidate(normalized) {
		return &identifiers.ValidationError{Code: "invalid_format", Message: "UK national account must be 14 digits (6 sort code + 8 account number)"}
	}
	sortCode := normalized[:6]
	account := normalized[6:]
	if sortCode == "000000" {
		return &identifiers.ValidationError{Code: "invalid_sort_code", Message: "sort code cannot be all zeros"}
	}
	if account == "00000000" {
		return &identifiers.ValidationError{Code: "invalid_account_number", Message: "account number cannot be all zeros"}
	}
	return nil
}

func (m *Module) Parse(normalized string) (map[string]string, error) {
	if err := m.Validate(normalized); err != nil {
		return nil, err
	}
	return map[string]string{
		"country_code":    "GB",
		"sort_code":       normalized[:6],
		"account_number":  normalized[6:],
		"composite_number": normalized,
		"raw":             normalized,
	}, nil
}
