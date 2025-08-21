package mocks

import (
	"gormgoskeleton/src/domain/models"
)

type MockPasswordRepository struct {
	MockRepositoryBase[models.PasswordCreate, models.PasswordUpdate, models.Password, models.PasswordInDB]
}

func (m *MockPasswordRepository) GetActivePassword(userEmail string) (*models.Password, error) {
	args := m.Called(userEmail)
	return args.Get(0).(*models.Password), args.Error(1)
}
