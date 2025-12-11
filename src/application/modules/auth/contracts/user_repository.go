package authcontracts

import (
	application_errors "github.com/simon3640/goprojectskeleton/src/application/shared/errors"
	"github.com/simon3640/goprojectskeleton/src/domain/models"
)

// IUserRepository is the interface for the user repository
type IUserRepository interface {
	// GetUserWithRole gets a user with their role
	GetUserWithRole(id uint) (*models.UserWithRole, *application_errors.ApplicationError)
	// GetByEmailOrPhone gets a user by email or phone
	GetByEmailOrPhone(emailOrPhone string) (*models.User, *application_errors.ApplicationError)
}
