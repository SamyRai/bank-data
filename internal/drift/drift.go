package drift

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// DriftReport represents the output JSON for a drift report.
type DriftReport struct {
	CountryChanges   map[string]int      `json:"country_changes"`
	AddedBankCodes   map[string][]string `json:"added_bank_codes"`
	RemovedBankCodes map[string][]string `json:"removed_bank_codes"`
}

func CalculateDrift(oldData, newData map[string]map[string]bool) DriftReport {
	report := DriftReport{
		CountryChanges:   make(map[string]int),
		AddedBankCodes:   make(map[string][]string),
		RemovedBankCodes: make(map[string][]string),
	}

	// Compare old to new to find removed codes
	for country, codes := range oldData {
		for code := range codes {
			if _, ok := newData[country]; !ok {
				report.RemovedBankCodes[country] = append(report.RemovedBankCodes[country], code)
				report.CountryChanges[country]--
			} else if _, ok := newData[country][code]; !ok {
				report.RemovedBankCodes[country] = append(report.RemovedBankCodes[country], code)
				report.CountryChanges[country]--
			}
		}
	}

	// Compare new to old to find added codes
	for country, codes := range newData {
		if _, ok := oldData[country]; !ok {
			for code := range codes {
				report.AddedBankCodes[country] = append(report.AddedBankCodes[country], code)
				report.CountryChanges[country]++
			}
		} else {
			for code := range codes {
				if _, ok := oldData[country][code]; !ok {
					report.AddedBankCodes[country] = append(report.AddedBankCodes[country], code)
					report.CountryChanges[country]++
				}
			}
		}
	}

	// Calculate net country changes (remove entries with 0 change)
	for country, change := range report.CountryChanges {
		if change == 0 {
			delete(report.CountryChanges, country)
		}
	}

	// Sort slices for deterministic output
	for _, v := range report.AddedBankCodes {
		sort.Strings(v)
	}
	for _, v := range report.RemovedBankCodes {
		sort.Strings(v)
	}

	return report
}

func ParseCSV(path string, defaultCountry string) (map[string]map[string]bool, error) {
	safePath, err := sanitizeRelativePath(path)
	if err != nil {
		return nil, err
	}

	// #nosec G304 -- input path is sanitized to a repository-relative path.
	file, err := os.Open(safePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data := make(map[string]map[string]bool)
	scanner := bufio.NewScanner(file)

	// Assume first line is header
	if scanner.Scan() {
		_ = scanner.Text()
	}

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ",")

		// Typically our csv looks like: BankCode, BIC, BankName
		// If it's a generic CSV with Country first, then use it
		var country, code string
		if len(parts) >= 3 {
			// This could be our standard registry format, or something else.
			// Let's assume standard format BankCode is parts[0]
			code = strings.TrimSpace(parts[0])
			country = defaultCountry

			// If it looks like it has a country prefix: "AT,12345,..."
			if len(parts[0]) == 2 && parts[0] == strings.ToUpper(parts[0]) {
				country = strings.TrimSpace(parts[0])
				code = strings.TrimSpace(parts[1])
			}

			if country != "" && code != "" {
				if data[country] == nil {
					data[country] = make(map[string]bool)
				}
				data[country][code] = true
			}
		}
	}

	return data, scanner.Err()
}

func sanitizeRelativePath(path string) (string, error) {
	trimmed := strings.TrimSpace(path)
	if trimmed == "" {
		return "", fmt.Errorf("path is empty")
	}

	clean := filepath.Clean(trimmed)
	if filepath.IsAbs(clean) {
		return "", fmt.Errorf("absolute paths are not allowed")
	}
	if clean == "." || clean == ".." || strings.HasPrefix(clean, ".."+string(os.PathSeparator)) {
		return "", fmt.Errorf("path must stay within repository")
	}

	return clean, nil
}
