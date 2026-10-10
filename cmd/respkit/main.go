// Command respkit supplies deterministic AI-readable catalogs and diagnostics.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/ZihxS/RespKit/resp"
)

func printJSON(v any) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(v); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func main() {
	if len(os.Args) < 2 {
		if err := usage(); err != nil {
			os.Exit(1)
		}
		return
	}
	switch os.Args[1] {
	case "version":
		fmt.Println(resp.Version)
	case "catalog":
		printJSON(resp.Catalog())
	case "explain":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "usage: respkit explain ERROR_CODE")
			os.Exit(2)
		}
		for _, item := range resp.Catalog() {
			if item.Code == strings.ToUpper(os.Args[2]) {
				printJSON(item)
				return
			}
		}
		fmt.Fprintln(os.Stderr, "unknown error code")
		os.Exit(1)
	case "doctor":
		// Honest: no hidden dependency scanning or runtime adapter detection.
		printJSON(map[string]any{
			"tool": "respkit", "version": resp.Version,
			"default_naming": "snake_case", "default_locale": "id",
			"built_in_locales": []string{"id", "en"},
			"supported_contexts": []string{
				"net/http (HTTP wrapper)",
				"Chi (HTTP wrapper)",
				"Gorilla/mux (HTTP wrapper)",
				"Beego v2 (HTTP wrapper)",
				"Gin (structural)",
				"Echo v4/v5 (structural)",
				"Fiber v2/v3 (structural)",
			},
			"note": "framework integration requires runtime verification within the consuming project; " +
				"this command does not inspect the process",
		})
	default:
		if err := usage(); err != nil {
			os.Exit(1)
		}
		os.Exit(2)
	}
}

func usage() error {
	const usageText = "usage: respkit [version|doctor|catalog|explain ERROR_CODE]"
	_, err := fmt.Fprintln(os.Stdout, "respkit v"+resp.Version+"\n"+usageText)
	return err
}
