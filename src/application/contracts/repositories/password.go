package contracts_repositories

import (
	dtos "goprojectskeleton/src/application/shared/DTOs"
	application_errors "goprojectskeleton/src/application/shared/errors"
	"goprojectskeleton/src/domain/models"
)

type IPasswordRepository interface {
	IRepositoryBase[dtos.PasswordCreate, dtos.PasswordUpdate, models.Password, models.PasswordInDB]
	GetActivePassword(userEmail string) (*models.Password, *application_errors.ApplicationError)
}
