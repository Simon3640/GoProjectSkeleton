package repositoriesmocks

import (
	contracts_repositories "github.com/simon3640/goprojectskeleton/src/application/contracts/repositories"
	dtos "github.com/simon3640/goprojectskeleton/src/application/shared/DTOs"
	application_errors "github.com/simon3640/goprojectskeleton/src/application/shared/errors"
	sharedmodels "github.com/simon3640/goprojectskeleton/src/domain/shared/models"
)

type MockOneTimeTokenRepository struct {
	MockRepositoryBase[dtos.OneTimeTokenCreate, dtos.OneTimeTokenUpdate, sharedmodels.OneTimeToken, sharedmodels.OneTimeToken]
}

// GetByTokenHash retrieves a one-time token by its hash
func (m *MockOneTimeTokenRepository) GetByTokenHash(tokenHash []byte) (*sharedmodels.OneTimeToken, *application_errors.ApplicationError) {
	args := m.Called(tokenHash)
	errorArg := args.Get(1)
	if errorArg != nil {
		return nil, errorArg.(*application_errors.ApplicationError)
	}
	if args.Get(0) == nil {
		return nil, nil
	}
	return args.Get(0).(*sharedmodels.OneTimeToken), nil
}

var _ contracts_repositories.IOneTimeTokenRepository = (*MockOneTimeTokenRepository)(nil)
