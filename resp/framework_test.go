package resp

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ZihxS/RespKit/internal/transport"
)

// These are contract doubles with the method signatures of Gin, Echo, and
// Fiber. Real upstream integration builds should also run in a CI matrix.
type ginDouble struct {
	Request *http.Request
	Writer  http.ResponseWriter
	code    int
}

func (c *ginDouble) JSON(code int, v any) {
	c.code = code
	c.Writer.WriteHeader(code)
	b, _ := json.Marshal(v)
	_, _ = c.Writer.Write(b)
}

func (c *ginDouble) Header(key, value string) {
	c.Writer.Header().Set(key, value)
}

func (c *ginDouble) GetHeader(key string) string {
	return c.Request.Header.Get(key)
}

func (c *ginDouble) Status(code int) {
	c.code = code
	c.Writer.WriteHeader(code)
}

type echoResponseDouble struct{ header http.Header }

func (r *echoResponseDouble) Header() http.Header {
	return r.header
}

type echoDouble struct {
	request  *http.Request
	response *echoResponseDouble
	code     int
	body     []byte
	ctype    string
}

func (c *echoDouble) Blob(code int, ctype string, buf []byte) error {
	c.code = code
	c.ctype = ctype
	c.body = buf
	return nil
}

func (c *echoDouble) Request() *http.Request {
	return c.request
}

func (c *echoDouble) Response() *echoResponseDouble {
	return c.response
}

type fiberDoubleV2 struct {
	req    http.Header
	res    http.Header
	status int
	body   []byte
}

func (c *fiberDoubleV2) Get(key string, defaultValue ...string) string {
	return c.req.Get(key)
}

func (c *fiberDoubleV2) Set(key, value string) {
	c.res.Set(key, value)
}

func (c *fiberDoubleV2) Send(body []byte) error {
	c.body = body
	return nil
}

func (c *fiberDoubleV2) Status(status int) *fiberDoubleV2 {
	c.status = status
	return c
}

func (c *fiberDoubleV2) Method(_ ...string) string {
	return "GET"
}

func (c *fiberDoubleV2) GetRespHeader(key string, _ ...string) string {
	return c.res.Get(key)
}

type fiberDoubleV3 struct{ fiberDoubleV2 }

func (c *fiberDoubleV3) Status(status int) *fiberDoubleV3 {
	c.status = status
	return c
}

func (c *fiberDoubleV3) Method() string {
	return "GET"
}

func TestGinStructuralIntegration(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Accept-Language", "en")
	recorder := httptest.NewRecorder()
	recorder.Header().Set("Vary", "Origin")
	c := &ginDouble{Request: req, Writer: recorder}
	if err := OK(c, 123); err != nil {
		t.Fatal(err)
	}
	vary := strings.Join(recorder.Header().Values("Vary"), ",")
	if recorder.Code != 200 || !transport.HasVary(vary, "Origin") || !transport.HasVary(vary, "Accept-Language") {
		t.Fatal(recorder.Code, recorder.Header())
	}
	c = &ginDouble{Request: req, Writer: httptest.NewRecorder()}
	if err := NoContent(c); err != nil || c.code != 204 {
		t.Fatal(err, c.code)
	}
}

func TestEchoStructuralIntegration(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Accept-Language", "en")
	c := &echoDouble{request: req, response: &echoResponseDouble{header: http.Header{}}}
	if err := Fail(c, ErrConflict); err != nil {
		t.Fatal(err)
	}
	if c.code != 409 || c.ctype != "application/json; charset=utf-8" || c.response.header.Get("Content-Language") != "en" {
		t.Fatal(c)
	}
	var b map[string]any
	if err := json.Unmarshal(c.body, &b); err != nil {
		t.Fatal(err)
	}
	if b["error"].(map[string]any)["code"] != "RESOURCE_CONFLICT" {
		t.Fatal(b)
	}
}

func TestFiberStructuralIntegration(t *testing.T) {
	c := &fiberDoubleV2{req: http.Header{}, res: http.Header{}}
	c.req.Set("Accept-Language", "en")
	c.res.Set("Vary", "Origin")
	if err := Fail(c, ErrNotFound); err != nil {
		t.Fatal(err)
	}
	if c.status != 404 || c.res.Get("Content-Language") != "en" || !transport.HasVary(c.res.Get("Vary"), "Origin") {
		t.Fatal(c)
	}
	d := &fiberDoubleV3{fiberDoubleV2{req: http.Header{}, res: http.Header{}}}
	if err := OK(d, map[string]any{"hello": "world"}); err != nil {
		t.Fatal(err)
	}
	if d.status != 200 || len(d.body) == 0 {
		t.Fatal(d)
	}
	if err := NoContent(d); err != nil || d.status != 204 || len(d.body) != 0 {
		t.Fatal(err, d)
	}
}
