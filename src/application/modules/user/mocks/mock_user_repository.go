// Package usermocks contains mock implementations of the contracts/repositories/ interfaces
package usermocks

import (
	usercontracts "github.com/simon3640/goprojectskeleton/src/application/modules/user/contracts"
	userdtos "github.com/simon3640/goprojectskeleton/src/application/modules/user/dtos"
	applicationerror "github.com/simon3640/goprojectskeleton/src/application/shared/errors"
	repositoriesmocks "github.com/simon3640/goprojectskeleton/src/application/shared/mocks/repositories"
	usermodels "github.com/simon3640/goprojectskeleton/src/domain/user/models"
)

// MockUserRepository is the mock implementation of the UserRepository interface
type MockUserRepository struct {
	repositoriesmocks.MockRepositoryBase[userdtos.UserCreate, userdtos.UserUpdate, usermodels.User, usermodels.User]
}

var _ usercontracts.IUserRepository = (*MockUserRepository)(nil)

// CreateWithPassword creates a new user with a password
func (m *MockUserRepository) CreateWithPassword(input userdtos.UserAndPasswordCreate) (*usermodels.User, *applicationerror.ApplicationError) {
	args := m.Called(input)
	errorArg := args.Get(1)
	if errorArg != nil {
		return args.Get(0).(*usermodels.User), errorArg.(*applicationerror.ApplicationError)
	}
	return args.Get(0).(*usermodels.User), nil
}

// GetUserWithRole gets a user with their role
func (m *MockUserRepository) GetUserWithRole(id uint) (*usermodels.UserWithRole, *applicationerror.ApplicationError) {
	args := m.Called(id)
	errorArg := args.Get(1)
	if errorArg != nil {
		return args.Get(0).(*usermodels.UserWithRole), errorArg.(*applicationerror.ApplicationError)
	}
	return args.Get(0).(*usermodels.UserWithRole), nil
}

// GetByEmailOrPhone gets a user by email or phone
func (m *MockUserRepository) GetByEmailOrPhone(emailOrPhone string) (*usermodels.User, *applicationerror.ApplicationError) {
	args := m.Called(emailOrPhone)
	errorArg := args.Get(1)
	if errorArg != nil {
		return nil, errorArg.(*applicationerror.ApplicationError)
	}
	return args.Get(0).(*usermodels.User), nil
}
