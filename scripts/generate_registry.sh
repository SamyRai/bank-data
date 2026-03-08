#!/usr/bin/env bash
set -euo pipefail

# Scripts generates the IBAN registry
echo "Downloading SWIFT IBAN registry..."
go run gen/cmd/download.go

echo "Regenerating internal/countrymeta/registry.go..."
go generate ./internal/countrymeta

echo "Updating manifest..."
# This is a bit tricky to do purely in bash, better to just let users know
echo "Please remember to update datasets/manifest.json with new checksums if datasets/iban-registry.txt has changed."

echo "Done."
