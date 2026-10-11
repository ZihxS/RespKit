package resolver

// Stable public error codes. Keep each wire value defined here once; changing
// a value breaks clients and the embedded translation catalogs.
const (
	CodeBadRequest            = "BAD_REQUEST"
	CodeUnauthorized          = "UNAUTHORIZED"
	CodeForbidden             = "FORBIDDEN"
	CodeResourceNotFound      = "RESOURCE_NOT_FOUND"
	CodeResourceConflict      = "RESOURCE_CONFLICT"
	CodeResourceAlreadyExists = "RESOURCE_ALREADY_EXISTS"
	CodeValidationFailed      = "VALIDATION_FAILED"
	CodeRateLimited           = "RATE_LIMITED"
	CodeServiceUnavailable    = "SERVICE_UNAVAILABLE"
	CodeInternalError         = "INTERNAL_ERROR"
	CodeInvalidEmail          = "INVALID_EMAIL"
	CodeFieldRequired         = "FIELD_REQUIRED"
	CodeMinLength             = "MIN_LENGTH"
)
