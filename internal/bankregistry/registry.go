package bankregistry

import (
	"bufio"
	"errors"
	"fmt"
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
