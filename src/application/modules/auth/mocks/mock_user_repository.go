// Package authmocks contains mock implementations of the contracts/auth/ interfaces
package authmocks

import (
	authcontracts "github.com/simon3640/goprojectskeleton/src/application/modules/auth/contracts"
	appllicationerrors "github.com/simon3640/goprojectskeleton/src/application/shared/errors"
	"github.com/simon3640/goprojectskeleton/src/domain/models"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository is the mock implementation of the UserRepository interface
type MockUserRepository struct {
	mock.Mock
}

var _ authcontracts.IUserRepository = (*MockUserRepository)(nil)

// GetUserWithRole gets a user with their role
func (m *MockUserRepository) GetUserWithRole(id uint) (*models.UserWithRole, *appllicationerrors.ApplicationError) {
	args := m.Called(id)
	errorArg := args.Get(1)
	if errorArg != nil {
		return args.Get(0).(*models.UserWithRole), errorArg.(*appllicationerrors.ApplicationError)
	}
	return args.Get(0).(*models.UserWithRole), nil
}

// GetByEmailOrPhone gets a user by email or phone
func (m *MockUserRepository) GetByEmailOrPhone(emailOrPhone string) (*models.User, *appllicationerrors.ApplicationError) {
	args := m.Called(emailOrPhone)
	errorArg := args.Get(1)
	if errorArg != nil {
		return nil, errorArg.(*appllicationerrors.ApplicationError)
	}
	return args.Get(0).(*models.User), nil
}
