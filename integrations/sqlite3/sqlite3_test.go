package sqlite3mapper

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/ZihxS/RespKit/resp"
	"github.com/mattn/go-sqlite3"
)

func TestMapSQLiteCodes(t *testing.T) {
	tests := []struct {
		name     string
		primary  sqlite3.ErrNo
		extended sqlite3.ErrNoExtended
		code     string
		status   int
	}{
		{
			name:     "unique",
			primary:  sqlite3.ErrConstraint,
			extended: sqlite3.ErrConstraintUnique,
			code:     "RESOURCE_ALREADY_EXISTS",
			status:   http.StatusConflict,
		},
		{
			name:     "primary key",
			primary:  sqlite3.ErrConstraint,
			extended: sqlite3.ErrConstraintPrimaryKey,
			code:     "RESOURCE_ALREADY_EXISTS",
			status:   http.StatusConflict,
		},
		{
			name:     "rowid",
			primary:  sqlite3.ErrConstraint,
			extended: sqlite3.ErrConstraintRowID,
			code:     "RESOURCE_ALREADY_EXISTS",
			status:   http.StatusConflict,
		},
		{
			name:     "not null",
			primary:  sqlite3.ErrConstraint,
			extended: sqlite3.ErrConstraintNotNull,
			code:     "VALIDATION_FAILED",
			status:   http.StatusUnprocessableEntity,
		},
		{
			name:     "foreign key",
			primary:  sqlite3.ErrConstraint,
			extended: sqlite3.ErrConstraintForeignKey,
			code:     "VALIDATION_FAILED",
			status:   http.StatusUnprocessableEntity,
		},
		{
			name:     "check",
			primary:  sqlite3.ErrConstraint,
			extended: sqlite3.ErrConstraintCheck,
			code:     "VALIDATION_FAILED",
			status:   http.StatusUnprocessableEntity,
		},
		{
			name:     "strict datatype",
			primary:  sqlite3.ErrConstraint,
			extended: sqlite3.ErrConstraint.Extend(12),
			code:     "VALIDATION_FAILED",
			status:   http.StatusUnprocessableEntity,
		},
		{
			name:     "mismatch",
			primary:  sqlite3.ErrMismatch,
			extended: 0,
			code:     "BAD_REQUEST",
			status:   http.StatusBadRequest,
		},
		{
			name:     "too big",
			primary:  sqlite3.ErrTooBig,
			extended: 0,
			code:     "VALIDATION_FAILED",
			status:   http.StatusUnprocessableEntity,
		},
		{
			name:     "busy",
			primary:  sqlite3.ErrBusy,
			extended: 0,
			code:     "RESOURCE_CONFLICT",
			status:   http.StatusConflict,
		},
		{
			name:     "busy recovery",
			primary:  sqlite3.ErrBusy,
			extended: sqlite3.ErrBusyRecovery,
			code:     "RESOURCE_CONFLICT",
			status:   http.StatusConflict,
		},
		{
			name:     "busy snapshot",
			primary:  sqlite3.ErrBusy,
			extended: sqlite3.ErrBusySnapshot,
			code:     "RESOURCE_CONFLICT",
			status:   http.StatusConflict,
		},
		{
			name:     "busy timeout",
			primary:  sqlite3.ErrBusy,
			extended: sqlite3.ErrBusy.Extend(3),
			code:     "RESOURCE_CONFLICT",
			status:   http.StatusConflict,
		},
		{
			name:     "locked",
			primary:  sqlite3.ErrLocked,
			extended: 0,
			code:     "RESOURCE_CONFLICT",
			status:   http.StatusConflict,
		},
		{
			name:     "locked shared cache",
			primary:  sqlite3.ErrLocked,
			extended: sqlite3.ErrLockedSharedCache,
			code:     "RESOURCE_CONFLICT",
			status:   http.StatusConflict,
		},
		{
			name:     "memory",
			primary:  sqlite3.ErrNomem,
			extended: 0,
			code:     "SERVICE_UNAVAILABLE",
			status:   http.StatusServiceUnavailable,
		},
		{
			name:     "full",
			primary:  sqlite3.ErrFull,
			extended: 0,
			code:     "SERVICE_UNAVAILABLE",
			status:   http.StatusServiceUnavailable,
		},
		{
			name:     "io",
			primary:  sqlite3.ErrIoErr,
			extended: 0,
			code:     "INTERNAL_ERROR",
			status:   http.StatusInternalServerError,
		},
		{
			name:     "io read",
			primary:  sqlite3.ErrIoErr,
			extended: sqlite3.ErrIoErrRead,
			code:     "INTERNAL_ERROR",
			status:   http.StatusInternalServerError,
		},
		{
			name:     "io write",
			primary:  sqlite3.ErrIoErr,
			extended: sqlite3.ErrIoErrWrite,
			code:     "INTERNAL_ERROR",
			status:   http.StatusInternalServerError,
		},
		{
			name:     "open",
			primary:  sqlite3.ErrCantOpen,
			extended: 0,
			code:     "SERVICE_UNAVAILABLE",
			status:   http.StatusServiceUnavailable,
		},
		{
			name:     "open directory",
			primary:  sqlite3.ErrCantOpen,
			extended: sqlite3.ErrCantOpenIsDir,
			code:     "SERVICE_UNAVAILABLE",
			status:   http.StatusServiceUnavailable,
		},
		{
			name:     "extended only",
			primary:  0,
			extended: sqlite3.ErrConstraintUnique,
			code:     "RESOURCE_ALREADY_EXISTS",
			status:   http.StatusConflict,
		},
		{
			name:     "extended primary",
			primary:  sqlite3.ErrNo(sqlite3.ErrBusySnapshot),
			extended: 0,
			code:     "RESOURCE_CONFLICT",
			status:   http.StatusConflict,
		},
		{
			name:     "generic constraint",
			primary:  sqlite3.ErrConstraint,
			extended: 0,
			code:     "INTERNAL_ERROR",
			status:   http.StatusInternalServerError,
		},
		{
			name:     "commit hook",
			primary:  sqlite3.ErrConstraint,
			extended: sqlite3.ErrConstraintCommitHook,
			code:     "INTERNAL_ERROR",
			status:   http.StatusInternalServerError,
		},
		{
			name:     "function",
			primary:  sqlite3.ErrConstraint,
			extended: sqlite3.ErrConstraintFunction,
			code:     "INTERNAL_ERROR",
			status:   http.StatusInternalServerError,
		},
		{
			name:     "trigger",
			primary:  sqlite3.ErrConstraint,
			extended: sqlite3.ErrConstraintTrigger,
			code:     "INTERNAL_ERROR",
			status:   http.StatusInternalServerError,
		},
		{
			name:     "virtual table",
			primary:  sqlite3.ErrConstraint,
			extended: sqlite3.ErrConstraintVTab,
			code:     "INTERNAL_ERROR",
			status:   http.StatusInternalServerError,
		},
		{
			name:     "pinned",
			primary:  sqlite3.ErrConstraint,
			extended: sqlite3.ErrConstraint.Extend(11),
			code:     "INTERNAL_ERROR",
			status:   http.StatusInternalServerError,
		},
		{
			name:     "future constraint",
			primary:  sqlite3.ErrConstraint,
			extended: sqlite3.ErrConstraint.Extend(99),
			code:     "INTERNAL_ERROR",
			status:   http.StatusInternalServerError,
		},
		{
			name:     "future primary",
			primary:  99,
			extended: 0,
			code:     "INTERNAL_ERROR",
			status:   http.StatusInternalServerError,
		},
		{
			name:     "zero value",
			primary:  0,
			extended: 0,
			code:     "INTERNAL_ERROR",
			status:   http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertSQLiteMapping(t, sqlite3.Error{
				Code:         tt.primary,
				ExtendedCode: tt.extended,
			}, tt.code, tt.status)
		})
	}
}

