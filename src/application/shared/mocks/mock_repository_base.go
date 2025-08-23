package mocks

import (
	contracts_repositories "gormgoskeleton/src/application/contracts/repositories"
	application_errors "gormgoskeleton/src/application/shared/errors"

	"github.com/stretchr/testify/mock"
)

type MockRepositoryBase[CreateDomainModel any, UpdateDomainModel any, DomainModel any, DBModel any] struct {
	mock.Mock
}

var _ contracts_repositories.IRepositoryBase[any, any, any, any] = (*MockRepositoryBase[any, any, any, any])(nil)

func (m *MockRepositoryBase[CreateDomainModel, UpdateDomainModel, DomainModel, DBModel]) Create(entity CreateDomainModel) (*DomainModel, *application_errors.ApplicationError) {
	args := m.Called(entity)
	return args.Get(0).(*DomainModel), args.Get(1).(*application_errors.ApplicationError)
}

func (m *MockRepositoryBase[CreateDomainModel, UpdateDomainModel, DomainModel, DBModel]) GetByID(id uint) (*DomainModel, *application_errors.ApplicationError) {
	args := m.Called(id)
	return args.Get(0).(*DomainModel), args.Get(1).(*application_errors.ApplicationError)
}

func (m *MockRepositoryBase[CreateDomainModel, UpdateDomainModel, DomainModel, DBModel]) Update(id uint, entity UpdateDomainModel) (*DomainModel, *application_errors.ApplicationError) {
	args := m.Called(id, entity)
	return args.Get(0).(*DomainModel), args.Get(1).(*application_errors.ApplicationError)
}

func (m *MockRepositoryBase[CreateDomainModel, UpdateDomainModel, DomainModel, DBModel]) Delete(id uint) *application_errors.ApplicationError {
	args := m.Called(id)
	return args.Get(0).(*application_errors.ApplicationError)
}

func (m *MockRepositoryBase[CreateDomainModel, UpdateDomainModel, DomainModel, DBModel]) GetAll(payload *map[string]string, skip *int, limit *int) ([]DomainModel, *application_errors.ApplicationError) {
	args := m.Called()
	return args.Get(0).([]DomainModel), args.Get(1).(*application_errors.ApplicationError)
}
