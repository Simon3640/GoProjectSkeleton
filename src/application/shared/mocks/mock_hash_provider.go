package mocks

import "github.com/stretchr/testify/mock"

type MockHashProvider struct {
	mock.Mock
}

func (mhp *MockHashProvider) HashPassword(password string) (string, error) {
	args := mhp.Called(password)
	return args.String(0), args.Error(1)
}

func (mhp *MockHashProvider) VerifyPassword(hashedPassword, password string) (bool, error) {
	args := mhp.Called(hashedPassword, password)
	return args.Bool(0), args.Error(1)
}
