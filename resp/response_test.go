package resp

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/ZihxS/RespKit/internal/transport"
)

func request(method, lang string) (*httptest.ResponseRecorder, *http.Request, *HTTPContext) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(method, "/users", nil)
	if lang != "" {
		r.Header.Set("Accept-Language", lang)
	}
	return w, r, HTTP(w, r)
}

func decode(t *testing.T, rec *httptest.ResponseRecorder) map[string]any {
	t.Helper()
	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("bad JSON %q: %v", rec.Body.String(), err)
	}
	return body
}

func TestDefaultSuccessAndLocalization(t *testing.T) {
	w, r, c := request("GET", "en-US,en;q=0.9")
	r.Header.Set("X-Request-ID", "req_ok.123")
	if err := OK(c, struct {
		FullName string `json:"full_name"`
	}{"Saleh"}); err != nil {
		t.Fatal(err)
	}
	if w.Code != 200 {
		t.Fatal(w.Code)
	}
	got := decode(t, w)
	if got["success"] != true {
		t.Fatal(got)
	}
	if got["data"].(map[string]any)["full_name"] != "Saleh" {
		t.Fatal(got)
	}
	if got["meta"].(map[string]any)["request_id"] != "req_ok.123" {
		t.Fatal(got)
	}
	if w.Header().Get("Content-Language") != "en" {
		t.Fatal(w.Header())
	}
	if w.Header().Get("Vary") != "Accept-Language" {
		t.Fatal(w.Header())
	}
}

func TestNamingDoesNotChangeBusinessData(t *testing.T) {
	for _, tc := range []struct {
		name       Naming
		requestKey string
		successKey string
	}{
		{SnakeCase, "request_id", "success"}, {CamelCase, "requestId", "success"}, {PascalCase, "RequestId", "Success"},
	} {
		t.Run(string(tc.name), func(t *testing.T) {
			e, err := New(WithNaming(tc.name))
			if err != nil {
				t.Fatal(err)
			}
			w, r, c := request("GET", "")
			r.Header.Set("X-Request-ID", "req_123")
			if err := e.OK(c, map[string]any{"custom_field": "hello"}); err != nil {
				t.Fatal(err)
			}
			b := decode(t, w)
			if b[tc.successKey] != true {
				t.Fatal(b)
			}
			dataKey := "data"
			metaKey := "meta"
			if tc.name == PascalCase {
				dataKey = "Data"
				metaKey = "Meta"
			}
			if b[dataKey].(map[string]any)["custom_field"] != "hello" {
				t.Fatal(b)
			}
			if b[metaKey].(map[string]any)[tc.requestKey] != "req_123" {
				t.Fatal(b)
			}
		})
	}
}

func TestErrorMappingAndNoLeak(t *testing.T) {
	tests := []struct {
		name   string
		err    error
		status int
		code   string
		mapper string
	}{
		{"sentinel", fmt.Errorf("lookup: %w", ErrNotFound), 404, "RESOURCE_NOT_FOUND", "builtin.sentinel"},
		{"syntax", &json.SyntaxError{Offset: 3}, 400, "BAD_REQUEST", "builtin.json"},
		{"unique", fmt.Errorf("db: %w", sqlstateError{"23505"}), 409, "RESOURCE_ALREADY_EXISTS", "builtin.sqlstate_unique"},
		{"ambiguous sqlstate", sqlstateError{"23503"}, 500, "INTERNAL_ERROR", "fallback"},
		{"sql no rows", sql.ErrNoRows, 500, "INTERNAL_ERROR", "fallback"},
		{"deadline", context.DeadlineExceeded, 500, "INTERNAL_ERROR", "fallback"},
		{"unknown", errors.New("db password secret=xyz"), 500, "INTERNAL_ERROR", "fallback"},
	}
	e, _ := New()
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w, _, c := request("GET", "id-ID")
			if err := e.Fail(c, tc.err); err != nil {
				t.Fatal(err)
			}
			if w.Code != tc.status {
				t.Fatalf("status=%d want %d", w.Code, tc.status)
			}
			b := decode(t, w)
			if code := b["error"].(map[string]any)["code"]; code != tc.code {
				t.Fatal(b)
			}
			if strings.Contains(w.Body.String(), "secret=") || strings.Contains(w.Body.String(), "23505") {
				t.Fatal("internal error leaked")
			}
			d := e.Inspect(tc.err)
			if d.Mapper != tc.mapper {
				t.Fatalf("mapper = %s want %s", d.Mapper, tc.mapper)
			}
		})
	}
}

type sqlstateError struct{ code string }

func (e sqlstateError) Error() string {
	return "SECRET SQLSTATE=" + e.code
}

func (e sqlstateError) SQLState() string {
	return e.code
}

