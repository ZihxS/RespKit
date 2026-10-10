package resp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func TestExposeInternalErrorOptIn(t *testing.T) {
	secret := "sql: user password=supersecret"
	for _, tc := range []struct {
		name    string
		opts    []Option
		enabled bool
		format  ErrorFormat
	}{
		{"default", nil, false, Envelope},
		{"explicit-false", []Option{WithExposeInternalErrors(false)}, false, Envelope},
		{"enabled", []Option{WithExposeInternalErrors(true)}, true, Envelope},
		{"problem", []Option{WithExposeInternalErrors(true), WithErrorFormat(ProblemDetails)}, true, ProblemDetails},
	} {
		t.Run(tc.name, func(t *testing.T) {
			e, err := New(tc.opts...)
			if err != nil {
				t.Fatal(err)
			}
			w, _, c := request("GET", "en")
			if err := e.Fail(c, fmt.Errorf("write failed: %w", errors.New(secret))); err != nil {
				t.Fatal(err)
			}
			b := decode(t, w)
			target := b
			if tc.format == Envelope {
				target = b["error"].(map[string]any)
			}
			raw, ok := target["debug_message"].(string)
			if ok != tc.enabled || (tc.enabled && !strings.Contains(raw, secret)) {
				t.Fatalf("%s: %v", tc.name, b)
			}
			if tc.enabled && target["message"] == secret {
				t.Fatal("localized message overwritten")
			}
			if !tc.enabled && strings.Contains(w.Body.String(), secret) {
				t.Fatal("internal data leak")
			}
			if w.Code != 500 {
				t.Fatal(w.Code)
			}
		})
	}
	e, _ := New(WithExposeInternalErrors(true))
	w, _, c := request("GET", "id")
	if err := e.Fail(c, errors.New(strings.Repeat("A", 9000))); err != nil {
		t.Fatal(err)
	}
	b := decode(t, w)
	if len(b["error"].(map[string]any)["debug_message"].(string)) != 8192 {
		t.Fatal("debug message must be truncated")
	}
	for _, enabled := range []bool{false, true} {
		e, _ := New(WithExposeInternalErrors(enabled))
		w, _, c := request("GET", "en")
		if err := e.OK(c, map[string]any{"invalid": make(chan int)}); err != nil {
			t.Fatal(err)
		}
		b := decode(t, w)
		_, ok := b["error"].(map[string]any)["debug_message"]
		if ok != enabled {
			t.Fatalf("marshal debug mismatch: %v", b)
		}
	}
	e, _ = New(WithExposeInternalErrors(true), WithErrorFormat(ProblemDetails))
	w, _, c = request("HEAD", "en")
	if err := e.OK(c, make(chan int)); err != nil {
		t.Fatal(err)
	}
	if w.Code != 500 || w.Body.Len() != 0 {
		t.Fatal(w.Code, w.Body)
	}
	w, _, c = request("GET", "en")
	if err := e.OK(c, make(chan int)); err != nil {
		t.Fatal(err)
	}
	b = decode(t, w)
	if !strings.Contains(b["debug_message"].(string), "unsupported type") {
		t.Fatal(b)
	}
}

func TestOptionsAndErrorValuesComplete(t *testing.T) {
	if _, err := New(nil); err != nil {
		t.Fatal(err)
	}
	if _, err := New(func(e *Engine) error {
		return errors.New("boom")
	}); err == nil {
		t.Fatal("expected option error")
	}
	if _, err := New(WithDefaultLocale("  ")); err == nil {
		t.Fatal("empty locale")
	}
	if _, err := New(WithTranslations(map[string]map[string]string{"??": {"AAA": "abc"}})); err == nil {
		t.Fatal("bad locale")
	}
	if _, err := New(WithTranslations(map[string]map[string]string{"fr": {"a": "abc"}})); err == nil {
		t.Fatal("bad error code")
	}
	if _, err := New(WithTranslations(map[string]map[string]string{"fr": {"FOO": "  "}})); err == nil {
		t.Fatal("empty message")
	}
	e, err := New(
		WithTranslations(map[string]map[string]string{
			"ja-JP": {"INTERNAL_ERROR": "error-ja"},
		}),
		WithDefaultLocale("JA-jp"),
	)
	if err != nil || e.defaultLocale != "ja" {
		t.Fatal(err, e)
	}
	if _, err := New(WithMapper("", func(error) (Definition, bool) {
		return Definition{}, false
	})); err == nil {
		t.Fatal("empty name")
	}
	if _, err := New(WithMapper("foo", nil)); err == nil {
		t.Fatal("nil mapper")
	}
	mapper := func(error) (Definition, bool) {
		return Definition{}, false
	}
	if _, err := New(WithMapper("dupe", mapper), WithMapper("dupe", mapper)); err == nil {
		t.Fatal("duplicate mapper")
	}
	if err := Configure(WithNaming("invalid")); err == nil {
		t.Fatal("invalid configure")
	}
	if err := Configure(WithNaming(CamelCase), WithExposeInternalErrors(true)); err != nil {
		t.Fatal(err)
	}
	if current().naming != CamelCase || !current().exposeInternal {
		t.Fatal("config not used")
	}
	if err := Configure(); err != nil {
		t.Fatal(err)
	}
	if current().naming != SnakeCase || current().exposeInternal {
		t.Fatal("default reset failed")
	}
	if current() == nil {
		t.Fatal("no default")
	}

	var nilErr *AppError
	if nilErr.Error() != "internal error" || nilErr.Unwrap() != nil {
		t.Fatal("nil app error")
	}
	if app := NewError("CODE", 409); app.Error() != "CODE" || app.Unwrap() != nil {
		t.Fatal("app without cause")
	}
	app := NewError("CODE", 409, WithCause(errors.New("detail")), nil)
	if app.Error() != "CODE: detail" || app.Unwrap() == nil {
		t.Fatal("app with cause")
	}
	if got := Inspect(ErrNotFound); got.ErrorCode != "RESOURCE_NOT_FOUND" {
		t.Fatal(got)
	}
	if !strings.Contains(Explain(ErrNotFound), "RESOURCE_NOT_FOUND") {
		t.Fatal("explain")
	}
	if got := Inspect(ErrNotFound).String(); !strings.Contains(got, "404") {
		t.Fatal(got)
	}
	w, _, c := request("GET", "")
	if err := NotFound(c); err != nil || w.Code != 404 {
		t.Fatal(err, w.Code)
	}
}

