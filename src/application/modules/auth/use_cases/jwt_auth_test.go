package authusecases

import (
	"context"
	"testing"
	"time"

	dtos "goprojectskeleton/src/application/shared/DTOs"
	"goprojectskeleton/src/application/shared/locales"
	"goprojectskeleton/src/application/shared/mocks"
	dtomocks "goprojectskeleton/src/application/shared/mocks/dtos"
	"goprojectskeleton/src/application/shared/status"
	"goprojectskeleton/src/domain/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuthenticationUseCase(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()

	testLogger := new(mocks.MockLoggerProvider)
	testJWTProvider := new(mocks.MockJWTProvider)
	testHashProvider := new(mocks.MockHashProvider)
	testPasswordRepository := new(mocks.MockPasswordRepository)
	testUserRepository := new(mocks.MockUserRepository)
	testOTPRepository := new(mocks.MockOneTimePasswordRepository)

	uc := NewAuthenticateUseCase(testLogger, testPasswordRepository, testUserRepository, testOTPRepository, testHashProvider, testJWTProvider)

	// Valid User Authentication
	userCredentials := dtos.UserCredentials{
		Email:    "user@example.com",
		Password: "plainPassword",
	}
	passwordExpiresAt := time.Now().Add(30 * 24 * time.Hour)
	passwordBase := models.PasswordBase{
		UserID:    uint(1),
		ExpiresAt: &passwordExpiresAt,
		IsActive:  true,
		Hash:      "hashedPassword123",
	}
	testPasswordRepository.On("GetActivePassword", "user@example.com").Return(&models.Password{
		PasswordBase: passwordBase,
		ID:           uint(1),
	}, nil)
	testHashProvider.On("VerifyPassword", passwordBase.Hash, userCredentials.Password).Return(true, nil)
	testJWTProvider.On("GenerateAccessToken", ctx, "1", mock.Anything).Return("accessToken", time.Now().Add(1*time.Hour), nil)
	testJWTProvider.On("GenerateRefreshToken", ctx, "1").Return("refreshToken", time.Now().Add(24*time.Hour), nil)
	testUserRepository.On("GetUserWithRole", uint(1)).Return(&dtomocks.UserWithRole, nil)
	result := uc.Execute(ctx, locales.EN_US, userCredentials)
	assert.NotNil(result)
	assert.True(result.IsSuccess())
	assert.Equal("accessToken", result.Data.AccessToken)
	assert.Equal("refreshToken", result.Data.RefreshToken)
}

//TODO: Add test for OTP when enabled

func TestAuthenticationUseCase_InvalidCredentials(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()

	testLogger := new(mocks.MockLoggerProvider)
	testJWTProvider := new(mocks.MockJWTProvider)
	testHashProvider := new(mocks.MockHashProvider)
	testPasswordRepository := new(mocks.MockPasswordRepository)
	testUserRepository := new(mocks.MockUserRepository)
	testOTPRepository := new(mocks.MockOneTimePasswordRepository)

	uc := NewAuthenticateUseCase(testLogger, testPasswordRepository, testUserRepository, testOTPRepository, testHashProvider, testJWTProvider)

	// Invalid User Authentication
	userCredentials := dtos.UserCredentials{
		Email:    "user@example.com",
		Password: "wrongPassword",
	}
	passwordExpiresAt := time.Now().Add(30 * 24 * time.Hour)
	passwordBase := models.PasswordBase{
		UserID:    uint(1),
		ExpiresAt: &passwordExpiresAt,
		IsActive:  true,
		Hash:      "hashedPassword123",
	}
	testPasswordRepository.On("GetActivePassword", "user@example.com").Return(&models.Password{
		PasswordBase: passwordBase,
		ID:           uint(1),
	}, nil)
	testHashProvider.On("VerifyPassword", passwordBase.Hash, userCredentials.Password).Return(false, nil)
	testUserRepository.On("GetUserWithRole", uint(1)).Return(&dtomocks.UserWithRole, nil)
	result := uc.Execute(ctx, locales.EN_US, userCredentials)
	assert.NotNil(result)
	assert.False(result.IsSuccess())
	assert.Equal(status.NotFound, result.StatusCode)
	assert.True(result.HasError())
	testPasswordRepository.AssertExpectations(t)
	testHashProvider.AssertExpectations(t)
	testJWTProvider.AssertExpectations(t)
}
