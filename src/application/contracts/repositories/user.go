package contracts_repositories

import (
	dtos "github.com/simon3640/goprojectskeleton/src/application/shared/DTOs"
	application_errors "github.com/simon3640/goprojectskeleton/src/application/shared/errors"
	"github.com/simon3640/goprojectskeleton/src/domain/models"
)

type IUserRepository interface {
	IRepositoryBase[dtos.UserCreate, dtos.UserUpdate, models.User, models.UserInDB]
	CreateWithPassword(input dtos.UserAndPasswordCreate) (*models.User, *application_errors.ApplicationError)
	GetUserWithRole(id uint) (*models.UserWithRole, *application_errors.ApplicationError)
	GetByEmailOrPhone(emailOrPhone string) (*models.User, *application_errors.ApplicationError)
}
