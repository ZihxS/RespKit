package resp

import (
	"sort"

	"github.com/ZihxS/RespKit/internal/resolver"
)

// Catalog returns deterministic machine-readable built-in mappings without
// requiring an HTTP framework or a running server.
type CatalogItem struct {
	Code       string            `json:"code"`
	HTTPStatus int               `json:"http_status"`
	Messages   map[string]string `json:"messages"`
}

func (e *Engine) Catalog() []CatalogItem {
	seen := resolver.Builtins()
	codes := make([]string, 0, len(seen))
	for c := range seen {
		codes = append(codes, c)
	}
	sort.Strings(codes)
	result := make([]CatalogItem, 0, len(codes))
	for _, c := range codes {
		messages := map[string]string{
			"id": e.translate("id", c, nil),
			"en": e.translate("en", c, nil),
		}
		result = append(result, CatalogItem{
			Code:       c,
			HTTPStatus: seen[c],
			Messages:   messages,
		})
	}
	return result
}

func Catalog() []CatalogItem {
	return current().Catalog()
}
