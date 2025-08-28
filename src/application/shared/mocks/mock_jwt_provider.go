package mocks

import (
	"context"
	contracts_providers "gormgoskeleton/src/application/contracts/providers"
	application_errors "gormgoskeleton/src/application/shared/errors"
	"time"

	"github.com/stretchr/testify/mock"
)

type MockJWTProvider struct {
	mock.Mock
}

var _ contracts_providers.IJWTProvider = (*MockJWTProvider)(nil)

func (m *MockJWTProvider) GenerateAccessToken(ctx context.Context, userID string, claims contracts_providers.JWTCLaims) (string, time.Time, *application_errors.ApplicationError) {
	args := m.Called(ctx, userID, claims)
	errorArg := args.Get(2)
	if errorArg != nil {
		return args.String(0), time.Time{}, errorArg.(*application_errors.ApplicationError)
	}
	return args.String(0), args.Get(1).(time.Time), nil
}

func (m *MockJWTProvider) GenerateRefreshToken(ctx context.Context, userID string) (string, time.Time, *application_errors.ApplicationError) {
	args := m.Called(ctx, userID)
	errorArg := args.Get(2)
	if errorArg != nil {
		return args.String(0), time.Time{}, errorArg.(*application_errors.ApplicationError)
	}
	return args.String(0), args.Get(1).(time.Time), nil
}

func (m *MockJWTProvider) ParseTokenAndValidate(tokenString string) (contracts_providers.JWTCLaims, *application_errors.ApplicationError) {
	args := m.Called(tokenString)
	errorArg := args.Get(1)
	if errorArg != nil {
		return args.Get(0).(contracts_providers.JWTCLaims), errorArg.(*application_errors.ApplicationError)
	}
	return args.Get(0).(contracts_providers.JWTCLaims), nil
}
