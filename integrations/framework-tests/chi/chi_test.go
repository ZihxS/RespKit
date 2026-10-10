package chi

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ZihxS/RespKit/integrations/framework-tests/internal/testkit"
	"github.com/go-chi/chi/v5"
)

func TestReal(t *testing.T) {
	app := chi.NewRouter()
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
		app.Method(route.method, route.path, testkit.StandardHandler(t, "chi", route.kind))
	}
	testkit.TestResponse(t, "chi", func(r *http.Request) (*http.Response, error) {
		w := httptest.NewRecorder()
		app.ServeHTTP(w, r)
		return w.Result(), nil
	})
}
