package resp

import "github.com/ZihxS/RespKit/internal/resolver"

// Built-in public codes are untyped string constants so existing Definition,
// NewError, FieldError, and translation APIs remain source-compatible.
const (
	// CodeBadRequest identifies an invalid request (HTTP 400).
	CodeBadRequest = resolver.CodeBadRequest
	// CodeUnauthorized identifies a request requiring authentication (HTTP 401).
	CodeUnauthorized = resolver.CodeUnauthorized
	// CodeForbidden identifies a request without the required permission (HTTP 403).
	CodeForbidden = resolver.CodeForbidden
	// CodeResourceNotFound identifies a missing application resource (HTTP 404).
	CodeResourceNotFound = resolver.CodeResourceNotFound
	// CodeResourceConflict identifies a conflict with resource state (HTTP 409).
	CodeResourceConflict = resolver.CodeResourceConflict
	// CodeResourceAlreadyExists identifies a duplicate resource (HTTP 409).
	CodeResourceAlreadyExists = resolver.CodeResourceAlreadyExists
	// CodeValidationFailed identifies invalid submitted data (HTTP 422).
	CodeValidationFailed = resolver.CodeValidationFailed
	// CodeRateLimited identifies an exceeded request rate limit (HTTP 429).
	CodeRateLimited = resolver.CodeRateLimited
	// CodeServiceUnavailable identifies an unavailable service (HTTP 503).
	CodeServiceUnavailable = resolver.CodeServiceUnavailable
	// CodeInternalError identifies a safely redacted internal failure (HTTP 500).
	CodeInternalError = resolver.CodeInternalError
	// CodeInvalidEmail identifies an invalid email field in validation details.
	CodeInvalidEmail = resolver.CodeInvalidEmail
	// CodeFieldRequired identifies a required field in validation details.
	CodeFieldRequired = resolver.CodeFieldRequired
	// CodeMinLength identifies a field shorter than its minimum length.
	CodeMinLength = resolver.CodeMinLength
)
