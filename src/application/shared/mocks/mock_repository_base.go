package mocks

import (
	contracts_repositories "goprojectskeleton/src/application/contracts/repositories"
	application_errors "goprojectskeleton/src/application/shared/errors"
	domain_utils "goprojectskeleton/src/domain/utils"

	"github.com/stretchr/testify/mock"
)

type MockRepositoryBase[CreateDomainModel any, UpdateDomainModel any, DomainModel any, DBModel any] struct {
	mock.Mock
}

var _ contracts_repositories.IRepositoryBase[any, any, any, any] = (*MockRepositoryBase[any, any, any, any])(nil)

func (m *MockRepositoryBase[CreateDomainModel, UpdateDomainModel, DomainModel, DBModel]) Create(entity CreateDomainModel) (*DomainModel, *application_errors.ApplicationError) {
	args := m.Called(entity)
	errorArg := args.Get(1)
	if errorArg != nil {
		return args.Get(0).(*DomainModel), errorArg.(*application_errors.ApplicationError)
	}
	return args.Get(0).(*DomainModel), nil
}

func (m *MockRepositoryBase[CreateDomainModel, UpdateDomainModel, DomainModel, DBModel]) GetByID(id uint) (*DomainModel, *application_errors.ApplicationError) {
	args := m.Called(id)
	errorArg := args.Get(1)
	if errorArg != nil {
		return args.Get(0).(*DomainModel), errorArg.(*application_errors.ApplicationError)
	}
	return args.Get(0).(*DomainModel), nil
}

func (m *MockRepositoryBase[CreateDomainModel, UpdateDomainModel, DomainModel, DBModel]) Update(id uint, entity UpdateDomainModel) (*DomainModel, *application_errors.ApplicationError) {
	args := m.Called(id, entity)
	errorArg := args.Get(1)
	if errorArg != nil {
		return args.Get(0).(*DomainModel), errorArg.(*application_errors.ApplicationError)
	}
	return args.Get(0).(*DomainModel), nil
}

func (m *MockRepositoryBase[CreateDomainModel, UpdateDomainModel, DomainModel, DBModel]) Delete(id uint) *application_errors.ApplicationError {
	args := m.Called(id)
	errorArg := args.Get(0)
	if errorArg != nil {
		return errorArg.(*application_errors.ApplicationError)
	}
	return nil
}

func (m *MockRepositoryBase[CreateDomainModel, UpdateDomainModel, DomainModel, DBModel]) SoftDelete(id uint) *application_errors.ApplicationError {
	args := m.Called(id)
	errorArg := args.Get(0)
	if errorArg != nil {
		return errorArg.(*application_errors.ApplicationError)
	}
	return nil
}

func (m *MockRepositoryBase[CreateDomainModel, UpdateDomainModel, DomainModel, DBModel]) GetAll(payload *domain_utils.QueryPayloadBuilder[DomainModel], skip int, limit int) ([]DomainModel, int64, *application_errors.ApplicationError) {
	args := m.Called(payload, skip, limit)
	errorArg := args.Get(2)
	if errorArg != nil {
		return args.Get(0).([]DomainModel), args.Get(1).(int64), errorArg.(*application_errors.ApplicationError)
	}
	return args.Get(0).([]DomainModel), args.Get(1).(int64), nil
}
