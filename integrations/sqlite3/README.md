# RespKit SQLite adapter

This optional module maps typed errors from [mattn/go-sqlite3](https://github.com/mattn/go-sqlite3) to safe RespKit error definitions. It uses the core Indonesian/English catalogs and supports application-provided languages. It does not open database connections or change global RespKit configuration.

## Requirements and installation

The adapter supports Go 1.22 and pins `go-sqlite3 v1.14.52`. The [driver requires cgo and a C compiler](https://github.com/mattn/go-sqlite3#installation): build with `CGO_ENABLED=1`. Cross-compiling requires a C compiler for the target platform. The root RespKit module remains standard-library-only and builds without cgo.

Install the adapter alongside core RespKit v1.1.0 or later:

```bash
go get github.com/ZihxS/RespKit/resp@v1.1.0
go get github.com/ZihxS/RespKit/integrations/sqlite3@v1.0.0
```

For local adapter development, use a `replace github.com/ZihxS/RespKit/integrations/sqlite3 => /path/to/RespKit/integrations/sqlite3` in your application's `go.mod`. Nested modules use their own release tags: [integrations/sqlite3/v1.0.0](https://github.com/ZihxS/RespKit/releases/tag/integrations/sqlite3/v1.0.0) versions this adapter, while [v1.1.0](https://github.com/ZihxS/RespKit/releases/tag/v1.1.0) versions the core.

## Choose a SQLite driver

Use this adapter when your application already uses `mattn/go-sqlite3` or can build with cgo and a C compiler. It handles that driver's error types and does not require another SQLite driver.

| Driver | Build requirements | RespKit integration |
| --- | --- | --- |
| [mattn/go-sqlite3](https://github.com/mattn/go-sqlite3) | cgo and a C compiler; target C toolchain required for cross-compilation | `integrations/sqlite3`, provided here |
| [modernc.org/sqlite](https://pkg.go.dev/modernc.org/sqlite) | Pure Go, no cgo; check the selected driver's Go/platform requirements | Separate planned adapter tracked in [issue #7](https://github.com/ZihxS/RespKit/issues/7) |

Choose modernc when a cgo-free build is a requirement. Its typed errors differ from mattn's, so `sqlite3mapper.Option()` does not map modernc errors. The modernc adapter is a separate task and has not shipped in this workspace. Install the integration for the driver you actually use; neither driver belongs in RespKit's core module.

## Configure and handle errors

Create an engine once at startup and reuse it in handlers:

```go
import (
	"net/http"

	sqlite3mapper "github.com/ZihxS/RespKit/integrations/sqlite3"
	"github.com/ZihxS/RespKit/resp"
)

func newResponseEngine() (*resp.Engine, error) {
	return resp.New(sqlite3mapper.Option())
}

func writeDatabaseError(api *resp.Engine, w http.ResponseWriter, r *http.Request, dbErr error) error {
	return api.Fail(resp.HTTP(w, r), dbErr)
}
```

Pass the error returned by `db.ExecContext`, `QueryRowContext().Scan`, or transaction operations to `api.Fail`. Wrapped errors (`fmt.Errorf("insert: %w", err)`) and joined errors are supported via `errors.As`. The mapper accepts `sqlite3.Error` values and pointers, and standalone `sqlite3.ErrNo` / `sqlite3.ErrNoExtended` errors. It never classifies ordinary errors by their text.

For Gin, Echo, and Fiber pass the native context to `api.Fail(ctx, dbErr)`. For `net/http`, Chi, Gorilla/mux, and Beego use `resp.HTTP(w, r)`. Applications using the global helpers can explicitly register the adapter at startup with `resp.Configure(sqlite3mapper.Option())`.

## Error mappings

Mappings follow [SQLite's primary and extended result codes](https://www.sqlite.org/rescode.html). Known extended constraint codes take precedence over the generic constraint category.

| SQLite condition | Driver code | Public code | HTTP |
| --- | --- | --- | ---: |
| Unique, primary-key, or rowid constraint | `ErrConstraintUnique`, `ErrConstraintPrimaryKey`, `ErrConstraintRowID` | `RESOURCE_ALREADY_EXISTS` | 409 |
| Not-null, foreign-key, or check constraint | `ErrConstraintNotNull`, `ErrConstraintForeignKey`, `ErrConstraintCheck` | `VALIDATION_FAILED` | 422 |
| STRICT table datatype constraint | `SQLITE_CONSTRAINT_DATATYPE` (3091) | `VALIDATION_FAILED` | 422 |
| Datatype mismatch, including invalid integer rowid | `ErrMismatch` | `BAD_REQUEST` | 400 |
| Text/BLOB exceeds the SQLite size limit | `ErrTooBig` | `VALIDATION_FAILED` | 422 |
| Database busy or table locked; reviewed recovery/snapshot/timeout/shared-cache codes | `ErrBusy`, `ErrBusyRecovery`, `ErrBusySnapshot`, `SQLITE_BUSY_TIMEOUT`, `ErrLocked`, `ErrLockedSharedCache` | `RESOURCE_CONFLICT` | 409 |
| Out of memory, database/disk full, or database cannot open; listed open variants | `ErrNomem`, `ErrFull`, `ErrCantOpen`, `ErrCantOpenNoTempDir`, `ErrCantOpenIsDir`, `ErrCantOpenFullPath`, `ErrCantOpenConvPath` | `SERVICE_UNAVAILABLE` | 503 |
| I/O failures (all extended codes), SQL/schema bugs, read-only or database permission failures, corruption, invalid database files, misuse, and other conditions | `ErrIoErr` and other typed SQLite errors | `INTERNAL_ERROR` | 500 |
| Generic/unknown constraints or application-specific trigger/function/commit-hook/virtual-table conditions | Other `ErrConstraint` extended codes | `INTERNAL_ERROR` | 500 |

Every recognized SQLite error receives a definition. Unknown current or future primary and extended codes fall back to `INTERNAL_ERROR` (500), even when the primary category is otherwise mapped. Unrelated errors and nil pointers are left to the core resolver. SQLite's `ErrNotFound` refers to an internal opcode condition, not a missing application resource; database authorization errors do not imply HTTP 401/403. Constraint/conflict semantics follow the policy in [issue #8](https://github.com/ZihxS/RespKit/issues/8) and the planned modernc adapter: reviewed conflicts use 409, selected validation constraints use 422, and I/O/corruption/SQL/schema/unknown conditions remain safe internal errors.

The mapper does not retry statements, roll back transactions, or set `Retry-After`. Handle those decisions in your application based on transaction state and whether the operation can safely be repeated. `api.Fail` preserves RespKit's cancellation behavior: errors wrapping `context.Canceled` return `resp.ErrRequestCanceled` without writing a response.

`sql.ErrNoRows` has application-specific meaning and stays unmapped. For a resource lookup, translate it explicitly:

```go
if errors.Is(dbErr, sql.ErrNoRows) {
	dbErr = resp.NewError(resp.CodeResourceNotFound, http.StatusNotFound, resp.WithCause(dbErr))
}
return api.Fail(resp.HTTP(w, r), dbErr)
```

Explicit `resp.AppError` definitions take precedence over the adapter. To override an individual mapping, register your own `resp.WithMapper` before `sqlite3mapper.Option()`.

## Localized messages

The default locale is Indonesian (`id`). `Accept-Language: en` selects English at request time; regional tags and quality values work with RespKit's locale matcher. Codes and HTTP statuses stay the same across languages. Every mapped public code already has both built-in translations:

| Public code | Indonesian | English |
| --- | --- | --- |
| `RESOURCE_ALREADY_EXISTS` | Data sudah terdaftar. | The resource already exists. |
| `VALIDATION_FAILED` | Data yang dikirim tidak valid. | The submitted data is invalid. |
| `BAD_REQUEST` | Permintaan tidak valid. | The request is invalid. |
| `RESOURCE_CONFLICT` | Terjadi konflik dengan data yang ada. | The request conflicts with the current resource state. |
| `SERVICE_UNAVAILABLE` | Layanan sementara tidak tersedia. | The service is temporarily unavailable. |
| `INTERNAL_ERROR` | Terjadi kesalahan pada server. | An internal server error occurred. |

A unique constraint failure with `Accept-Language: en` returns HTTP 409, `Content-Language: en`, and `Vary: Accept-Language`:

```json
{
  "success": false,
  "error": {
    "code": "RESOURCE_ALREADY_EXISTS",
    "message": "The resource already exists."
  }
}
```

Use `resp.Locale("en")` to override the locale for one response. To change the default language, add `resp.WithDefaultLocale("en")` when creating the engine. Additional languages use `resp.WithTranslations`; provide all six mapped codes for a complete catalog. Partial catalogs fall back to the default language when a message is missing.

```go
api, err := resp.New(
	sqlite3mapper.Option(),
	resp.WithTranslations(map[string]map[string]string{
		"fr": {
			resp.CodeResourceAlreadyExists: "La ressource existe déjà.",
			resp.CodeValidationFailed: "Les données soumises sont invalides.",
			resp.CodeBadRequest: "La requête est invalide.",
			resp.CodeResourceConflict: "La requête est en conflit avec l'état actuel de la ressource.",
			resp.CodeServiceUnavailable: "Le service est temporairement indisponible.",
			resp.CodeInternalError: "Une erreur interne du serveur s'est produite.",
		},
	}),
)
```

Both the default envelope and `resp.WithErrorFormat(resp.ProblemDetails)` use localized public messages. Raw driver messages, table/column names, database paths, and OS errors are hidden by default. The mapper does not generate field details from database schema names. Add trusted field details through an application error when needed. `resp.WithExposeInternalErrors(true)` remains an explicitly unsafe option for isolated local development.

`api.Inspect(dbErr)` identifies accepted mappings as `custom.sqlite3`. The driver does not implement `SQLState()`, so the core diagnostic `source_code` is empty for ordinary SQLite errors. Inspect the typed error privately when you need its primary/extended code; keep it out of public responses.

## Verify the adapter

From the repository root:

```bash
make integration-mappers
```

Or run its checks directly:

```bash
cd integrations/sqlite3
CGO_ENABLED=1 go test -race -count=1 -cover ./...
CGO_ENABLED=1 go vet ./...
```

The tests use the actual driver to create constraint failures, lock contention, read-only errors, capacity limits, and invalid/open failures. They verify Indonesian/English messages, request locale negotiation, custom-language fallback, both response formats, mapper precedence, cancellation, and confidentiality. Synthetic typed errors cover conditions such as out-of-memory and I/O failures without disrupting the host environment. CI is configured to test the adapter on Linux, macOS, and Windows, plus the Go 1.22 minimum on Linux.
