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
	case sqlStateNumericValueOutOfRange,
		sqlStateInvalidDatetimeFormat,
		sqlStateDatetimeFieldOverflow,
		sqlStateInvalidCharacterValueForCast,
		sqlStateCharacterNotInRepertoire,
		sqlStateInvalidParameterValue,
		sqlStateInvalidTextRepresentation,
		sqlStateInvalidBinaryRepresentation:
		return resp.Definition{
			Code:   resp.CodeBadRequest,
			Status: http.StatusBadRequest,
		}, true
	case sqlStateStringDataRightTruncation:
		return resp.Definition{
			Code:   resp.CodeValidationFailed,
			Status: http.StatusUnprocessableEntity,
		}, true

	// Constraint failures are request validation errors, except uniqueness and
	// exclusion failures which represent a conflict with existing resources.
	case sqlStateNotNullViolation,
		sqlStateForeignKeyViolation,
		sqlStateCheckViolation:
		return resp.Definition{
			Code:   resp.CodeValidationFailed,
			Status: http.StatusUnprocessableEntity,
		}, true
	case sqlStateUniqueViolation:
		return resp.Definition{
			Code:   resp.CodeResourceAlreadyExists,
			Status: http.StatusConflict,
		}, true
	case sqlStateIntegrityConstraintViolation,
		sqlStateRestrictViolation,
		sqlStateExclusionViolation:
		return resp.Definition{
			Code:   resp.CodeResourceConflict,
			Status: http.StatusConflict,
		}, true

	// Transaction failures commonly require retrying the request in a new
	// transaction. Keep that decision with the caller by reporting a conflict.
	case sqlStateTransactionRollback,
		sqlStateSerializationFailure,
		sqlStateDeadlockDetected:
		return resp.Definition{
			Code:   resp.CodeResourceConflict,
			Status: http.StatusConflict,
		}, true

	// These SQLSTATEs indicate that the database dependency cannot currently
	// serve the request. Credentials, SQL syntax, and protocol errors are not
	// mapped here because they generally require application/operator action.
	case sqlStateConnectionException,
		sqlStateSQLClientUnableToEstablishSQLConnection,
		sqlStateConnectionDoesNotExist,
		sqlStateSQLServerRejectedEstablishmentOfSQLConnection,
		sqlStateConnectionFailure,
		sqlStateTransactionResolutionUnknown,
		sqlStateInsufficientResources,
		sqlStateDiskFull,
		sqlStateOutOfMemory,
		sqlStateTooManyConnections,
		sqlStateConfigurationLimitExceeded,
		sqlStateAdminShutdown,
		sqlStateCrashShutdown,
		sqlStateCannotConnectNow:
		return resp.Definition{
			Code:   resp.CodeServiceUnavailable,
			Status: http.StatusServiceUnavailable,
		}, true

	case sqlStateLockNotAvailable:
		return resp.Definition{
			Code:   resp.CodeResourceConflict,
			Status: http.StatusConflict,
		}, true
	default:
		return resp.Definition{
			Code:   resp.CodeInternalError,
			Status: http.StatusInternalServerError,
		}, true
	}
}

// Option registers the pgx mapper on a RespKit engine.
func Option() resp.Option {
	return resp.WithMapper(Name, Map)
}
