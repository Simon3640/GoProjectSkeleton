package auth

import (
	"context"
	// "errors"

	contractsProviders "gormgoskeleton/src/application/contracts/providers"
	"gormgoskeleton/src/application/shared/locales"
	"gormgoskeleton/src/application/shared/mocks"
	dtomocks "gormgoskeleton/src/application/shared/mocks/dtos"
	"gormgoskeleton/src/application/shared/status"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAuthUserCase(t *testing.T) {
	asswert := assert.New(t)

	testLogger := new(mocks.MockLoggerProvider)
	testUserRepository := new(mocks.MockUserRepository)
	testJWTProvider := new(mocks.MockJWTProvider)

	authUserUseCase := NewAuthUserUseCase(
		testLogger,
		testUserRepository,
		testJWTProvider,
	)

	validToken := "validToken.123"
	claimsReturn := contractsProviders.JWTCLaims{
		"sub": "1",
		"typ": "access",
		"exp": float64(time.Now().Add(1 * time.Hour).Unix()),
	}

	testJWTProvider.On("ParseTokenAndValidate", validToken).Return(claimsReturn, nil)
	testUserRepository.On("GetUserWithRole", uint(1)).Return(&dtomocks.UserWithRole, nil)

	result := authUserUseCase.Execute(context.Background(), locales.EN_US, validToken)

	asswert.NotNil(result)
	asswert.True(result.IsSuccess())
	asswert.Equal(uint(1), result.Data.ID)
	asswert.Equal(dtomocks.UserWithRole.Name, result.Data.Name)
	asswert.Equal(dtomocks.UserWithRole.Email, result.Data.Email)
	asswert.Equal(dtomocks.UserWithRole.Phone, result.Data.Phone)
	asswert.Equal(dtomocks.UserWithRole.Status, result.Data.Status)
	asswert.Equal(dtomocks.UserWithRole.RoleID, result.Data.RoleID)

}

func TestAuthUserCase_InvalidToken(t *testing.T) {
	assert := assert.New(t)

	testLogger := new(mocks.MockLoggerProvider)
	testUserRepository := new(mocks.MockUserRepository)
	testJWTProvider := new(mocks.MockJWTProvider)

	authUserUseCase := NewAuthUserUseCase(
		testLogger,
		testUserRepository,
		testJWTProvider,
	)

	invalidToken := "invalidToken"

	result := authUserUseCase.Execute(context.Background(), locales.EN_US, invalidToken)

	assert.NotNil(result)
	assert.True(result.HasError())
	assert.Equal(status.Unauthorized, result.StatusCode)
	assert.NotNil(result.Error)
}

func TestAuthUserCase_ExpiredToken(t *testing.T) {
	assert := assert.New(t)

	testLogger := new(mocks.MockLoggerProvider)
	testUserRepository := new(mocks.MockUserRepository)
	testJWTProvider := new(mocks.MockJWTProvider)

	authUserUseCase := NewAuthUserUseCase(
		testLogger,
		testUserRepository,
		testJWTProvider,
	)

	expiredToken := "expiredToken.123"
	claimsReturn := contractsProviders.JWTCLaims{
		"sub": "1",
		"typ": "access",
		"exp": float64(time.Now().Add(-1 * time.Hour).Unix()),
	}

	testJWTProvider.On("ParseTokenAndValidate", expiredToken).Return(claimsReturn, nil)

	result := authUserUseCase.Execute(context.Background(), locales.EN_US, expiredToken)

	assert.NotNil(result)
	assert.True(result.HasError())
	assert.Equal(status.Unauthorized, result.StatusCode)
	assert.NotNil(result.Error)
}

func TestAuthUserCase_NoAccessToken(t *testing.T) {
	assert := assert.New(t)

	testLogger := new(mocks.MockLoggerProvider)
	testUserRepository := new(mocks.MockUserRepository)
	testJWTProvider := new(mocks.MockJWTProvider)

	authUserUseCase := NewAuthUserUseCase(
		testLogger,
		testUserRepository,
		testJWTProvider,
	)

	noAccessToken := "noAccessToken.123"
	claimsReturn := contractsProviders.JWTCLaims{
		"sub": "1",
		"typ": "refresh",
		"exp": float64(time.Now().Add(1 * time.Hour).Unix()),
	}

	testJWTProvider.On("ParseTokenAndValidate", noAccessToken).Return(claimsReturn, nil)

	result := authUserUseCase.Execute(context.Background(), locales.EN_US, noAccessToken)

	assert.NotNil(result)
	assert.True(result.HasError())
	assert.Equal(status.Unauthorized, result.StatusCode)
	assert.NotNil(result.Error)
}
