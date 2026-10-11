package sqlite3mapper

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ZihxS/RespKit/resp"
	"github.com/mattn/go-sqlite3"
)

func newDatabase(t *testing.T, dsn string) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		t.Fatal(err)
	}
	// In-memory databases belong to a connection, and lock tests explicitly
	// open a second pool. Keep each pool on one connection for determinism.
	db.SetMaxOpenConns(1)
	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Error(err)
		}
	})
	return db
}

func openDatabase(t *testing.T, dsn string) *sql.DB {
	t.Helper()
	db := newDatabase(t, dsn)
	if err := db.Ping(); err != nil {
		t.Fatal(err)
	}
	return db
}

func beginTransaction(t *testing.T, db *sql.DB) *sql.Tx {
	t.Helper()
	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
			t.Error(err)
		}
	})
	return tx
}

func execSQL(t *testing.T, db *sql.DB, query string) {
	t.Helper()
	if _, err := db.Exec(query); err != nil {
		t.Fatal(err)
	}
}

func requireSQLiteError(t *testing.T, err error, primary sqlite3.ErrNo, extended sqlite3.ErrNoExtended) error {
	t.Helper()
	var source sqlite3.Error
	if !errors.As(err, &source) || source.Code != primary || source.ExtendedCode != extended {
		t.Fatalf("driver error = %v (%+v), want SQLite codes %d/%d", err, source, primary, extended)
	}
	return fmt.Errorf("private operation: %w", err)
}

var publicMessages = map[string]map[string]string{
	"id": {
		"RESOURCE_ALREADY_EXISTS": "Data sudah terdaftar.",
		"VALIDATION_FAILED":       "Data yang dikirim tidak valid.",
		"BAD_REQUEST":             "Permintaan tidak valid.",
		"RESOURCE_CONFLICT":       "Terjadi konflik dengan data yang ada.",
		"SERVICE_UNAVAILABLE":     "Layanan sementara tidak tersedia.",
		"INTERNAL_ERROR":          "Terjadi kesalahan pada server.",
	},
	"en": {
		"RESOURCE_ALREADY_EXISTS": "The resource already exists.",
		"VALIDATION_FAILED":       "The submitted data is invalid.",
		"BAD_REQUEST":             "The request is invalid.",
		"RESOURCE_CONFLICT":       "The request conflicts with the current resource state.",
		"SERVICE_UNAVAILABLE":     "The service is temporarily unavailable.",
		"INTERNAL_ERROR":          "An internal server error occurred.",
	},
}

func assertResponse(t *testing.T, api *resp.Engine, err error, code string, status int, accept, locale, message string, opts ...resp.ResponseOption) {
	t.Helper()
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/accounts", nil)
	r.Header.Set("Accept-Language", accept)
	if writeErr := api.Fail(resp.HTTP(w, r), err, opts...); writeErr != nil {
		t.Fatal(writeErr)
	}
	if w.Code != status || w.Header().Get("Content-Language") != locale || w.Header().Get("Vary") != "Accept-Language" {
		t.Fatalf("response = %d, %v, want status %d and locale %s", w.Code, w.Header(), status, locale)
	}
	var body struct {
		Success bool `json:"success"`
		Error   struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
		Code   string `json:"code"`
		Status int    `json:"status"`
		Detail string `json:"detail"`
		Type   string `json:"type"`
	}
	if decodeErr := json.Unmarshal(w.Body.Bytes(), &body); decodeErr != nil {
		t.Fatal(decodeErr)
	}
	if w.Header().Get("Content-Type") == "application/problem+json" {
		if body.Code != code || body.Status != status || body.Detail != message || body.Type != "about:blank" {
			t.Fatalf("problem response = %s, want %s/%d/%q", w.Body.String(), code, status, message)
		}
	} else if w.Header().Get("Content-Type") != "application/json; charset=utf-8" || body.Success || body.Error.Code != code || body.Error.Message != message {
		t.Fatalf("envelope response = %s, want %s/%q", w.Body.String(), code, message)
	}
	for _, private := range []string{
		err.Error(),
		"private_",
		"debug_message",
		"SystemErrno",
		"ExtendedCode",
	} {
		if strings.Contains(w.Body.String(), private) {
			t.Fatalf("response disclosed %q: %s", private, w.Body.String())
		}
	}
}

