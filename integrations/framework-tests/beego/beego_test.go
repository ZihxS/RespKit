package beego

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ZihxS/RespKit/integrations/framework-tests/internal/testkit"
	"github.com/ZihxS/RespKit/resp"
	web "github.com/beego/beego/v2/server/web"
	beecontext "github.com/beego/beego/v2/server/web/context"
)

func TestReal(t *testing.T) {
	app := web.NewControllerRegister()
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
		handler := func(c *beecontext.Context) {
			if err := testkit.Execute(kind, resp.HTTP(c.ResponseWriter, c.Request)); err != nil {
				t.Errorf("beego: %v", err)
			}
		}
		app.AddMethod(strings.ToLower(route.method), route.path, handler)
	}
	testkit.TestResponse(t, "beego", func(r *http.Request) (*http.Response, error) {
		w := httptest.NewRecorder()
		app.ServeHTTP(w, r)
		return w.Result(), nil
	})
}
