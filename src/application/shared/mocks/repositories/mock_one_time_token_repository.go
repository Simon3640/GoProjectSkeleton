package repositoriesmocks

import (
	contracts_repositories "github.com/simon3640/goprojectskeleton/src/application/contracts/repositories"
	dtos "github.com/simon3640/goprojectskeleton/src/application/shared/DTOs"
	application_errors "github.com/simon3640/goprojectskeleton/src/application/shared/errors"
	"github.com/simon3640/goprojectskeleton/src/domain/models"
)

type MockOneTimeTokenRepository struct {
	MockRepositoryBase[dtos.OneTimeTokenCreate, dtos.OneTimeTokenUpdate, models.OneTimeToken, models.OneTimeToken]
}

func (m *MockOneTimeTokenRepository) GetByTokenHash(tokenHash []byte) (*models.OneTimeToken, *application_errors.ApplicationError) {
	args := m.Called(tokenHash)
	errorArg := args.Get(1)
	if errorArg != nil {
		return nil, errorArg.(*application_errors.ApplicationError)
	}
	return args.Get(0).(*models.OneTimeToken), nil
}

var _ contracts_repositories.IOneTimeTokenRepository = (*MockOneTimeTokenRepository)(nil)
