<h1 align="center">RespKit — Consistent HTTP Responses for Go</h1>

<p align="center">
  <img src="assets/respkit.png" alt="RespKit" width="560">
</p>

RespKit gives HTTP handlers a shared response format, maps Go errors to stable status codes,
and selects translated messages for each request. It works with Go's standard HTTP types and
popular web frameworks, while the core module stays dependency-free.

[Documentation](https://zihxs.github.io/RespKit/) ·
[Examples](https://zihxs.github.io/RespKit/quickstart/) ·
[Report a security issue](SECURITY.md)

## Features

- Zero-configuration helpers for successful and failed responses.
- Stable, language-independent error codes and conservative Go error mapping.
- JSON translation catalogs, with Indonesian and English included and support for additional
  languages.
- English and Indonesian messages selected from each request's Accept-Language header.
- RFC 9457 Problem Details for error responses.
- Safe defaults that keep raw internal error messages out of client responses.
- Optional adapters that keep database and framework dependencies out of the core module.

## Requirements and installation

RespKit supports Go 1.22 and later. From an existing Go module, install the response package
and Gin with:

```bash
go get github.com/gin-gonic/gin github.com/ZihxS/RespKit/resp
```

## Quick start with Gin

Save the following as main.go, then run it with go run main.go. It starts a Gin server with
a success route and an error route:

```go
package main

import (
	"log"

	"github.com/ZihxS/RespKit/resp"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/users/123", func(c *gin.Context) {
		user := gin.H{
			"user_id": 123,
			"name":    "Ada",
		}
		resp.OK(c, user)
	})

	router.GET("/users/missing", func(c *gin.Context) {
		resp.Fail(c, resp.ErrNotFound)
	})

	log.Fatal(router.Run(":8080"))
}
```

Gin, Echo, and Fiber accept their native contexts. Standard net/http handlers, Chi, Gorilla/mux,
and Beego pass the response writer and request to RespKit. See the
[framework guide](https://zihxs.github.io/RespKit/frameworks/) for examples for each router.

## Response formats

A 200 OK response uses the default success envelope. A 201 Created response uses the same
envelope and can also set a Location header. Application data appears under the data field:

```json
{
  "success": true,
  "data": {
    "user_id": 123,
    "name": "Ada"
  }
}
```

The error examples use English messages selected with Accept-Language: en. With the default
envelope format, a missing resource returns status 404, a stable code, and a localized message:

```json
{
  "success": false,
  "error": {
    "code": "RESOURCE_NOT_FOUND",
    "message": "The requested resource was not found."
  }
}
```

Validation failures return status 422 and can include field-level details for clients:

```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_FAILED",
    "message": "The submitted data is invalid.",
    "details": [
      {
        "field": "email",
        "code": "FIELD_REQUIRED",
        "message": "The email field is required."
      }
    ]
  }
}
```

Applications can select RFC 9457 Problem Details for errors through configuration.
A 404 response then looks like this:

```json
{
  "type": "about:blank",
  "title": "Not Found",
  "status": 404,
  "detail": "The requested resource was not found.",
  "code": "RESOURCE_NOT_FOUND"
}
```

Problem Details changes the representation of errors only; successful responses keep the
envelope shown above. A No Content response has status 204 and no body. Requests using the
HEAD method also receive no body. RespKit does not rewrite application data keys.

Read more about [response formats](https://zihxs.github.io/RespKit/responses/),
[error mapping](https://zihxs.github.io/RespKit/errors/), and
[localization](https://zihxs.github.io/RespKit/localization/).

## Optional integrations

The core module recognizes PostgreSQL's unique-constraint SQLSTATE without importing a
database driver. Optional adapters add mappings for a specific database or driver and live in
separate Go modules, so applications only include the integrations they use.

The [pgx adapter](integrations/pgx/README.md) maps additional PostgreSQL errors. The
[SQLite adapter](integrations/sqlite3/README.md) maps typed `mattn/go-sqlite3` errors,
including constraints, contention, and database availability failures. Both use RespKit's
Indonesian/English messages and support additional translation catalogs. See the
[PostgreSQL guide](https://zihxs.github.io/RespKit/postgres/) and
[SQLite guide](https://zihxs.github.io/RespKit/sqlite/) for mappings and setup.

## Testing and compatibility

The repository includes unit tests and runtime tests against upstream Gin, Echo v4/v5, Fiber
v2/v3, Chi, Gorilla/mux, and Beego. Run the checks from the repository root:

```bash
make test
make race
make vet
make coverage
make integration
make integration-mappers
```

The main module supports Go 1.22 and later. GitHub Actions runs manually: the default `core`
suite checks the core on Ubuntu, while `full` also checks macOS, Windows, and integrations
with Go 1.27.x. See the [contribution guide](CONTRIBUTING.md#run-github-actions-manually) for
manual runs and the [verification report](docs/VERIFICATION.md) for pinned framework versions
and test results.

## License

RespKit is available under the [MIT License](LICENSE).

## Contributors

Thanks to everyone who contributes to RespKit. See the
[contributors graph](https://github.com/ZihxS/RespKit/graphs/contributors) for the current
list.

To contribute, read the [contribution guide](CONTRIBUTING.md) or choose a task from the
[optional integrations roadmap](https://github.com/ZihxS/RespKit/issues/4).

## Star History

[![Star History Chart](https://api.star-history.com/svg?repos=ZihxS/RespKit&type=Date)](https://star-history.com/#ZihxS/RespKit&Date)
