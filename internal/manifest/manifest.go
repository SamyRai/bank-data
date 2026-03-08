package manifest

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

type Manifest struct {
	PolicyVersion int       `json:"policy_version"`
	UpdatedAt     string    `json:"updated_at"`
	Datasets      []Dataset `json:"datasets"`
}

type Dataset struct {
	ID                  string `json:"id"`
	Path                string `json:"path"`
	Format              string `json:"format"`
	Source              string `json:"source"`
	SourceURL           string `json:"source_url"`
	VersionDate         string `json:"version_date"`
	License             string `json:"license"`
	Checksum            string `json:"checksum"`
	GenerationTimestamp string `json:"generation_timestamp"`
}

var checksumSHA256Regex = regexp.MustCompile(`^[a-f0-9]{64}$`)

func ValidateManifest(data []byte, calculateChecksumFunc func(string) (string, error)) []error {
	var errs []error
	var manifest Manifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return []error{fmt.Errorf("error parsing manifest JSON: %v", err)}
	}

	if manifest.PolicyVersion <= 0 {
		errs = append(errs, fmt.Errorf("Missing or invalid policy_version"))
	}
	if manifest.UpdatedAt == "" {
		errs = append(errs, fmt.Errorf("Missing updated_at"))
	} else if _, err := time.Parse(time.RFC3339, manifest.UpdatedAt); err != nil {
		errs = append(errs, fmt.Errorf("Invalid updated_at format, must be RFC3339"))
	}

	seenDatasetIDs := make(map[string]struct{}, len(manifest.Datasets))
	for i, ds := range manifest.Datasets {
		if ds.ID == "" {
			errs = append(errs, fmt.Errorf("Dataset %d: Missing id", i))
		} else {
			if _, exists := seenDatasetIDs[ds.ID]; exists {
				errs = append(errs, fmt.Errorf("Dataset %s: Duplicate id", ds.ID))
			}
			seenDatasetIDs[ds.ID] = struct{}{}
		}
		if strings.TrimSpace(ds.Path) == "" {
			errs = append(errs, fmt.Errorf("Dataset %s: Missing path", ds.ID))
		}
		if strings.TrimSpace(ds.Format) == "" {
			errs = append(errs, fmt.Errorf("Dataset %s: Missing format", ds.ID))
		}
		if strings.TrimSpace(ds.Source) == "" {
			errs = append(errs, fmt.Errorf("Dataset %s: Missing source", ds.ID))
		}
		if ds.SourceURL == "" {
			errs = append(errs, fmt.Errorf("Dataset %s: Missing source_url", ds.ID))
		}
		if ds.License == "" {
			errs = append(errs, fmt.Errorf("Dataset %s: Missing license", ds.ID))
		}
		if ds.VersionDate == "" {
			errs = append(errs, fmt.Errorf("Dataset %s: Missing version_date", ds.ID))
		} else if _, err := time.Parse("2006-01-02", ds.VersionDate); err != nil {
			errs = append(errs, fmt.Errorf("Dataset %s: Invalid version_date format, must be YYYY-MM-DD", ds.ID))
		}
		if !checksumSHA256Regex.MatchString(ds.Checksum) {
			errs = append(errs, fmt.Errorf("Dataset %s: Invalid checksum '%s'", ds.ID, ds.Checksum))
		}
		if ds.GenerationTimestamp == "" {
			errs = append(errs, fmt.Errorf("Dataset %s: Missing generation_timestamp", ds.ID))
		} else {
			if _, err := time.Parse(time.RFC3339, ds.GenerationTimestamp); err != nil {
				errs = append(errs, fmt.Errorf("Dataset %s: Invalid generation_timestamp format, must be RFC3339", ds.ID))
			}
		}

		if strings.TrimSpace(ds.Path) != "" && checksumSHA256Regex.MatchString(ds.Checksum) {
			// Verify checksum matches actual files only when both path and checksum shape are valid.
			actualChecksum, err := calculateChecksumFunc(ds.Path)
			if err != nil {
				errs = append(errs, fmt.Errorf("Dataset %s: Could not calculate checksum for %s: %v", ds.ID, ds.Path, err))
			} else if actualChecksum != ds.Checksum {
				errs = append(errs, fmt.Errorf("Dataset %s: Checksum mismatch for %s. Expected %s, got %s", ds.ID, ds.Path, ds.Checksum, actualChecksum))
			}
		}
	}
	return errs
}

func CalculateChecksum(pattern string) (string, error) {
	if err := validateSafeRelativePath(pattern); err != nil {
		return "", fmt.Errorf("invalid checksum path pattern: %w", err)
	}

	matches, err := filepath.Glob(pattern)
	if err != nil {
		return "", err
	}
	if len(matches) == 0 {
		return "", fmt.Errorf("no files found matching pattern")
	}
	sort.Strings(matches)

	// For consistent hashing of multiple files, sort by path and include file boundaries.
	hash := sha256.New()
	for _, match := range matches {
		if err := validateSafeRelativePath(match); err != nil {
			return "", fmt.Errorf("invalid matched dataset path %q: %w", match, err)
		}

		_, _ = io.WriteString(hash, match)
		_, _ = hash.Write([]byte{0})

		// #nosec G304 -- match values are constrained to sanitized repository-relative paths.
		file, err := os.Open(match)
		if err != nil {
			return "", err
		}
		if _, err := io.Copy(hash, file); err != nil {
			closeErr := file.Close()
			if closeErr != nil {
				return "", fmt.Errorf("copy failed: %v (close failed: %v)", err, closeErr)
			}
			return "", err
		}
		if err := file.Close(); err != nil {
			return "", fmt.Errorf("close failed for %s: %w", match, err)
		}
		_, _ = hash.Write([]byte{0})
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

func validateSafeRelativePath(path string) error {
	trimmed := strings.TrimSpace(path)
	if trimmed == "" {
		return fmt.Errorf("path is empty")
	}

	clean := filepath.Clean(trimmed)
	if filepath.IsAbs(clean) {
		return fmt.Errorf("absolute paths are not allowed")
	}
	if clean == "." || clean == ".." || strings.HasPrefix(clean, ".."+string(os.PathSeparator)) {
		return fmt.Errorf("path must stay within repository")
	}

	return nil
}
