package contracts_repositories

import (
	dtos "gormgoskeleton/src/application/shared/DTOs"
	application_errors "gormgoskeleton/src/application/shared/errors"
	"gormgoskeleton/src/domain/models"
)

type IUserRepository interface {
	IRepositoryBase[dtos.UserCreate, dtos.UserUpdate, models.User, models.UserInDB]
	CreateWithPassword(input dtos.UserAndPasswordCreate) (*models.User, *application_errors.ApplicationError)
	GetUserWithRole(id uint) (*models.UserWithRole, *application_errors.ApplicationError)
}
