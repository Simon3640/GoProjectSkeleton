package repositories

import (
	application_errors "gormgoskeleton/src/application/shared/errors"
	"gormgoskeleton/src/application/shared/locales/messages"
	"gormgoskeleton/src/application/shared/status"

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
	gorm.ErrDuplicatedKey: application_errors.NewApplicationError(
		status.Conflict,
		messages.MessageKeysInstance.RESOURCE_EXISTS,
		"Resource already exists",
	),
	gorm.ErrForeignKeyViolated: application_errors.NewApplicationError(
		status.Conflict,
		messages.MessageKeysInstance.RESOURCE_NOT_FOUND,
		"Foreign key violation",
	),
}

var DefaultORMError = application_errors.NewApplicationError(
	status.InternalError,
	messages.MessageKeysInstance.SOMETHING_WENT_WRONG,
	"Unexpected database error",
)
