package shared

import (
	"errors"

	applicationerrors "github.com/simon3640/goprojectskeleton/src/application/shared/errors"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales/messages"
	"github.com/simon3640/goprojectskeleton/src/application/shared/status"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

// ORMErrorMapping is the mapping of GORM errors to application errors
var ORMErrorMapping = map[error]*applicationerrors.ApplicationError{
	gorm.ErrDuplicatedKey: applicationerrors.NewApplicationError(
		status.Conflict,
		messages.MessageKeysInstance.RESOURCE_EXISTS,
		"Resource already exists",
	),
	gorm.ErrInvalidData: applicationerrors.NewApplicationError(
		status.InvalidInput,
		messages.MessageKeysInstance.SOMETHING_WENT_WRONG,
		"Invalid data",
	),
	gorm.ErrRecordNotFound: applicationerrors.NewApplicationError(
		status.NotFound,
		messages.MessageKeysInstance.RESOURCE_NOT_FOUND,
		"Resource not found",
	),
	gorm.ErrForeignKeyViolated: applicationerrors.NewApplicationError(
		status.Conflict,
		messages.MessageKeysInstance.RESOURCE_NOT_FOUND,
		"Foreign key violation",
	),
}

// PostgresErrorMapping is the mapping of Postgres errors to application errors
var PostgresErrorMapping = map[string]*applicationerrors.ApplicationError{
	"23505": applicationerrors.NewApplicationError(
		status.Conflict,
		messages.MessageKeysInstance.RESOURCE_EXISTS,
		"Unique constraint violated",
	),
	"23503": applicationerrors.NewApplicationError(
		status.Conflict,
		messages.MessageKeysInstance.RESOURCE_NOT_FOUND,
		"Foreign key violation",
	),
	"23502": applicationerrors.NewApplicationError(
		status.InvalidInput,
		messages.MessageKeysInstance.SOMETHING_WENT_WRONG,
		"Null value violation",
	),
	"22001": applicationerrors.NewApplicationError(
		status.InvalidInput,
		messages.MessageKeysInstance.SOMETHING_WENT_WRONG,
		"Value too long for column",
	),
	"22P02": applicationerrors.NewApplicationError(
		status.InvalidInput,
		messages.MessageKeysInstance.SOMETHING_WENT_WRONG,
		"Invalid input syntax",
	),
	"23514": applicationerrors.NewApplicationError(
		status.InvalidInput,
		messages.MessageKeysInstance.SOMETHING_WENT_WRONG,
		"Check constraint violation",
	),
	"40001": applicationerrors.NewApplicationError(
		status.Conflict,
		messages.MessageKeysInstance.SOMETHING_WENT_WRONG,
		"Serialization failure",
	),
}

// DefaultORMError is the default application error for database errors
var DefaultORMError = applicationerrors.NewApplicationError(
	status.InternalError,
	messages.MessageKeysInstance.SOMETHING_WENT_WRONG,
	"Unexpected database error",
)

// MapOrmError maps a GORM error to an application error
func MapOrmError(err error) *applicationerrors.ApplicationError {
	if err == nil {
		return nil
	}

	// 1. Errores conocidos de GORM
	for gormErr, appErr := range ORMErrorMapping {
		if errors.Is(err, gormErr) {
			return appErr
		}
	}

	// 2. Errores de Postgres
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if appErr, ok := PostgresErrorMapping[pgErr.Code]; ok {
			return appErr
		}
	}

	// 3. Por defecto
	return DefaultORMError
}
