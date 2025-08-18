package mocks

import (
	"context"
	"gormgoskeleton/src/application/contracts"
	"time"

	"github.com/stretchr/testify/mock"
)

type MockJWTProvider struct {
	mock.Mock
}

var _ contracts.IJWTProvider = (*MockJWTProvider)(nil)

func (m *MockJWTProvider) GenerateAccessToken(ctx context.Context, userID string, claims contracts.JWTCLaims) (string, time.Time, error) {
	args := m.Called(ctx, userID, claims)
	return args.String(0), args.Get(1).(time.Time), args.Error(2)
}

func (m *MockJWTProvider) GenerateRefreshToken(ctx context.Context, userID string) (string, time.Time, error) {
	args := m.Called(ctx, userID)
	return args.String(0), args.Get(1).(time.Time), args.Error(2)
}

func (m *MockJWTProvider) ParseTokenAndValidate(tokenString string) (contracts.JWTCLaims, error) {
	args := m.Called(tokenString)
	if claims, ok := args.Get(0).(contracts.JWTCLaims); ok {
		return claims, args.Error(1)
	}
	return nil, args.Error(1)
}
