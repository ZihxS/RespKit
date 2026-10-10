package i18n

import (
	"sort"
	"strconv"
	"strings"
)

func NormalizeLocale(lang string) string {
	lang = strings.ToLower(strings.TrimSpace(strings.ReplaceAll(lang, "_", "-")))
	if lang == "" || len(lang) > 64 {
		return ""
	}
	for _, r := range lang {
		if (r < 'a' || r > 'z') && (r < '0' || r > '9') && r != '-' {
			return ""
		}
	}
	return strings.SplitN(lang, "-", 2)[0]
}

func ChooseLocale(acceptLanguage, forced, defaultLocale string, translations map[string]map[string]string) string {
	if forced != "" {
		l := NormalizeLocale(forced)
		if _, exists := translations[l]; exists {
			return l
		}
		return defaultLocale
	}
	type pref struct {
		lang   string
		weight float64
		order  int
	}
	var prefs []pref
	for i, part := range strings.Split(acceptLanguage, ",") {
		if i >= 32 {
			break
		}
		fields := strings.Split(strings.TrimSpace(part), ";")
		language := NormalizeLocale(fields[0])
		if strings.TrimSpace(fields[0]) == "*" {
			language = "*"
		}
		if language == "" {
			continue
		}
		q := 1.0
		for _, attr := range fields[1:] {
			pair := strings.SplitN(strings.TrimSpace(attr), "=", 2)
			if len(pair) == 2 && strings.EqualFold(strings.TrimSpace(pair[0]), "q") {
				n, err := strconv.ParseFloat(strings.TrimSpace(pair[1]), 64)
				if err != nil || n < 0 || n > 1 {
					q = 0
				} else {
					q = n
				}
				break
			}
		}
		if q > 0 {
			prefs = append(prefs, pref{language, q, i})
		}
	}
	sort.SliceStable(prefs, func(i, j int) bool { return prefs[i].weight > prefs[j].weight })
	for _, p := range prefs {
		if p.lang == "*" {
			return defaultLocale
		}
		if _, ok := translations[p.lang]; ok {
			return p.lang
		}
	}
	return defaultLocale
}
