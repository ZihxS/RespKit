package resp

import (
	"net/http"

	"github.com/ZihxS/RespKit/internal/transport"
)

// HTTPContext combines net/http's two handler arguments. Use it with net/http,
// Chi, Gorilla/mux, Beego, and other handlers that expose both values.
type HTTPContext = transport.HTTPContext

// HTTP adapts a standard net/http handler request to the unified resp API.
func HTTP(w http.ResponseWriter, r *http.Request) *HTTPContext {
	return transport.HTTP(w, r)
}
