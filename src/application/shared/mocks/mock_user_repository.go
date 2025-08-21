package mocks

import (
	"gormgoskeleton/src/domain/models"
)

type MockUserRepository struct {
	MockRepositoryBase[models.UserCreate, models.UserUpdate, models.User, models.User]
}

func (m *MockUserRepository) CreateWithPassword(input models.UserAndPasswordCreate) (*models.User, error) {
	args := m.Called(input)
	return args.Get(0).(*models.User), args.Error(1)
}
