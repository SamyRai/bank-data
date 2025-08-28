// Package bic provides BIC (SWIFT) code validation utilities for banking applications.
package bic

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	bicmap "github.com/SamyRai/bank-data/internal/bic/map"
	"github.com/SamyRai/bank-data/internal/countrymeta"
)

var (
	// BIC: 8 or 11 uppercase characters; optional branch code
	bicRegex = regexp.MustCompile(`^[A-Z]{4}[A-Z]{2}[A-Z0-9]{2}([A-Z0-9]{3})?$`)
)

// Meta holds metadata for a BIC, e.g., active/inactive status.
type Meta struct {
	Active        bool
	InactiveSince string // ISO date if inactive
}

// Directory provides lookup for BIC metadata.
var Directory = map[string]Meta{
	// Example entries; production version should load from SWIFTRef
	"DEUTDEFF": {Active: true},
	"NEDSZAJJ": {Active: false, InactiveSince: "2022-01-01"},
	"DABADEBB": {Active: false, InactiveSince: "2010-12-31"},
}

// Errors for typed checking (future-friendly)
var (
	ErrBICFormat      = errors.New("invalid BIC format")
	ErrBICLength      = errors.New("invalid BIC length")
	ErrInvalidCountry = errors.New("invalid country code")
	ErrInactiveBIC    = errors.New("BIC is inactive")
)

// isValidISOCountryCode checks if the code is a valid ISO 3166-1 alpha-2.
func isValidISOCountryCode(code string) bool {
	_, ok := countrymeta.Registry[code]
	return ok
}

// validateBICSyntax checks the basic syntax of a BIC.
// It validates the length (8 or 11 characters), format (alphanumeric),
// and ensures the country code corresponds to a valid ISO 3166-1 alpha-2 code.
func validateBICSyntax(bic string) error {
	bic = strings.ToUpper(strings.TrimSpace(bic))
	if len(bic) != 8 && len(bic) != 11 {
		return ErrBICLength
	}
	if !bicRegex.MatchString(bic) {
		return ErrBICFormat
	}
	country := bic[4:6]
	if !isValidISOCountryCode(country) {
		return ErrInvalidCountry
	}
	return nil
}

// validateBICSemantic performs semantic validation of a BIC.
// It checks if the bank code in the BIC maps to a known bank in the provided bank code map.
// This helps to ensure that the BIC belongs to a legitimate financial institution.
func validateBICSemantic(bic string) error {
	country := bic[4:6]
	if bicmapData != nil {
		bankCode := bic[:4]
		countryMap, countryExists := bicmapData[country]
		if countryExists {
			entry, ok := countryMap[bankCode]
			if ok {
				if entry.BIC != bic[:8] && entry.BIC != bic {
					return fmt.Errorf("BIC does not match known bank code mapping")
				}
			}
		}
	}
	return nil
}

// validateBICDirectory checks a BIC against a directory of known BICs.
// If strict mode is enabled, the BIC must be present and active in the directory.
// If strict mode is disabled, it only checks for inactive status if the BIC is found.
func validateBICDirectory(bic string, strict bool) error {
	meta, found := Directory[bic[:8]]
	if strict {
		if !found {
			return fmt.Errorf("BIC not found in directory")
		}
		if !meta.Active {
			return fmt.Errorf("%w: %s", ErrInactiveBIC, meta.InactiveSince)
		}
	} else {
		if found && !meta.Active {
			return fmt.Errorf("%w: %s", ErrInactiveBIC, meta.InactiveSince)
		}
	}
	return nil
}

// Validate performs a comprehensive validation of a BIC (SWIFT code).
// It checks for syntactic correctness, semantic validity against a bank code map,
// and presence and status in a directory of known BICs.
//
// The strictWhitelist parameter controls the directory validation behavior:
//   - If true, the BIC must be present and active in the directory.
//   - If false (default), the validation only fails if the BIC is found but is inactive.
//
// The function returns an error if the BIC is invalid, otherwise nil.
func Validate(bic string, strictWhitelist ...bool) error {
	bic = strings.ToUpper(strings.TrimSpace(bic))
	if err := validateBICSyntax(bic); err != nil {
		return err
	}
	if err := validateBICSemantic(bic); err != nil {
		return err
	}
	strict := len(strictWhitelist) > 0 && strictWhitelist[0]
	if err := validateBICDirectory(bic, strict); err != nil {
		return err
	}
	// If not found in mapping or directory, reject.
	// This ensures that we only accept BICs that are known to us.
	foundInMapping := false
	if bicmapData != nil {
		bankCode := bic[:4]
		country := bic[4:6]
		countryMap, countryExists := bicmapData[country]
		if countryExists {
			_, ok := countryMap[bankCode]
			if ok {
				foundInMapping = true
			}
		}
	}
	_, foundInDirectory := Directory[bic[:8]]
	if !foundInMapping && !foundInDirectory {
		return fmt.Errorf("BIC not found in mapping or directory")
	}
	return nil
}

var bicmapData bicmap.BankBICMap

// SetBankBICMap wires the country-agnostic bank code to BIC mapping for validation.
func SetBankBICMap(m bicmap.BankBICMap) {
	bicmapData = m
}

// BankInfo holds IBAN, BIC, and bank name for enrichment.
type BankInfo struct {
	IBAN     string
	BIC      string
	BankName string
}

// EnrichWithBIC extracts the bank code from IBAN (country-specific), looks up BIC and bank name.
// Currently supports Germany (DE) using BLZ → BIC mapping.
func EnrichWithBIC(iban string) (BankInfo, error) {
	iban = strings.ReplaceAll(iban, " ", "")
	if len(iban) < 12 {
		return BankInfo{}, fmt.Errorf("IBAN too short")
	}
	country := iban[:2]
	switch country {
	case "DE":
		if len(iban) < 12 {
			return BankInfo{}, fmt.Errorf("IBAN too short for DE")
		}
		blz := iban[4:12] // BLZ is 8 digits at pos 5-12 (0-based)
		if bicmapData == nil {
			return BankInfo{}, fmt.Errorf("bank code to BIC mapping not loaded")
		}
		countryMap, ok := bicmapData["DE"]
		if !ok {
			return BankInfo{}, fmt.Errorf("no mapping for DE")
		}
		entry, ok := countryMap[blz]
		if !ok {
			return BankInfo{}, fmt.Errorf("BLZ %s not found in mapping", blz)
		}
		return BankInfo{
			IBAN:     iban,
			BIC:      entry.BIC,
			BankName: entry.BankName,
		}, nil
	default:
		return BankInfo{}, fmt.Errorf("unsupported country: %s", country)
	}
}
