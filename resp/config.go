package resp

import (
	"errors"
	"fmt"
	"strings"
	"sync/atomic"

	"github.com/ZihxS/RespKit/internal/resolver"
)

const Version = "1.0.0"

type Naming string

const (
	SnakeCase  Naming = "snake_case"
	CamelCase  Naming = "camelCase"
	PascalCase Naming = "PascalCase"
)

type ErrorFormat string

const (
	Envelope       ErrorFormat = "envelope"
	ProblemDetails ErrorFormat = "problem_details"
)

// Definition is the safe public representation used by error mappers.
type Definition = resolver.Definition

// Mapper matches an error and optionally provides a public HTTP classification.
type Mapper = resolver.Mapper

type Option func(*Engine) error

type Engine struct {
	naming         Naming
	format         ErrorFormat
	defaultLocale  string
	exposeInternal bool
	translations   map[string]map[string]string
	mappers        []resolver.NamedMapper
}

func New(opts ...Option) (*Engine, error) {
	translations, err := builtInTranslations()
	if err != nil {
		return nil, fmt.Errorf("resp: load built-in translations: %w", err)
	}

	e := &Engine{
		naming:        SnakeCase,
		format:        Envelope,
		defaultLocale: "id",
		translations:  translations,
	}
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if err := opt(e); err != nil {
			return nil, err
		}
	}
	if _, ok := e.translations[e.defaultLocale]; !ok {
		return nil, fmt.Errorf("resp: no translation catalog for default locale %q", e.defaultLocale)
	}
	return e, nil
}

func WithNaming(n Naming) Option {
	return func(e *Engine) error {
		if n != SnakeCase && n != CamelCase && n != PascalCase {
			return fmt.Errorf("resp: invalid naming convention %q", n)
		}
		e.naming = n
		return nil
	}
}

func WithErrorFormat(f ErrorFormat) Option {
	return func(e *Engine) error {
		if f != Envelope && f != ProblemDetails {
			return fmt.Errorf("resp: invalid error format %q", f)
		}
		e.format = f
		return nil
	}
}

// WithExposeInternalErrors includes the original err.Error() in error.debug_message.
// DANGEROUS: use ONLY for local development, never for internet-accessible apps.
// The default is false; no request header, env var, or debug flag enables it.
// Debug messages are never localized and may contain credentials or PII.
func WithExposeInternalErrors(enabled bool) Option {
	return func(e *Engine) error {
		e.exposeInternal = enabled
		return nil
	}
}

func WithDefaultLocale(locale string) Option {
	return func(e *Engine) error {
		locale = normalizeLocale(locale)
		if locale == "" {
			return errors.New("resp: invalid default locale")
		}
		e.defaultLocale = locale
		return nil
	}
}

// WithTranslations adds or replaces safe public messages by locale and error code.
// Example: map[string]map[string]string{"en":{"USER_SUSPENDED":"Account suspended."}}
func WithTranslations(catalog map[string]map[string]string) Option {
	return func(e *Engine) error {
		for locale, values := range catalog {
			locale = normalizeLocale(locale)
			if locale == "" {
				return errors.New("resp: invalid locale in translation catalog")
			}
			if e.translations[locale] == nil {
				e.translations[locale] = map[string]string{}
			}
			for code, message := range values {
				if !validCode(code) || strings.TrimSpace(message) == "" {
					return fmt.Errorf("resp: invalid translation entry %q", code)
				}
				e.translations[locale][code] = message
			}
		}
		return nil
	}
}

func WithMapper(name string, mapper Mapper) Option {
	return func(e *Engine) error {
		if name == "" || mapper == nil {
			return errors.New("resp: mapper name and function required")
		}
		for _, m := range e.mappers {
			if m.Name == name {
				return fmt.Errorf("resp: duplicate mapper %q", name)
			}
		}
		e.mappers = append(e.mappers, resolver.NamedMapper{Name: name, Fn: mapper})
		return nil
	}
}

// Configuring is optional. Configure swaps in an immutable snapshot atomically.
// Call at startup for deterministic behavior; prefer Engine instances for isolation.
var defaultEngine atomic.Pointer[Engine]

func init() {
	e, err := New()
	if err != nil {
		panic(fmt.Errorf("resp: initialize default engine: %w", err))
	}
	defaultEngine.Store(e)
}

func current() *Engine {
	return defaultEngine.Load()
}

func Configure(opts ...Option) error {
	e, err := New(opts...)
	if err != nil {
		return err
	}
	defaultEngine.Store(e)
	return nil
}
