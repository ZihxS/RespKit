package echov4

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ZihxS/RespKit/integrations/framework-tests/internal/testkit"
	echo "github.com/labstack/echo/v4"
)

func TestReal(t *testing.T) {
	app := echo.New()
	app.GET("/ok", func(c echo.Context) error { return testkit.Execute("ok", c) })
	app.GET("/error", func(c echo.Context) error { return testkit.Execute("error", c) })
	app.GET("/raw", func(c echo.Context) error { return testkit.Execute("raw", c) })
	app.GET("/debug", func(c echo.Context) error { return testkit.Execute("debug", c) })
	app.POST("/created", func(c echo.Context) error { return testkit.Execute("created", c) })
	app.GET("/empty", func(c echo.Context) error { return testkit.Execute("empty", c) })
	testkit.TestResponse(t, "echo-v4", func(r *http.Request) (*http.Response, error) {
		w := httptest.NewRecorder()
		app.ServeHTTP(w, r)
		return w.Result(), nil
	})
}
