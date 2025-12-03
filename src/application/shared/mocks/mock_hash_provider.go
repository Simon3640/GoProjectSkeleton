package mocks

import (
	contractsProviders "goprojectskeleton/src/application/contracts/providers"
	application_errors "goprojectskeleton/src/application/shared/errors"

	"github.com/stretchr/testify/mock"
)

type MockHashProvider struct {
	mock.Mock
}

var _ contractsProviders.IHashProvider = (*MockHashProvider)(nil)

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

func (mhp *MockHashProvider) OneTimeToken() (string, []byte, *application_errors.ApplicationError) {
	args := mhp.Called()
	errorArg := args.Get(2)
	if errorArg != nil {
		return args.String(0), args.Get(1).([]byte), errorArg.(*application_errors.ApplicationError)
	}
	return args.String(0), args.Get(1).([]byte), nil
}

func (mhp *MockHashProvider) HashOneTimeToken(token string) []byte {
	args := mhp.Called(token)
	return args.Get(0).([]byte)
}

func (mhp *MockHashProvider) ValidateOneTimeToken(hashedToken []byte, token string) bool {
	args := mhp.Called(hashedToken, token)
	return args.Bool(0)
}

func (mhp *MockHashProvider) GenerateOTP() (string, []byte, *application_errors.ApplicationError) {
	args := mhp.Called()
	errorArg := args.Get(2)
	if errorArg != nil {
		return args.String(0), args.Get(1).([]byte), errorArg.(*application_errors.ApplicationError)
	}
	return args.String(0), args.Get(1).([]byte), nil
}

func (mhp *MockHashProvider) ValidateOTP(hashedOTP []byte, otp string) bool {
	args := mhp.Called(hashedOTP, otp)
	return args.Bool(0)
}
