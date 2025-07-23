package assets

import (
	"embed"
)

// AssetsFS provides access to embedded assets (binary, text, etc.).
//
//go:embed *
var AssetsFS embed.FS