func assertLocalizedResponses(t *testing.T, err error, code string, status int) {
	t.Helper()
	for _, format := range []resp.ErrorFormat{
		resp.Envelope,
		resp.ProblemDetails,
	} {
		t.Run(string(format), func(t *testing.T) {
			api, newErr := resp.New(Option(), resp.WithErrorFormat(format))
			if newErr != nil {
				t.Fatal(newErr)
			}
			for _, language := range []struct {
				name   string
				accept string
				locale string
			}{
				{
					name:   "default",
					locale: "id",
				},
				{
					name:   "Indonesian",
					accept: "id-ID",
					locale: "id",
				},
				{
					name:   "English",
					accept: "en-US",
					locale: "en",
				},
				{
					name:   "weighted",
					accept: "id;q=0.2, en;q=0.9",
					locale: "en",
				},
				{
					name:   "unsupported",
					accept: "fr-FR",
					locale: "id",
				},
				{
					name:   "excluded",
					accept: "en;q=0",
					locale: "id",
				},
			} {
				t.Run(language.name, func(t *testing.T) {
					assertResponse(t, api, err, code, status, language.accept, language.locale, publicMessages[language.locale][code])
				})
			}
			t.Run("locale override", func(t *testing.T) {
				assertResponse(t, api, err, code, status, "en", "id", publicMessages["id"][code], resp.Locale("id"))
			})
		})
	}
}

func TestRealSQLiteConstraintsAndFallback(t *testing.T) {
	db := openDatabase(t, "file:constraints?mode=memory&cache=private&_foreign_keys=on")
	execSQL(t, db, `
		CREATE TABLE private_accounts (id INTEGER PRIMARY KEY, email TEXT NOT NULL UNIQUE, age INTEGER CHECK (age >= 18));
		CREATE TABLE private_links (account_id INTEGER REFERENCES private_accounts(id));
		CREATE TABLE private_rowids (name TEXT);
		CREATE TABLE private_strict (value INTEGER) STRICT;
		CREATE TABLE private_blocked (value TEXT);
		CREATE TRIGGER private_guard BEFORE INSERT ON private_blocked BEGIN SELECT RAISE(ABORT, 'private_secret'); END;
		INSERT INTO private_accounts VALUES (1, 'private@example.test', 20);
		INSERT INTO private_rowids (rowid, name) VALUES (1, 'private');
	`)
	for _, tt := range []struct {
		name     string
		query    string
		primary  sqlite3.ErrNo
		extended sqlite3.ErrNoExtended
		code     string
		status   int
	}{
		{
			name:     "unique",
			query:    `INSERT INTO private_accounts VALUES (2, 'private@example.test', 20)`,
			primary:  sqlite3.ErrConstraint,
			extended: sqlite3.ErrConstraintUnique,
			code:     "RESOURCE_ALREADY_EXISTS",
			status:   http.StatusConflict,
		},
		{
			name:     "primary key",
			query:    `INSERT INTO private_accounts VALUES (1, 'other@example.test', 20)`,
			primary:  sqlite3.ErrConstraint,
			extended: sqlite3.ErrConstraintPrimaryKey,
			code:     "RESOURCE_ALREADY_EXISTS",
			status:   http.StatusConflict,
		},
		{
			name:     "rowid",
			query:    `INSERT INTO private_rowids (rowid, name) VALUES (1, 'other')`,
			primary:  sqlite3.ErrConstraint,
			extended: sqlite3.ErrConstraintRowID,
			code:     "RESOURCE_ALREADY_EXISTS",
			status:   http.StatusConflict,
		},
		{
			name:     "not null",
			query:    `INSERT INTO private_accounts VALUES (2, NULL, 20)`,
			primary:  sqlite3.ErrConstraint,
			extended: sqlite3.ErrConstraintNotNull,
			code:     "VALIDATION_FAILED",
			status:   http.StatusUnprocessableEntity,
		},
		{
			name:     "foreign key",
			query:    `INSERT INTO private_links VALUES (999)`,
			primary:  sqlite3.ErrConstraint,
			extended: sqlite3.ErrConstraintForeignKey,
			code:     "VALIDATION_FAILED",
			status:   http.StatusUnprocessableEntity,
		},
		{
			name:     "check",
			query:    `INSERT INTO private_accounts VALUES (2, 'other@example.test', 10)`,
			primary:  sqlite3.ErrConstraint,
			extended: sqlite3.ErrConstraintCheck,
			code:     "VALIDATION_FAILED",
			status:   http.StatusUnprocessableEntity,
		},
		{
			name:     "strict datatype",
			query:    `INSERT INTO private_strict VALUES ('private_text')`,
			primary:  sqlite3.ErrConstraint,
			extended: sqlite3.ErrConstraint.Extend(12),
			code:     "VALIDATION_FAILED",
			status:   http.StatusUnprocessableEntity,
		},
		{
			name:     "mismatch",
			query:    `INSERT INTO private_accounts VALUES ('private_text', 'other@example.test', 20)`,
			primary:  sqlite3.ErrMismatch,
			extended: sqlite3.ErrNoExtended(sqlite3.ErrMismatch),
			code:     "BAD_REQUEST",
			status:   http.StatusBadRequest,
		},
		{
			name:     "syntax",
			query:    `INSERT INTO private_missing VALUES (1)`,
			primary:  sqlite3.ErrError,
			extended: sqlite3.ErrNoExtended(sqlite3.ErrError),
			code:     "INTERNAL_ERROR",
			status:   http.StatusInternalServerError,
		},
		{
			name:     "ambiguous trigger",
			query:    `INSERT INTO private_blocked VALUES ('value')`,
			primary:  sqlite3.ErrConstraint,
			extended: sqlite3.ErrConstraintTrigger,
			code:     "INTERNAL_ERROR",
			status:   http.StatusInternalServerError,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			_, driverErr := db.Exec(tt.query)
			err := requireSQLiteError(t, driverErr, tt.primary, tt.extended)
			assertLocalizedResponses(t, err, tt.code, tt.status)
		})
	}
	// SQLite's ErrNotFound is an internal opcode error, while sql.ErrNoRows is
	// domain-specific. The adapter deliberately does not turn either into 404.
	t.Run("no rows", func(t *testing.T) {
		var id int
		err := db.QueryRow(`SELECT id FROM private_accounts WHERE id = 999`).Scan(&id)
		if !errors.Is(err, sql.ErrNoRows) {
			t.Fatalf("QueryRow() = %v, want sql.ErrNoRows", err)
		}
		assertLocalizedResponses(t, err, "INTERNAL_ERROR", http.StatusInternalServerError)
	})
}

