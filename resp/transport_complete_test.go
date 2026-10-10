package resp

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ZihxS/RespKit/internal/transport"
)

type ginMinimal struct {
	header  http.Header
	request http.Header
	code    int
	payload any
}

func (g *ginMinimal) JSON(status int, body any) {
	g.code = status
	g.payload = body
}

func (g *ginMinimal) Header(key, value string) {
	g.header.Set(key, value)
}

func (g *ginMinimal) GetHeader(key string) string {
	return g.request.Get(key)
}

type ginNoRequest struct {
	Writer  http.ResponseWriter
	request http.Header
	status  int
}

func (g *ginNoRequest) JSON(code int, body any) {
	g.status = code
}

func (g *ginNoRequest) Header(key, value string) {
	g.Writer.Header().Set(key, value)
}

func (g *ginNoRequest) GetHeader(key string) string {
	return g.request.Get(key)
}

type echoBadResponse1 struct{ request *http.Request }

func (e *echoBadResponse1) Request() *http.Request {
	return e.request
}

func (e *echoBadResponse1) Blob(int, string, []byte) error {
	return nil
}

func (e *echoBadResponse1) Response() int {
	return 1
}

type echoBadResponse2 struct{ echoBadResponse1 }

func (e *echoBadResponse2) Response() *echoResponseDouble {
	return nil
}

type echoBadResponse3 struct{ echoBadResponse1 }

func (e *echoBadResponse3) Response() *badResponseHeader {
	return &badResponseHeader{}
}

type badResponseHeader struct{}

func (*badResponseHeader) Header() int {
	return 0
}

type echoSendError struct {
	request  *http.Request
	response *echoResponseDouble
	last     string
	code     int
}

func (e *echoSendError) Request() *http.Request {
	return e.request
}

func (e *echoSendError) Response() *echoResponseDouble {
	return e.response
}

func (e *echoSendError) Blob(code int, ctype string, body []byte) error {
	e.code = code
	e.last = ctype
	return errors.New("echo: fail write")
}

type fiberNoMethod struct{ fiberDoubleV2 }

func (f *fiberNoMethod) Method(int) int {
	return 2
}

type fiberMissingStatus struct{ fiberDoubleV2 }

// shadowing inherited Status with incorrect signature makes adaptation fail
func (f *fiberMissingStatus) Status(string) *fiberMissingStatus {
	return f
}

type fiberNoStatus struct {
	req http.Header
	res http.Header
}

func (f *fiberNoStatus) Get(key string, _ ...string) string {
	return f.req.Get(key)
}

func (f *fiberNoStatus) Set(key, value string) {
	f.res.Set(key, value)
}

func (f *fiberNoStatus) Send([]byte) error {
	return nil
}

type fiberInvalidStatus struct{ fiberNoStatus }

func (f *fiberInvalidStatus) Status(i int, other int) *fiberInvalidStatus {
	return f
}

type fiberSendErr struct{ fiberDoubleV2 }

func (f *fiberSendErr) Send([]byte) error {
	return errors.New("fiber send failed")
}

type ginStructAbsentWriter struct{ Request *http.Request }

func (g *ginStructAbsentWriter) JSON(int, any) {
}

func (g *ginStructAbsentWriter) Header(string, string) {
}

func (g *ginStructAbsentWriter) GetHeader(string) string {
	return ""
}

type fakeNonStructGin int

func (fakeNonStructGin) JSON(int, any) {
}

func (fakeNonStructGin) Header(string, string) {
}

func (fakeNonStructGin) GetHeader(string) string {
	return ""
}

