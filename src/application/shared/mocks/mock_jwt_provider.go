package mocks

import (
	"context"
	"gormgoskeleton/src/application/contracts"
	application_errors "gormgoskeleton/src/application/shared/errors"
	"time"

	"github.com/stretchr/testify/mock"
)

type MockJWTProvider struct {
	mock.Mock
}

var _ contracts.IJWTProvider = (*MockJWTProvider)(nil)

func (m *MockJWTProvider) GenerateAccessToken(ctx context.Context, userID string, claims contracts.JWTCLaims) (string, time.Time, *application_errors.ApplicationError) {
	args := m.Called(ctx, userID, claims)
	return args.String(0), args.Get(1).(time.Time), args.Get(2).(*application_errors.ApplicationError)
}

func (m *MockJWTProvider) GenerateRefreshToken(ctx context.Context, userID string) (string, time.Time, *application_errors.ApplicationError) {
	args := m.Called(ctx, userID)
	return args.String(0), args.Get(1).(time.Time), args.Get(2).(*application_errors.ApplicationError)
}

func (m *MockJWTProvider) ParseTokenAndValidate(tokenString string) (contracts.JWTCLaims, *application_errors.ApplicationError) {
	args := m.Called(tokenString)
	return args.Get(0).(contracts.JWTCLaims), args.Get(1).(*application_errors.ApplicationError)
}
