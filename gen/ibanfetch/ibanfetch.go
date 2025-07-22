// Package ibanfetch: Fetches the official SWIFT IBAN registry CSV file
// download_iban_registry.go: Download the official SWIFT IBAN registry CSV for code generation
package ibanfetch

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/SamyRai/bank-data/internal/log"
)

const SwiftURL = "https://www.swift.com/swift-resource/11971/download"
const OutFile = "datasets/iban_registry.txt"

// progressReader wraps an io.Reader and prints progress to stdout.
type progressReader struct {
	reader     io.Reader
	total      int64
	read       int64
	lastUpdate time.Time
}

func (p *progressReader) Read(b []byte) (int, error) {
	n, err := p.reader.Read(b)
	p.read += int64(n)
	now := time.Now()
	if now.Sub(p.lastUpdate) > 500*time.Millisecond || err == io.EOF {
		p.printProgress()
		p.lastUpdate = now
	}
	return n, err
}

func (p *progressReader) printProgress() {
	if p.total > 0 {
		percent := float64(p.read) / float64(p.total) * 100
		fmt.Printf("\rDownloading: %.2f%% (%d/%d bytes)", percent, p.read, p.total)
	} else {
		fmt.Printf("\rDownloading: %d bytes", p.read)
	}
	if p.read == p.total {
		fmt.Println()
	}
}

// DownloadIBANRegistry downloads the official SWIFT IBAN registry CSV to OutFile with progress.
func DownloadIBANRegistry() error {
	// If the file already exists, skip download
	if _, err := os.Stat(OutFile); err == nil {
		log.Info("IBAN registry file already exists, skipping download", log.Fields{"file": OutFile})
		return nil
	}
	resp, err := http.Get(SwiftURL)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			log.Warn("Failed to close response body", log.Fields{"err": cerr.Error()})
		}
	}()
	f, err := os.Create(OutFile)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := f.Close(); cerr != nil {
			log.Warn("Failed to close output file", log.Fields{"err": cerr.Error()})
		}
	}()
	total := resp.ContentLength
	pr := &progressReader{reader: resp.Body, total: total}
	_, err = io.Copy(f, pr)
	if err != nil {
		return err
	}
	pr.printProgress() // ensure 100% is printed
	log.Info("Downloaded IBAN registry", log.Fields{"file": OutFile})
	return nil
}
