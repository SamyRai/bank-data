// Package dataset provides compact encoding and decoding utilities for key-value datasets.

package dataset

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/SamyRai/bank-data/internal/log"
)

// RegistryDecodeError represents an error during registry decoding.
type RegistryDecodeError struct {
	Reason string
	Data   []byte
}

func (e *RegistryDecodeError) Error() string {
	return fmt.Sprintf("registry decode error: %s", e.Reason)
}

// ChecksumError represents a checksum mismatch error.
type ChecksumError struct {
	Expected string
	Actual   string
}

func (e *ChecksumError) Error() string {
	return fmt.Sprintf("checksum mismatch: expected %s, got %s", e.Expected, e.Actual)
}

// Decode reads encoded records and returns a slice of parsed maps.
// reverseKeyMap is the inverse of the encoder map: e.g., {"N": "Name"}
// delimiter is the separator between key-value tokens (default: Sep)
// If checksum is non-empty, validates SHA256 of all input lines.
func Decode(r io.Reader, reverseKeyMap map[string]string, delimiter string, checksum string) ([]map[string]string, error) {
	if delimiter == "" {
		delimiter = Sep
	}
	var records []map[string]string
	scanner := bufio.NewScanner(r)
	var allLines []byte
	lineNum := 0
	for scanner.Scan() {
		line := scanner.Bytes()
		allLines = append(allLines, line...)
		allLines = append(allLines, '\n')
		lineNum++
		tokens := bytes.Split(line, []byte(delimiter))
		if len(tokens)%2 != 0 {
			log.Warn("Skipping malformed line", log.Fields{"lineNum": lineNum, "line": string(line)})
			continue
		}
		record := make(map[string]string)
		for i := 0; i < len(tokens)-1; i += 2 {
			shortKey := string(tokens[i])
			value := string(tokens[i+1])
			longKey, ok := reverseKeyMap[shortKey]
			if !ok {
				log.Warn("Unknown short key", log.Fields{"shortKey": shortKey, "lineNum": lineNum})
				continue
			}
			record[longKey] = value
		}
		if len(record) > 0 {
			records = append(records, record)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Error("Decode error", log.Fields{"lineNum": lineNum, "error": err})
		return nil, &RegistryDecodeError{Reason: fmt.Sprintf("decode error at line %d: %v", lineNum, err), Data: allLines}
	}
	if checksum != "" {
		h := sha256.Sum256(allLines)
		actual := hex.EncodeToString(h[:])
		if actual != checksum {
			log.Error("Checksum mismatch", log.Fields{"expected": checksum, "actual": actual})
			return nil, &ChecksumError{Expected: checksum, Actual: actual}
		}
	}
	log.Info("Decode completed", log.Fields{"records": len(records), "lines": lineNum})
	return records, nil
}

// DecodeStream streams records via callback, for memory efficiency.
func DecodeStream(r io.Reader, reverseKeyMap map[string]string, delimiter string, onRecord func(map[string]string, int)) error {
	if delimiter == "" {
		delimiter = Sep
	}
	scanner := bufio.NewScanner(r)
	lineNum := 0
	for scanner.Scan() {
		line := scanner.Bytes()
		lineNum++
		tokens := bytes.Split(line, []byte(delimiter))
		if len(tokens)%2 != 0 {
			log.Warn("Skipping malformed line", log.Fields{"lineNum": lineNum, "line": string(line)})
			continue
		}
		record := make(map[string]string)
		for i := 0; i < len(tokens)-1; i += 2 {
			shortKey := string(tokens[i])
			value := string(tokens[i+1])
			longKey, ok := reverseKeyMap[shortKey]
			if !ok {
				log.Warn("Unknown short key", log.Fields{"shortKey": shortKey, "lineNum": lineNum})
				continue
			}
			record[longKey] = value
		}
		if len(record) > 0 {
			onRecord(record, lineNum)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Error("DecodeStream error", log.Fields{"lineNum": lineNum, "error": err})
		return fmt.Errorf("decode error at line %d: %w", lineNum, err)
	}
	log.Info("DecodeStream completed", log.Fields{"lines": lineNum})
	return nil
}
