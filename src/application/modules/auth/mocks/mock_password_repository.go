package authmocks

import (
	authcontracts "github.com/simon3640/goprojectskeleton/src/application/modules/auth/contracts"
	applicationerrors "github.com/simon3640/goprojectskeleton/src/application/shared/errors"
	"github.com/simon3640/goprojectskeleton/src/domain/models"
	"github.com/stretchr/testify/mock"
)

// MockPasswordRepository is the mock implementation of the PasswordRepository interface
type MockPasswordRepository struct {
	mock.Mock
}

var _ authcontracts.IPasswordRepository = (*MockPasswordRepository)(nil)

// GetActivePassword gets the active password for a user
func (m *MockPasswordRepository) GetActivePassword(userEmail string) (*models.Password, *applicationerrors.ApplicationError) {
	args := m.Called(userEmail)
	errorArg := args.Get(1)
	if errorArg != nil {
		return nil, errorArg.(*applicationerrors.ApplicationError)
	}
	return args.Get(0).(*models.Password), nil
}