func TestRealSQLiteContention(t *testing.T) {
	t.Run("busy", func(t *testing.T) {
		dsn := (&url.URL{
			Scheme:   "file",
			Path:     filepath.ToSlash(filepath.Join(t.TempDir(), "private.db")),
			RawQuery: "_busy_timeout=0",
		}).String()
		first := openDatabase(t, dsn)
		execSQL(t, first, `CREATE TABLE private_items (id INTEGER)`)
		second := openDatabase(t, dsn)
		tx := beginTransaction(t, first)
		if _, err := tx.Exec(`INSERT INTO private_items VALUES (1)`); err != nil {
			t.Fatal(err)
		}
		_, driverErr := second.Exec(`INSERT INTO private_items VALUES (2)`)
		err := requireSQLiteError(t, driverErr, sqlite3.ErrBusy, sqlite3.ErrNoExtended(sqlite3.ErrBusy))
		assertLocalizedResponses(t, err, "RESOURCE_CONFLICT", http.StatusConflict)
	})
	t.Run("locked shared cache", func(t *testing.T) {
		dsn := "file:locked?mode=memory&cache=shared&_busy_timeout=0"
		first := openDatabase(t, dsn)
		execSQL(t, first, `CREATE TABLE private_items (id INTEGER)`)
		second := openDatabase(t, dsn)
		tx := beginTransaction(t, first)
		if _, err := tx.Exec(`INSERT INTO private_items VALUES (1)`); err != nil {
			t.Fatal(err)
		}
		_, driverErr := second.Exec(`INSERT INTO private_items VALUES (2)`)
		err := requireSQLiteError(t, driverErr, sqlite3.ErrLocked, sqlite3.ErrLockedSharedCache)
		assertLocalizedResponses(t, err, "RESOURCE_CONFLICT", http.StatusConflict)
	})
}

