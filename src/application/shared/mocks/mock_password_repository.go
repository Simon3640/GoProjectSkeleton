package mocks

import (
	contracts_repositories "gormgoskeleton/src/application/contracts/repositories"
	application_errors "gormgoskeleton/src/application/shared/errors"
	"gormgoskeleton/src/domain/models"
)

type MockPasswordRepository struct {
	MockRepositoryBase[models.PasswordCreate, models.PasswordUpdate, models.Password, models.PasswordInDB]
}

var _ contracts_repositories.IPasswordRepository = (*MockPasswordRepository)(nil)

func (m *MockPasswordRepository) GetActivePassword(userEmail string) (*models.Password, *application_errors.ApplicationError) {
	args := m.Called(userEmail)
	errorArg := args.Get(1)
	if errorArg != nil {
		return args.Get(0).(*models.Password), errorArg.(*application_errors.ApplicationError)
	}
	return args.Get(0).(*models.Password), nil
}
