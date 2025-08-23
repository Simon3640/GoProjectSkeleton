package mocks

import (
	"gormgoskeleton/src/application/contracts"
	application_errors "gormgoskeleton/src/application/shared/errors"

	"github.com/stretchr/testify/mock"
)

type MockHashProvider struct {
	mock.Mock
}

var _ contracts.IHashProvider = (*MockHashProvider)(nil)

func (mhp *MockHashProvider) HashPassword(password string) (string, *application_errors.ApplicationError) {
	args := mhp.Called(password)
	return args.String(0), args.Get(1).(*application_errors.ApplicationError)
}

func (mhp *MockHashProvider) VerifyPassword(hashedPassword string, password string) (bool, *application_errors.ApplicationError) {
	args := mhp.Called(hashedPassword, password)
	return args.Bool(0), args.Get(1).(*application_errors.ApplicationError)
}
