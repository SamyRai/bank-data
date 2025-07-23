// cmd/download.go: Standalone downloader for IBAN registry
package main

import (
	"github.com/SamyRai/bank-data/gen/ibanfetch"
)

func main() {
	if err := ibanfetch.DownloadIBANRegistry(); err != nil {
		panic(err)
	}
}
