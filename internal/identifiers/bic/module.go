package bic

import (
	"regexp"
	"strings"

	"github.com/SamyRai/bank-data/internal/countrymeta"
	"github.com/SamyRai/bank-data/internal/identifiers"
)

var bicRegex = regexp.MustCompile(`^[A-Z]{4}[A-Z]{2}[A-Z0-9]{2}([A-Z0-9]{3})?$`)

// Module validates and parses BIC/SWIFT codes.
type Module struct{}

func New() *Module { return &Module{} }

func (m *Module) Normalize(input string) string {
	return strings.ToUpper(strings.ReplaceAll(strings.TrimSpace(input), " ", ""))
}

func (m *Module) DetectCandidate(normalized string) bool {
	l := len(normalized)
	if l != 8 && l != 11 {
		return false
	}
	return bicRegex.MatchString(normalized)
}

func (m *Module) Validate(normalized string) error {
	if !m.DetectCandidate(normalized) {
		return &identifiers.ValidationError{Code: "invalid_format", Message: "BIC must be 8 or 11 chars and match SWIFT format"}
	}
	if _, ok := countrymeta.Registry[normalized[4:6]]; !ok {
		return &identifiers.ValidationError{Code: "invalid_country", Message: "BIC country code is not recognized"}
	}
	return nil
}

func (m *Module) Parse(normalized string) (map[string]string, error) {
	if err := m.Validate(normalized); err != nil {
		return nil, err
	}
	branch := "XXX"
	if len(normalized) == 11 {
		branch = normalized[8:11]
	}
	return map[string]string{
		"institution_code": normalized[0:4],
		"country_code":     normalized[4:6],
		"location_code":    normalized[6:8],
		"branch_code":      branch,
		"raw":              normalized,
	}, nil
}