func TestValidationLocalizationAndHeaders(t *testing.T) {
	e, _ := New()
	app := NewError("VALIDATION_FAILED", 422, WithDetails(
		FieldError{Field: "email", Code: "INVALID_EMAIL"},
		FieldError{Field: "password", Code: "MIN_LENGTH", Params: map[string]string{"min": "8"}},
	))
	w, _, ctx := request("POST", "en-AU, en;q=0.9")
	if err := e.Fail(ctx, app); err != nil {
		t.Fatal(err)
	}
	b := decode(t, w)
	details := b["error"].(map[string]any)["details"].([]any)
	if details[0].(map[string]any)["message"] != "The email address is invalid." {
		t.Fatal(details)
	}
	if details[1].(map[string]any)["message"] != "The password field must be at least 8 characters." {
		t.Fatal(details)
	}
	if w.Header().Get("Content-Language") != "en" {
		t.Fatal(w.Header())
	}
}

func TestLocalePriorityAndFallback(t *testing.T) {
	e, _ := New(WithTranslations(map[string]map[string]string{"ja": {"INTERNAL_ERROR": "サーバーエラー"}}))
	tests := []struct{ header, force, want string }{
		{"en-US,en;q=0.9", "", "en"}, {"en;q=0, id;q=1", "", "id"},
		{"ja-JP,ja;q=0.9", "", "ja"}, {"es-ES", "", "id"},
		{"id;q=0.2,en;q=0.8", "", "en"}, {"en", "id", "id"},
	}
	for _, tc := range tests {
		w, _, c := request("GET", tc.header)
		if err := e.Fail(c, errors.New("internal"), Locale(tc.force)); err != nil {
			t.Fatal(err)
		}
		if have := w.Header().Get("Content-Language"); have != tc.want {
			t.Errorf("%q: %q want %q", tc.header, have, tc.want)
		}
	}
}

func TestRFC9457(t *testing.T) {
	e, _ := New(WithErrorFormat(ProblemDetails), WithNaming(PascalCase))
	w, r, c := request("GET", "en")
	r.Header.Set("X-Request-ID", "req_1")
	if err := e.Fail(c, ErrNotFound); err != nil {
		t.Fatal(err)
	}
	if w.Code != 404 || w.Header().Get("Content-Type") != "application/problem+json" {
		t.Fatal(w.Code, w.Header())
	}
	b := decode(t, w)
	if b["type"] != "about:blank" || b["status"] != float64(404) || b["request_id"] != "req_1" {
		t.Fatal(b)
	}
	if _, exists := b["Success"]; exists {
		t.Fatal("envelope incorrectly included")
	}
}

func TestHTTPInvariants(t *testing.T) {
	w, _, c := request("GET", "")
	if err := NoContent(c); err != nil {
		t.Fatal(err)
	}
	if w.Code != 204 || w.Body.Len() != 0 {
		t.Fatalf("got %d %s", w.Code, w.Body.String())
	}
	w, _, c = request("HEAD", "en")
	if err := OK(c, 123); err != nil {
		t.Fatal(err)
	}
	if w.Body.Len() != 0 {
		t.Fatal("HEAD has body")
	}
	w, _, c = request("POST", "")
	if err := Created(c, 123, Location("/users/123")); err != nil {
		t.Fatal(err)
	}
	if w.Code != 201 || w.Header().Get("Location") != "/users/123" {
		t.Fatal(w.Code, w.Header())
	}
	w, _, c = request("GET", "")
	if err := Fail(c, ErrUnauthorized); err != nil {
		t.Fatal(err)
	}
	if w.Code != 401 || w.Header().Get("WWW-Authenticate") != "Bearer" {
		t.Fatal(w.Code, w.Header())
	}
	w, _, c = request("GET", "")
	err := Fail(c, NewError("UNAUTHORIZED", 401, WithHeader("WWW-Authenticate", "Basic realm=\"example\"")))
	if err != nil || w.Header().Get("WWW-Authenticate") != "Basic realm=\"example\"" {
		t.Fatal(err, w.Header())
	}
}

func TestInvalidDataFallback(t *testing.T) {
	w, _, c := request("GET", "")
	if err := OK(c, map[string]any{"invalid": make(chan int)}); err != nil {
		t.Fatal(err)
	}
	if w.Code != 500 || strings.Contains(w.Body.String(), "invalid") {
		t.Fatal(w.Code, w.Body.String())
	}
	if b := decode(t, w); b["error"].(map[string]any)["code"] != "INTERNAL_ERROR" {
		t.Fatal(b)
	}
}

