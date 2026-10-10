# RespKit

[Documentation](https://zihxs.github.io/RespKit/) · [GitHub](https://github.com/ZihxS/RespKit) · [License](LICENSE)

RespKit gives Go applications consistent HTTP responses and maps errors to HTTP status codes and client messages. Use it with net/http and popular Go frameworks. Indonesian and English translations are included, and applications can add their own languages.

## Quick start

RespKit supports Go 1.22 and later. Install the current version from the main branch:

```bash
go get github.com/ZihxS/RespKit/resp@main
```

In Gin, Echo, or Fiber, pass the framework context to a response helper:

```go
resp.OK(c, user)
resp.Created(c, user, resp.Location("/users/123"))
resp.Fail(c, err)
```

With net/http, Chi, or Gorilla/mux, pass the response writer and request through RespKit:

```go
if err := resp.OK(resp.HTTP(w, r), user); err != nil {
	log.Print(err)
}
```

Beego v2 uses the same approach with its response writer and request. See the [framework guide](https://zihxs.github.io/RespKit/frameworks/) for complete examples.

## Responses and errors

A successful response looks like this:

```json
{"success":true,"data":{"user_id":123}}
```

RespKit keeps application data unchanged, selects a message language for each request, and hides raw internal error messages from clients by default. It recognizes wrapped Go errors and PostgreSQL's unique-constraint error. Error responses can use RFC 9457 Problem Details.

The built-in translation files cover Indonesian and English. Add more translations through configuration. Read about [response formats](https://zihxs.github.io/RespKit/responses/), [error mapping](https://zihxs.github.io/RespKit/errors/), and [localization](https://zihxs.github.io/RespKit/localization/).

## PostgreSQL

The core recognizes the PostgreSQL unique-constraint code without importing a database driver. The optional [pgx adapter](integrations/pgx/README.md) maps additional PostgreSQL errors. It has a separate Go module, so applications only add it when needed.

## Tests

Run the checks from the repository root:

```bash
make test
make race
make vet
make coverage
make integration
make integration-mappers
```

The integration suite sends requests through Gin, Echo v4/v5, Fiber v2/v3, Chi, Gorilla/mux, and Beego. CI uses Go 1.27.x on Ubuntu, macOS, and Windows. See the [verification report](docs/VERIFICATION.md) for tested versions and results.

## Version

The module path is github.com/ZihxS/RespKit. The source is available on GitHub; a GitHub release for v1.0.0 has not been published yet. The [release checklist](docs/RELEASE_CHECKLIST.md) tracks the remaining steps.
