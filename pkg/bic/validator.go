// Package bic provides tools for validating and parsing Bank Identifier Codes (BIC), also known as SWIFT codes.
package bic

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	bicmap "github.com/SamyRai/bank-data/internal/bic/map"
	"github.com/SamyRai/bank-data/internal/countrymeta"
	"github.com/SamyRai/bank-data/pkg/validation"
)

// Errors for typed checking
var (
	ErrBICFormat      = errors.New("invalid BIC format")
	ErrBICLength      = errors.New("invalid BIC length")
	ErrInvalidCountry = errors.New("invalid country code")
	ErrInactiveBIC    = errors.New("BIC is inactive")
	ErrNotFound       = errors.New("BIC not found in mapping or directory")
)

// Meta holds metadata for a BIC.
type Meta struct {
	Active        bool
	InactiveSince string // ISO date if inactive
}

// Directory provides lookup for BIC metadata.
var Directory = map[string]Meta{
	"DEUTDEFF": {Active: true},
	"NEDSZAJJ": {Active: false, InactiveSince: "2022-01-01"},
}

var bicmapData *bicmap.BankBICMap

// SetBankBICMap wires the country-agnostic bank code to BIC mapping for validation.
func SetBankBICMap(m *bicmap.BankBICMap) {
	bicmapData = m
}

// BIC length constants
const (
	BIC8  = 8
	BIC11 = 11
)

// BIC segments
type BICInfo struct {
	Institution string // 4 letters
	Country     string // 2 letters
	Location    string // 2 alphanumeric
	Branch      string // 3 alphanumeric (optional, 'XXX' means primary)
	Raw         string
}

// bicRegex validates the basic structure of a BIC
var bicRegex = regexp.MustCompile(`^[A-Z]{4}[A-Z]{2}[A-Z0-9]{2}([A-Z0-9]{3})?$`)

// Validator implements the validation.Validator interface for BIC.
type Validator struct{}

// NewValidator returns a new BIC Validator.
func NewValidator() *Validator {
	return &Validator{}
}

// Validate checks if the string is a valid BIC (format, country, mapping, and directory).
func (v *Validator) Validate(input string) validation.ValidationResult {
	res := validation.ValidationResult{
		Input: input,
		Valid: true,
	}

	norm := strings.ToUpper(strings.TrimSpace(input))

	// 1. Syntax check
	if len(norm) != BIC8 && len(norm) != BIC11 {
		res.Valid = false
		res.Error = ErrBICLength
		return res
	}
	if !bicRegex.MatchString(norm) {
		res.Valid = false
		res.Error = ErrBICFormat
		return res
	}

	// 2. Country check
	country := norm[4:6]
	if _, ok := countrymeta.Registry[country]; !ok {
		res.Valid = false
		res.Error = ErrInvalidCountry
		return res
	}

	// 3. Semantic check (mapping)
	foundInMapping := false
	if bicmapData != nil {
		bankCode := norm[:4]
		if entry, ok := bicmapData.LookupBIC(country, bankCode); ok {
			if entry.BIC != norm[:8] && entry.BIC != norm {
				res.Valid = false
				res.Error = fmt.Errorf("BIC does not match known bank code mapping")
				return res
			}
			foundInMapping = true
		}
	}

	// 4. Directory check
	meta, foundInDirectory := Directory[norm[:8]]
	if foundInDirectory && !meta.Active {
		res.Valid = false
		res.Error = fmt.Errorf("%w: %s", ErrInactiveBIC, meta.InactiveSince)
		return res
	}

	// 5. Final existence check
	if !foundInMapping && !foundInDirectory {
		res.Valid = false
		res.Error = ErrNotFound
		return res
	}

	return res
}

// Parse decompose a BIC string into its segments.
func Parse(bicStr string) (*BICInfo, error) {
	norm := strings.ToUpper(strings.ReplaceAll(bicStr, " ", ""))
	if !bicRegex.MatchString(norm) {
		return nil, fmt.Errorf("invalid BIC format")
	}

	info := &BICInfo{
		Institution: norm[0:4],
		Country:     norm[4:6],
		Location:    norm[6:8],
		Raw:         norm,
	}

	if len(norm) == 11 {
		info.Branch = norm[8:11]
	} else {
		info.Branch = "XXX" // Standard default
	}

	return info, nil
}
