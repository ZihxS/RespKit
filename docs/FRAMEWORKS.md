# HTTP framework integration

RespKit uses one public response package across Go HTTP frameworks. The core module does not import framework packages, so applications install only the framework they already use.

Gin, Echo v4/v5, and Fiber v2/v3 accept their native context values. Standard `net/http`, Chi, Gorilla/mux, and Beego v2 handlers use `resp.HTTP` to pair the response writer with the request.

## Standard `net/http`, Chi, and Gorilla/mux

Chi and Gorilla/mux register regular `net/http` handlers. This example works with either router:

```go
package handlers

import (
	"log"
	"net/http"

	"github.com/ZihxS/RespKit/resp"
)

func GetUser(w http.ResponseWriter, r *http.Request) {
	user := map[string]any{
		"user_id": 7,
		"name":    "Ada",
	}
	if err := resp.OK(resp.HTTP(w, r), user); err != nil {
		log.Printf("write response: %v", err)
	}
}
```

Register the handler with Chi or Gorilla/mux:

```go
router.Get("/users", handlers.GetUser) // Chi
router.HandleFunc("/users", handlers.GetUser).Methods(http.MethodGet) // Gorilla/mux
```

## Beego v2

Beego exposes the underlying response writer and request through its context:

```go
func GetUser(c *beecontext.Context) {
	user := map[string]any{
		"user_id": 7,
		"name":    "Ada",
	}
	if err := resp.OK(resp.HTTP(c.ResponseWriter, c.Request), user); err != nil {
		log.Printf("write response: %v", err)
	}
}
```

Import the context type from `github.com/beego/beego/v2/server/web/context`.

## Gin

Pass the Gin context directly. Return helper errors to Gin's error list for application middleware to inspect:

```go
func GetUser(c *gin.Context) {
	user := map[string]any{
		"user_id": 7,
		"name":    "Ada",
	}
	if err := resp.OK(c, user); err != nil {
		_ = c.Error(err)
	}
}
```

## Echo v4 and v5

Return helper errors so Echo's HTTP error handler can process them. Echo v4 uses `echo.Context`:

```go
func GetUserV4(c echo.Context) error {
	return resp.Fail(c, resp.ErrNotFound)
}
```

Echo v5 uses `*echo.Context`:

```go
func GetUserV5(c *echo.Context) error {
	return resp.Fail(c, resp.ErrNotFound)
}
```

## Fiber v2 and v3

Return the helper error to Fiber. Fiber v2 uses `*fiber.Ctx`:

```go
func GetUserV2(c *fiber.Ctx) error {
	return resp.OK(c, map[string]any{
		"user_id": 7,
		"name":    "Ada",
	})
}
```

Fiber v3 uses `fiber.Ctx`:

```go
func GetUserV3(c fiber.Ctx) error {
	return resp.OK(c, map[string]any{
		"user_id": 7,
		"name":    "Ada",
	})
}
```

## Compatibility details

- Gin writes the pre-encoded payload through `JSON(int, any)` and provides `GetHeader` and `Header` methods.
- Echo writes through `Blob(int, string, []byte)`, reads the request with `Request()`, and uses the response header map.
- Fiber provides `Get`, `Set`, and `Send`. RespKit reflects on `Status(int)` because the return type differs between Fiber v2 and v3.
- Chi, Gorilla/mux, and Beego use the standard HTTP writer and request; the core module has no dependency on those frameworks.
- Middleware that has already written a status or body can prevent RespKit from setting the response. Call a response helper at the end of the handler.

## Run the upstream tests

`integrations/framework-tests` contains runtime tests against the actual upstream routers. They check status codes, response JSON, localization headers, error confidentiality, `Location`, and 204 responses. See [VERIFICATION.md](VERIFICATION.md) for the pinned versions and results. Local interface doubles do not replace these integration tests.
