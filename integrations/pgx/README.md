# RespKit pgx mapper

This optional module maps common [`pgx/v5`](https://github.com/jackc/pgx) `*pgconn.PgError` SQLSTATEs to safe HTTP error definitions. It does not open connections or change global RespKit configuration.

## Installation

Import `github.com/ZihxS/RespKit/integrations/pgx` alongside RespKit. This adapter is pinned to pgx/v5 `v5.11.0`, which requires Go 1.25. It is a separate module; the core RespKit module remains Go 1.22-compatible and has no pgx dependency.

## Enable the mapper

```go
import (
	pgxmapper "github.com/ZihxS/RespKit/integrations/pgx"
	"github.com/ZihxS/RespKit/resp"
)

func writeDatabaseError(ctx any, dbErr error) error {
	api, err := resp.New(pgxmapper.Option())
	if err != nil {
		return err
	}
	return api.Fail(ctx, dbErr)
}
```

The mapper recognizes wrapped pgx errors with `errors.As`:

| SQLSTATEs | PostgreSQL conditions | RespKit code | HTTP |
| --- | --- | --- | ---: |
| `22003`, `22007`, `22008`, `22018`, `22021`, `22023`, `22P02`, `22P03` | Invalid numeric, date/time, text, or binary request values | `BAD_REQUEST` | 400 |
| `22001` | String data right truncation | `VALIDATION_FAILED` | 422 |
| `23502`, `23503`, `23514` | Not-null, foreign-key, or check constraint violation | `VALIDATION_FAILED` | 422 |
| `23505` | Unique violation | `RESOURCE_ALREADY_EXISTS` | 409 |
| `23000`, `23001`, `23P01` | Integrity, restrict, or exclusion conflict | `RESOURCE_CONFLICT` | 409 |
| `40000`, `40001`, `40P01` | Transaction rollback, serialization failure, or deadlock | `RESOURCE_CONFLICT` | 409 |
| `08000`, `08001`, `08003`, `08004`, `08006`, `08007` | Connection exception or failure | `SERVICE_UNAVAILABLE` | 503 |
| `53000`, `53100`, `53200`, `53300`, `53400` | Insufficient resources, disk full, out of memory, too many connections, or configuration limit | `SERVICE_UNAVAILABLE` | 503 |
| `57P01`, `57P02`, `57P03` | Database shutdown or cannot connect now | `SERVICE_UNAVAILABLE` | 503 |
| `55P03` | Lock not available | `RESOURCE_CONFLICT` | 409 |

The mapper accepts every `*pgconn.PgError` SQLSTATE. Codes in the table receive the listed HTTP definition; every other current or future PostgreSQL code receives `INTERNAL_ERROR` (500), with its SQLSTATE retained in diagnostics. This avoids assigning 4xx responses to SQL bugs, authentication/configuration problems, protocol failures, and internal database conditions. No raw PostgreSQL message is exposed. The core also maps `23505` from any error implementing `SQLState() string`; enabling this adapter identifies the mapping as `custom.pgx` in diagnostics.

## Run tests

From the repository root:

```bash
cd integrations/pgx
go test -race -count=1 ./...
```
