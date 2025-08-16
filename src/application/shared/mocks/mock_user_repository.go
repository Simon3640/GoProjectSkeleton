package mocks

import (
	contracts_repositories "gormgoskeleton/src/application/contracts/repositories"
	"gormgoskeleton/src/domain/models"

	"github.com/stretchr/testify/mock"
)

type MockUserRespository struct {
	mock.Mock
}

type DBUser struct {
	Name   string
	Email  string
	Phone  string
	Status string
	ID     int
}

var _ contracts_repositories.IUserRepository = (*MockUserRespository)(nil)

func (m *MockUserRespository) Create(input models.UserCreate) (*models.User, error) {
	args := m.Called(input)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRespository) Update(id int, input models.UserUpdate) (*models.User, error) {
	args := m.Called(id, input)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRespository) GetByID(id int) (*models.User, error) {
	args := m.Called(id)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRespository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserRespository) GetAll(payload *map[string]string, skip *int, limit *int) ([]models.User, error) {
	args := m.Called()
	return args.Get(0).([]models.User), args.Error(1)
}
