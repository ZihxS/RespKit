package resp

import (
	"errors"
	"testing"

	"github.com/ZihxS/RespKit/internal/diagnostics"
)

// Regression coverage for isolated internal components after modularization.
func TestNilDefinitionAndDiagnosticEncodingFailure(t *testing.T) {
	var e *AppError
	def := e.RespKitDefinition()
	if def.Code != "" || def.Status != 0 {
		t.Fatalf("nil error definition: %+v", def)
	}
	if got := diagnostics.JSON(make(chan int)); got != "{}" {
		t.Fatalf("unexpected failure fallback: %q", got)
	}
	if got := diagnostics.JSON(map[string]string{"code": "OK"}); got == "{}" {
		t.Fatal("diagnostic serialization unexpectedly failed")
	}
	if got := Inspect(errors.New("secret")); got.ErrorCode != "INTERNAL_ERROR" {
		t.Fatal(got)
	}
}
