package transport

import (
	"net/http"
)

// Echo Context provides Blob(code, contentType, body), avoiding double encoding.
type echoLike interface {
	Blob(int, string, []byte) error
	Request() *http.Request
}

type echoAdapter struct {
	c   echoLike
	raw any
}

func (a echoAdapter) Header(key string) string {
	return a.c.Request().Header.Get(key)
}

func (a echoAdapter) SetHeader(k, v string) {
	if h, ok := EchoResponseHeader(a.raw); ok {
		h.Set(k, v)
	}
}

func (a echoAdapter) AddVary(v string) {
	if h, ok := EchoResponseHeader(a.raw); ok {
		AddVary(h, v)
	}
}

func (a echoAdapter) Method() string {
	return a.c.Request().Method
}

func (a echoAdapter) Write(code int, body []byte, ctype string) error {
	if body == nil { // RFC 9110: HEAD and 204 must not contain a response body.
		if h, ok := EchoResponseHeader(a.raw); ok {
			h.Del("Content-Type")
		}
		return a.c.Blob(code, "", nil)
	}
	return a.c.Blob(code, ctype, body)
}
