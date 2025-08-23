package repositories

import (
	"errors"
	application_errors "gormgoskeleton/src/application/shared/errors"
	"gormgoskeleton/src/application/shared/locales/messages"
	"gormgoskeleton/src/application/shared/status"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

var ORMErrorMapping = map[error]*application_errors.ApplicationError{
	gorm.ErrDuplicatedKey: application_errors.NewApplicationError(
		status.Conflict,
		messages.MessageKeysInstance.RESOURCE_EXISTS,
		"Resource already exists",
	),
	gorm.ErrInvalidData: application_errors.NewApplicationError(
		status.InvalidInput,
		messages.MessageKeysInstance.SOMETHING_WENT_WRONG,
		"Invalid data",
	),
	gorm.ErrRecordNotFound: application_errors.NewApplicationError(
		status.NotFound,
		messages.MessageKeysInstance.RESOURCE_NOT_FOUND,
		"Resource not found",
	),
	gorm.ErrForeignKeyViolated: application_errors.NewApplicationError(
		status.Conflict,
		messages.MessageKeysInstance.RESOURCE_NOT_FOUND,
		"Foreign key violation",
	),
}

// Errores Postgres por código
var PostgresErrorMapping = map[string]*application_errors.ApplicationError{
	"23505": application_errors.NewApplicationError(
		status.Conflict,
		messages.MessageKeysInstance.RESOURCE_EXISTS,
		"Unique constraint violated",
	),
	"23503": application_errors.NewApplicationError(
		status.Conflict,
		messages.MessageKeysInstance.RESOURCE_NOT_FOUND,
		"Foreign key violation",
	),
	"23502": application_errors.NewApplicationError(
		status.InvalidInput,
		messages.MessageKeysInstance.SOMETHING_WENT_WRONG,
		"Null value violation",
	),
	"22001": application_errors.NewApplicationError(
		status.InvalidInput,
		messages.MessageKeysInstance.SOMETHING_WENT_WRONG,
		"Value too long for column",
	),
	"22P02": application_errors.NewApplicationError(
		status.InvalidInput,
		messages.MessageKeysInstance.SOMETHING_WENT_WRONG,
		"Invalid input syntax",
	),
	"23514": application_errors.NewApplicationError(
		status.InvalidInput,
		messages.MessageKeysInstance.SOMETHING_WENT_WRONG,
		"Check constraint violation",
	),
	"40001": application_errors.NewApplicationError(
		status.Conflict,
		messages.MessageKeysInstance.SOMETHING_WENT_WRONG,
		"Serialization failure",
	),
}

// Error por defecto
var DefaultORMError = application_errors.NewApplicationError(
	status.InternalError,
	messages.MessageKeysInstance.SOMETHING_WENT_WRONG,
	"Unexpected database error",
)

// Función de mapeo unificada
func MapOrmError(err error) *application_errors.ApplicationError {
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
