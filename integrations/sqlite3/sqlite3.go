// Package sqlite3mapper provides opt-in SQLite error mapping for mattn/go-sqlite3.
// It requires cgo and uses RespKit's localized public message catalogs.
package sqlite3mapper

import (
	"errors"
	"net/http"

	"github.com/ZihxS/RespKit/resp"
	"github.com/mattn/go-sqlite3"
)

// Name is the mapper name used in RespKit diagnostics.
const Name = "sqlite3"

// These stable SQLite result codes are not exported by go-sqlite3.
// https://www.sqlite.org/rescode.html
const (
	busyTimeout        sqlite3.ErrNoExtended = 773
	constraintDatatype sqlite3.ErrNoExtended = 3091
)

// Map classifies typed SQLite errors, including wrapped values and pointers.
// Known constraints and contention receive stable public definitions. Database
// permissions, query bugs, and ambiguous or unknown conditions remain internal
// errors. Raw messages, schema names, and OS errors are never copied to a
// definition. The caller retains control of transaction rollback and retries.
func Map(err error) (resp.Definition, bool) {
	sqliteErr, ok := sqliteError(err)
	if !ok {
		return resp.Definition{}, false
	}

	code := sqliteErr.ExtendedCode
	if code == 0 {
		code = sqlite3.ErrNoExtended(sqliteErr.Code)
	}

	// Match the complete result code rather than reducing an unknown extended
	// code to a recognized primary category. Future codes must safely fall back.
	switch code {
	case sqlite3.ErrConstraintUnique,
		sqlite3.ErrConstraintPrimaryKey,
		sqlite3.ErrConstraintRowID:
		return resp.Definition{
			Code:   resp.CodeResourceAlreadyExists,
			Status: http.StatusConflict,
		}, true
	case sqlite3.ErrConstraintNotNull,
		sqlite3.ErrConstraintForeignKey,
		sqlite3.ErrConstraintCheck,
		constraintDatatype:
		return resp.Definition{
			Code:   resp.CodeValidationFailed,
			Status: http.StatusUnprocessableEntity,
		}, true
	case sqlite3.ErrNoExtended(sqlite3.ErrMismatch):
		return resp.Definition{
			Code:   resp.CodeBadRequest,
			Status: http.StatusBadRequest,
		}, true
	case sqlite3.ErrNoExtended(sqlite3.ErrTooBig):
		return resp.Definition{
			Code:   resp.CodeValidationFailed,
			Status: http.StatusUnprocessableEntity,
		}, true
	case sqlite3.ErrNoExtended(sqlite3.ErrBusy),
		sqlite3.ErrBusyRecovery,
		sqlite3.ErrBusySnapshot,
		busyTimeout,
		sqlite3.ErrNoExtended(sqlite3.ErrLocked),
		sqlite3.ErrLockedSharedCache:
		return resp.Definition{
			Code:   resp.CodeResourceConflict,
			Status: http.StatusConflict,
		}, true
	case sqlite3.ErrNoExtended(sqlite3.ErrNomem),
		sqlite3.ErrNoExtended(sqlite3.ErrFull),
		sqlite3.ErrNoExtended(sqlite3.ErrCantOpen),
		sqlite3.ErrCantOpenNoTempDir,
		sqlite3.ErrCantOpenIsDir,
		sqlite3.ErrCantOpenFullPath,
		sqlite3.ErrCantOpenConvPath:
		return resp.Definition{
			Code:   resp.CodeServiceUnavailable,
			Status: http.StatusServiceUnavailable,
		}, true
	}
	// Includes all I/O errors, corruption, SQL/schema failures, and unknown codes.
	return resp.Definition{
		Code:   resp.CodeInternalError,
		Status: http.StatusInternalServerError,
	}, true
}

func sqliteError(err error) (sqlite3.Error, bool) {
	// database/sql returns sqlite3.Error values; application wrappers can also
	// contain pointers or the driver's standalone numeric error types.
	var value sqlite3.Error
	if errors.As(err, &value) {
		return value, true
	}
	var pointer *sqlite3.Error
	if errors.As(err, &pointer) {
		if pointer == nil {
			return sqlite3.Error{}, false
		}
		return *pointer, true
	}
	var primary sqlite3.ErrNo
	if errors.As(err, &primary) {
		return sqlite3.Error{Code: primary}, true
	}
	var extended sqlite3.ErrNoExtended
	if errors.As(err, &extended) {
		return sqlite3.Error{ExtendedCode: extended}, true
	}
	return sqlite3.Error{}, false
}

// Option registers the SQLite mapper on a RespKit engine. It does not open a
// database connection or change RespKit's global configuration.
func Option() resp.Option {
	return resp.WithMapper(Name, Map)
}
