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
	esLoader := NewESLoader()
	err = esLoader.Load("../../testdata/bankregistry/es_sample.csv", r, SourceMetadata{Source: "fixture", VersionDate: "2026-03-08", License: "test"})
	if err != nil {
		t.Fatalf("load ES fixture: %v", err)
	}
	itLoader := NewITLoader()
	err = itLoader.Load("../../testdata/bankregistry/it_sample.csv", r, SourceMetadata{Source: "fixture", VersionDate: "2026-03-08", License: "test"})
	if err != nil {
		t.Fatalf("load IT fixture: %v", err)
	}
	plLoader := NewPLLoader()
	err = plLoader.Load("../../testdata/bankregistry/pl_sample.csv", r, SourceMetadata{Source: "fixture", VersionDate: "2026-03-08", License: "test"})
	if err != nil {
		t.Fatalf("load PL fixture: %v", err)
	}
	seLoader := NewSELoader()
	err = seLoader.Load("../../testdata/bankregistry/se_sample.csv", r, SourceMetadata{Source: "fixture", VersionDate: "2026-03-08", License: "test"})
	if err != nil {
		t.Fatalf("load SE fixture: %v", err)
	}
	chLoader := NewCHLoader()
	err = chLoader.Load("../../testdata/bankregistry/ch_sample.csv", r, SourceMetadata{Source: "fixture", VersionDate: "2026-03-08", License: "test"})
	if err != nil {
		t.Fatalf("load CH fixture: %v", err)
	}
	ukLoader := NewUKLoader()
	err = ukLoader.Load("../../testdata/bankregistry/uk_sample.csv", r, SourceMetadata{Source: "fixture", VersionDate: "2026-03-08", License: "test"})
	if err != nil {
		t.Fatalf("load GB fixture: %v", err)
	}

	err = r.LoadDEFixedWidth("../../testdata/bankregistry/de_sample.txt", SourceMetadata{Source: "fixture", VersionDate: "2026-03-08", License: "test"})
	if err != nil {
		t.Fatalf("load DE fixture: %v", err)
	}

	meta, err := r.Metadata("DE")
	if err != nil || meta.VersionDate != "2026-03-08" {
		t.Fatalf("metadata for DE failed: %v, %v", err, meta)
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
	if bi, ok := r.LookupBank("ES", "0128"); !ok || bi.BIC != "BKTRESMM" {
		t.Fatalf("unexpected ES lookup: ok=%v bi=%+v", ok, bi)
	}
	if bi, ok := r.LookupBank("IT", "02008"); !ok || bi.BIC != "UNCRITM1" {
		t.Fatalf("unexpected IT lookup: ok=%v bi=%+v", ok, bi)
	}
	if bi, ok := r.LookupBank("PL", "1010"); !ok || bi.BIC != "NBPLPLPW" {
		t.Fatalf("unexpected PL lookup: ok=%v bi=%+v", ok, bi)
	}
	if bi, ok := r.LookupBank("SE", "5000"); !ok || bi.BIC != "ESSEBESS" {
		t.Fatalf("unexpected SE lookup: ok=%v bi=%+v", ok, bi)
	}
	if bi, ok := r.LookupBank("CH", "08000"); !ok || bi.BIC != "UBSWCHZH" {
		t.Fatalf("unexpected CH lookup: ok=%v bi=%+v", ok, bi)
	}
	if bi, ok := r.LookupBank("GB", "200000"); !ok || bi.BIC != "BARCGB22" {
		t.Fatalf("unexpected GB lookup: ok=%v bi=%+v", ok, bi)
	}
	if bi, ok := r.LookupBank("DE", "10000000"); !ok || bi.BIC != "DEUTDEFF" {
		t.Fatalf("unexpected DE lookup: ok=%v bi=%+v", ok, bi)
	}

	countries := r.Countries()
	if len(countries) != 10 {
		t.Fatalf("unexpected countries list length: %d", len(countries))
	}
}
