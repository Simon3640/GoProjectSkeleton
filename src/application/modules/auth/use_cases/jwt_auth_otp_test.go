package authusecases

import (
	"context"
	"strconv"
	"testing"
	"time"

	authmocks "github.com/simon3640/goprojectskeleton/src/application/modules/auth/mocks"
	app_context "github.com/simon3640/goprojectskeleton/src/application/shared/context"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales"
	dtomocks "github.com/simon3640/goprojectskeleton/src/application/shared/mocks/dtos"
	providersmocks "github.com/simon3640/goprojectskeleton/src/application/shared/mocks/providers"
	"github.com/simon3640/goprojectskeleton/src/application/shared/status"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuthenticateOTPUseCase_Valid(t *testing.T) {
	assert := assert.New(t)
	ctx := &app_context.AppContext{Context: context.Background()}

	testUserRepository := new(authmocks.MockUserRepository)
	testOTPRepository := new(authmocks.MockOneTimePasswordRepository)
	testJWTProvider := new(authmocks.MockJWTProvider)
	testHashProvider := new(providersmocks.MockHashProvider)

	authOTPUseCase := NewAuthenticateOTPUseCase(
		testUserRepository,
		testOTPRepository,
		testHashProvider,
		testJWTProvider,
	)

	// Mocking Methods
	testHashProvider.On("HashOneTimeToken", "validOTP").Return(authmocks.OneTimePassword.Hash)

	testOTPRepository.On(
		"GetByPasswordHash", authmocks.OneTimePassword.Hash,
	).Return(&authmocks.OneTimePassword, nil)

	testUserRepository.On(
		"GetUserWithRole",
		authmocks.OneTimePassword.UserID,
	).Return(&dtomocks.UserWithRole, nil)

	testJWTProvider.On(
		"GenerateAccessToken",
		ctx,
		strconv.FormatUint(uint64(authmocks.OneTimePassword.UserID), 10),
		mock.Anything,
	).Return("newAccessToken", time.Now().Add(1*time.Hour), nil)

	testJWTProvider.On(
		"GenerateRefreshToken",
		ctx,
		strconv.FormatUint(uint64(authmocks.OneTimePassword.UserID), 10),
	).Return("newRefreshToken", time.Now().Add(24*time.Hour), nil)

	testOTPRepository.On(
		"Update",
		mock.Anything,
		mock.AnythingOfType("authdtos.OneTimePasswordUpdate"),
	).Return(&authmocks.OneTimePassword, nil)

	result := authOTPUseCase.Execute(ctx, locales.EN_US, "validOTP")

	assert.NotNil(result)
	assert.True(result.IsSuccess())

}

func TestAuthenticateOTPUseCase_InvalidOTP(t *testing.T) {
	assert := assert.New(t)
	ctx := &app_context.AppContext{Context: context.Background()}

	testUserRepository := new(authmocks.MockUserRepository)
	testOTPRepository := new(authmocks.MockOneTimePasswordRepository)
	testJWTProvider := new(authmocks.MockJWTProvider)
	testHashProvider := new(providersmocks.MockHashProvider)

	authOTPUseCase := NewAuthenticateOTPUseCase(
		testUserRepository,
		testOTPRepository,
		testHashProvider,
		testJWTProvider,
	)

	// Mocking Methods
	testHashProvider.On("HashOneTimeToken", "invalidOTP").Return(authmocks.ExpiredOneTimePassword.Hash)

	testOTPRepository.On(
		"GetByPasswordHash", authmocks.ExpiredOneTimePassword.Hash,
	).Return(&authmocks.ExpiredOneTimePassword, nil)

	result := authOTPUseCase.Execute(ctx, locales.EN_US, "invalidOTP")

	assert.NotNil(result)
	assert.False(result.IsSuccess())
	assert.Equal(result.GetStatusCode(), status.Unauthorized)
	assert.NotNil(result.Details)
}
