// Package bicmap provides a streaming BIC dataset reader.
package bicmap

import (
	"fmt"
	"io"
	"log"
	"regexp"
	"strings"

	"github.com/SamyRai/bank-data/internal/core"
)

type BICEntryReader struct {
	entries []BankBICEntry
	idx     int
}

var (
	bankCodeRe = regexp.MustCompile(`^(\d{9})`)
)

// NewBICEntryReader opens a file and returns a robust reader using core.RobustLineReader.
func NewBICEntryReader(path string) (*BICEntryReader, error) {
	entries, total, skipped, err := core.RobustLineReader(path, func(line string) (BankBICEntry, error) {
		codeMatch := bankCodeRe.FindStringSubmatch(line)
		if len(codeMatch) == 2 {
			bankCode := codeMatch[1]
			rest := line[len(bankCode):]
			fields := strings.Fields(rest)
			if len(fields) >= 4 {
				bic := fields[len(fields)-2]
				return BankBICEntry{
					BankCode: bankCode,
					BankName: fields[0],
					BIC:      bic,
					Country:  "DE",
				}, nil
			}
		}
		return BankBICEntry{}, logError("malformed BIC record", line)
	})
	if err != nil {
		return nil, err
	}
	log.Printf("BICReader stats: total=%d, parsed=%d, skipped=%d", total, len(entries), skipped)
	return &BICEntryReader{entries: entries, idx: 0}, nil
}

// Next returns the next BankBICEntry or io.EOF.
func (r *BICEntryReader) Next() (BankBICEntry, error) {
	if r.idx >= len(r.entries) {
		return BankBICEntry{}, io.EOF
	}
	entry := r.entries[r.idx]
	r.idx++
	return entry, nil
}

func (r *BICEntryReader) Close() error {
	return nil // nothing to close
}

func logError(msg, line string) error {
	log.Printf("%s: %s", msg, line)
	return fmt.Errorf("%s: %s", msg, line)
}

var _ core.DatasetReader[BankBICEntry] = (*BICEntryReader)(nil)
