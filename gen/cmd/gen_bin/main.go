package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/SamyRai/bank-data/internal/bankdata"
	bicmap "github.com/SamyRai/bank-data/internal/bic/map"
	"github.com/SamyRai/bank-data/internal/core"
	"github.com/SamyRai/bank-data/internal/dataset"
)

func decodeAndPrint(input string, from, to, lines int) error {
	if filepath.IsAbs(input) || strings.Contains(input, "..") {
		return fmt.Errorf("invalid input file path: %s", input)
	}
	f, err := os.Open(input)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer f.Close()

	count := 0
	printed := 0

	r := bufio.NewReader(f)
	for {
		line, err := r.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		count++
		if (from > 0 && count < from) || (to > 0 && count > to) {
			continue
		}
		fmt.Print(line)
		printed++
		if lines > 0 && printed >= lines {
			break
		}
		if to > 0 && count >= to {
			break
		}
	}
	return nil
}

func encodeWithStats[T any](w io.Writer, reader func() (core.DatasetReader[T], error), encodeFunc func(T) ([]byte, error)) error {
	r, err := reader()
	if err != nil {
		return err
	}
	defer r.Close()
	var total, encoded, skipped int
	for {
		record, err := r.Next()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			skipped++
			continue
		}
		total++
		data, err := encodeFunc(record)
		if err != nil || len(data) == 0 {
			skipped++
			continue
		}
		if _, err := w.Write(data); err != nil {
			return err
		}
		encoded++
	}
	fmt.Printf("Encoding stats: total=%d, encoded=%d, skipped=%d\n", total, encoded, skipped)
	return nil
}

func main() {
	datasetType := flag.String("type", "iban", "Dataset type: iban|mfi|bic")
	input := flag.String("input", "datasets/iban-registry.txt", "Path to input file")
	output := flag.String("output", "assets/iban_registry.bin", "Path to output binary file")
	decode := flag.Bool("decode", false, "Decode and print binary file")
	lines := flag.Int("lines", 0, "Number of lines to print")
	from := flag.Int("from", 0, "Start line (inclusive)")
	to := flag.Int("to", 0, "End line (inclusive)")
	flag.Parse()

	if *decode {
		if err := decodeAndPrint(*input, *from, *to, *lines); err != nil {
			fmt.Fprintf(os.Stderr, "Decoding failed: %v\n", err)
			os.Exit(1)
		}
		return
	}

	if err := os.MkdirAll(filepath.Dir(*output), 0750); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create output directory: %v\n", err)
		os.Exit(1)
	}

	outFile, err := os.Create(*output)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create output file: %v\n", err)
		os.Exit(1)
	}
	defer outFile.Close()

	switch *datasetType {
	case "mfi":
		err = encodeWithStats(outFile, func() (core.DatasetReader[bankdata.MFIRecord], error) {
			return bankdata.NewMFIRecordReader(*input)
		}, func(rec bankdata.MFIRecord) ([]byte, error) {
			line := fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
				rec.RIADCode, rec.LEI, rec.CountryOfRegistration, rec.Name, rec.Box, rec.Address, rec.Postal, rec.City, rec.Category, rec.HeadCountryOfRegistration, rec.HeadName, rec.HeadRIADCode, rec.HeadLEI, rec.Report)
			return []byte(line), nil
		})
	case "bic":
		err = encodeWithStats(outFile, func() (core.DatasetReader[bicmap.BankBICEntry], error) {
			return bicmap.NewBICEntryReader(*input)
		}, func(rec bicmap.BankBICEntry) ([]byte, error) {
			line := fmt.Sprintf("%s,%s,%s,%s\n", rec.BankCode, rec.BankName, rec.BIC, rec.Country)
			return []byte(line), nil
		})
	case "iban":
		if filepath.IsAbs(*input) || strings.Contains(*input, "..") {
			fmt.Fprintf(os.Stderr, "invalid IBAN registry file path: %s\n", *input)
			os.Exit(1)
		}
		f, ibanErr := os.Open(*input)
		if ibanErr != nil {
			fmt.Fprintf(os.Stderr, "Failed to open IBAN registry: %v\n", ibanErr)
			os.Exit(1)
		}
		defer f.Close()
		// Use streaming parser for IBANMeta
		var metas []dataset.IBANMeta
		ibanErr = dataset.ParseIBANMetaStream(f, ',', func(meta dataset.IBANMeta, _ int) {
			metas = append(metas, meta)
		})
		if ibanErr != nil {
			fmt.Fprintf(os.Stderr, "Failed to parse IBAN registry: %v\n", ibanErr)
			os.Exit(1)
		}
		binData, err := dataset.EncodeIBANRegistry(metas)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to encode IBAN registry: %v\n", err)
			os.Exit(1)
		}
		if _, err := outFile.Write(binData); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to write IBAN binary: %v\n", err)
			os.Exit(1)
		}
	default:
		f, defaultErr := os.Open(*input)
		if defaultErr != nil {
			fmt.Fprintf(os.Stderr, "Failed to open input file: %v\n", defaultErr)
			os.Exit(1)
		}
		defer f.Close()
		// Use streaming parser for IBANMeta
		var metas []dataset.IBANMeta
		defaultErr = dataset.ParseIBANMetaStream(f, ',', func(meta dataset.IBANMeta, _ int) {
			metas = append(metas, meta)
		})
		if defaultErr != nil {
			fmt.Fprintf(os.Stderr, "Failed to parse IBAN registry: %v\n", defaultErr)
			os.Exit(1)
		}
		binData, err := dataset.EncodeIBANRegistry(metas)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to encode IBAN registry: %v\n", err)
			os.Exit(1)
		}
		if _, err := outFile.Write(binData); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to write IBAN binary: %v\n", err)
			os.Exit(1)
		}
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Encoding failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully generated binary asset: %s\n", *output)
}
