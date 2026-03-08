package bankregistry

import "testing"

func TestLoadCSVFixtures(t *testing.T) {
	r := New()
	schema := CSVSchema{Delimiter: ',', HasHeader: true, BankCodeColumn: 0, BankNameColumn: 1, BICColumn: 2}

	err := r.LoadCSV("../../testdata/bankregistry/fr_sample.csv", "FR", schema, SourceMetadata{Source: "fixture", VersionDate: "2026-03-08", License: "test"})
	if err != nil {
		t.Fatalf("load FR fixture: %v", err)
	}
	err = r.LoadCSV("../../testdata/bankregistry/at_sample.csv", "AT", schema, SourceMetadata{Source: "fixture", VersionDate: "2026-03-08", License: "test"})
	if err != nil {
		t.Fatalf("load AT fixture: %v", err)
	}
	err = r.LoadCSV("../../testdata/bankregistry/nl_sample.csv", "NL", schema, SourceMetadata{Source: "fixture", VersionDate: "2026-03-08", License: "test"})
	if err != nil {
		t.Fatalf("load NL fixture: %v", err)
	}

	if bi, ok := r.LookupBank("FR", "30004"); !ok || bi.BIC != "BNPAFRPP" {
		t.Fatalf("unexpected FR lookup: ok=%v bi=%+v", ok, bi)
	}
	if bi, ok := r.LookupBank("AT", "12000"); !ok || bi.BIC != "BKAUATWW" {
		t.Fatalf("unexpected AT lookup: ok=%v bi=%+v", ok, bi)
	}
	if bi, ok := r.LookupBank("NL", "ABNA"); !ok || bi.BIC != "ABNANL2A" {
		t.Fatalf("unexpected NL lookup: ok=%v bi=%+v", ok, bi)
	}

	countries := r.Countries()
	if len(countries) != 3 || countries[0] != "AT" || countries[1] != "FR" || countries[2] != "NL" {
		t.Fatalf("unexpected countries list: %+v", countries)
	}
}
