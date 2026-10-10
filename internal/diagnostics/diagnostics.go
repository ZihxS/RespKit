// Package diagnostics produces structured reports with no raw error strings.
package diagnostics

import "encoding/json"

// JSON serializes a known-safe diagnostic value for AI/CLI consumption.
func JSON(v any) string {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "{}"
	}
	return string(b)
}
