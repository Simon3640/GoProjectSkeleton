package authusecases

import (
	"context"
	"testing"
	"time"

	shareddtos "github.com/simon3640/goprojectskeleton/src/application/shared/DTOs"
	application_errors "github.com/simon3640/goprojectskeleton/src/application/shared/errors"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales/messages"
	providersmocks "github.com/simon3640/goprojectskeleton/src/application/shared/mocks/providers"
	email_service "github.com/simon3640/goprojectskeleton/src/application/shared/services/emails"
	email_models "github.com/simon3640/goprojectskeleton/src/application/shared/services/emails/models"
	"github.com/simon3640/goprojectskeleton/src/application/shared/status"
	"github.com/simon3640/goprojectskeleton/src/domain/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetResetPasswordSendEmailUseCase_Execute_Success(t *testing.T) {
	assert := assert.New(t)

	ctx := context.Background()
	testLogger := new(providersmocks.MockLoggerProvider)
	mockRenderProvider := new(providersmocks.MockRenderProvider[email_models.ResetPasswordEmailData])
	mockEmailProvider := new(providersmocks.MockEmailProvider)

	userStatus := models.UserStatusActive
	testUser := models.User{
		UserBase: models.UserBase{
			Name:   "Test User",
			Email:  "test@example.com",
			Phone:  "1234567890",
			Status: &userStatus,
		},
		DBBaseModel: models.DBBaseModel{
			ID:        1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	token := "test-reset-token-123"
	input := shareddtos.OneTimeTokenUser{
		User:  testUser,
		Token: token,
	}

	// Mock email service
	mockRenderProvider.On("Render", mock.Anything, mock.Anything).Return("test-rendered-email", nil)
	mockEmailProvider.On("SendEmail", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	email_service.ResetPasswordEmailServiceInstance.SetUp(
		mockRenderProvider,
		mockEmailProvider,
	)

	uc := NewGetResetPasswordSendEmailUseCase(testLogger)

	result := uc.Execute(ctx, locales.EN_US, input)

	assert.NotNil(result)
	assert.True(result.IsSuccess())
	assert.Equal(true, *result.Data)
	assert.Equal(status.Success, result.StatusCode)

	mockRenderProvider.AssertExpectations(t)
	mockEmailProvider.AssertExpectations(t)
}

func TestGetResetPasswordSendEmailUseCase_Execute_ErrorRenderingEmail(t *testing.T) {
	assert := assert.New(t)

	ctx := context.Background()
	testLogger := new(providersmocks.MockLoggerProvider)
	mockRenderProvider := new(providersmocks.MockRenderProvider[email_models.ResetPasswordEmailData])
	mockEmailProvider := new(providersmocks.MockEmailProvider)

	userStatus := models.UserStatusActive
	testUser := models.User{
		UserBase: models.UserBase{
			Name:   "Test User",
			Email:  "test@example.com",
			Phone:  "1234567890",
			Status: &userStatus,
		},
		DBBaseModel: models.DBBaseModel{
			ID:        1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	token := "test-reset-token-123"
	input := shareddtos.OneTimeTokenUser{
		User:  testUser,
		Token: token,
	}

	appError := application_errors.NewApplicationError(
		status.InternalError,
		messages.MessageKeysInstance.SOMETHING_WENT_WRONG,
		"Error rendering email template",
	)

	// Mock email service - error al renderizar
	mockRenderProvider.On("Render", mock.Anything, mock.Anything).Return("", appError)

	email_service.ResetPasswordEmailServiceInstance.SetUp(
		mockRenderProvider,
		mockEmailProvider,
	)

	uc := NewGetResetPasswordSendEmailUseCase(testLogger)

	result := uc.Execute(ctx, locales.EN_US, input)

	assert.NotNil(result)
	assert.True(result.HasError())
	assert.Equal(status.InternalError, result.StatusCode)

	mockRenderProvider.AssertExpectations(t)
}

func TestGetResetPasswordSendEmailUseCase_Execute_ErrorSendingEmail(t *testing.T) {
	assert := assert.New(t)

	ctx := context.Background()
	testLogger := new(providersmocks.MockLoggerProvider)
	mockRenderProvider := new(providersmocks.MockRenderProvider[email_models.ResetPasswordEmailData])
	mockEmailProvider := new(providersmocks.MockEmailProvider)

	userStatus := models.UserStatusActive
	testUser := models.User{
		UserBase: models.UserBase{
			Name:   "Test User",
			Email:  "test@example.com",
			Phone:  "1234567890",
			Status: &userStatus,
		},
		DBBaseModel: models.DBBaseModel{
			ID:        1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	token := "test-reset-token-123"
	input := shareddtos.OneTimeTokenUser{
		User:  testUser,
		Token: token,
	}

	appError := application_errors.NewApplicationError(
		status.ProviderError,
		messages.MessageKeysInstance.SOMETHING_WENT_WRONG,
		"Error sending email",
	)

	// Mock email service - error al enviar
	mockRenderProvider.On("Render", mock.Anything, mock.Anything).Return("test-rendered-email", nil)
	mockEmailProvider.On("SendEmail", mock.Anything, mock.Anything, mock.Anything).Return(appError)

	email_service.ResetPasswordEmailServiceInstance.SetUp(
		mockRenderProvider,
		mockEmailProvider,
	)

	uc := NewGetResetPasswordSendEmailUseCase(testLogger)

	result := uc.Execute(ctx, locales.EN_US, input)

	assert.NotNil(result)
	assert.True(result.HasError())
	assert.Equal(status.ProviderError, result.StatusCode)

	mockRenderProvider.AssertExpectations(t)
	mockEmailProvider.AssertExpectations(t)
}

func TestGetResetPasswordSendEmailUseCase_buildEmailData(t *testing.T) {
	assert := assert.New(t)

	testLogger := new(providersmocks.MockLoggerProvider)
	uc := NewGetResetPasswordSendEmailUseCase(testLogger)

	userStatus := models.UserStatusActive
	testUser := models.User{
		UserBase: models.UserBase{
			Name:   "Test User",
			Email:  "test@example.com",
			Phone:  "1234567890",
			Status: &userStatus,
		},
		DBBaseModel: models.DBBaseModel{
			ID:        1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	token := "test-reset-token-123"
	input := shareddtos.OneTimeTokenUser{
		User:  testUser,
		Token: token,
	}

	emailData := uc.buildEmailData(input)

	assert.Equal(testUser.Name, emailData.Name)
	assert.NotEmpty(emailData.ResetLink)
	assert.Contains(emailData.ResetLink, token)
	// Verificamos que los campos estén presentes (pueden tener valores por defecto)
	assert.NotNil(emailData)
	assert.GreaterOrEqual(emailData.ExpirationMinutes, int64(0))
}
