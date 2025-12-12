// Package repositoriesmocks contains mock implementations of the contracts/repositories/ interfaces
package repositoriesmocks

import (
	usercontracts "github.com/simon3640/goprojectskeleton/src/application/modules/user/contracts"
	userdtos "github.com/simon3640/goprojectskeleton/src/application/modules/user/dtos"
	applicationerror "github.com/simon3640/goprojectskeleton/src/application/shared/errors"
	repositoriesmocks "github.com/simon3640/goprojectskeleton/src/application/shared/mocks/repositories"
	"github.com/simon3640/goprojectskeleton/src/domain/models"
)

// MockUserRepository is the mock implementation of the UserRepository interface
type MockUserRepository struct {
	repositoriesmocks.MockRepositoryBase[userdtos.UserCreate, userdtos.UserUpdate, models.User, models.User]
}

var _ usercontracts.IUserRepository = (*MockUserRepository)(nil)

// CreateWithPassword creates a new user with a password
func (m *MockUserRepository) CreateWithPassword(input userdtos.UserAndPasswordCreate) (*models.User, *applicationerror.ApplicationError) {
	args := m.Called(input)
	errorArg := args.Get(1)
	if errorArg != nil {
		return args.Get(0).(*models.User), errorArg.(*applicationerror.ApplicationError)
	}
	return args.Get(0).(*models.User), nil
}

// GetUserWithRole gets a user with their role
func (m *MockUserRepository) GetUserWithRole(id uint) (*models.UserWithRole, *applicationerror.ApplicationError) {
	args := m.Called(id)
	errorArg := args.Get(1)
	if errorArg != nil {
		return args.Get(0).(*models.UserWithRole), errorArg.(*applicationerror.ApplicationError)
	}
	return args.Get(0).(*models.UserWithRole), nil
}

// GetByEmailOrPhone gets a user by email or phone
func (m *MockUserRepository) GetByEmailOrPhone(emailOrPhone string) (*models.User, *applicationerror.ApplicationError) {
	args := m.Called(emailOrPhone)
	errorArg := args.Get(1)
	if errorArg != nil {
		return nil, errorArg.(*applicationerror.ApplicationError)
	}
	return args.Get(0).(*models.User), nil
}
