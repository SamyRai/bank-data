package main

import (
	"bytes"
	"fmt"
	"os"
	"sort"
	"text/template"

	"github.com/SamyRai/bank-data/internal/bankregistry"
)

type CoverageRow struct {
	CountryName         string
	CountryCode         string
	BankCodeSupport     bool
	LocalAccountSupport bool
	LastDatasetVersion  string
}

func main() {
	r := bankregistry.New()

	if err := loadCoverageFixtures(r); err != nil {
		fmt.Fprintf(os.Stderr, "fixture load error: %v\n", err)
		os.Exit(1)
	}

	var rows []CoverageRow
	countries := r.Countries()

	countryNames := map[string]string{
		"FR": "France",
		"AT": "Austria",
		"NL": "Netherlands",
		"ES": "Spain",
		"IT": "Italy",
		"PL": "Poland",
		"SE": "Sweden",
		"CH": "Switzerland",
		"GB": "United Kingdom",
		"DE": "Germany",
	}
	localAccountSupport := map[string]bool{
		"GB": true,
	}

	for _, cc := range countries {
		meta, err := r.Metadata(cc)
		vd := "Unknown"
		if err == nil && meta.VersionDate != "" {
			vd = meta.VersionDate
		}

		name, ok := countryNames[cc]
		if !ok {
			name = cc
		}

		rows = append(rows, CoverageRow{
			CountryName:         name,
			CountryCode:         cc,
			BankCodeSupport:     true,
			LocalAccountSupport: localAccountSupport[cc],
			LastDatasetVersion:  vd,
		})
	}

	sort.Slice(rows, func(i, j int) bool {
		return rows[i].CountryCode < rows[j].CountryCode
	})

	const tmpl = `# Country Coverage Matrix

This document outlines the level of support currently available for different countries within the ` + "`bank-data`" + ` library.

| Country | Code | Bank-Code Support | Local Account Support | Last Dataset Version |
| ------- | ---- | ----------------- | --------------------- | -------------------- |
{{- range .Rows }}
| {{ .CountryName }} | {{ .CountryCode }} | {{ if .BankCodeSupport }}✅{{ else }}❌{{ end }} | {{ if .LocalAccountSupport }}✅{{ else }}❌{{ end }} | {{ .LastDatasetVersion }} |
{{- end }}

_Note: This list is continually expanding. For feature requests or updating source datasets, refer to ` + "`CONTRIBUTING.md`" + `._
`

	t, err := template.New("coverage").Parse(tmpl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Template error: %v\n", err)
		os.Exit(1)
	}

	data := struct {
		Rows []CoverageRow
	}{
		Rows: rows,
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		fmt.Fprintf(os.Stderr, "Execute error: %v\n", err)
		os.Exit(1)
	}

	if err := os.WriteFile("docs/coverage_matrix.md", buf.Bytes(), 0600); err != nil {
		fmt.Fprintf(os.Stderr, "Write error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Generated docs/coverage_matrix.md")
}

func loadCoverageFixtures(r *bankregistry.Registry) error {
	meta := bankregistry.SourceMetadata{Source: "fixture", VersionDate: "2026-03-08", License: "test"}
	baseSchema := bankregistry.CSVSchema{Delimiter: ',', HasHeader: true, BankCodeColumn: 0, BankNameColumn: 1, BICColumn: 2}

	if err := r.LoadCSV("testdata/bankregistry/fr_sample.csv", "FR", baseSchema, meta); err != nil {
		return fmt.Errorf("load FR fixture: %w", err)
	}
	if err := r.LoadCSV("testdata/bankregistry/at_sample.csv", "AT", baseSchema, meta); err != nil {
		return fmt.Errorf("load AT fixture: %w", err)
	}
	if err := r.LoadCSV("testdata/bankregistry/nl_sample.csv", "NL", baseSchema, meta); err != nil {
		return fmt.Errorf("load NL fixture: %w", err)
	}
	if err := bankregistry.NewESLoader().Load("testdata/bankregistry/es_sample.csv", r, meta); err != nil {
		return fmt.Errorf("load ES fixture: %w", err)
	}
	if err := bankregistry.NewITLoader().Load("testdata/bankregistry/it_sample.csv", r, meta); err != nil {
		return fmt.Errorf("load IT fixture: %w", err)
	}
	if err := bankregistry.NewPLLoader().Load("testdata/bankregistry/pl_sample.csv", r, meta); err != nil {
		return fmt.Errorf("load PL fixture: %w", err)
	}
	if err := bankregistry.NewSELoader().Load("testdata/bankregistry/se_sample.csv", r, meta); err != nil {
		return fmt.Errorf("load SE fixture: %w", err)
	}
	if err := bankregistry.NewCHLoader().Load("testdata/bankregistry/ch_sample.csv", r, meta); err != nil {
		return fmt.Errorf("load CH fixture: %w", err)
	}
	if err := bankregistry.NewUKLoader().Load("testdata/bankregistry/uk_sample.csv", r, meta); err != nil {
		return fmt.Errorf("load GB fixture: %w", err)
	}
	if err := r.LoadDEFixedWidth("testdata/bankregistry/de_sample.txt", meta); err != nil {
		return fmt.Errorf("load DE fixture: %w", err)
	}

	return nil
}
