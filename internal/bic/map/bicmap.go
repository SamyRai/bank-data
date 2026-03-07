// Package bicmap provides country-agnostic bank code to BIC mapping utilities.
package bicmap

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/SamyRai/bank-data/internal/log"
	"github.com/SamyRai/bank-data/pkg/bank"
)

type BankBICEntry struct {
	BankCode string // e.g. BLZ for DE, ABI/CAB for IT, etc.
	BIC      string
	BankName string
	Country  string // ISO 3166-1 alpha-2
}

// BankBICMap manages bank-to-BIC mappings with lazy loading.
type BankBICMap struct {
	datasetsDir  string
	countryFiles map[string]string
	cache        map[string]map[string]BankBICEntry // country -> bankcode -> entry
	mu           sync.RWMutex
}

// NewBankBICMap initializes a new map manager.
func NewBankBICMap(datasetsDir string, countryFiles map[string]string) *BankBICMap {
	return &BankBICMap{
		datasetsDir:  datasetsDir,
		countryFiles: countryFiles,
		cache:        make(map[string]map[string]BankBICEntry),
	}
}

// DefaultLoader returns a BankBICMap manager with default configurations.
func DefaultLoader() (*BankBICMap, error) {
	// Root-relative path discovery
	datasetsDir := "datasets"
	if _, err := os.Stat(datasetsDir); os.IsNotExist(err) {
		datasetsDir = "../datasets"
		if _, err := os.Stat(datasetsDir); os.IsNotExist(err) {
			datasetsDir = "../../datasets"
		}
	}

	countryFiles := map[string]string{
		"DE": "blz_bic.csv",
	}
	return NewBankBICMap(datasetsDir, countryFiles), nil
}

// LookupBIC returns the BIC for a given country and bank code, loading the dataset if needed.
func (m *BankBICMap) LookupBIC(country, bankCode string) (BankBICEntry, bool) {
	m.mu.RLock()
	countryMap, ok := m.cache[country]
	m.mu.RUnlock()

	if !ok {
		// Try to load
		if err := m.loadCountry(country); err != nil {
			return BankBICEntry{}, false
		}
		m.mu.RLock()
		countryMap = m.cache[country]
		m.mu.RUnlock()
		if countryMap == nil {
			return BankBICEntry{}, false
		}
	}

	entry, ok := countryMap[bankCode]
	return entry, ok
}

func (m *BankBICMap) loadCountry(country string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Double check
	if _, ok := m.cache[country]; ok {
		return nil
	}

	file, ok := m.countryFiles[country]
	if !ok {
		return fmt.Errorf("no dataset for country %s", country)
	}

	path := fmt.Sprintf("%s/%s", m.datasetsDir, file)
	log.Debug("Lazy loading BIC map", log.Fields{"country": country, "path": path})

	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	entries := make(map[string]BankBICEntry)
	if country == "DE" {
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := scanner.Text()
			if len(line) < 150 {
				continue
			}
			blz := strings.TrimSpace(line[0:8])
			name := strings.TrimSpace(line[9:67])
			bic := strings.TrimSpace(line[139:150])
			if blz != "" && bic != "" {
				entries[blz] = BankBICEntry{
					BankCode: blz,
					BIC:      bic,
					BankName: name,
					Country:  country,
				}
			}
		}
	} else {
		r := csv.NewReader(f)
		for {
			rec, err := r.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}
			if len(rec) >= 6 {
				entries[rec[0]] = BankBICEntry{
					BankCode: rec[0],
					BIC:      rec[5],
					BankName: rec[3],
					Country:  country,
				}
			}
		}
	}

	m.cache[country] = entries
	return nil
}

// ToBankInfo converts a BankBICEntry to a global BankInfo entity.
func (e BankBICEntry) ToBankInfo() *bank.BankInfo {
	return &bank.BankInfo{
		CountryCode: e.Country,
		BankCode:    e.BankCode,
		BIC:         e.BIC,
		BankName:    e.BankName,
	}
}

// LookupBankInfo returns BankInfo for a given country and bank code.
func (m *BankBICMap) LookupBankInfo(country, bankCode string) (*bank.BankInfo, bool) {
	entry, ok := m.LookupBIC(country, bankCode)
	if !ok {
		return nil, false
	}
	return entry.ToBankInfo(), true
}
