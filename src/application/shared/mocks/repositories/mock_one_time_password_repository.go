package repositoriesmocks

import (
	contracts_repositories "goprojectskeleton/src/application/contracts/repositories"
	dtos "goprojectskeleton/src/application/shared/DTOs"
	application_errors "goprojectskeleton/src/application/shared/errors"
	"goprojectskeleton/src/domain/models"
)

type MockOneTimePasswordRepository struct {
	MockRepositoryBase[dtos.OneTimePasswordCreate, dtos.OneTimePasswordUpdate, models.OneTimePassword, models.OneTimePassword]
}

func (m *MockOneTimePasswordRepository) GetByPasswordHash(tokenHash []byte) (*models.OneTimePassword, *application_errors.ApplicationError) {
	args := m.Called(tokenHash)
	errorArg := args.Get(1)
	if errorArg != nil {
		return nil, errorArg.(*application_errors.ApplicationError)
	}
	return args.Get(0).(*models.OneTimePassword), nil
}

var _ contracts_repositories.IOneTimePasswordRepository = (*MockOneTimePasswordRepository)(nil)
