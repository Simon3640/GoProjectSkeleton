package mocks

import (
	contracts_repositories "gormgoskeleton/src/application/contracts/repositories"
	dtos "gormgoskeleton/src/application/shared/DTOs"
	application_errors "gormgoskeleton/src/application/shared/errors"
	"gormgoskeleton/src/domain/models"
)

type MockUserRepository struct {
	MockRepositoryBase[dtos.UserCreate, dtos.UserUpdate, models.User, models.User]
}

var _ contracts_repositories.IUserRepository = (*MockUserRepository)(nil)

func (m *MockUserRepository) CreateWithPassword(input dtos.UserAndPasswordCreate) (*models.User, *application_errors.ApplicationError) {
	args := m.Called(input)
	errorArg := args.Get(1)
	if errorArg != nil {
		return args.Get(0).(*models.User), errorArg.(*application_errors.ApplicationError)
	}
	return args.Get(0).(*models.User), nil
}

func (m *MockUserRepository) GetUserWithRole(id uint) (*models.UserWithRole, *application_errors.ApplicationError) {
	args := m.Called(id)
	errorArg := args.Get(1)
	if errorArg != nil {
		return args.Get(0).(*models.UserWithRole), errorArg.(*application_errors.ApplicationError)
	}
	return args.Get(0).(*models.UserWithRole), nil
}

func (m *MockUserRepository) GetByEmailOrPhone(emailOrPhone string) (*models.User, *application_errors.ApplicationError) {
	args := m.Called(emailOrPhone)
	errorArg := args.Get(1)
	if errorArg != nil {
		return args.Get(0).(*models.User), errorArg.(*application_errors.ApplicationError)
	}
	return args.Get(0).(*models.User), nil
}
