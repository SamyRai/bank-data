package manifest

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

func mockCalculateChecksum(expectedChecksums map[string]string) func(string) (string, error) {
	return func(pattern string) (string, error) {
		if sum, ok := expectedChecksums[pattern]; ok {
			return sum, nil
		}
		return "", fmt.Errorf("checksum for pattern %s not configured", pattern)
	}
}

const testChecksum = "abcd1234abcd1234abcd1234abcd1234abcd1234abcd1234abcd1234abcd1234"

func validDataset() Dataset {
	return Dataset{
		ID:                  "test-dataset",
		Path:                "test/path/*.csv",
		Format:              "csv",
		Source:              "example source",
		SourceURL:           "http://example.com",
		License:             "MIT",
		VersionDate:         "2023-01-01",
		Checksum:            testChecksum,
		GenerationTimestamp: "2023-01-01T12:00:00Z",
	}
}

func validManifest() Manifest {
	return Manifest{
		PolicyVersion: 1,
		UpdatedAt:     "2026-03-08T12:00:00Z",
		Datasets:      []Dataset{validDataset()},
	}
}

func assertContains(t *testing.T, errs []error, contains string) {
	t.Helper()
	for _, err := range errs {
		if strings.Contains(err.Error(), contains) {
			return
		}
	}
	t.Fatalf("expected error containing %q, got: %v", contains, errs)
}

func TestValidateManifest_Valid(t *testing.T) {
	manifest := validManifest()
	data, _ := json.Marshal(manifest)

	errs := ValidateManifest(data, mockCalculateChecksum(map[string]string{
		"test/path/*.csv": testChecksum,
	}))

	if len(errs) > 0 {
		t.Fatalf("Expected no errors, got: %v", errs)
	}
}

func TestValidateManifest_MissingFields(t *testing.T) {
	manifest := Manifest{
		UpdatedAt: "not-a-date",
		Datasets: []Dataset{
			{
				// Intentionally sparse to assert strict required fields.
			},
		},
	}
	data, _ := json.Marshal(manifest)
	errs := ValidateManifest(data, mockCalculateChecksum(nil))

	assertContains(t, errs, "missing or invalid policy_version")
	assertContains(t, errs, "invalid updated_at format")
	assertContains(t, errs, "Missing id")
	assertContains(t, errs, "Missing path")
	assertContains(t, errs, "Missing format")
	assertContains(t, errs, "Missing source")
	assertContains(t, errs, "Missing source_url")
	assertContains(t, errs, "Missing license")
	assertContains(t, errs, "Missing version_date")
	assertContains(t, errs, "Invalid checksum")
	assertContains(t, errs, "Missing generation_timestamp")
}

func TestValidateManifest_InvalidDate(t *testing.T) {
	manifest := validManifest()
	manifest.Datasets[0].GenerationTimestamp = "2023-01-01" // Not RFC3339

	data, _ := json.Marshal(manifest)
	errs := ValidateManifest(data, mockCalculateChecksum(map[string]string{
		"test/path/*.csv": testChecksum,
	}))

	if len(errs) != 1 {
		t.Fatalf("Expected 1 error, got: %v", errs)
	}
	if !strings.Contains(errs[0].Error(), "Invalid generation_timestamp format") {
		t.Errorf("Unexpected error message: %v", errs[0])
	}
}

func TestValidateManifest_ChecksumMismatch(t *testing.T) {
	manifest := validManifest()
	data, _ := json.Marshal(manifest)
	errs := ValidateManifest(data, mockCalculateChecksum(map[string]string{
		"test/path/*.csv": "differentchecksum",
	}))

	if len(errs) != 1 {
		t.Fatalf("Expected 1 error, got: %v", errs)
	}
	if !strings.Contains(errs[0].Error(), "Checksum mismatch") {
		t.Errorf("Unexpected error message: %v", errs[0])
	}
}

func TestValidateManifest_DuplicateDatasetID(t *testing.T) {
	ds1 := validDataset()
	ds2 := validDataset()
	ds2.Path = "another/path/*.csv"

	manifest := Manifest{
		PolicyVersion: 1,
		UpdatedAt:     "2026-03-08T12:00:00Z",
		Datasets:      []Dataset{ds1, ds2},
	}
	data, _ := json.Marshal(manifest)

	errs := ValidateManifest(data, mockCalculateChecksum(map[string]string{
		ds1.Path: testChecksum,
		ds2.Path: testChecksum,
	}))

	assertContains(t, errs, "Duplicate id")
}

func TestValidateManifest_InvalidVersionDate(t *testing.T) {
	manifest := validManifest()
	manifest.Datasets[0].VersionDate = "2023/01/01"

	data, _ := json.Marshal(manifest)
	errs := ValidateManifest(data, mockCalculateChecksum(map[string]string{
		"test/path/*.csv": testChecksum,
	}))

	assertContains(t, errs, "Invalid version_date format")
}
