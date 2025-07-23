// Package bicmap provides country-agnostic bank code to BIC mapping utilities.
package bicmap

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"

	"github.com/SamyRai/bank-data/pkg/bank"
)

type BankBICEntry struct {
	BankCode string // e.g. BLZ for DE, ABI/CAB for IT, etc.
	BIC      string
	BankName string
	Country  string // ISO 3166-1 alpha-2
}

type BankBICMap map[string]map[string]BankBICEntry // country -> bankcode -> entry

// LoadBankBICMap loads all country datasets from the datasets directory.
func LoadBankBICMap(datasetsDir string, countryFiles map[string]string) (BankBICMap, error) {
	m := make(BankBICMap)
	for country, file := range countryFiles {
		f, err := os.Open(fmt.Sprintf("%s/%s", datasetsDir, file))
		if err != nil {
			return nil, err
		}
		r := csv.NewReader(f)
		m[country] = make(map[string]BankBICEntry)
		for {
			rec, err := r.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				return nil, err
			}
			// For DE: rec[0]=BLZ, rec[5]=BIC, rec[3]=BankName
			entry := BankBICEntry{
				BankCode: rec[0],
				BIC:      rec[5],
				BankName: rec[3],
				Country:  country,
			}
			m[country][entry.BankCode] = entry
		}
		if err := f.Close(); err != nil {
			return nil, err
		}
	}
	return m, nil
}

// LookupBIC returns the BIC for a given country and bank code.
func (m BankBICMap) LookupBIC(country, bankCode string) (BankBICEntry, bool) {
	countryMap, ok := m[country]
	if !ok {
		return BankBICEntry{}, false
	}
	entry, ok := countryMap[bankCode]
	return entry, ok
}

// ToBankInfo converts a BankBICEntry to a global BankInfo entity.
func (e BankBICEntry) ToBankInfo() *bank.Info {
	return &bank.Info{
		CountryCode: e.Country,
		BankCode:    e.BankCode,
		BIC:         e.BIC,
		BankName:    e.BankName,
	}
}

// LookupBankInfo returns BankInfo for a given country and bank code.
func (m BankBICMap) LookupBankInfo(country, bankCode string) (*bank.Info, bool) {
	entry, ok := m.LookupBIC(country, bankCode)
	if !ok {
		return nil, false
	}
	return entry.ToBankInfo(), true
}
