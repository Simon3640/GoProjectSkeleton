package authmocks

import (
	"context"
	"time"

	authcontracts "github.com/simon3640/goprojectskeleton/src/application/modules/auth/contracts"
	application_errors "github.com/simon3640/goprojectskeleton/src/application/shared/errors"

	"github.com/stretchr/testify/mock"
)

// MockJWTProvider is the mock implementation of the JWTProvider interface
type MockJWTProvider struct {
	mock.Mock
}

var _ authcontracts.IJWTProvider = (*MockJWTProvider)(nil)

// GenerateAccessToken generates an access token
func (m *MockJWTProvider) GenerateAccessToken(ctx context.Context, userID string, claims authcontracts.JWTCLaims) (string, time.Time, *application_errors.ApplicationError) {
	args := m.Called(ctx, userID, claims)
	errorArg := args.Get(2)
	if errorArg != nil {
		return args.String(0), time.Time{}, errorArg.(*application_errors.ApplicationError)
	}
	return args.String(0), args.Get(1).(time.Time), nil
}

// GenerateRefreshToken generates a refresh token
func (m *MockJWTProvider) GenerateRefreshToken(ctx context.Context, userID string) (string, time.Time, *application_errors.ApplicationError) {
	args := m.Called(ctx, userID)
	errorArg := args.Get(2)
	if errorArg != nil {
		return args.String(0), time.Time{}, errorArg.(*application_errors.ApplicationError)
	}
	return args.String(0), args.Get(1).(time.Time), nil
}

// ParseTokenAndValidate parses and validates a token
func (m *MockJWTProvider) ParseTokenAndValidate(tokenString string) (authcontracts.JWTCLaims, *application_errors.ApplicationError) {
	args := m.Called(tokenString)
	errorArg := args.Get(1)
	if errorArg != nil {
		return args.Get(0).(authcontracts.JWTCLaims), errorArg.(*application_errors.ApplicationError)
	}
	return args.Get(0).(authcontracts.JWTCLaims), nil
}
