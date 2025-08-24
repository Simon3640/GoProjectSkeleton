package contracts_repositories

import (
	application_errors "gormgoskeleton/src/application/shared/errors"
	"gormgoskeleton/src/domain/models"
)

type IUserRepository interface {
	IRepositoryBase[models.UserCreate, models.UserUpdate, models.User, models.UserInDB]
	CreateWithPassword(input models.UserAndPasswordCreate) (*models.User, *application_errors.ApplicationError)
	GetUserWithRole(id uint) (*models.UserWithRole, *application_errors.ApplicationError)
}
