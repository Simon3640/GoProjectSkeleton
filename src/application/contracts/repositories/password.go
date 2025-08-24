package contracts_repositories

import (
	application_errors "gormgoskeleton/src/application/shared/errors"
	"gormgoskeleton/src/domain/models"
)

type IPasswordRepository interface {
	IRepositoryBase[models.PasswordCreate, models.PasswordUpdate, models.Password, models.PasswordInDB]
	GetActivePassword(userEmail string) (*models.Password, *application_errors.ApplicationError)
}
