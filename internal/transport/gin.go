package transport

import (
	"encoding/json"
)

type ginLike interface {
	JSON(int, any)
	Header(string, string)
	GetHeader(string) string
}

type ginAdapter struct {
	c   ginLike
	raw any
}

func (a ginAdapter) Header(key string) string {
	return a.c.GetHeader(key)
}

func (a ginAdapter) SetHeader(key, value string) {
	a.c.Header(key, value)
}

func (a ginAdapter) AddVary(v string) {
	if h, ok := GinWriterHeader(a.raw); ok {
		AddVary(h, v)
		return
	}
	a.c.Header("Vary", v)
}

func (a ginAdapter) Method() string {
	if r, ok := GinRequest(a.raw); ok {
		return r.Method
	}
	return ""
}

func (a ginAdapter) Write(code int, body []byte, ctype string) error {
	if body == nil {
		if status, ok := a.raw.(interface{ Status(int) }); ok {
			status.Status(code)
			return nil
		}
	}
	if ctype != "" {
		a.c.Header("Content-Type", ctype)
	}
	if body == nil {
		body = []byte{}
	}
	a.c.JSON(code, json.RawMessage(body))
	return nil
}
