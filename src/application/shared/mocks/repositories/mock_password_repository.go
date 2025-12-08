package repositoriesmocks

import (
	contracts_repositories "github.com/simon3640/goprojectskeleton/src/application/contracts/repositories"
	dtos "github.com/simon3640/goprojectskeleton/src/application/shared/DTOs"
	application_errors "github.com/simon3640/goprojectskeleton/src/application/shared/errors"
	"github.com/simon3640/goprojectskeleton/src/domain/models"
)

type MockPasswordRepository struct {
	MockRepositoryBase[dtos.PasswordCreate, dtos.PasswordUpdate, models.Password, models.PasswordInDB]
}

var _ contracts_repositories.IPasswordRepository = (*MockPasswordRepository)(nil)

func (m *MockPasswordRepository) GetActivePassword(userEmail string) (*models.Password, *application_errors.ApplicationError) {
	args := m.Called(userEmail)
	errorArg := args.Get(1)
	if errorArg != nil {
		return nil, errorArg.(*application_errors.ApplicationError)
	}
	return args.Get(0).(*models.Password), nil
}
