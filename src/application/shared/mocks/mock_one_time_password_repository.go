package mocks

import (
	contracts_repositories "gormgoskeleton/src/application/contracts/repositories"
	dtos "gormgoskeleton/src/application/shared/DTOs"
	application_errors "gormgoskeleton/src/application/shared/errors"
	"gormgoskeleton/src/domain/models"
)

type MockOneTimePasswordRepository struct {
	MockRepositoryBase[dtos.OneTimePasswordCreate, dtos.OneTimePasswordUpdate, models.OneTimePassword, models.OneTimePassword]
}

func (m *MockOneTimePasswordRepository) GetByPasswordHash(tokenHash []byte) (*models.OneTimePassword, *application_errors.ApplicationError) {
	args := m.Called(tokenHash)
	errorArg := args.Get(1)
	if errorArg != nil {
		return args.Get(0).(*models.OneTimePassword), errorArg.(*application_errors.ApplicationError)
	}
	return args.Get(0).(*models.OneTimePassword), nil
}

var _ contracts_repositories.IOneTimePasswordRepository = (*MockOneTimePasswordRepository)(nil)
