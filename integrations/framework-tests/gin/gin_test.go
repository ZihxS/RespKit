package gin

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ZihxS/RespKit/integrations/framework-tests/internal/testkit"
	"github.com/gin-gonic/gin"
)

func TestReal(t *testing.T) {
	gin.SetMode(gin.TestMode)
	app := gin.New()
	routes := []struct {
		path, method, kind string
	}{
		{"/ok", "GET", "ok"},
		{"/error", "GET", "error"},
		{"/raw", "GET", "raw"},
		{"/debug", "GET", "debug"},
		{"/created", "POST", "created"},
		{"/empty", "GET", "empty"},
	}
	for _, route := range routes {
		kind := route.kind
		handler := func(c *gin.Context) {
			if err := testkit.Execute(kind, c); err != nil {
				t.Errorf("gin: %v", err)
			}
		}
		app.Handle(route.method, route.path, handler)
	}
	testkit.TestResponse(t, "gin", func(r *http.Request) (*http.Response, error) {
		w := httptest.NewRecorder()
		app.ServeHTTP(w, r)
		return w.Result(), nil
	})
}
