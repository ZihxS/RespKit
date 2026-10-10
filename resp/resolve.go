package resp

import (
	"github.com/ZihxS/RespKit/internal/diagnostics"
	"github.com/ZihxS/RespKit/internal/resolver"
)

// Diagnostic contains safe, machine-readable information about how an error was resolved.
type Diagnostic = resolver.Diagnostic

func validCode(code string) bool {
	return resolver.ValidCode(code)
}

func (e *Engine) Resolve(err error) (Definition, Diagnostic) {
	return resolver.Resolve(err, e.mappers)
}

func (e *Engine) Inspect(err error) Diagnostic {
	_, diagnostic := e.Resolve(err)
	return diagnostic
}

func Inspect(err error) Diagnostic {
	return current().Inspect(err)
}

func Explain(err error) string {
	return diagnostics.JSON(Inspect(err))
}
