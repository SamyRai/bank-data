// Package core provides shared abstractions for dataset streaming.
package core

// DatasetReader is a generic streaming interface for reading records from a dataset.
type DatasetReader[T any] interface {
	// Next returns the next record, or io.EOF when done.
	Next() (T, error)
	// Close releases any resources held by the reader.
	Close() error
}