func TestRealSQLiteDatabaseFailures(t *testing.T) {
	t.Run("read only", func(t *testing.T) {
		db := openDatabase(t, ":memory:")
		execSQL(t, db, `CREATE TABLE private_items (id INTEGER); PRAGMA query_only = ON`)
		_, driverErr := db.Exec(`INSERT INTO private_items VALUES (1)`)
		err := requireSQLiteError(t, driverErr, sqlite3.ErrReadonly, sqlite3.ErrNoExtended(sqlite3.ErrReadonly))
		assertLocalizedResponses(t, err, "INTERNAL_ERROR", http.StatusInternalServerError)
	})
	t.Run("cannot open", func(t *testing.T) {
		db := newDatabase(t, filepath.Join(t.TempDir(), "private_missing", "private.db"))
		err := requireSQLiteError(t, db.Ping(), sqlite3.ErrCantOpen, sqlite3.ErrNoExtended(sqlite3.ErrCantOpen))
		assertLocalizedResponses(t, err, "SERVICE_UNAVAILABLE", http.StatusServiceUnavailable)
	})
	t.Run("not a database", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "private.db")
		if err := os.WriteFile(path, []byte(strings.Repeat("private_secret", 100)), 0600); err != nil {
			t.Fatal(err)
		}
		db := newDatabase(t, path)
		err := requireSQLiteError(t, db.Ping(), sqlite3.ErrNotADB, sqlite3.ErrNoExtended(sqlite3.ErrNotADB))
		assertLocalizedResponses(t, err, "INTERNAL_ERROR", http.StatusInternalServerError)
	})
	t.Run("full", func(t *testing.T) {
		db := openDatabase(t, ":memory:")
		execSQL(t, db, `CREATE TABLE private_blobs (value BLOB); PRAGMA max_page_count = 2`)
		_, driverErr := db.Exec(`INSERT INTO private_blobs VALUES (zeroblob(1048576))`)
		err := requireSQLiteError(t, driverErr, sqlite3.ErrFull, sqlite3.ErrNoExtended(sqlite3.ErrFull))
		assertLocalizedResponses(t, err, "SERVICE_UNAVAILABLE", http.StatusServiceUnavailable)
	})
	t.Run("too big", func(t *testing.T) {
		db := openDatabase(t, ":memory:")
		execSQL(t, db, `CREATE TABLE private_blobs (value BLOB)`)
		conn, err := db.Conn(context.Background())
		if err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() {
			if err := conn.Close(); err != nil {
				t.Error(err)
			}
		})
		if err := conn.Raw(func(driverConn any) error {
			driverConn.(*sqlite3.SQLiteConn).SetLimit(sqlite3.SQLITE_LIMIT_LENGTH, 1024)
			return nil
		}); err != nil {
			t.Fatal(err)
		}
		_, driverErr := conn.ExecContext(context.Background(), `INSERT INTO private_blobs VALUES (zeroblob(2048))`)
		err = requireSQLiteError(t, driverErr, sqlite3.ErrTooBig, sqlite3.ErrNoExtended(sqlite3.ErrTooBig))
		assertLocalizedResponses(t, err, "VALIDATION_FAILED", http.StatusUnprocessableEntity)
	})
}

func TestCustomLocalesAndCancellation(t *testing.T) {
	db := openDatabase(t, ":memory:")
	execSQL(t, db, `CREATE TABLE private_items (id INTEGER PRIMARY KEY); INSERT INTO private_items VALUES (1)`)
	_, driverErr := db.Exec(`INSERT INTO private_items VALUES (1)`)
	err := requireSQLiteError(t, driverErr, sqlite3.ErrConstraint, sqlite3.ErrConstraintPrimaryKey)
	for _, format := range []resp.ErrorFormat{
		resp.Envelope,
		resp.ProblemDetails,
	} {
		t.Run(string(format), func(t *testing.T) {
			api, newErr := resp.New(
				Option(),
				resp.WithErrorFormat(format),
				resp.WithDefaultLocale("en"),
				resp.WithTranslations(map[string]map[string]string{
					"fr": {
						"RESOURCE_ALREADY_EXISTS": "La ressource existe déjà.",
					},
					"es": {
						"INTERNAL_ERROR": "Error interno del servidor.",
					},
				}),
			)
			if newErr != nil {
				t.Fatal(newErr)
			}
			assertResponse(t, api, err, "RESOURCE_ALREADY_EXISTS", http.StatusConflict, "", "en", publicMessages["en"]["RESOURCE_ALREADY_EXISTS"])
			assertResponse(t, api, err, "RESOURCE_ALREADY_EXISTS", http.StatusConflict, "fr-FR", "fr", "La ressource existe déjà.")
			assertResponse(t, api, err, "RESOURCE_ALREADY_EXISTS", http.StatusConflict, "es", "en", publicMessages["en"]["RESOURCE_ALREADY_EXISTS"])
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/accounts", nil)
			if writeErr := api.Fail(resp.HTTP(w, r), errors.Join(context.Canceled, err)); !errors.Is(writeErr, resp.ErrRequestCanceled) || w.Body.Len() != 0 || len(w.Header()) != 0 {
				t.Fatalf("canceled request wrote a response: %v, %s, %v", writeErr, w.Body.String(), w.Header())
			}
		})
	}
}
