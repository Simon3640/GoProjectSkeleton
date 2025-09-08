package auth

import (
	"context"
	"testing"
	"time"

	"gormgoskeleton/src/application/shared/locales"
	"gormgoskeleton/src/application/shared/mocks"
	dtomocks "gormgoskeleton/src/application/shared/mocks/dtos"
	"gormgoskeleton/src/application/shared/status"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuthenticateOTPUseCase_Valid(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()

	testLogger := new(mocks.MockLoggerProvider)
	testUserRepository := new(mocks.MockUserRepository)
	testOTPRepository := new(mocks.MockOneTimePasswordRepository)
	testJWTProvider := new(mocks.MockJWTProvider)
	testHashProvider := new(mocks.MockHashProvider)

	authOTPUseCase := NewAuthenticateOTPUseCase(
		testLogger,
		testUserRepository,
		testOTPRepository,
		testHashProvider,
		testJWTProvider,
	)

	// Mocking Methods
	testHashProvider.On("HashOneTimeToken", "validOTP").Return(dtomocks.OneTimePassword.Hash)

	testOTPRepository.On(
		"GetByPasswordHash", dtomocks.OneTimePassword.Hash,
	).Return(&dtomocks.OneTimePassword, nil)

	testUserRepository.On(
		"GetUserWithRole",
		dtomocks.OneTimePassword.UserID,
	).Return(&dtomocks.UserWithRole, nil)

	testJWTProvider.On(
		"GenerateAccessToken",
		ctx,
		string(rune(dtomocks.OneTimePassword.UserID)),
		mock.Anything,
	).Return("newAccessToken", time.Now().Add(1*time.Hour), nil)

	testJWTProvider.On(
		"GenerateRefreshToken",
		ctx,
		string(rune(dtomocks.OneTimePassword.UserID)),
	).Return("newRefreshToken", time.Now().Add(24*time.Hour), nil)

	testOTPRepository.On(
		"Update",
		mock.Anything,
		mock.AnythingOfType("dtos.OneTimePasswordUpdate"),
	).Return(&dtomocks.OneTimePassword, nil)

	result := authOTPUseCase.Execute(ctx, locales.EN_US, "validOTP")

	assert.NotNil(result)
	assert.True(result.IsSuccess())

}

func TestAuthenticateOTPUseCase_InvalidOTP(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()

	testLogger := new(mocks.MockLoggerProvider)
	testUserRepository := new(mocks.MockUserRepository)
	testOTPRepository := new(mocks.MockOneTimePasswordRepository)
	testJWTProvider := new(mocks.MockJWTProvider)
	testHashProvider := new(mocks.MockHashProvider)

	authOTPUseCase := NewAuthenticateOTPUseCase(
		testLogger,
		testUserRepository,
		testOTPRepository,
		testHashProvider,
		testJWTProvider,
	)

	// Mocking Methods
	testHashProvider.On("HashOneTimeToken", "invalidOTP").Return(dtomocks.ExpiredOneTimePassword.Hash)

	testOTPRepository.On(
		"GetByPasswordHash", dtomocks.ExpiredOneTimePassword.Hash,
	).Return(&dtomocks.ExpiredOneTimePassword, nil)

	result := authOTPUseCase.Execute(ctx, locales.EN_US, "invalidOTP")

	assert.NotNil(result)
	assert.False(result.IsSuccess())
	assert.Equal(result.GetStatusCode(), status.Unauthorized)
	assert.NotNil(result.Details)
}
