package i18n

import (
	"strings"

	"github.com/ZihxS/RespKit/internal/resolver"
)

func Translate(
	translations map[string]map[string]string,
	defaultLocale string,
	locale string,
	code string,
	params map[string]string,
) string {
	msg := translations[locale][code]
	if msg == "" {
		msg = translations[defaultLocale][code]
	}
	if msg == "" {
		msg = translations[locale][resolver.CodeInternalError]
		if msg == "" {
			msg = translations[defaultLocale][resolver.CodeInternalError]
		}
	}
	for k, v := range params {
		// Parameters come from trusted application code; only known placeholders are substituted.
		msg = strings.ReplaceAll(msg, "{"+k+"}", v)
	}
	return msg
}