func TestAllAdapterEdgeCases(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	if _, err := transport.Adapt(nil); !errors.Is(err, ErrUnsupportedContext) {
		t.Fatal(err)
	}
	var ctx *HTTPContext
	if _, err := transport.Adapt(ctx); !errors.Is(err, ErrUnsupportedContext) {
		t.Fatal(err)
	}
	for _, bad := range []*HTTPContext{{Writer: nil, Request: r}, {Writer: w, Request: nil}} {
		if _, err := transport.Adapt(bad); !errors.Is(err, ErrUnsupportedContext) {
			t.Fatal(err)
		}
	}
	var gnil *ginDouble
	if _, err := transport.Adapt(gnil); !errors.Is(err, ErrUnsupportedContext) {
		t.Fatal(err)
	}
	badContexts := []any{
		&echoBadResponse1{r},
		&echoBadResponse2{echoBadResponse1{r}},
		&echoBadResponse3{echoBadResponse1{r}},
		&echoBadResponse1{nil},
		&fiberNoStatus{http.Header{}, http.Header{}},
		&fiberMissingStatus{},
		&fiberInvalidStatus{},
	}
	for _, bad := range badContexts {
		if _, err := transport.Adapt(bad); !errors.Is(err, ErrUnsupportedContext) {
			t.Fatalf("%T: %v", bad, err)
		}
	}
	ginContexts := []any{
		&ginMinimal{header: http.Header{}, request: http.Header{}},
		fakeNonStructGin(1),
		&ginStructAbsentWriter{Request: r},
	}
	for _, c := range ginContexts {
		if err := OK(c, map[string]any{"foo": "bar"}); err != nil {
			t.Fatalf("%T: %v", c, err)
		}
	}
	g := &ginMinimal{header: http.Header{}, request: http.Header{}}
	if err := NoContent(g); err != nil || g.code != 204 {
		t.Fatal(err, g.code)
	}
	if err := Fail(g, ErrNotFound); err != nil || !transport.HasVary(g.header.Get("Vary"), "Accept-Language") {
		t.Fatal(err, g.header)
	}
	g2 := &ginNoRequest{Writer: httptest.NewRecorder(), request: http.Header{}}
	if err := OK(g2, 1); err != nil || g2.status != 200 {
		t.Fatal(err)
	}
	if _, ok := transport.GinRequest(fakeNonStructGin(1)); ok {
		t.Fatal("invalid")
	}
	if _, ok := transport.GinRequest(&ginStructAbsentWriter{Request: nil}); !ok {
		t.Fatal("present nil field reports shape")
	}
	if _, ok := transport.GinRequest(&ginMinimal{}); ok {
		t.Fatal("no Request field")
	}
	if _, ok := transport.GinWriterHeader(fakeNonStructGin(1)); ok {
		t.Fatal("invalid")
	}
	if _, ok := transport.GinWriterHeader(&ginStructAbsentWriter{}); ok {
		t.Fatal("no Writer field")
	}
	if _, ok := transport.GinWriterHeader(&ginDouble{Writer: nil}); ok {
		t.Fatal("nil Writer")
	}
	if _, ok := transport.GinWriterHeader(&ginDouble{Writer: httptest.NewRecorder()}); !ok {
		t.Fatal("valid Writer")
	}
	var nilGin *ginDouble
	if _, ok := transport.GinWriterHeader(nilGin); ok {
		t.Fatal("nil gin")
	}
	if _, ok := transport.GinRequest(nilGin); ok {
		t.Fatal("nil gin request")
	}
	if _, ok := transport.GinRequest(23); ok {
		t.Fatal("gin request int")
	}
	e := &echoDouble{request: r, response: &echoResponseDouble{header: http.Header{}}}
	if err := NoContent(e); err != nil || e.code != 204 || len(e.body) != 0 {
		t.Fatal(err)
	}
	e2 := &echoSendError{request: r, response: &echoResponseDouble{header: http.Header{}}}
	if err := OK(e2, 23); err == nil {
		t.Fatal("echo must propagate error")
	}
	if err := NoContent(e2); err == nil || e2.last != "" {
		t.Fatal("echo no-content error", err)
	}
	echoContexts := []any{
		nil,
		0,
		&echoBadResponse1{r},
		&echoBadResponse2{echoBadResponse1{r}},
		&echoBadResponse3{echoBadResponse1{r}},
	}
	for _, v := range echoContexts {
		if _, ok := transport.EchoResponseHeader(v); ok {
			t.Fatalf("unexpected echo response %T", v)
		}
	}
	f := &fiberNoMethod{fiberDoubleV2{req: http.Header{}, res: http.Header{}}}
	if err := OK(f, 4); err != nil {
		t.Fatal(err)
	}
	if _, err := transport.Adapt(&fiberNoStatus{req: http.Header{}, res: http.Header{}}); err == nil {
		t.Fatal("bad status accepted")
	}
	ferr := &fiberSendErr{fiberDoubleV2{req: http.Header{}, res: http.Header{}}}
	if err := OK(ferr, 23); err == nil || !strings.Contains(err.Error(), "response write") {
		t.Fatal("fiber failure", err)
	}
	for _, val := range []string{"*", "Origin, Accept-Language"} {
		if !transport.HasVary(val, "Accept-Language") {
			t.Fatal(val)
		}
	}
}

func TestDefaultEngineColdStart(t *testing.T) {
	old := defaultEngine.Load()
	defer defaultEngine.Store(old)
	defaultEngine.Store(old)
	p := current()
	if p == nil || p.naming != SnakeCase {
		t.Fatal(p)
	}
	if current() != p {
		t.Fatal("default not memoized")
	}
}

// Pointer-to-non-struct values exercise reflective fallback without framework imports.
type pointerNonStruct int

func (*pointerNonStruct) JSON(int, any) {
}

func (*pointerNonStruct) Header(string, string) {
}

func (*pointerNonStruct) GetHeader(string) string {
	return ""
}

type fiberNoMethodAtAll struct {
	req http.Header
	res http.Header
}

func (f *fiberNoMethodAtAll) Get(key string, _ ...string) string {
	return f.req.Get(key)
}

func (f *fiberNoMethodAtAll) Set(key, value string) {
	f.res.Set(key, value)
}

func (f *fiberNoMethodAtAll) Send([]byte) error {
	return nil
}

func (f *fiberNoMethodAtAll) Status(int) *fiberNoMethodAtAll {
	return f
}

func TestReflectiveUnreachableViaAdapterDirectly(t *testing.T) {
	ptr := pointerNonStruct(7)
	if _, ok := transport.GinWriterHeader(&ptr); ok {
		t.Fatal("pointer int shouldn't be writer")
	}
	if _, ok := transport.GinRequest(&ptr); ok {
		t.Fatal("pointer int shouldn't be request")
	}
	f := &fiberNoMethodAtAll{req: http.Header{}, res: http.Header{}}
	if m := (transport.FiberAdapter{C: f, Raw: f}).Method(); m != "" {
		t.Fatal(m)
	}
	badStatus := &fiberNoStatus{req: http.Header{}, res: http.Header{}}
	adapter := transport.FiberAdapter{C: f, Raw: badStatus}
	if err := adapter.Write(204, nil, ""); !errors.Is(err, ErrUnsupportedContext) {
		t.Fatal(err)
	}
}
