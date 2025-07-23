package core

import (
	"bufio"
	"log"
	"os"
)

// RobustLineReader reads a file line-by-line and applies a user-supplied parse function.
// It logs and skips malformed lines, and returns all successfully parsed records.
func RobustLineReader[T any](path string, parseFunc func(string) (T, error)) ([]T, int, int, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, 0, 0, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var records []T
	total, skipped := 0, 0
	for scanner.Scan() {
		total++
		rec, err := parseFunc(scanner.Text())
		if err != nil {
			log.Printf("Skipping line %d: %v", total, err)
			skipped++
			continue
		}
		records = append(records, rec)
	}
	return records, total, skipped, scanner.Err()
}
