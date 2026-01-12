package authcontracts

import (
	applicationerrors "github.com/simon3640/goprojectskeleton/src/application/shared/errors"
	usermodels "github.com/simon3640/goprojectskeleton/src/domain/user/models"
)

// IUserRepository is the interface for the user repository
type IUserRepository interface {
	// GetUserWithRole gets a user with their role
	GetUserWithRole(id uint) (*usermodels.UserWithRole, *applicationerrors.ApplicationError)
	// GetByEmailOrPhone gets a user by email or phone
	GetByEmailOrPhone(emailOrPhone string) (*usermodels.User, *applicationerrors.ApplicationError)
}