func TestCustomMapperAndInvalidConfiguration(t *testing.T) {
	e, err := New(WithMapper("domain", func(err error) (Definition, bool) {
		if errors.Is(err, ErrConflict) {
			return Definition{Code: "CUSTOM_CONFLICT", Status: 409}, true
		}
		return Definition{}, false
	}), WithTranslations(map[string]map[string]string{"en": {"CUSTOM_CONFLICT": "Domain conflict."}}))
	if err != nil {
		t.Fatal(err)
	}
	w, _, c := request("GET", "en")
	if err := e.Fail(c, fmt.Errorf("wrapped: %w", ErrConflict)); err != nil {
		t.Fatal(err)
	}
	if b := decode(t, w); b["error"].(map[string]any)["message"] != "Domain conflict." {
		t.Fatal(b)
	}
	if _, err := New(WithNaming("invalid")); err == nil {
		t.Fatal("invalid naming accepted")
	}
	if _, err := New(WithErrorFormat("invalid")); err == nil {
		t.Fatal("invalid format accepted")
	}
	if _, err := New(WithDefaultLocale("ja")); err == nil {
		t.Fatal("missing catalog accepted")
	}
}

func TestCanceledAndUnsupported(t *testing.T) {
	w, _, c := request("GET", "")
	if err := Fail(c, fmt.Errorf("request: %w", context.Canceled)); !errors.Is(err, ErrRequestCanceled) {
		t.Fatal(err)
	}
	if w.Body.Len() != 0 {
		t.Fatal("sent canceled response")
	}
	if err := OK(struct{}{}, 123); !errors.Is(err, ErrUnsupportedContext) {
		t.Fatal(err)
	}
}

func TestRequestIDSanitization(t *testing.T) {
	w, r, c := request("GET", "")
	r.Header.Set("X-Request-ID", "\"inject\nnot-safe")
	_ = OK(c, 5)
	if _, ok := decode(t, w)["meta"]; ok {
		t.Fatal("untrusted request id leaked")
	}
}

func TestConcurrentEngineRequests(t *testing.T) {
	e, _ := New()
	var wg sync.WaitGroup
	for i := range 50 {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			locale := "en"
			if i%2 == 0 {
				locale = "id"
			}
			w, _, c := request("GET", locale)
			if err := e.Fail(c, ErrNotFound); err != nil {
				t.Error(err)
				return
			}
			if w.Header().Get("Content-Language") != locale {
				t.Error("language leaked between requests")
			}
		}(i)
	}
	wg.Wait()
}

func TestCatalog(t *testing.T) {
	c := Catalog()
	if len(c) < 8 {
		t.Fatal(c)
	}
	codes := map[string]bool{}
	for _, v := range c {
		if codes[v.Code] {
			t.Fatal("duplicate code")
		}
		codes[v.Code] = true
	}
	if !codes["RESOURCE_NOT_FOUND"] {
		t.Fatal(c)
	}
}

func TestReflectHasVary(t *testing.T) {
	if !transport.HasVary("Origin, Accept-Language", "Accept-Language") || transport.HasVary("Accept", "Accept-Language") {
		t.Fatal("vary comparison")
	}
}

func TestJoinedErrorsAndInvalidMapper(t *testing.T) {
	e, err := New(WithMapper("bad", func(err error) (Definition, bool) {
		return Definition{Code: "lowercase", Status: 200}, true
	}))
	if err != nil {
		t.Fatal(err)
	}
	d, diag := e.Resolve(errors.Join(ErrNotFound, errors.New("sensitive")))
	if d.Code != "INTERNAL_ERROR" || diag.Detected || diag.Note == "" {
		t.Fatal(d, diag)
	}
	normal, _ := New()
	d, _ = normal.Resolve(errors.Join(ErrNotFound, errors.New("sensitive")))
	if d.Code != "RESOURCE_NOT_FOUND" {
		t.Fatal(d)
	}
}

func TestPartialLocaleCatalogFallsBackForEntireError(t *testing.T) {
	e, _ := New(WithTranslations(map[string]map[string]string{"ja": {"INTERNAL_ERROR": "サーバーエラー"}}))
	w, _, c := request("GET", "ja-JP")
	if err := e.Fail(c, ErrNotFound); err != nil {
		t.Fatal(err)
	}
	if w.Header().Get("Content-Language") != "id" {
		t.Fatal(w.Header())
	}
	b := decode(t, w)
	if b["error"].(map[string]any)["message"] != "Data tidak ditemukan." {
		t.Fatal(b)
	}
}

func TestUnknownDomainErrorIsSafelyLocalizedByHTTPStatus(t *testing.T) {
	w, _, c := request("GET", "en")
	if err := Fail(c, NewError("NEW_BUSINESS_RULE", 403, WithCause(errors.New("sensitive account token")))); err != nil {
		t.Fatal(err)
	}
	b := decode(t, w)
	if w.Code != 403 || b["error"].(map[string]any)["code"] != "NEW_BUSINESS_RULE" {
		t.Fatal(b)
	}
	message := b["error"].(map[string]any)["message"]
	if strings.Contains(w.Body.String(), "token") || message == "An internal server error occurred." {
		t.Fatal(b)
	}
}

func BenchmarkNetHTTP_OK(b *testing.B) {
	e, _ := New()
	data := struct {
		UserID int `json:"user_id"`
	}{UserID: 123}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/users", nil)
		if err := e.OK(HTTP(w, r), data); err != nil {
			b.Fatal(err)
		}
	}
}
