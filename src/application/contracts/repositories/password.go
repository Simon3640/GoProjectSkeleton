package contracts_repositories

import (
	dtos "gormgoskeleton/src/application/shared/DTOs"
	application_errors "gormgoskeleton/src/application/shared/errors"
	"gormgoskeleton/src/domain/models"
)

type IPasswordRepository interface {
	IRepositoryBase[dtos.PasswordCreate, dtos.PasswordUpdate, models.Password, models.PasswordInDB]
	GetActivePassword(userEmail string) (*models.Password, *application_errors.ApplicationError)
}
