package auth

import (
	"context"
	"testing"
	"time"

	"gormgoskeleton/src/application/shared/locales"
	"gormgoskeleton/src/application/shared/mocks"
	dto_mocks "gormgoskeleton/src/application/shared/mocks/dtos"
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
	testHashProvider.On("HashOneTimeToken", "validOTP").Return(dto_mocks.OneTimePassword.Hash)

	testOTPRepository.On(
		"GetByPasswordHash", dto_mocks.OneTimePassword.Hash,
	).Return(&dto_mocks.OneTimePassword, nil)

	testUserRepository.On(
		"GetUserWithRole",
		dto_mocks.OneTimePassword.UserID,
	).Return(&dto_mocks.UserWithRole, nil)

	testJWTProvider.On(
		"GenerateAccessToken",
		ctx,
		string(rune(dto_mocks.OneTimePassword.UserID)),
		mock.Anything,
	).Return("newAccessToken", time.Now().Add(1*time.Hour), nil)

	testJWTProvider.On(
		"GenerateRefreshToken",
		ctx,
		string(rune(dto_mocks.OneTimePassword.UserID)),
	).Return("newRefreshToken", time.Now().Add(24*time.Hour), nil)

	testOTPRepository.On(
		"Update",
		mock.Anything,
		mock.AnythingOfType("dtos.OneTimePasswordUpdate"),
	).Return(&dto_mocks.OneTimePassword, nil)

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
	testHashProvider.On("HashOneTimeToken", "invalidOTP").Return(dto_mocks.ExpiredOneTimePassword.Hash)

	testOTPRepository.On(
		"GetByPasswordHash", dto_mocks.ExpiredOneTimePassword.Hash,
	).Return(&dto_mocks.ExpiredOneTimePassword, nil)

	result := authOTPUseCase.Execute(ctx, locales.EN_US, "invalidOTP")

	assert.NotNil(result)
	assert.False(result.IsSuccess())
	assert.Equal(result.GetStatusCode(), status.Unauthorized)
	assert.NotNil(result.Details)
}
