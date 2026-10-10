package pgxmapper

import (
	"errors"
	"fmt"
	"testing"

	"github.com/ZihxS/RespKit/resp"
	"github.com/jackc/pgx/v5/pgconn"
)

func TestMapWrappedPostgresErrors(t *testing.T) {
	tests := []struct {
		name       string
		code       string
		wantCode   string
		wantStatus int
		wantMapped bool
	}{
		{name: "unique violation", code: "23505", wantCode: "RESOURCE_ALREADY_EXISTS", wantStatus: 409, wantMapped: true},
		{name: "exclusion violation", code: "23P01", wantCode: "RESOURCE_CONFLICT", wantStatus: 409, wantMapped: true},
		{name: "foreign key violation", code: "23503", wantCode: "VALIDATION_FAILED", wantStatus: 422, wantMapped: true},
		{name: "not null violation", code: "23502", wantCode: "VALIDATION_FAILED", wantStatus: 422, wantMapped: true},
		{name: "connection failure", code: "08006", wantCode: "SERVICE_UNAVAILABLE", wantStatus: 503, wantMapped: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := fmt.Errorf("query failed: %w", &pgconn.PgError{Code: tt.code})
			got, ok := Map(err)
			if ok != tt.wantMapped {
				t.Fatalf("Map() mapped = %v, want %v", ok, tt.wantMapped)
			}
			if !ok {
				return
			}
			if got.Code != tt.wantCode || got.Status != tt.wantStatus {
				t.Fatalf("Map() = (%q, %d), want (%q, %d)", got.Code, got.Status, tt.wantCode, tt.wantStatus)
			}
		})
	}
}

func TestMapAbstainsForNonPgxErrors(t *testing.T) {
	errorsToCheck := []error{
		nil,
		errors.New("duplicate key value violates unique constraint"),
		sqlStateError("23P01"),
	}
	for _, err := range errorsToCheck {
		if got, ok := Map(err); ok {
			t.Errorf("Map(%v) = (%+v, true), want no mapping", err, got)
		}
	}
}

func TestOptionUsesRealPgxErrorAndPreservesFallback(t *testing.T) {
	e, err := resp.New(Option())
	if err != nil {
		t.Fatal(err)
	}

	d, diag := e.Resolve(fmt.Errorf("insert event: %w", &pgconn.PgError{Code: "23505"}))
	if d.Code != "RESOURCE_ALREADY_EXISTS" || d.Status != 409 {
		t.Fatalf("unique violation definition = (%q, %d), want (RESOURCE_ALREADY_EXISTS, 409)", d.Code, d.Status)
	}
	if diag.Mapper != "custom.pgx" || !diag.Detected {
		t.Fatalf("unique violation diagnostic = %+v, want detected custom.pgx mapper", diag)
	}

	d, diag = e.Resolve(fmt.Errorf("insert event: %w", &pgconn.PgError{Code: "23P01"}))
	if d.Code != "RESOURCE_CONFLICT" || d.Status != 409 {
		t.Fatalf("Resolve() definition = (%q, %d), want (RESOURCE_CONFLICT, 409)", d.Code, d.Status)
	}
	if diag.Mapper != "custom.pgx" || !diag.Detected {
		t.Fatalf("Resolve() diagnostic = %+v, want detected custom.pgx mapper", diag)
	}

	d, diag = e.Resolve(fmt.Errorf("insert event: %w", &pgconn.PgError{Code: "42601"}))
	if d.Code != "INTERNAL_ERROR" || d.Status != 500 {
		t.Fatalf("unclassified SQLSTATE definition = (%q, %d), want (INTERNAL_ERROR, 500)", d.Code, d.Status)
	}
	if diag.Mapper != "custom.pgx" || diag.SourceCode != "42601" {
		t.Fatalf("unclassified SQLSTATE diagnostic = %+v, want pgx generic fallback and source code", diag)
	}
}

type sqlStateError string

func (e sqlStateError) Error() string {
	return "database error"
}

func (e sqlStateError) SQLState() string {
	return string(e)
}
