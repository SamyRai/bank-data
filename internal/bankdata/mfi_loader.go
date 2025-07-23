// Package bankdata provides utilities for loading and querying MFI (Monetary Financial Institution) datasets.
package bankdata

// MFIRecord represents a row from the MFI dataset.
type MFIRecord struct {
	RIADCode                  string
	LEI                       string
	CountryOfRegistration     string
	Name                      string
	Box                       string
	Address                   string
	Postal                    string
	City                      string
	Category                  string
	HeadCountryOfRegistration string
	HeadName                  string
	HeadRIADCode              string
	HeadLEI                   string
	Report                    string
}

// LoadMFIDataset loads the MFI dataset from a CSV/TSV file using streaming.
func LoadMFIDataset(path string) ([]MFIRecord, error) {
	r, err := NewMFIRecordReader(path)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	var out []MFIRecord
	for {
		rec, err := r.Next()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return nil, err
		}
		if rec.RIADCode == "" && rec.LEI == "" {
			continue // skip incomplete/empty
		}
		out = append(out, rec)
	}
	return out, nil
}

// FindByLEI returns the first MFIRecord with the given LEI.
func FindByLEI(records []MFIRecord, lei string) *MFIRecord {
	for _, r := range records {
		if r.LEI == lei {
			return &r
		}
	}
	return nil
}

// FindByRIADCode returns the first MFIRecord with the given RIAD code.
func FindByRIADCode(records []MFIRecord, code string) *MFIRecord {
	for _, r := range records {
		if r.RIADCode == code {
			return &r
		}
	}
	return nil
}
