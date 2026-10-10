// Package pgxmapper provides opt-in PostgreSQL error mapping for pgx/v5.
package pgxmapper

import (
	"errors"
	"net/http"

	"github.com/ZihxS/RespKit/resp"
	"github.com/jackc/pgx/v5/pgconn"
)

// Name is the mapper name used in RespKit diagnostics.
const Name = "pgx"

// Map classifies PostgreSQL errors with a reasonably stable HTTP meaning.
// Other SQLSTATEs are handled as generic internal errors because they usually
// indicate a query bug, database configuration problem, or application-specific
// condition rather than a client error.
func Map(err error) (resp.Definition, bool) {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) || pgErr == nil {
		return resp.Definition{}, false
	}

	switch pgErr.Code {
	// Invalid request values reported while converting PostgreSQL parameters.
	case "22003", // numeric_value_out_of_range
		"22007", // invalid_datetime_format
		"22008", // datetime_field_overflow
		"22018", // invalid_character_value_for_cast
		"22021", // character_not_in_repertoire
		"22023", // invalid_parameter_value
		"22P02", // invalid_text_representation
		"22P03": // invalid_binary_representation
		return resp.Definition{Code: "BAD_REQUEST", Status: http.StatusBadRequest}, true
	case "22001": // string_data_right_truncation
		return resp.Definition{Code: "VALIDATION_FAILED", Status: http.StatusUnprocessableEntity}, true

	// Constraint failures are request validation errors, except uniqueness and
	// exclusion failures which represent a conflict with existing resources.
	case "23502", // not_null_violation
		"23503", // foreign_key_violation
		"23514": // check_violation
		return resp.Definition{Code: "VALIDATION_FAILED", Status: http.StatusUnprocessableEntity}, true
	case "23505": // unique_violation
		return resp.Definition{Code: "RESOURCE_ALREADY_EXISTS", Status: http.StatusConflict}, true
	case "23000", // integrity_constraint_violation
		"23001", // restrict_violation
		"23P01": // exclusion_violation
		return resp.Definition{Code: "RESOURCE_CONFLICT", Status: http.StatusConflict}, true

	// Transaction failures commonly require retrying the request in a new
	// transaction. Keep that decision with the caller by reporting a conflict.
	case "40000", // transaction_rollback
		"40001", // serialization_failure
		"40P01": // deadlock_detected
		return resp.Definition{Code: "RESOURCE_CONFLICT", Status: http.StatusConflict}, true

	// These SQLSTATEs indicate that the database dependency cannot currently
	// serve the request. Credentials, SQL syntax, and protocol errors are not
	// mapped here because they generally require application/operator action.
	case "08000", // connection_exception
		"08001", // sqlclient_unable_to_establish_sqlconnection
		"08003", // connection_does_not_exist
		"08004", // sqlserver_rejected_establishment_of_sqlconnection
		"08006", // connection_failure
		"08007", // transaction_resolution_unknown
		"53000", // insufficient_resources
		"53100", // disk_full
		"53200", // out_of_memory
		"53300", // too_many_connections
		"53400", // configuration_limit_exceeded
		"57P01", // admin_shutdown
		"57P02", // crash_shutdown
		"57P03": // cannot_connect_now
		return resp.Definition{Code: "SERVICE_UNAVAILABLE", Status: http.StatusServiceUnavailable}, true

	case "55P03": // lock_not_available
		return resp.Definition{Code: "RESOURCE_CONFLICT", Status: http.StatusConflict}, true
	default:
		return resp.Definition{Code: "INTERNAL_ERROR", Status: http.StatusInternalServerError}, true
	}
}

// Option registers the pgx mapper on a RespKit engine.
func Option() resp.Option {
	return resp.WithMapper(Name, Map)
}
