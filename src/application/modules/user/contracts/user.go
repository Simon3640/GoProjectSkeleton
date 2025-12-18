// Package usercontracts contains the interfaces for the user module
package usercontracts

import (
	contractsrepositories "github.com/simon3640/goprojectskeleton/src/application/contracts/repositories"
	userdtos "github.com/simon3640/goprojectskeleton/src/application/modules/user/dtos"
	applicationerrors "github.com/simon3640/goprojectskeleton/src/application/shared/errors"
	"github.com/simon3640/goprojectskeleton/src/domain/models"
)

// IUserRepository is the interface for the user repository
type IUserRepository interface {
	contractsrepositories.IRepositoryBase[userdtos.UserCreate, userdtos.UserUpdate, models.User, models.UserInDB]
	// CreateWithPassword creates a new user with a password
	CreateWithPassword(input userdtos.UserAndPasswordCreate) (*models.User, *applicationerrors.ApplicationError)
	// GetUserWithRole gets a user with their role
	GetUserWithRole(id uint) (*models.UserWithRole, *applicationerrors.ApplicationError)
	// GetByEmailOrPhone gets a user by email or phone
	GetByEmailOrPhone(emailOrPhone string) (*models.User, *applicationerrors.ApplicationError)
}
