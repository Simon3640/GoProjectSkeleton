package auth

import (
	"context"
	"testing"
	"time"

	contractsProviders "goprojectskeleton/src/application/contracts/providers"
	"goprojectskeleton/src/application/shared/locales"
	"goprojectskeleton/src/application/shared/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuthenticationRefreshUseCase(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()

	testLogger := new(mocks.MockLoggerProvider)
	testJWTProvider := new(mocks.MockJWTProvider)

	uc := NewAuthenticationRefreshUseCase(testLogger, testJWTProvider)

	// Valid Token Refresh
	validToken := "validAccessToken.123"
	claimsReturn := contractsProviders.JWTCLaims{
		"sub": "1",
		"typ": "refresh",
		"exp": float64(time.Now().Add(1 * time.Hour).Unix()),
	}
	testJWTProvider.On("ParseTokenAndValidate", validToken).Return(claimsReturn, nil)

	testJWTProvider.On("GenerateAccessToken", ctx, "1", mock.Anything).Return("newAccessToken", time.Now().Add(1*time.Hour), nil)

	testJWTProvider.On("GenerateRefreshToken", ctx, "1").Return("newRefreshToken", time.Now().Add(24*time.Hour), nil)
	result := uc.Execute(ctx, locales.EN_US, validToken)

	assert.NotNil(result)
	assert.True(result.IsSuccess())
	assert.Equal("newAccessToken", result.Data.AccessToken)
	assert.Equal("newRefreshToken", result.Data.RefreshToken)
	assert.NotNil(result.Data.AccessTokenExpiresAt)
	assert.NotNil(result.Data.RefreshTokenExpiresAt)
}
