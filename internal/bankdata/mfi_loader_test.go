package bankdata

import (
	"os"
	"strings"
	"testing"
)

const mfiSample = `RIADCode	LEI	CountryOfRegistration	Name	Box	Address	Postal	City	Category	HeadCountryOfRegistration	HeadName	HeadRIADCode	HeadLEI	Report
12345	LEI123	DE	Test Bank	Box1	Addr1	12345	Berlin	CAT1	DE	Head Bank	54321	LEI543	2025-07-23
67890	LEI678	FR	Other Bank	Box2	Addr2	67890	Paris	CAT2	FR	Head Other	09876	LEI098	2025-07-23
`

func writeTempMFIFile(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "mfi_test_*.tsv")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		if cerr := f.Close(); cerr != nil {
			t.Errorf("failed to close temp file: %v", cerr)
		}
		t.Fatalf("failed to write temp file: %v", err)
	}
	if err := f.Close(); err != nil {
		t.Errorf("failed to close temp file: %v", err)
	}
	return f.Name()
}

func TestLoadMFIDataset(t *testing.T) {
	path := writeTempMFIFile(t, mfiSample)
	defer os.Remove(path)
	recs, err := LoadMFIDataset(path)
	if err != nil {
		t.Fatalf("LoadMFIDataset failed: %v", err)
	}
	// Skip header line
	var filtered []MFIRecord
	for _, r := range recs {
		if r.Name != "Name" { // header field
			filtered = append(filtered, r)
		}
	}
	if len(filtered) != 2 {
		t.Fatalf("expected 2 records, got %d", len(filtered))
	}
	if filtered[0].Name != "Test Bank" || filtered[1].City != "Paris" {
		t.Errorf("unexpected record data: %+v", filtered)
	}
}

func TestMFIRecordReader(t *testing.T) {
	path := writeTempMFIFile(t, mfiSample)
	defer os.Remove(path)
	r, err := NewMFIRecordReader(path)
	if err != nil {
		t.Fatalf("NewMFIRecordReader failed: %v", err)
	}
	defer r.Close()
	var names []string
	for {
		rec, err := r.Next()
		if err != nil {
			if strings.Contains(err.Error(), "EOF") {
				break
			}
			t.Fatalf("Next failed: %v", err)
		}
		if rec.Name != "Name" && rec.Name != "" { // skip header
			names = append(names, rec.Name)
		}
	}
	if len(names) != 2 || names[0] != "Test Bank" || names[1] != "Other Bank" {
		t.Errorf("unexpected names: %v", names)
	}
}
