package transport

import (
	"errors"
	"net/http"
	"strings"
)

var ErrUnsupportedContext = errors.New("unsupported HTTP context")

// HTTP combines the ResponseWriter and Request used by net/http handlers and
// compatible routers/frameworks such as Chi, Gorilla/mux, and Beego.
func HTTP(w http.ResponseWriter, r *http.Request) *HTTPContext {
	return &HTTPContext{Writer: w, Request: r}
}

type HTTPContext struct {
	Writer  http.ResponseWriter
	Request *http.Request
}

type Adapter interface {
	Header(key string) string
	SetHeader(key, value string)
	AddVary(value string)
	Method() string
	Write(status int, body []byte, mediaType string) error
}

type netAdapter struct {
	w http.ResponseWriter
	r *http.Request
}

func (a netAdapter) Header(key string) string {
	return a.r.Header.Get(key)
}

func (a netAdapter) SetHeader(key, value string) {
	a.w.Header().Set(key, value)
}

func (a netAdapter) AddVary(value string) {
	AddVary(a.w.Header(), value)
}

func (a netAdapter) Method() string {
	return a.r.Method
}

func (a netAdapter) Write(code int, body []byte, ctype string) error {
	if ctype != "" {
		a.w.Header().Set("Content-Type", ctype)
	}
	a.w.WriteHeader(code)
	if body == nil {
		return nil
	}
	_, err := a.w.Write(body)
	return err
}

func HasVary(vary, token string) bool {
	for _, t := range strings.Split(vary, ",") {
		if strings.EqualFold(strings.TrimSpace(t), token) || strings.TrimSpace(t) == "*" {
			return true
		}
	}
	return false
}

func AddVary(h http.Header, value string) {
	if !HasVary(h.Get("Vary"), value) {
		h.Add("Vary", value)
	}
}
