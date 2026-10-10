package gorillamux

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ZihxS/RespKit/integrations/framework-tests/internal/testkit"
	"github.com/gorilla/mux"
)

func TestReal(t *testing.T) {
	app := mux.NewRouter()
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
		app.Handle(route.path, testkit.StandardHandler(t, "gorilla/mux", route.kind)).Methods(route.method)
	}
	testkit.TestResponse(t, "gorilla-mux", func(r *http.Request) (*http.Response, error) {
		w := httptest.NewRecorder()
		app.ServeHTTP(w, r)
		return w.Result(), nil
	})
}
