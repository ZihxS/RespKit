# RespKit

RespKit provides consistent HTTP responses, error mapping, and configurable localized messages for Go applications. You can use the same API with net/http, Gin, Echo, Fiber, Chi, Gorilla/mux, and Beego.

[Documentation](https://zihxs.github.io/RespKit/) · [GitHub](https://github.com/ZihxS/RespKit)

## Install

The module supports Go 1.22 and later. Install the current default branch with:

```bash
go get github.com/ZihxS/RespKit/resp@main
```

The root module uses only the Go standard library. Framework test dependencies and the optional pgx adapter are kept in separate modules.

## Use

With Gin, Echo, or Fiber, pass the framework context directly:

```go
resp.OK(c, user)
resp.Created(c, user, resp.Location("/users/123"))
resp.NoContent(c)
resp.Fail(c, err)
resp.NotFound(c)
```

For net/http, Chi, and Gorilla/mux, pass the response writer and request to RespKit:

```go
func handler(w http.ResponseWriter, r *http.Request) {
	if err := resp.OK(resp.HTTP(w, r), map[string]any{"user_id": 123}); err != nil {
		log.Print(err)
	}
}
```

Beego v2 uses the same approach with its response writer and request. Each response helper returns an error. Return it from Echo and Fiber handlers, or handle it in Gin and standard HTTP handlers.

These integrations are tested through requests served by the upstream routers. The pinned versions and test results are listed in the [verification report](docs/VERIFICATION.md).

## Defaults

Successful responses use this envelope:

```json
{"success":true,"data":{"user_id":123}}
```

An unknown error produces a generic HTTP 500 response. Raw error messages are not sent to clients by default. Other defaults include:

- RespKit's JSON keys use snake_case by default; application data is left unchanged.
- Indonesian is the default locale, and English is included. Applications can add more locales, and RespKit reads the Accept-Language header on each request.
- Built-in translations are embedded in the binary from one JSON file per locale.
- Wrapped errors are classified with errors.Is and errors.As. PostgreSQL SQLSTATE 23505 maps to HTTP 409.
- Created responses can set a Location header. HEAD and 204 No Content responses have no body.
- Error responses can use the RFC 9457 Problem Details format.

## Configuration

The response helpers work with the default settings. Create an Engine when an application needs its own settings:

```go
api, err := resp.New(
	resp.WithDefaultLocale("en"),
	resp.WithErrorFormat(resp.ProblemDetails),
)
if err != nil {
	return err
}

return api.Fail(c, resp.ErrNotFound)
```

See [configuration](https://zihxs.github.io/RespKit/configuration/) and [response formats](https://zihxs.github.io/RespKit/responses/) for available options.

## PostgreSQL and pgx

The core recognizes SQLSTATE 23505 without a database dependency. The optional [pgx adapter module](integrations/pgx/README.md) maps other common PostgreSQL error codes without adding pgx to the core module.

## Run checks

```bash
make test
make race
make vet
make coverage
make integration
make integration-mappers
```

CI uses Go 1.27.x on Ubuntu, macOS, and Windows. The framework and pgx checks fetch their dependencies from upstream.

## Repository layout

```text
resp/                   Public API and response pipeline
internal/transport/     HTTP and framework context adapters
internal/resolver/      Error classification
internal/i18n/          Translation catalogs and locale matching
internal/formatter/     JSON envelopes and Problem Details
internal/diagnostics/   Structured error information
cmd/respkit/            Version and diagnostic commands
examples/nethttp/       Runnable net/http example
integrations/           Framework tests and optional adapters
book/                   Documentation site source
```

## Version status

The module path is github.com/ZihxS/RespKit. The repository contains the v1.0.0 source, but a GitHub release for that version has not been published. See the [release checklist](docs/RELEASE_CHECKLIST.md) for the remaining release steps.
