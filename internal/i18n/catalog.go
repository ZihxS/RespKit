package i18n

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"path"
	"strings"
)

//go:embed locales
var localeCatalogs embed.FS

// BuiltInTranslations loads the bundled JSON catalog for each supported locale.
func BuiltInTranslations() (map[string]map[string]string, error) {
	return loadTranslations(localeCatalogs)
}

func loadTranslations(catalogs fs.FS) (map[string]map[string]string, error) {
	files, err := fs.Glob(catalogs, "locales/*.json")
	if err != nil {
		return nil, fmt.Errorf("i18n: list embedded locale catalogs: %w", err)
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("i18n: no embedded locale catalogs found")
	}

	translations := make(map[string]map[string]string, len(files))
	for _, filename := range files {
		data, err := fs.ReadFile(catalogs, filename)
		if err != nil {
			return nil, fmt.Errorf("i18n: read embedded catalog %q: %w", filename, err)
		}

		locale := strings.TrimSuffix(path.Base(filename), ".json")
		messages := make(map[string]string)
		if err := json.Unmarshal(data, &messages); err != nil {
			return nil, fmt.Errorf("i18n: decode embedded catalog %q: %w", filename, err)
		}
		if len(messages) == 0 {
			return nil, fmt.Errorf("i18n: embedded catalog %q is empty", filename)
		}
		for code, message := range messages {
			if strings.TrimSpace(code) == "" || strings.TrimSpace(message) == "" {
				return nil, fmt.Errorf("i18n: embedded catalog %q has an empty code or message", filename)
			}
		}
		translations[locale] = messages
	}

	return translations, nil
}
