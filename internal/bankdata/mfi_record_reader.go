// Package bankdata provides a streaming MFI dataset reader.
package bankdata

import (
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/SamyRai/bank-data/internal/core"
)

type MFIRecordReader struct {
	records []MFIRecord
	idx     int
}

// NewMFIRecordReader opens a TSV file and returns a robust reader using core.RobustLineReader.
func NewMFIRecordReader(path string) (*MFIRecordReader, error) {
	records, total, skipped, err := core.RobustLineReader(path, func(line string) (MFIRecord, error) {
		fields := strings.Split(line, "\t")
		if len(fields) < 14 {
			return MFIRecord{}, logError("malformed MFI record", line)
		}
		return MFIRecord{
			RIADCode:                  fields[0],
			LEI:                       fields[1],
			CountryOfRegistration:     fields[2],
			Name:                      fields[3],
			Box:                       fields[4],
			Address:                   fields[5],
			Postal:                    fields[6],
			City:                      fields[7],
			Category:                  fields[8],
			HeadCountryOfRegistration: fields[9],
			HeadName:                  fields[10],
			HeadRIADCode:              fields[11],
			HeadLEI:                   fields[12],
			Report:                    fields[13],
		}, nil
	})
	if err != nil {
		return nil, err
	}
	log.Printf("MFIReader stats: total=%d, parsed=%d, skipped=%d", total, len(records), skipped)
	return &MFIRecordReader{records: records, idx: 0}, nil
}

// Next returns the next MFIRecord or io.EOF.
func (r *MFIRecordReader) Next() (MFIRecord, error) {
	if r.idx >= len(r.records) {
		return MFIRecord{}, io.EOF
	}
	record := r.records[r.idx]
	r.idx++
	return record, nil
}

func (r *MFIRecordReader) Close() error {
	return nil // nothing to close
}

func logError(msg, line string) error {
	log.Printf("%s: %s", msg, line)
	return fmt.Errorf("%s: %s", msg, line)
}

var _ core.DatasetReader[MFIRecord] = (*MFIRecordReader)(nil)