func assertSQLiteMapping(t *testing.T, value sqlite3.Error, code string, status int) {
	t.Helper()
	for _, tt := range []struct {
		name string
		err  error
	}{
		{
			name: "value",
			err:  value,
		},
		{
			name: "pointer",
			err:  &value,
		},
		{
			name: "wrapped value",
			err:  fmt.Errorf("private query: %w", value),
		},
		{
			name: "wrapped pointer",
			err:  fmt.Errorf("private query: %w", &value),
		},
		{
			name: "joined",
			err:  errors.Join(errors.New("private query"), value),
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			def, ok := Map(tt.err)
			if !ok || def.Code != code || def.Status != status {
				t.Fatalf("Map() = (%+v, %v), want (%s, %d, true)", def, ok, code, status)
			}
			if len(def.Details) != 0 || len(def.Headers) != 0 {
				t.Fatalf("Map() included database details or headers: %+v", def)
			}
		})
	}
}

func TestMapInternalSQLiteCodes(t *testing.T) {
	for _, primary := range []sqlite3.ErrNo{
		sqlite3.ErrError,
		sqlite3.ErrInternal,
		sqlite3.ErrPerm,
		sqlite3.ErrAbort,
		sqlite3.ErrReadonly,
		sqlite3.ErrInterrupt,
		sqlite3.ErrCorrupt,
		sqlite3.ErrNotFound,
		sqlite3.ErrProtocol,
		sqlite3.ErrEmpty,
		sqlite3.ErrSchema,
		sqlite3.ErrMisuse,
		sqlite3.ErrNoLFS,
		sqlite3.ErrAuth,
		sqlite3.ErrFormat,
		sqlite3.ErrRange,
		sqlite3.ErrNotADB,
		sqlite3.ErrNotice,
		sqlite3.ErrWarning,
	} {
		t.Run(fmt.Sprintf("primary_%d", primary), func(t *testing.T) {
			assertSQLiteMapping(t, sqlite3.Error{Code: primary}, "INTERNAL_ERROR", http.StatusInternalServerError)
		})
	}
}

