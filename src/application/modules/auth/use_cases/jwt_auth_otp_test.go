package authusecases

import (
	"context"
	"strconv"
	"testing"
	"time"

	"goprojectskeleton/src/application/shared/locales"
	dtomocks "goprojectskeleton/src/application/shared/mocks/dtos"
	providersmocks "goprojectskeleton/src/application/shared/mocks/providers"
	repositoriesmocks "goprojectskeleton/src/application/shared/mocks/repositories"
	"goprojectskeleton/src/application/shared/status"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuthenticateOTPUseCase_Valid(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()

	testLogger := new(providersmocks.MockLoggerProvider)
	testUserRepository := new(repositoriesmocks.MockUserRepository)
	testOTPRepository := new(repositoriesmocks.MockOneTimePasswordRepository)
	testJWTProvider := new(providersmocks.MockJWTProvider)
	testHashProvider := new(providersmocks.MockHashProvider)

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
		strconv.FormatUint(uint64(dtomocks.OneTimePassword.UserID), 10),
		mock.Anything,
	).Return("newAccessToken", time.Now().Add(1*time.Hour), nil)

	testJWTProvider.On(
		"GenerateRefreshToken",
		ctx,
		strconv.FormatUint(uint64(dtomocks.OneTimePassword.UserID), 10),
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

	testLogger := new(providersmocks.MockLoggerProvider)
	testUserRepository := new(repositoriesmocks.MockUserRepository)
	testOTPRepository := new(repositoriesmocks.MockOneTimePasswordRepository)
	testJWTProvider := new(providersmocks.MockJWTProvider)
	testHashProvider := new(providersmocks.MockHashProvider)

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
