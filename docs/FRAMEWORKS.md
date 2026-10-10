# Framework integration — single import

These examples use the selected module path `github.com/ZihxS/RespKit`. They will resolve after the repository is published at that path.

## Standard net/http, Chi, and Gorilla/mux

Chi and Gorilla/mux use ordinary `net/http` handlers. Pass the handler's writer and request through `resp.HTTP`:

```go
import (
    "net/http"
    "github.com/ZihxS/RespKit/resp"
)
func handler(w http.ResponseWriter,r *http.Request) {
    if err:=resp.OK(resp.HTTP(w,r),map[string]any{"user_id":7});err!=nil {
        // Write failure has already occurred; log it, don't write again.
    }
}
```

The same handler body works when registered with Chi or Gorilla/mux:

```go
router.Get("/users", Handler) // Chi
router.HandleFunc("/users", Handler).Methods(http.MethodGet) // Gorilla/mux
```

## Beego v2

Beego's context exposes the underlying HTTP response writer and request. Wrap those values with `resp.HTTP`:

```go
import (
    beecontext "github.com/beego/beego/v2/server/web/context"
    "github.com/ZihxS/RespKit/resp"
)

func handler(c *beecontext.Context) {
    if err := resp.OK(resp.HTTP(c.ResponseWriter, c.Request), map[string]any{"user_id": 7}); err != nil {
        // Write failure has already occurred; log it, don't write again.
    }
}
```

## Gin

```go
import (
    "github.com/gin-gonic/gin"
    "github.com/ZihxS/RespKit/resp"
)
func handler(c *gin.Context) {
    if err:=resp.OK(c,map[string]any{"user_id":7});err!=nil { _=c.Error(err) }
}
```

## Echo v4

```go
import (
    "github.com/labstack/echo/v4"
    "github.com/ZihxS/RespKit/resp"
)
func handler(c echo.Context) error {
    return resp.Fail(c,resp.ErrNotFound)
}
```

## Echo v5

```go
import (
    "github.com/labstack/echo/v5"
    "github.com/ZihxS/RespKit/resp"
)
func handler(c *echo.Context) error {
    return resp.Fail(c, resp.ErrNotFound)
}
```

## Fiber v2

```go
import (
    "github.com/gofiber/fiber/v2"
    "github.com/ZihxS/RespKit/resp"
)
func handler(c *fiber.Ctx) error {
    return resp.OK(c,map[string]any{"user_id":7})
}
```

## Fiber v3

```go
import (
    "github.com/gofiber/fiber/v3"
    "github.com/ZihxS/RespKit/resp"
)
func handler(c fiber.Ctx) error {
    return resp.OK(c,map[string]any{"user_id":7})
}
```

## Important compatibility details

- Gin adapter detects `JSON(int, any)`, `GetHeader(string)`, `Header(string,string)`; the built-in `c.JSON` is used to send the pre-encoded payload.
- Echo adapter detects `Blob(int,string,[]byte)` and `Request() *http.Request`; the Echo response headers are resolved through `Response().Header()`.
- Fiber adapter detects `Get(string,...string) string`, `Set(string,string)`, `Send([]byte) error`, plus a callable `Status(int)` method. Fiber v2/v3 Status return types differ, so just this operation is reflected.
- Chi, Gorilla/mux, and Beego v2 use the standard `net/http` response writer and request through `resp.HTTP`; the core package has no framework dependency for them.
- None of the above frameworks is linked into RespKit core. Their version compatibility must be checked against the version used by the consumer.
- Framework-specific middleware which writes a response first can prevent helpers from setting headers or status. Place RespKit at the terminal response point.

### How to verify real-framework compatibility before publishing

The independent `integrations/framework-tests` Go module runs end-to-end requests against pinned upstream framework versions, with a separate folder per framework/version and shared assertions in `internal/testkit`. It verifies status, JSON, `Content-Language`, `Vary`, `Content-Type`, error confidentiality, `Location`, and 204 behavior for every supported integration family. Contract-double tests in `resp/framework_test.go` are necessary but not sufficient.
