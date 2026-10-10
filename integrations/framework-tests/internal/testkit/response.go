package testkit

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ZihxS/RespKit/resp"
)

type Performer func(*http.Request) (*http.Response, error)

// TestResponse runs the shared response contract against a real framework router.
func TestResponse(t *testing.T, framework string, call Performer) {
	t.Helper()
	for _, tt := range []struct {
		path, method, locale string
		status               int
		code                 string
	}{
		{"/ok", "GET", "en-US", 200, ""},
		{"/error", "GET", "en-US", 409, "RESOURCE_ALREADY_EXISTS"},
		{"/error", "GET", "id-ID", 409, "RESOURCE_ALREADY_EXISTS"},
		{"/raw", "GET", "en", 500, "INTERNAL_ERROR"},
		{"/debug", "GET", "en", 500, "INTERNAL_ERROR"},
		{"/created", "POST", "en", 201, ""},
		{"/empty", "GET", "id", 204, ""},
	} {
		t.Run(framework+tt.path+tt.locale, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "http://localhost"+tt.path, nil)
			req.Header.Set("Accept-Language", tt.locale)
			req.Header.Set("X-Request-ID", "req_integration")
			r, err := call(req)
			if err != nil {
				t.Fatal(err)
			}
			defer r.Body.Close()
			body, err := io.ReadAll(r.Body)
			if err != nil {
				t.Fatal(err)
			}
			if r.StatusCode != tt.status {
				t.Fatalf("status=%d want=%d body=%s", r.StatusCode, tt.status, body)
			}
			if tt.status == 204 {
				if len(body) != 0 {
					t.Fatalf("204 with body %q", body)
				}
				return
			}
			if r.Header.Get("Content-Language") != strings.Split(tt.locale, "-")[0] {
				t.Fatalf("locale: %#v", r.Header)
			}
			if !strings.Contains(r.Header.Get("Vary"), "Accept-Language") {
				t.Fatalf("missing vary header: %v", r.Header)
			}
			if r.Header.Get("Content-Type") == "" {
				t.Fatalf("missing content type: %v", r.Header)
			}
			var obj map[string]any
			if err := json.Unmarshal(body, &obj); err != nil {
				t.Fatalf("invalid JSON: %s: %v", body, err)
			}
			if tt.code != "" {
				got := obj["error"].(map[string]any)["code"]
				if got != tt.code {
					t.Fatalf("code %v want %v", got, tt.code)
				}
				if tt.path == "/debug" {
					if obj["error"].(map[string]any)["debug_message"] != "local-debug-only" {
						t.Fatalf("missing local debug message: %s", body)
					}
				} else if strings.Contains(string(body), "secret-pg-password") || strings.Contains(string(body), "debug_message") {
					t.Fatal("leaked raw error")
				}
			} else if obj["success"] != true {
				t.Fatal(obj)
			}
			if tt.status == 201 && r.Header.Get("Location") != "/users/99" {
				t.Fatalf("missing location: %v", r.Header)
			}
		})
	}
}

// Execute runs one shared response scenario using the given framework context.
func Execute(kind string, c any) error {
	switch kind {
	case "ok":
		return resp.OK(c, map[string]any{"user_id": 99})
	case "error":
		return resp.Fail(c, structuralSQLState{"23505"})
	case "raw":
		return resp.Fail(c, errors.New("secret-pg-password"))
	case "debug":
		api, err := resp.New(resp.WithExposeInternalErrors(true))
		if err != nil {
			return err
		}
		return api.Fail(c, errors.New("local-debug-only"))
	case "created":
		return resp.Created(c, map[string]int{"id": 99}, resp.Location("/users/99"))
	case "empty":
		return resp.NoContent(c)
	default:
		return errors.New("unknown route")
	}
}

// StandardHandler adapts a shared scenario for net/http-compatible routers.
func StandardHandler(t *testing.T, framework, kind string) http.HandlerFunc {
	t.Helper()
	return func(w http.ResponseWriter, r *http.Request) {
		if err := Execute(kind, resp.HTTP(w, r)); err != nil {
			t.Errorf("%s: %v", framework, err)
		}
	}
}

type structuralSQLState struct{ code string }

func (e structuralSQLState) Error() string {
	return fmt.Sprintf("sqlstate %s", e.code)
}

func (e structuralSQLState) SQLState() string {
	return e.code
}
