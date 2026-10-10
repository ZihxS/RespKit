# RespKit v1.0.0 (source candidate)

**Zero-config REST response formatting, safe Go error normalization, Indonesian/English localization and structured diagnostics. One `resp` API across Go HTTP frameworks.**

**Dokumentasi:** [RespKit Book](https://zihxs.github.io/RespKit/)

> **Release state:** Code is an unpublished `v1.0.0` source candidate, not a tagged GitHub release. The selected module path is `github.com/ZihxS/RespKit`. Runtime unit tests pass with **100% combined statement coverage** for `resp` and `internal/...`. Upstream framework tests run from `integrations/framework-tests`; see [verification](docs/VERIFICATION.md) for tested versions.

## Zero config

```go
import "github.com/ZihxS/RespKit/resp"

resp.OK(c, user)
resp.Created(c, user, resp.Location("/users/123"))
resp.NoContent(c)
resp.Fail(c, err)
resp.NotFound(c)
```

No `New()`, `Init()`, or explicit registration is needed. Each method returns `error`; return it from Echo/Fiber handlers, or handle it explicitly in Gin.

For standard `net/http`, Chi, and Gorilla/mux, use the ordinary handler writer and request:

```go
func Handler(w http.ResponseWriter, r *http.Request) {
    if err := resp.OK(resp.HTTP(w, r), map[string]any{"user_id": 123}); err != nil {
        log.Print(err)
    }
}
```

Gin, Echo v4/v5, Fiber v2/v3 use the native context with the same import and `resp.OK(c, data)` surface. Detection uses structural interfaces and small amounts of reflection, **not** an imported framework-specific wrapper.

Beego v2 uses the standard HTTP wrapper around its context fields: `resp.HTTP(c.ResponseWriter, c.Request)`. Chi, Gorilla/mux, and Beego are tested against upstream routers in `integrations/framework-tests`.

## Default behavior

Success HTTP 200:

```json
{"success":true,"data":{"user_id":123}}
```

Wrapped `resp.ErrNotFound` on an Indonesian request produces real HTTP 404:

```json
{"success":false,"error":{"code":"RESOURCE_NOT_FOUND","message":"Data tidak ditemukan."}}
```

- Default `snake_case` for **RespKit keys only**, leaving the application's JSON tags and `data` untouched.
- Default locale `id`, supported `id/en`, `Accept-Language` detection, region fallback and `Content-Language`/`Vary` headers.
- Built-in translations live in one JSON file per locale and are embedded in the binary; no locale file is read from the deployment filesystem at runtime.
- Domain sentinel errors (`errors.Is/As`), JSON syntax/type errors and PostgreSQL SQLSTATE `23505` structural matching. Ambiguous DB errors and unknown errors safely yield HTTP 500 unless explicitly mapped.
- Error code is stable across languages. Raw `err.Error()` is hidden by default.
- Optional pgx/v5 mapping is available from [`integrations/pgx`](integrations/pgx/README.md); the core stays free of pgx dependencies.
- Optional Problem Details (RFC 9457) for error responses, `HEAD`/204 no body, 401 Bearer challenge (overrideable), and 201 Location header.

## Optional configuration

```go
err := resp.Configure(
    resp.WithNaming(resp.CamelCase),
    resp.WithDefaultLocale("en"),
    resp.WithTranslations(map[string]map[string]string{
        "en": {"USER_SUSPENDED": "Your account is suspended."},
        "id": {"USER_SUSPENDED": "Akun Anda ditangguhkan."},
    }),
)
if err != nil { panic(err) }
```

Or create an isolated engine:

```go
api, err := resp.New(resp.WithErrorFormat(resp.ProblemDetails))
if err != nil { panic(err) }
_ = api.Fail(c, resp.ErrConflict)
```

Use `resp.Locale("en")` to override the locale for one response. `resp.NewError("USER_SUSPENDED", 403, resp.WithCause(err))` creates an explicit domain error with an internal cause. `resp.WithMapper` supports extensions that abstain when classification is ambiguous.

The built-in catalogs are maintained in `internal/i18n/locales/id.json` and `internal/i18n/locales/en.json`. `WithTranslations` adds application-specific messages or replaces a built-in message for a locale.

### Unsafe local-only internal error messages

**Default: OFF.** For isolated local/dev usage only:

```go
if err := resp.Configure(resp.WithExposeInternalErrors(true)); err != nil { panic(err) }
```

This adds an `error.debug_message` (or the corresponding naming convention) while preserving the safe public `message`:

```json
{"success":false,"error":{"code":"INTERNAL_ERROR","message":"Terjadi kesalahan pada server.","debug_message":"db: connection refused"}}
```

The option may expose SQL queries, tokens or personal information. **Never enable it on internet-accessible services**, and never let a request header turn it on. You can also use `resp.New(resp.WithExposeInternalErrors(true))` for an isolated instance.

## Tests and source layout

```bash
make test       # root tests
make race       # race detector
make vet
make coverage   # 100.0% combined statement coverage for resp + internal runtime
make integration  # genuine upstream Gin/Echo/Fiber, Chi, Gorilla/mux, Beego; Go 1.26+ and network required
make integration-mappers  # optional pgx mapper tests; Go 1.25+ and network required
```

```text
resp/                   # stable public API and response pipeline
internal/transport/     # net/http adapter + Gin / Echo / Fiber bridges and detection
internal/resolver/      # Go-idiomatic errors, custom and built-in mapping
internal/i18n/          # built-in catalogs, locale matching, fallback
internal/formatter/     # stable JSON envelope, casing, RFC 9457
internal/diagnostics/   # machine-readable, safe diagnostics
cmd/respkit/            # version, doctor, catalog, explain
examples/nethttp/       # runnable HTTP server
integrations/           # integration test harnesses and optional adapters
  framework-tests/      # independent module, upstream tests, folders by framework
    beego/ chi/ gin/ gorilla_mux/
    echo/v4/ echo/v5/ fiber/v2/ fiber/v3/
    internal/testkit/   # shared HTTP response assertions
  pgx/                  # optional PostgreSQL/pgx error mapper module
```

The root module has no external dependencies. Each `integrations/` subfolder has its own `go.mod`; ordinary RespKit users do not install test frameworks or optional adapters unless they choose to use them.

Diagnostics for coding assistants require no AI key: `go run ./cmd/respkit doctor`, `catalog`, `explain RESOURCE_NOT_FOUND`. Developer-specific instructions are in [`AGENTS.md`](AGENTS.md).

See [architecture](docs/ARCHITECTURE.md), [verification report](docs/VERIFICATION.md), [release checklist](docs/RELEASE_CHECKLIST.md), and [contributing](CONTRIBUTING.md). License: MIT.
