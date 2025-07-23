// Package dataset provides compact encoding and decoding utilities for key-value datasets.

package dataset

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"sort"
	"strconv"
	"strings"

	"github.com/SamyRai/bank-data/internal/core"
)

const Sep = "\x1F"

// EncodeRecord serializes a single map[string]string using a field mapping (e.g. "Name" -> "N")
func EncodeRecord(record map[string]string, keyMap map[string]string) string {
	var pairs []string
	keys := make([]string, 0, len(keyMap))
	for k := range keyMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, longKey := range keys {
		shortKey, ok := keyMap[longKey]
		if !ok {
			continue
		}
		val, exists := record[longKey]
		if exists && val != "" {
			pairs = append(pairs, shortKey+Sep+val)
		}
	}
	return strings.Join(pairs, Sep) + "\n"
}

// Encode writes multiple records to writer using the compact encoding
func Encode(w io.Writer, records []map[string]string, keyMap map[string]string) error {
	for _, r := range records {
		_, err := w.Write([]byte(EncodeRecord(r, keyMap)))
		if err != nil {
			return err
		}
	}
	return nil
}

// IBANMeta represents a single IBAN registry entry.
type IBANMeta struct {
	Country      string
	Length       int
	Regex        string
	BankStart    int
	BankEnd      int
	AccountStart int
	AccountEnd   int
}

// ParseIBANMeta parses a CSV record into IBANMeta.
func ParseIBANMeta(rec []string) (IBANMeta, error) {
	if len(rec) < 7 {
		return IBANMeta{}, fmt.Errorf("invalid record length: %d", len(rec))
	}
	length, err := strconv.Atoi(rec[1])
	if err != nil {
		return IBANMeta{}, fmt.Errorf("invalid length: %w", err)
	}
	bankStart, err := strconv.Atoi(rec[3])
	if err != nil {
		return IBANMeta{}, fmt.Errorf("invalid bankStart: %w", err)
	}
	bankEnd, err := strconv.Atoi(rec[4])
	if err != nil {
		return IBANMeta{}, fmt.Errorf("invalid bankEnd: %w", err)
	}
	accountStart, err := strconv.Atoi(rec[5])
	if err != nil {
		return IBANMeta{}, fmt.Errorf("invalid accountStart: %w", err)
	}
	accountEnd, err := strconv.Atoi(rec[6])
	if err != nil {
		return IBANMeta{}, fmt.Errorf("invalid accountEnd: %w", err)
	}
	return IBANMeta{
		Country:      rec[0],
		Length:       length,
		Regex:        rec[2],
		BankStart:    bankStart,
		BankEnd:      bankEnd,
		AccountStart: accountStart,
		AccountEnd:   accountEnd,
	}, nil
}

// EncodeIBANRegistry encodes a slice of IBANMeta to binary.
func EncodeIBANRegistry(metas []IBANMeta) ([]byte, error) {
	var buf bytes.Buffer
	for _, meta := range metas {
		line := fmt.Sprintf("%s,%d,%s,%d,%d,%d,%d\n", meta.Country, meta.Length, meta.Regex, meta.BankStart, meta.BankEnd, meta.AccountStart, meta.AccountEnd)
		buf.WriteString(line)
	}
	return buf.Bytes(), nil
}

// ParseIBANMetaStream parses IBANMeta records from any io.Reader (CSV or TXT), streaming and skipping malformed lines.
func ParseIBANMetaStream(r io.Reader, delimiter rune, onRecord func(IBANMeta, int)) error {
	scanner := bufio.NewScanner(r)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		var rec []string
		if delimiter == ',' {
			// Use encoding/csv for robust CSV parsing
			csvReader := csv.NewReader(strings.NewReader(line))
			csvReader.Comma = delimiter
			parsed, err := csvReader.Read()
			if err != nil {
				log.Printf("Skipping malformed CSV at line %d: %v", lineNum, err)
				continue
			}
			rec = parsed
		} else {
			rec = strings.Split(line, string(delimiter))
		}
		meta, err := ParseIBANMeta(rec)
		if err != nil {
			log.Printf("Skipping malformed record at line %d: %v", lineNum, err)
			continue
		}
		onRecord(meta, lineNum)
	}
	return scanner.Err()
}

// EncodeFromReader encodes records from a DatasetReader to binary.
func EncodeFromReader[T any](w io.Writer, reader core.DatasetReader[T], encodeFunc func(T) ([]byte, error)) error {
	for {
		record, err := reader.Next()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return err
		}
		data, err := encodeFunc(record)
		if err != nil {
			return err
		}
		if _, err := w.Write(data); err != nil {
			return err
		}
	}
	return reader.Close()
}
