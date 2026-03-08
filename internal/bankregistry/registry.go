package bankregistry

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/SamyRai/bank-data/pkg/bank"
)

// SourceMetadata describes where dataset rows were sourced from.
type SourceMetadata struct {
	Source      string
	VersionDate string
	Checksum    string
	License     string
}

// Record is a normalized bank registry row.
type Record struct {
	CountryCode string
	BankCode    string
	BIC         string
	BankName    string
}

// CountryLoader loads bank records for one country into Registry.
type CountryLoader interface {
	Load(path string, registry *Registry, meta SourceMetadata) error
}

// CSVSchema describes column positions for CSV-based country datasets.
type CSVSchema struct {
	Delimiter      rune
	HasHeader      bool
	BankCodeColumn int
	BICColumn      int
	BankNameColumn int
}

// Registry stores bank lookup records by country and code.
type Registry struct {
	mu       sync.RWMutex
	records  map[string]map[string]Record
	metaByCC map[string]SourceMetadata
}

func New() *Registry {
	return &Registry{
		records:  make(map[string]map[string]Record),
		metaByCC: make(map[string]SourceMetadata),
	}
}

func (r *Registry) SetSourceMetadata(countryCode string, meta SourceMetadata) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.metaByCC[countryCode] = meta
}

func (r *Registry) Add(record Record) {
	r.mu.Lock()
	defer r.mu.Unlock()
	cc := strings.ToUpper(strings.TrimSpace(record.CountryCode))
	if _, ok := r.records[cc]; !ok {
		r.records[cc] = make(map[string]Record)
	}
	r.records[cc][strings.TrimSpace(record.BankCode)] = Record{
		CountryCode: cc,
		BankCode:    strings.TrimSpace(record.BankCode),
		BIC:         strings.ToUpper(strings.TrimSpace(record.BIC)),
		BankName:    strings.TrimSpace(record.BankName),
	}
}

func (r *Registry) LookupBank(countryCode, bankCode string) (*bank.BankInfo, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	rows, ok := r.records[strings.ToUpper(strings.TrimSpace(countryCode))]
	if !ok {
		return nil, false
	}
	row, ok := rows[strings.TrimSpace(bankCode)]
	if !ok {
		return nil, false
	}
	return &bank.BankInfo{
		CountryCode: row.CountryCode,
		BankCode:    row.BankCode,
		BIC:         row.BIC,
		BankName:    row.BankName,
	}, true
}

// LoadDEFixedWidth loads Bundesbank BLZ/BIC data from explicit path.
func (r *Registry) LoadDEFixedWidth(path string, meta SourceMetadata) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()

	s := bufio.NewScanner(f)
	loaded := 0
	for s.Scan() {
		line := s.Text()
		if len(line) < 150 {
			continue
		}
		blz := strings.TrimSpace(line[0:8])
		name := strings.TrimSpace(line[9:67])
		bic := strings.TrimSpace(line[139:150])
		if blz == "" || bic == "" {
			continue
		}
		r.Add(Record{CountryCode: "DE", BankCode: blz, BIC: bic, BankName: name})
		loaded++
	}
	if err := s.Err(); err != nil {
		return err
	}
	if loaded == 0 {
		return errors.New("no valid DE bank records loaded")
	}
	r.SetSourceMetadata("DE", meta)
	return nil
}

func (r *Registry) Metadata(countryCode string) (SourceMetadata, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	meta, ok := r.metaByCC[strings.ToUpper(strings.TrimSpace(countryCode))]
	if !ok {
		return SourceMetadata{}, fmt.Errorf("no metadata for country %s", countryCode)
	}
	return meta, nil
}

// Countries returns sorted list of countries currently loaded in the registry.
func (r *Registry) Countries() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]string, 0, len(r.records))
	for cc := range r.records {
		out = append(out, cc)
	}
	// small list; simple insertion sort avoids additional imports.
	for i := 1; i < len(out); i++ {
		for j := i; j > 0 && out[j] < out[j-1]; j-- {
			out[j], out[j-1] = out[j-1], out[j]
		}
	}
	return out
}

// LoadCSV loads country records from a deterministic local CSV fixture/file.
func (r *Registry) LoadCSV(path, countryCode string, schema CSVSchema, meta SourceMetadata) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()

	reader := csv.NewReader(f)
	if schema.Delimiter != 0 {
		reader.Comma = schema.Delimiter
	}

	row := 0
	loaded := 0
	for {
		rec, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		row++
		if schema.HasHeader && row == 1 {
			continue
		}
		maxCol := max(schema.BankCodeColumn, max(schema.BICColumn, schema.BankNameColumn))
		if len(rec) <= maxCol {
			continue
		}

		bankCode := strings.TrimSpace(rec[schema.BankCodeColumn])
		bic := strings.TrimSpace(rec[schema.BICColumn])
		name := strings.TrimSpace(rec[schema.BankNameColumn])
		if bankCode == "" || bic == "" {
			continue
		}
		r.Add(Record{
			CountryCode: countryCode,
			BankCode:    bankCode,
			BIC:         bic,
			BankName:    name,
		})
		loaded++
	}
	if loaded == 0 {
		return fmt.Errorf("no valid %s records loaded from %s", countryCode, path)
	}
	r.SetSourceMetadata(countryCode, meta)
	return nil
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
