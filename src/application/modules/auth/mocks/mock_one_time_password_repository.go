package authmocks

import (
	authcontracts "github.com/simon3640/goprojectskeleton/src/application/modules/auth/contracts"
	dtos "github.com/simon3640/goprojectskeleton/src/application/modules/auth/dtos"
	applicationerrors "github.com/simon3640/goprojectskeleton/src/application/shared/errors"
	repositoriesmocks "github.com/simon3640/goprojectskeleton/src/application/shared/mocks/repositories"
	"github.com/simon3640/goprojectskeleton/src/domain/models"
)

// MockOneTimePasswordRepository is the mock implementation of the one time password repository
type MockOneTimePasswordRepository struct {
	repositoriesmocks.MockRepositoryBase[dtos.OneTimePasswordCreate, dtos.OneTimePasswordUpdate, models.OneTimePassword, models.OneTimePassword]
}

// GetByPasswordHash gets a one time password by its hash
func (m *MockOneTimePasswordRepository) GetByPasswordHash(tokenHash []byte) (*models.OneTimePassword, *applicationerrors.ApplicationError) {
	args := m.Called(tokenHash)
	errorArg := args.Get(1)
	if errorArg != nil {
		return nil, errorArg.(*applicationerrors.ApplicationError)
	}
	return args.Get(0).(*models.OneTimePassword), nil
}

var _ authcontracts.IOneTimePasswordRepository = (*MockOneTimePasswordRepository)(nil)