func TestLocaleMatchingAndMissingTranslations(t *testing.T) {
	e, _ := New()
	cases := map[string]string{
		"en-US,en;q=0.8":                     "en",
		"id-ID,en;q=1":                       "id",
		"en;q=0,id;q=1":                      "id",
		"en;q=bad,id;q=0.6":                  "id",
		"en;q=1.3,id;q=0.5":                  "id",
		"x_invalid-language-string,en;q=0.8": "en",
		"*":                                  "id",
		"fr, en;q=0.8":                       "en",
		"xx,yy":                              "id",
		"en;q=0.3,id;q=0.9":                  "id",
		"en;q=-1,id":                         "id",
		"en;q=0.7, id;q=0.7":                 "en",
	}
	for header, want := range cases {
		if got := e.chooseLocale(header, ""); got != want {
			t.Errorf("%q got %q want %q", header, got, want)
		}
	}
	if got := e.chooseLocale("en", "de"); got != "id" {
		t.Fatal(got)
	}
	if got := e.chooseLocale("en", "EN_us"); got != "en" {
		t.Fatal(got)
	}
	if normalizeLocale(strings.Repeat("a", 65)) != "" || normalizeLocale("?bad") != "" || normalizeLocale("   ") != "" {
		t.Fatal("invalid locales accepted")
	}
	message := e.translate("en", "FIELD_REQUIRED", map[string]string{"field": "email"})
	if message != "The email field is required." {
		t.Fatal(message)
	}
	if got := e.translate("en", "NOT_IN_CATALOG", nil); got != "An internal server error occurred." {
		t.Fatal(got)
	}
	if got := e.translate("de", "RESOURCE_NOT_FOUND", nil); got != "Data tidak ditemukan." {
		t.Fatal(got)
	}
	if got := e.translate("de", "UNREGISTERED", nil); got != "Terjadi kesalahan pada server." {
		t.Fatal(got)
	}
	if got := e.chooseLocale(strings.Repeat("unknown,", 33)+"en", ""); got != "id" {
		t.Fatal(got)
	}
	if e.effectiveLocale("en", nil) != "en" || e.effectiveLocale("id", &Definition{Code: "NOT_IN_CATALOG"}) != "id" {
		t.Fatal("wrong locale")
	}
	d := &Definition{Code: "RESOURCE_NOT_FOUND", Details: []FieldError{{Code: "MISSING_JA_KEY"}}}
	if e.effectiveLocale("en", d) != "en" {
		t.Fatal("unknown detail should use generic english fallback")
	}
	ei, _ := New(WithTranslations(map[string]map[string]string{"id": {"ID_ONLY_DETAIL": "Hanya id"}}))
	d.Details[0].Code = "ID_ONLY_DETAIL"
	if ei.effectiveLocale("en", d) != "id" {
		t.Fatal("partial locale must fallback to id")
	}
	if got := e.publicMessage("en", "NEW_UNKNOWN", 502); got != "An internal server error occurred." {
		t.Fatal(got)
	}
	if got := e.publicMessage("en", "NEW_UNKNOWN", 400); got != "The request is invalid." {
		t.Fatal(got)
	}
	if got := e.fieldMessage("en", "NOT_EXIST", nil); got != "The submitted data is invalid." {
		t.Fatal(got)
	}
	message = e.fieldMessage("en", "FIELD_REQUIRED", map[string]string{"field": "name"})
	if message != "The name field is required." {
		t.Fatal(message)
	}
}

