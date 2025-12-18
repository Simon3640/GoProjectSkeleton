// Package passwordmocks contains mock implementations of the contracts/repositories/ interfaces
package passwordmocks

import (
	passwordcontracts "github.com/simon3640/goprojectskeleton/src/application/modules/password/contracts"
	dtos "github.com/simon3640/goprojectskeleton/src/application/modules/password/dtos"
	applicationerrors "github.com/simon3640/goprojectskeleton/src/application/shared/errors"
	repositoriesmocks "github.com/simon3640/goprojectskeleton/src/application/shared/mocks/repositories"
	"github.com/simon3640/goprojectskeleton/src/domain/models"
)

// MockPasswordRepository is the mock implementation of the PasswordRepository interface
type MockPasswordRepository struct {
	repositoriesmocks.MockRepositoryBase[dtos.PasswordCreate, dtos.PasswordUpdate, models.Password, models.PasswordInDB]
}

var _ passwordcontracts.IPasswordRepository = (*MockPasswordRepository)(nil)

// GetActivePassword gets the active password for a user
func (m *MockPasswordRepository) GetActivePassword(userEmail string) (*models.Password, *applicationerrors.ApplicationError) {
	args := m.Called(userEmail)
	errorArg := args.Get(1)
	if errorArg != nil {
		return nil, errorArg.(*applicationerrors.ApplicationError)
	}
	return args.Get(0).(*models.Password), nil
}
