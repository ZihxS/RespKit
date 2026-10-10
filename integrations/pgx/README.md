# RespKit pgx adapter

This optional module maps selected PostgreSQL SQLSTATEs from pgx/v5 to safe RespKit error definitions. It does not open database connections or change global RespKit configuration.

## Installation

Add the core response package and the adapter to the application:

```bash
go get github.com/ZihxS/RespKit/resp
go get github.com/ZihxS/RespKit/integrations/pgx
```

The adapter uses pgx/v5 v5.11.0, which requires Go 1.25. It is a separate module; the core module supports Go 1.22 and later and has no pgx dependency.

## Configure the mapper

Create the engine once during application startup, then reuse it when writing responses:

```go
package handlers

import (
	"net/http"

	pgxmapper "github.com/ZihxS/RespKit/integrations/pgx"
	"github.com/ZihxS/RespKit/resp"
)

func newResponseEngine() (*resp.Engine, error) {
	return resp.New(pgxmapper.Option())
}

func writeDatabaseError(
	api *resp.Engine,
	w http.ResponseWriter,
	r *http.Request,
	dbErr error,
) error {
	return api.Fail(resp.HTTP(w, r), dbErr)
}
```

The mapper recognizes wrapped pgx errors with `errors.As`:

| SQLSTATE | Common PostgreSQL condition | RespKit code | HTTP status |
| --- | --- | --- | ---: |
| `22003`, `22007`, `22008`, `22018`, `22021`, `22023`, `22P02`, `22P03` | Invalid numeric, date/time, text, or binary input | `BAD_REQUEST` | 400 |
| `22001` | Text value exceeds the column length | `VALIDATION_FAILED` | 422 |
| `23502`, `23503`, `23514` | Not-null, foreign-key, or check constraint violation | `VALIDATION_FAILED` | 422 |
| `23505` | Unique constraint violation | `RESOURCE_ALREADY_EXISTS` | 409 |
| `23000`, `23001`, `23P01` | Integrity, restrict, or exclusion conflict | `RESOURCE_CONFLICT` | 409 |
| `40000`, `40001`, `40P01` | Transaction rollback, serialization failure, or deadlock | `RESOURCE_CONFLICT` | 409 |
| `08000`, `08001`, `08003`, `08004`, `08006`, `08007` | Database connection failure | `SERVICE_UNAVAILABLE` | 503 |
| `53000`, `53100`, `53200`, `53300`, `53400` | Insufficient resources or a server configuration limit | `SERVICE_UNAVAILABLE` | 503 |
| `57P01`, `57P02`, `57P03` | Database shutdown or not accepting connections | `SERVICE_UNAVAILABLE` | 503 |
| `55P03` | Lock not available | `RESOURCE_CONFLICT` | 409 |

The mapper accepts every `*pgconn.PgError` SQLSTATE. Listed codes receive the definitions above; any other current or future PostgreSQL code maps to `INTERNAL_ERROR` (500), with its SQLSTATE retained in diagnostics. This avoids returning a client error for SQL bugs, authentication or configuration problems, protocol failures, and internal database conditions. Raw PostgreSQL messages are never exposed.

The core module also maps `23505` from any error that implements `SQLState() string`. When enabled, this adapter identifies its mapping as `custom.pgx` in diagnostics.

## Run the adapter tests

From the repository root:

```bash
cd integrations/pgx
go test -race -count=1 ./...
```
