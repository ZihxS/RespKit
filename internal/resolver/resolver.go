package resolver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
)

var codePattern = regexp.MustCompile(`^[A-Z][A-Z0-9_]{1,63}$`)

func ValidCode(c string) bool {
	return codePattern.MatchString(c)
}

type FieldError struct {
	Field  string            `json:"field"`
	Code   string            `json:"code"`
	Params map[string]string `json:"-"`
}

type Definition struct {
	Code    string
	Status  int
	Details []FieldError
	Headers http.Header
}

type Mapper func(error) (Definition, bool)

type NamedMapper struct {
	Name string
	Fn   Mapper
}

type Diagnostic struct {
	SchemaVersion string `json:"schema_version"`
	Detected      bool   `json:"detected"`
	Mapper        string `json:"mapper"`
	ErrorCode     string `json:"error_code"`
	HTTPStatus    int    `json:"http_status"`
	SourceCode    string `json:"source_code,omitempty"`
	Note          string `json:"note,omitempty"`
}

func (d Diagnostic) String() string {
	return fmt.Sprintf("%s (%d) via %s", d.ErrorCode, d.HTTPStatus, d.Mapper)
}

var (
	ErrBadRequest   = errors.New("bad request")
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
	ErrNotFound     = errors.New("not found")
	ErrConflict     = errors.New("conflict")
	ErrValidation   = errors.New("validation failed")
	ErrRateLimited  = errors.New("rate limited")
	ErrUnavailable  = errors.New("service unavailable")
	ErrInternal     = errors.New("internal error")
)

var builtins = []struct {
	sentinel error
	code     string
	status   int
}{
	{ErrBadRequest, "BAD_REQUEST", 400},
	{ErrUnauthorized, "UNAUTHORIZED", 401},
	{ErrForbidden, "FORBIDDEN", 403},
	{ErrNotFound, "RESOURCE_NOT_FOUND", 404},
	{ErrConflict, "RESOURCE_CONFLICT", 409},
	{ErrValidation, "VALIDATION_FAILED", 422},
	{ErrRateLimited, "RATE_LIMITED", 429},
	{ErrUnavailable, "SERVICE_UNAVAILABLE", 503},
	{ErrInternal, "INTERNAL_ERROR", 500},
}

func Resolve(err error, mappers []NamedMapper) (Definition, Diagnostic) {
	if err == nil {
		err = ErrInternal
	}
	def := Definition{Code: "INTERNAL_ERROR", Status: http.StatusInternalServerError}
	diagnostic := Diagnostic{SchemaVersion: "1", Mapper: "fallback", ErrorCode: def.Code, HTTPStatus: def.Status}
	// Record SQLSTATE before custom mappers run so a database adapter can provide
	// a safe HTTP definition without dropping the original diagnostic code.
	var sqlState interface{ SQLState() string }
	if errors.As(err, &sqlState) {
		diagnostic.SourceCode = sqlState.SQLState()
	}
	finish := func(d Definition, source string, detected bool) (Definition, Diagnostic) {
		if d.Status < 400 || d.Status > 599 || !ValidCode(d.Code) {
			diagnostic.Note = "invalid mapper definition; safe fallback applied"
			return def, diagnostic
		}
		diagnostic.Detected = detected
		diagnostic.Mapper = source
		diagnostic.ErrorCode = d.Code
		diagnostic.HTTPStatus = d.Status
		return d, diagnostic
	}
	var app interface{ RespKitDefinition() Definition }
	if errors.As(err, &app) && app != nil {
		return finish(app.RespKitDefinition(), "app_error", true)
	}
	for _, mapper := range mappers {
		if mapped, ok := mapper.Fn(err); ok {
			return finish(mapped, "custom."+mapper.Name, true)
		}
	}
	for _, b := range builtins {
		if errors.Is(err, b.sentinel) {
			d := Definition{Code: b.code, Status: b.status}
			if b.status == 401 {
				// Auth scheme is application-specific. The app should supply a
				// challenge via AppError.WithHeader / custom mapper.
				diagnostic.Note = "401 requires a WWW-Authenticate challenge configured by the application"
			}
			return finish(d, "builtin.sentinel", true)
		}
	}
	var syntax *json.SyntaxError
	var jsonType *json.UnmarshalTypeError
	if errors.As(err, &syntax) || errors.As(err, &jsonType) {
		return finish(Definition{Code: "BAD_REQUEST", Status: 400}, "builtin.json", true)
	}
	// pgx/v4 and pgx/v5 implement SQLState(). Inspect via structural interface
	// without pulling PostgreSQL drivers into the core module.
	if errors.As(err, &sqlState) {
		code := sqlState.SQLState()
		if code == "23505" {
			return finish(Definition{Code: "RESOURCE_ALREADY_EXISTS", Status: 409}, "builtin.sqlstate_unique", true)
		}
		// Foreign key violations and others have ambiguous REST semantics.
		diagnostic.Note = "SQLSTATE recognized but requires domain-specific mapping"
		return def, diagnostic
	}
	if errors.Is(err, context.Canceled) {
		diagnostic.Note = "request canceled; caller should not write a response"
		return def, diagnostic
	}
	if errors.Is(err, context.DeadlineExceeded) {
		diagnostic.Note = "deadline exceeded; HTTP 504 is only appropriate for gateway timeouts"
	}
	return def, diagnostic
}

// Builtins returns the stable, documented public error/status catalog.
func Builtins() map[string]int {
	out := make(map[string]int, len(builtins)+1)
	for _, b := range builtins {
		out[b.code] = b.status
	}
	out["RESOURCE_ALREADY_EXISTS"] = 409
	return out
}
