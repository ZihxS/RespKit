package resp

import "github.com/ZihxS/RespKit/internal/i18n"

func builtInTranslations() (map[string]map[string]string, error) {
	return i18n.BuiltInTranslations()
}

func normalizeLocale(lang string) string {
	return i18n.NormalizeLocale(lang)
}

func (e *Engine) chooseLocale(accept, forced string) string {
	return i18n.ChooseLocale(accept, forced, e.defaultLocale, e.translations)
}

func (e *Engine) translate(locale, code string, params map[string]string) string {
	return i18n.Translate(e.translations, e.defaultLocale, locale, code, params)
}

// hasTranslation is used to ensure Content-Language reflects the language of
// the entire error response. Partial locale catalogs safely fall back to the
// default catalog for the response, rather than mixing languages silently.
func (e *Engine) hasTranslation(locale, code string) bool {
	return e.translations[locale][code] != ""
}

func (e *Engine) effectiveLocale(locale string, failure *Definition) string {
	if failure == nil || locale == e.defaultLocale {
		return locale
	}
	if !e.hasTranslation(locale, failure.Code) && e.hasTranslation(e.defaultLocale, failure.Code) {
		return e.defaultLocale
	}
	for _, d := range failure.Details {
		if !e.hasTranslation(locale, d.Code) && e.hasTranslation(e.defaultLocale, d.Code) {
			return e.defaultLocale
		}
	}
	return locale
}