func TestMapNumericDriverErrors(t *testing.T) {
	for _, err := range []error{
		sqlite3.ErrBusy,
		fmt.Errorf("wrapped: %w", sqlite3.ErrBusySnapshot),
	} {
		def, ok := Map(err)
		if !ok || def.Code != "RESOURCE_CONFLICT" || def.Status != http.StatusConflict {
			t.Fatalf("Map() = (%+v, %v), want conflict", def, ok)
		}
	}
}

func TestMapUnknownExtendedCodes(t *testing.T) {
	for _, primary := range []sqlite3.ErrNo{
		sqlite3.ErrConstraint,
		sqlite3.ErrMismatch,
		sqlite3.ErrTooBig,
		sqlite3.ErrBusy,
		sqlite3.ErrLocked,
		sqlite3.ErrNomem,
		sqlite3.ErrFull,
		sqlite3.ErrCantOpen,
		sqlite3.ErrIoErr,
		sqlite3.ErrCorrupt,
		sqlite3.ErrError,
	} {
		t.Run(fmt.Sprintf("primary_%d", primary), func(t *testing.T) {
			extended := primary.Extend(99)
			assertSQLiteMapping(t, sqlite3.Error{
				Code:         primary,
				ExtendedCode: extended,
			}, "INTERNAL_ERROR", http.StatusInternalServerError)
			def, ok := Map(fmt.Errorf("private query: %w", extended))
			if !ok || def.Code != resp.CodeInternalError || def.Status != http.StatusInternalServerError {
				t.Fatalf("unknown numeric error = (%+v, %v), want internal error", def, ok)
			}
		})
	}
}

func TestKnownDatabaseCodesUseLocalizedSafeMessages(t *testing.T) {
	for _, tt := range []struct {
		name     string
		primary  sqlite3.ErrNo
		extended sqlite3.ErrNoExtended
	}{
		{
			name:    "io",
			primary: sqlite3.ErrIoErr,
		},
		{
			name:     "io read",
			primary:  sqlite3.ErrIoErr,
			extended: sqlite3.ErrIoErrRead,
		},
		{
			name:     "io write",
			primary:  sqlite3.ErrIoErr,
			extended: sqlite3.ErrIoErrWrite,
		},
		{
			name:    "corruption",
			primary: sqlite3.ErrCorrupt,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			err := fmt.Errorf("private database path: %w", sqlite3.Error{
				Code:         tt.primary,
				ExtendedCode: tt.extended,
			})
			assertLocalizedResponses(t, err, "INTERNAL_ERROR", http.StatusInternalServerError)
		})
	}
}

func TestMapAbstainsForUnrelatedErrors(t *testing.T) {
	var typedNil *sqlite3.Error
	for _, err := range []error{
		nil,
		typedNil,
		fmt.Errorf("wrapped: %w", typedNil),
		errors.New("UNIQUE constraint failed: private_users.email"),
		sql.ErrNoRows,
		context.Canceled,
		context.DeadlineExceeded,
		resp.ErrNotFound,
		resp.ErrForbidden,
	} {
		if def, ok := Map(err); ok || def.Code != "" || def.Status != 0 {
			t.Fatalf("Map() = (%+v, %v), want no mapping", def, ok)
		}
	}
}

func TestOptionAndResolverPrecedence(t *testing.T) {
	api, err := resp.New(Option())
	if err != nil {
		t.Fatal(err)
	}
	sqliteErr := sqlite3.Error{
		Code:         sqlite3.ErrConstraint,
		ExtendedCode: sqlite3.ErrConstraintUnique,
	}
	def, diag := api.Resolve(fmt.Errorf("insert: %w", sqliteErr))
	if def.Code != "RESOURCE_ALREADY_EXISTS" || def.Status != http.StatusConflict || diag.Mapper != "custom.sqlite3" || !diag.Detected {
		t.Fatalf("SQLite resolution = %+v, %+v", def, diag)
	}
	core, err := resp.New()
	if err != nil {
		t.Fatal(err)
	}
	def, diag = core.Resolve(sqliteErr)
	if def.Code != "INTERNAL_ERROR" || def.Status != http.StatusInternalServerError || diag.Mapper != "fallback" || diag.Detected {
		t.Fatalf("core unexpectedly enabled SQLite: %+v, %+v", def, diag)
	}
	def, diag = api.Resolve(resp.NewError("RESOURCE_NOT_FOUND", http.StatusNotFound, resp.WithCause(sqliteErr)))
	if def.Code != "RESOURCE_NOT_FOUND" || def.Status != http.StatusNotFound || diag.Mapper != "app_error" {
		t.Fatalf("application definition lost precedence: %+v, %+v", def, diag)
	}
	def, diag = api.Resolve(resp.ErrForbidden)
	if def.Code != "FORBIDDEN" || def.Status != http.StatusForbidden || diag.Mapper != "builtin.sentinel" {
		t.Fatalf("builtin fallback changed: %+v, %+v", def, diag)
	}
	if _, err := resp.New(Option(), Option()); err == nil {
		t.Fatal("duplicate SQLite mapper accepted")
	}
}
