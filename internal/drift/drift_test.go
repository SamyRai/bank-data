package drift

import (
	"reflect"
	"testing"
)

func TestCalculateDrift(t *testing.T) {
	oldData := map[string]map[string]bool{
		"AT": {
			"12345": true,
			"67890": true,
		},
		"FR": {
			"11111": true,
		},
	}

	newData := map[string]map[string]bool{
		"AT": {
			"12345": true, // Kept
			// "67890" removed
		},
		"FR": {
			"11111": true, // Kept
			"22222": true, // Added
		},
		"NL": {
			"99999": true, // Added
		},
	}

	expectedReport := DriftReport{
		CountryChanges: map[string]int{
			"AT": -1,
			"FR": 1,
			"NL": 1,
		},
		AddedBankCodes: map[string][]string{
			"FR": {"22222"},
			"NL": {"99999"},
		},
		RemovedBankCodes: map[string][]string{
			"AT": {"67890"},
		},
	}

	report := CalculateDrift(oldData, newData)

	if !reflect.DeepEqual(report, expectedReport) {
		t.Errorf("CalculateDrift() mismatch.\nExpected: %v\nGot: %v", expectedReport, report)
	}
}

func TestCalculateDrift_NoChange(t *testing.T) {
	data := map[string]map[string]bool{
		"AT": {
			"12345": true,
		},
	}

	expectedReport := DriftReport{
		CountryChanges:   make(map[string]int),
		AddedBankCodes:   make(map[string][]string),
		RemovedBankCodes: make(map[string][]string),
	}

	report := CalculateDrift(data, data)

	if !reflect.DeepEqual(report, expectedReport) {
		t.Errorf("CalculateDrift() mismatch.\nExpected: %v\nGot: %v", expectedReport, report)
	}
}
