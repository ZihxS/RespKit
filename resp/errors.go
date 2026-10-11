// Package resp provides zero-configuration, framework-adaptive HTTP responses,
// safe error mapping, localization, and structured diagnostics.
package resp

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/ZihxS/RespKit/internal/resolver"
	"github.com/ZihxS/RespKit/internal/transport"
)

// Sentinels intentionally represent application/transport semantics, not raw DB errors.
var (
	ErrBadRequest         = resolver.ErrBadRequest
	ErrUnauthorized       = resolver.ErrUnauthorized
	ErrForbidden          = resolver.ErrForbidden
	ErrNotFound           = resolver.ErrNotFound
	ErrConflict           = resolver.ErrConflict
	ErrValidation         = resolver.ErrValidation
	ErrRateLimited        = resolver.ErrRateLimited
	ErrUnavailable        = resolver.ErrUnavailable
	ErrInternal           = resolver.ErrInternal
	ErrRequestCanceled    = errors.New("request canceled; response not written")
	ErrUnsupportedContext = transport.ErrUnsupportedContext
)

// FieldError represents a trusted application-defined validation failure.
// Field is a JSON path; Code is stable and locale-independent. Params may
// contain only non-sensitive, developer-approved substitution values.
type FieldError = resolver.FieldError

// AppError defines an explicit public error. Cause is NEVER serialized.
// Status must be in the 400..599 range; unknown codes fall back to a safe message.
type AppError struct {
	Code    string
	Status  int
	Cause   error
	Details []FieldError
	Headers http.Header
}

func (e *AppError) Error() string {
	if e == nil {
		return "internal error"
	}
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Code, e.Cause)
	}
	return e.Code
}

// RespKitDefinition provides the safe, explicit public classification.
// The wrapped cause stays internal and is never serialized by the resolver.
func (e *AppError) RespKitDefinition() resolver.Definition {
	if e == nil {
		return resolver.Definition{}
	}
	return resolver.Definition{
		Code:    e.Code,
		Status:  e.Status,
		Details: e.Details,
		Headers: e.Headers,
	}
}

func (e *AppError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Cause
}

type ErrorOption func(*AppError)

func WithCause(err error) ErrorOption {
	return func(e *AppError) {
		e.Cause = err
	}
}

func WithDetails(details ...FieldError) ErrorOption {
	return func(e *AppError) {
		e.Details = append([]FieldError(nil), details...)
	}
}

// WithHeader allows the caller to set a response header (e.g. WWW-Authenticate).
func WithHeader(key, value string) ErrorOption {
	return func(e *AppError) {
		if e.Headers == nil {
			e.Headers = make(http.Header)
		}
		e.Headers.Set(key, value)
	}
}

func NewError(code string, status int, options ...ErrorOption) *AppError {
	e := &AppError{Code: code, Status: status}
	for _, opt := range options {
		if opt != nil {
			opt(e)
		}
	}
	return e
}