func TestResolutionAdditionalCases(t *testing.T) {
	e, _ := New(WithMapper("unused", func(error) (Definition, bool) {
		return Definition{}, false
	}))
	d, diag := e.Resolve(nil)
	if d.Code != "INTERNAL_ERROR" || diag.Mapper != "builtin.sentinel" {
		t.Fatal(d, diag)
	}
	d, diag = e.Resolve(&json.UnmarshalTypeError{Value: "array", Type: reflect.TypeFor[int]()})
	if d.Status != 400 || diag.Mapper != "builtin.json" {
		t.Fatal(d, diag)
	}
	d, diag = e.Resolve(context.Canceled)
	if d.Code != "INTERNAL_ERROR" || diag.Note == "" {
		t.Fatal(d, diag)
	}
	d, diag = e.Resolve(NewError("invalid", 200))
	if d.Code != "INTERNAL_ERROR" || diag.Note == "" {
		t.Fatal(d, diag)
	}
	e, _ = New(WithMapper("yes", func(error) (Definition, bool) {
		return Definition{Code: "CUSTOM_ERR", Status: 409}, true
	}))
	d, diag = e.Resolve(errors.New("anything"))
	if d.Code != "CUSTOM_ERR" || diag.Mapper != "custom.yes" {
		t.Fatal(d, diag)
	}
	e, _ = New(WithMapper("invalid", func(error) (Definition, bool) {
		return Definition{Code: "VALID", Status: 600}, true
	}))
	d, diag = e.Resolve(errors.New("anything"))
	if d.Code != "INTERNAL_ERROR" || diag.Detected {
		t.Fatal(d, diag)
	}
	d, diag = NewEngineResolve(t, sqlstateError{"23503"})
	if d.Status != 500 || diag.SourceCode != "23503" {
		t.Fatal(d, diag)
	}
}

func NewEngineResolve(t *testing.T, err error) (Definition, Diagnostic) {
	t.Helper()
	e, eerr := New()
	if eerr != nil {
		t.Fatal(eerr)
	}
	return e.Resolve(err)
}

func TestRemainingResponseBranches(t *testing.T) {
	e, _ := New(WithNaming(PascalCase), WithErrorFormat(ProblemDetails))
	w, r, c := request("GET", "en")
	r.Header.Set("X-Request-ID", "req_777")
	if err := e.OK(c, make(chan int)); err != nil {
		t.Fatal(err)
	}
	if b := decode(t, w); b["request_id"] != "req_777" {
		t.Fatal(b)
	}
	w, r, c = request("GET", "en")
	r.Header.Set("X-Request-ID", "req_777")
	if err := e.Fail(c, ErrNotFound); err != nil {
		t.Fatal(err)
	}
	if b := decode(t, w); b["request_id"] != "req_777" {
		t.Fatal(b)
	}
	e, _ = New()
	w, _, c = request("GET", "en")
	validationError := NewError(
		"VALIDATION_FAILED",
		422,
		WithDetails(FieldError{
			Field: "email",
			Code:  "FIELD_REQUIRED",
			Params: map[string]string{
				"field": "ignored",
				"ok":    strings.Repeat("x", 300),
				"other": "safe",
			},
		}),
	)
	err := e.Fail(c, validationError)
	if err != nil {
		t.Fatal(err)
	}
	b := decode(t, w)
	ds := b["error"].(map[string]any)["details"].([]any)
	if ds[0].(map[string]any)["message"] != "The email field is required." {
		t.Fatal(b)
	}
	w, _, c = request("GET", "")
	if err := e.OK(c, nil, nil); err != nil {
		t.Fatal(err)
	}
	if _, ok := decode(t, w)["data"]; !ok {
		t.Fatal("missing data")
	}
	w, r, c = request("GET", "en")
	r.Header.Set("X-Request-ID", "req_marshal")
	if err := e.OK(c, make(chan int)); err != nil {
		t.Fatal(err)
	}
	if b := decode(t, w); b["meta"].(map[string]any)["request_id"] != "req_marshal" {
		t.Fatal(b)
	}
	if got := e.field("_lead__tail_"); got != "_lead__tail_" {
		t.Fatal(got)
	}
	camel, _ := New(WithNaming(CamelCase))
	if got := camel.field("_lead__tail_"); got != "LeadTail" {
		t.Fatal(got)
	}
	w, _, c = request("GET", "en")
	if err := e.Fail(c, NewError("NOT_TRANSLATED", 418)); err != nil {
		t.Fatal(err)
	}
	if got := decode(t, w)["error"].(map[string]any)["message"]; got != "An internal server error occurred." {
		t.Fatal(got)
	}
}

type failedWriter struct {
	header http.Header
	status int
}

func (w *failedWriter) Header() http.Header {
	return w.header
}

func (w *failedWriter) WriteHeader(status int) {
	w.status = status
}

func (w *failedWriter) Write([]byte) (int, error) {
	return 0, errors.New("network write failed")
}

func TestTransportWriteError(t *testing.T) {
	w := &failedWriter{header: http.Header{}}
	r := httptest.NewRequest("GET", "/", nil)
	err := OK(HTTP(w, r), 123)
	if err == nil || !strings.Contains(err.Error(), "resp: response write") || w.status != 200 {
		t.Fatal(err, w.status)
	}
	w = &failedWriter{header: http.Header{}}
	err = OK(HTTP(w, r), make(chan int))
	if err == nil || w.status != 500 {
		t.Fatal(err, w.status)
	}
}
