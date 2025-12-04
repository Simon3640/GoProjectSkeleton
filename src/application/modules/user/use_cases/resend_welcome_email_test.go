package userusecases

import (
	"context"
	"testing"
	"time"

	dtos "goprojectskeleton/src/application/shared/DTOs"
	application_errors "goprojectskeleton/src/application/shared/errors"
	"goprojectskeleton/src/application/shared/locales"
	providersmocks "goprojectskeleton/src/application/shared/mocks/providers"
	repositoriesmocks "goprojectskeleton/src/application/shared/mocks/repositories"
	email_service "goprojectskeleton/src/application/shared/services/emails"
	email_models "goprojectskeleton/src/application/shared/services/emails/models"
	"goprojectskeleton/src/application/shared/status"
	"goprojectskeleton/src/domain/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestResendWelcomeEmailUseCase_Execute_Success(t *testing.T) {
	assert := assert.New(t)

	ctx := context.Background()

	testLogger := new(providersmocks.MockLoggerProvider)
	testHashProvider := new(providersmocks.MockHashProvider)
	testUserRepository := new(repositoriesmocks.MockUserRepository)
	testTokenRepository := new(repositoriesmocks.MockOneTimeTokenRepository)

	mockRenderProvider := new(providersmocks.MockRenderProvider[email_models.NewUserEmailData])
	mockEmailProvider := new(providersmocks.MockEmailProvider)

	email := "test@example.com"
	userStatus := models.UserStatusPending
	testUser := &models.User{
		UserBase: models.UserBase{
			Name:   "Test User",
			Email:  email,
			Phone:  "1234567890",
			Status: &userStatus,
			RoleID: 2,
		},
		DBBaseModel: models.DBBaseModel{
			ID:        1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	// Mock GetByEmailOrPhone
	testUserRepository.On("GetByEmailOrPhone", email).Return(testUser, nil)

	// Mock CreateOneTimeToken
	testHashProvider.On("OneTimeToken").Return("test-token", []byte("hash"), nil)
	testTokenRepository.On("Create", mock.Anything).Return(&models.OneTimeToken{
		OneTimeTokenBase: models.OneTimeTokenBase{
			UserID:  1,
			Purpose: models.OneTimeTokenPurposeEmailVerify,
			Hash:    []byte("hash"),
			IsUsed:  false,
			Expires: time.Now().Add(24 * time.Hour),
		},
		DBBaseModel: models.DBBaseModel{
			ID: 1,
		},
	}, nil)
	mockRenderProvider.On("Render", mock.Anything, mock.Anything).Return("test-rendered", nil)
	mockEmailProvider.On("SendEmail", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	email_service.RegisterUserEmailServiceInstance.SetUp(
		mockRenderProvider,
		mockEmailProvider,
	)

	uc := NewResendWelcomeEmailUseCase(
		testLogger,
		testHashProvider,
		testUserRepository,
		testTokenRepository,
	)

	input := dtos.ResendWelcomeEmailRequest{
		Email: email,
	}

	result := uc.Execute(ctx, locales.EN_US, input)

	assert.NotNil(result)
	assert.False(result.HasError())
	assert.NotNil(result.Data)
	assert.True(*result.Data)
	assert.Equal(status.Success, result.StatusCode)
}

func TestResendWelcomeEmailUseCase_Execute_InvalidEmail(t *testing.T) {
	assert := assert.New(t)

	ctx := context.Background()

	testLogger := new(providersmocks.MockLoggerProvider)
	testHashProvider := new(providersmocks.MockHashProvider)
	testUserRepository := new(repositoriesmocks.MockUserRepository)
	testTokenRepository := new(repositoriesmocks.MockOneTimeTokenRepository)

	uc := NewResendWelcomeEmailUseCase(
		testLogger,
		testHashProvider,
		testUserRepository,
		testTokenRepository,
	)

	input := dtos.ResendWelcomeEmailRequest{
		Email: "invalid-email",
	}

	result := uc.Execute(ctx, locales.EN_US, input)

	assert.NotNil(result)
	assert.True(result.HasError())
	assert.Equal(status.InvalidInput, result.StatusCode)
}

func TestResendWelcomeEmailUseCase_Execute_UserAlreadyVerified(t *testing.T) {
	assert := assert.New(t)

	ctx := context.Background()

	testLogger := new(providersmocks.MockLoggerProvider)
	testHashProvider := new(providersmocks.MockHashProvider)
	testUserRepository := new(repositoriesmocks.MockUserRepository)
	testTokenRepository := new(repositoriesmocks.MockOneTimeTokenRepository)

	email := "test@example.com"
	userStatus := models.UserStatusActive
	testUser := &models.User{
		UserBase: models.UserBase{
			Name:   "Test User",
			Email:  email,
			Phone:  "1234567890",
			Status: &userStatus,
			RoleID: 2,
		},
		DBBaseModel: models.DBBaseModel{
			ID:        1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	// Mock GetByEmailOrPhone
	testUserRepository.On("GetByEmailOrPhone", email).Return(testUser, nil)

	uc := NewResendWelcomeEmailUseCase(
		testLogger,
		testHashProvider,
		testUserRepository,
		testTokenRepository,
	)

	input := dtos.ResendWelcomeEmailRequest{
		Email: email,
	}

	result := uc.Execute(ctx, locales.EN_US, input)

	assert.NotNil(result)
	assert.True(result.HasError())
	assert.Equal(status.Conflict, result.StatusCode)
}

func TestResendWelcomeEmailUseCase_Execute_UserNotFound(t *testing.T) {
	assert := assert.New(t)

	ctx := context.Background()

	testLogger := new(providersmocks.MockLoggerProvider)
	testHashProvider := new(providersmocks.MockHashProvider)
	testUserRepository := new(repositoriesmocks.MockUserRepository)
	testTokenRepository := new(repositoriesmocks.MockOneTimeTokenRepository)

	email := "notfound@example.com"

	// Mock GetByEmailOrPhone returning nil (user not found)
	appErr := application_errors.NewApplicationError(
		status.NotFound,
		"RESOURCE_NOT_FOUND",
		"User not found",
	)
	testUserRepository.On("GetByEmailOrPhone", email).Return(nil, appErr)

	uc := NewResendWelcomeEmailUseCase(
		testLogger,
		testHashProvider,
		testUserRepository,
		testTokenRepository,
	)

	input := dtos.ResendWelcomeEmailRequest{
		Email: email,
	}

	result := uc.Execute(ctx, locales.EN_US, input)

	assert.NotNil(result)
	assert.True(result.HasError())
	assert.Equal(status.NotFound, result.StatusCode)
}

func TestResendWelcomeEmailUseCase_Execute_RepositoryError(t *testing.T) {
	assert := assert.New(t)

	ctx := context.Background()

	testLogger := new(providersmocks.MockLoggerProvider)
	testHashProvider := new(providersmocks.MockHashProvider)
	testUserRepository := new(repositoriesmocks.MockUserRepository)
	testTokenRepository := new(repositoriesmocks.MockOneTimeTokenRepository)

	email := "test@example.com"

	// Mock GetByEmailOrPhone returning error
	appErr := application_errors.NewApplicationError(
		status.InternalError,
		"RESOURCE_NOT_FOUND",
		"Database error",
	)
	testUserRepository.On("GetByEmailOrPhone", email).Return(nil, appErr)

	uc := NewResendWelcomeEmailUseCase(
		testLogger,
		testHashProvider,
		testUserRepository,
		testTokenRepository,
	)

	input := dtos.ResendWelcomeEmailRequest{
		Email: email,
	}

	result := uc.Execute(ctx, locales.EN_US, input)

	assert.NotNil(result)
	assert.True(result.HasError())
	assert.Equal(status.InternalError, result.StatusCode)
}
