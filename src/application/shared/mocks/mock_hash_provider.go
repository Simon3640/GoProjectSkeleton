package mocks

import (
	contracts_providers "gormgoskeleton/src/application/contracts/providers"
	application_errors "gormgoskeleton/src/application/shared/errors"

	"github.com/stretchr/testify/mock"
)

type MockHashProvider struct {
	mock.Mock
}

var _ contracts_providers.IHashProvider = (*MockHashProvider)(nil)

func (mhp *MockHashProvider) HashPassword(password string) (string, *application_errors.ApplicationError) {
	args := mhp.Called(password)
	errorArg := args.Get(1)
	if errorArg != nil {
		return args.String(0), errorArg.(*application_errors.ApplicationError)
	}
	return args.String(0), nil
}

func (mhp *MockHashProvider) VerifyPassword(hashedPassword string, password string) (bool, *application_errors.ApplicationError) {
	args := mhp.Called(hashedPassword, password)
	errorArg := args.Get(1)
	if errorArg != nil {
		return args.Bool(0), errorArg.(*application_errors.ApplicationError)
	}
	return args.Bool(0), nil
}
