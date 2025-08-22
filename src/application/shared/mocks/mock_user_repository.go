package mocks

import (
	contracts_repositories "gormgoskeleton/src/application/contracts/repositories"
	"gormgoskeleton/src/domain/models"

	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

type DBUser struct {
	Name   string
	Email  string
	Phone  string
	Status string
	ID     int
}

var _ contracts_repositories.IUserRepository = (*MockUserRepository)(nil)

func (m *MockUserRepository) Create(input models.UserCreate) (*models.User, error) {
	args := m.Called(input)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) Update(id uint, input models.UserUpdate) (*models.User, error) {
	args := m.Called(id, input)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetByID(id uint) (*models.User, error) {
	args := m.Called(id)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserRepository) GetAll(payload *map[string]string, skip *int, limit *int) ([]models.User, error) {
	args := m.Called()
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockUserRepository) CreateWithPassword(input models.UserAndPasswordCreate) (*models.User, error) {
	args := m.Called(input)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetUserWithRole(id uint) (*models.UserWithRole, error) {
	args := m.Called(id)
	return args.Get(0).(*models.UserWithRole), args.Error(1)
}
