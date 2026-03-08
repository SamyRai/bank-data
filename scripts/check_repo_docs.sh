#!/usr/bin/env bash
set -euo pipefail

has_pattern() {
  local pattern="$1"
  local file="$2"
  if command -v rg >/dev/null 2>&1; then
    rg -q "$pattern" "$file"
    return
  fi
  grep -Eq "$pattern" "$file"
}

required_files=(
  "README.md"
  "CHANGELOG.md"
  "SECURITY.md"
  "LICENSE"
  "CONTRIBUTING.md"
  "docs/api_reference.md"
  "docs/architecture.md"
  "docs/development.md"
  "docs/getting_started.md"
)

for f in "${required_files[@]}"; do
  if [[ ! -f "$f" ]]; then
    echo "missing required file: $f" >&2
    exit 1
  fi
done

if ! has_pattern '^## \[Unreleased\]' CHANGELOG.md; then
  echo "CHANGELOG.md must contain an [Unreleased] section" >&2
  exit 1
fi

if ! has_pattern 'Report a vulnerability' SECURITY.md; then
  echo "SECURITY.md must document private vulnerability reporting" >&2
  exit 1
fi

if ! has_pattern '^## Supported Identifier Types' README.md; then
  echo "README.md must document supported identifier types" >&2
  exit 1
fi

if ! has_pattern 'Canonical Public API: `pkg/financial`' docs/api_reference.md; then
  echo "docs/api_reference.md must document canonical pkg/financial API" >&2
  exit 1
fi

echo "doc checks passed"
