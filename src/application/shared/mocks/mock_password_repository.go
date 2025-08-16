package mocks

import (
	contracts_repositories "gormgoskeleton/src/application/contracts/repositories"
	"gormgoskeleton/src/domain/models"

	"github.com/stretchr/testify/mock"
)

type MockPasswordRespository struct {
	mock.Mock
}

type DBPassword struct {
	Name   string
	Email  string
	Phone  string
	Status string
	ID     int
}

var _ contracts_repositories.IPasswordRepository = (*MockPasswordRespository)(nil)

func (m *MockPasswordRespository) Create(input models.PasswordCreate) (*models.Password, error) {
	args := m.Called(input)
	return args.Get(0).(*models.Password), args.Error(1)
}

func (m *MockPasswordRespository) Update(id int, input models.PasswordUpdate) (*models.Password, error) {
	args := m.Called(id, input)
	return args.Get(0).(*models.Password), args.Error(1)
}

func (m *MockPasswordRespository) GetByID(id int) (*models.Password, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Password), args.Error(1)
}

func (m *MockPasswordRespository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockPasswordRespository) GetAll(payload *map[string]string, skip *int, limit *int) ([]models.Password, error) {
	args := m.Called()
	return args.Get(0).([]models.Password), args.Error(1)
}

func (m *MockPasswordRespository) CleanPasswords(userID int) error {
	args := m.Called(userID)
	return args.Error(0)
}
