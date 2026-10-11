package resp_test

import (
	"net/http"
	"testing"

	"github.com/ZihxS/RespKit/resp"
)

func TestPublicCodeConstantsMatchWireCatalog(t *testing.T) {
	api, err := resp.New()
	if err != nil {
		t.Fatal(err)
	}
	catalog := make(map[string]resp.CatalogItem)
	for _, item := range api.Catalog() {
		catalog[item.Code] = item
	}
	for _, tt := range []struct {
		code   string
		wire   string
		status int
	}{
		{
			code:   resp.CodeBadRequest,
			wire:   "BAD_REQUEST",
			status: http.StatusBadRequest,
		},
		{
			code:   resp.CodeUnauthorized,
			wire:   "UNAUTHORIZED",
			status: http.StatusUnauthorized,
		},
		{
			code:   resp.CodeForbidden,
			wire:   "FORBIDDEN",
			status: http.StatusForbidden,
		},
		{
			code:   resp.CodeResourceNotFound,
			wire:   "RESOURCE_NOT_FOUND",
			status: http.StatusNotFound,
		},
		{
			code:   resp.CodeResourceConflict,
			wire:   "RESOURCE_CONFLICT",
			status: http.StatusConflict,
		},
		{
			code:   resp.CodeResourceAlreadyExists,
			wire:   "RESOURCE_ALREADY_EXISTS",
			status: http.StatusConflict,
		},
		{
			code:   resp.CodeValidationFailed,
			wire:   "VALIDATION_FAILED",
			status: http.StatusUnprocessableEntity,
		},
		{
			code:   resp.CodeRateLimited,
			wire:   "RATE_LIMITED",
			status: http.StatusTooManyRequests,
		},
		{
			code:   resp.CodeServiceUnavailable,
			wire:   "SERVICE_UNAVAILABLE",
			status: http.StatusServiceUnavailable,
		},
		{
			code:   resp.CodeInternalError,
			wire:   "INTERNAL_ERROR",
			status: http.StatusInternalServerError,
		},
	} {
		t.Run(tt.wire, func(t *testing.T) {
			if tt.code != tt.wire {
				t.Fatalf("public code = %q, want stable wire value %q", tt.code, tt.wire)
			}
			item, ok := catalog[tt.code]
			if !ok || item.HTTPStatus != tt.status || item.Messages["id"] == "" || item.Messages["en"] == "" {
				t.Fatalf("catalog entry = %+v, want HTTP %d and both built-in translations", item, tt.status)
			}
		})
	}
}
