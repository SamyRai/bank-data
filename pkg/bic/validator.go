// Package bic provides compatibility APIs for BIC/SWIFT validation.
package bic

import (
	"errors"
	"regexp"
	"strings"
)

var (
	ErrBICFormat      = errors.New("invalid BIC format")
	ErrBICLength      = errors.New("invalid BIC length")
	ErrInvalidCountry = errors.New("invalid country code")
)

type BICInfo struct {
	Institution string
	Country     string
	Location    string
	Branch      string
	Raw         string
}

var bicRegex = regexp.MustCompile(`^[A-Z]{4}[A-Z]{2}[A-Z0-9]{2}([A-Z0-9]{3})?$`)

type ValidationResult struct {
	Input   string
	Valid   bool
	Code    string
	Message string
}

type Validator struct{}

func NewValidator() *Validator { return &Validator{} }

func (v *Validator) Validate(input string) ValidationResult {
	res := ValidationResult{Input: input, Valid: true}
	norm := strings.ToUpper(strings.TrimSpace(input))
	if len(norm) != 8 && len(norm) != 11 {
		res.Valid = false
		res.Code = "invalid_length"
		res.Message = ErrBICLength.Error()
		return res
	}
	if !bicRegex.MatchString(norm) {
		res.Valid = false
		res.Code = "invalid_format"
		res.Message = ErrBICFormat.Error()
		return res
	}
	// Country code must be letters.
	if norm[4] < 'A' || norm[4] > 'Z' || norm[5] < 'A' || norm[5] > 'Z' {
		res.Valid = false
		res.Code = "invalid_country"
		res.Message = ErrInvalidCountry.Error()
	}
	return res
}

func Parse(bicStr string) (*BICInfo, error) {
	norm := strings.ToUpper(strings.ReplaceAll(bicStr, " ", ""))
	if !bicRegex.MatchString(norm) {
		return nil, ErrBICFormat
	}
	branch := "XXX"
	if len(norm) == 11 {
		branch = norm[8:11]
	}
	return &BICInfo{
		Institution: norm[0:4],
		Country:     norm[4:6],
		Location:    norm[6:8],
		Branch:      branch,
		Raw:         norm,
	}, nil
}
