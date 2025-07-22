// main.go: Full pipeline CLI for IBAN registry update
package main

import (
	"os"

	"github.com/SamyRai/bank-data/gen/ibanfetch"
	"github.com/SamyRai/bank-data/gen/ibanmeta"
	"github.com/SamyRai/bank-data/internal/log"
)

func main() {
	if err := run(); err != nil {
		log.Error("Pipeline failed", log.Fields{"err": err.Error()})
		os.Exit(1)
	}
}

func run() error {
	if err := ibanfetch.DownloadIBANRegistry(); err != nil {
		log.Error("Failed to download IBAN registry", log.Fields{"err": err.Error()})
		return err
	}
	code, err := ibanmeta.GenerateRegistry("iban_registry.csv")
	if err != nil {
		log.Error("Registry generation failed", log.Fields{"err": err.Error()})
		return err
	}
	f, err := os.Create("../../internal/countrymeta/registry.go")
	if err != nil {
		log.Error("Failed to open registry.go for writing", log.Fields{"err": err.Error()})
		return err
	}
	defer func() {
		if cerr := f.Close(); cerr != nil {
			log.Warn("Failed to close registry.go", log.Fields{"err": cerr.Error()})
		}
	}()
	_, err = f.Write(code)
	if err != nil {
		log.Error("Failed to write registry.go", log.Fields{"err": err.Error()})
		return err
	}
	log.Info("Successfully updated registry.go from latest IBAN registry.", nil)
	return nil
}
