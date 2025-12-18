package userusecases

import (
	"context"
	"testing"
	"time"

	app_context "github.com/simon3640/goprojectskeleton/src/application/shared/context"
	applicationerror "github.com/simon3640/goprojectskeleton/src/application/shared/errors"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales/messages"
	providersmocks "github.com/simon3640/goprojectskeleton/src/application/shared/mocks/providers"
	repositoriesmocks "github.com/simon3640/goprojectskeleton/src/application/shared/mocks/repositories"
	emailservice "github.com/simon3640/goprojectskeleton/src/application/shared/services/emails"
	emailmodels "github.com/simon3640/goprojectskeleton/src/application/shared/services/emails/models"
	"github.com/simon3640/goprojectskeleton/src/application/shared/status"
	"github.com/simon3640/goprojectskeleton/src/domain/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateUserSendEmailUseCase_Execute_Success(t *testing.T) {
	assert := assert.New(t)

	ctx := &app_context.AppContext{Context: context.Background()}

	testHashProvider := new(providersmocks.MockHashProvider)
	testTokenRepository := new(repositoriesmocks.MockOneTimeTokenRepository)

	mockRenderProvider := new(providersmocks.MockRenderProvider[emailmodels.NewUserEmailData])
	mockEmailProvider := new(providersmocks.MockEmailProvider)

	userStatus := models.UserStatusPending
	testUser := models.User{
		UserBase: models.UserBase{
			Name:   "Test User",
			Email:  "test@example.com",
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

	// Mock OneTimeToken
	token := "test-token-123"
	tokenHash := []byte("hash")
	testHashProvider.On("OneTimeToken").Return(token, tokenHash, nil)

	// Mock Create token repository
	testTokenRepository.On("Create", mock.AnythingOfType("dtos.OneTimeTokenCreate")).Return(&models.OneTimeToken{
		OneTimeTokenBase: models.OneTimeTokenBase{
			UserID:  1,
			Purpose: models.OneTimeTokenPurposeEmailVerify,
			Hash:    tokenHash,
			IsUsed:  false,
			Expires: time.Now().Add(24 * time.Hour),
		},
		DBBaseModel: models.DBBaseModel{
			ID: 1,
		},
	}, nil)

	// Mock email service
	mockRenderProvider.On("Render", mock.Anything, mock.Anything).Return("test-rendered-email", nil)
	mockEmailProvider.On("SendEmail", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	emailservice.RegisterUserEmailServiceInstance.SetUp(
		mockRenderProvider,
		mockEmailProvider,
	)

	uc := NewCreateUserSendEmailUseCase(
		testHashProvider,
		testTokenRepository,
	)

	result := uc.Execute(ctx, locales.EN_US, testUser)

	assert.NotNil(result)
	assert.False(result.HasError())
	assert.True(result.IsSuccess())
	assert.NotNil(result.Data)
	assert.Equal(testUser.ID, result.Data.ID)
	assert.Equal(testUser.Email, result.Data.Email)
	assert.Equal(testUser.Name, result.Data.Name)
	assert.Equal(status.Created, result.StatusCode)
}

func TestCreateUserSendEmailUseCase_Execute_OneTimeTokenError(t *testing.T) {
	assert := assert.New(t)

	ctx := &app_context.AppContext{Context: context.Background()}

	testHashProvider := new(providersmocks.MockHashProvider)
	testTokenRepository := new(repositoriesmocks.MockOneTimeTokenRepository)

	userStatus := models.UserStatusPending
	testUser := models.User{
		UserBase: models.UserBase{
			Name:   "Test User",
			Email:  "test@example.com",
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

	// Mock OneTimeToken returning error
	appErr := applicationerror.NewApplicationError(
		status.InternalError,
		messages.MessageKeysInstance.SOMETHING_WENT_WRONG,
		"Failed to generate token",
	)
	testHashProvider.On("OneTimeToken").Return("", []byte(nil), appErr)

	uc := NewCreateUserSendEmailUseCase(
		testHashProvider,
		testTokenRepository,
	)

	result := uc.Execute(ctx, locales.EN_US, testUser)

	assert.NotNil(result)
	assert.True(result.HasError())
	assert.False(result.IsSuccess())
	assert.Equal(status.InternalError, result.StatusCode)
}

func TestCreateUserSendEmailUseCase_Execute_TokenRepositoryCreateError(t *testing.T) {
	assert := assert.New(t)

	ctx := &app_context.AppContext{Context: context.Background()}

	testHashProvider := new(providersmocks.MockHashProvider)
	testTokenRepository := new(repositoriesmocks.MockOneTimeTokenRepository)

	userStatus := models.UserStatusPending
	testUser := models.User{
		UserBase: models.UserBase{
			Name:   "Test User",
			Email:  "test@example.com",
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

	// Mock OneTimeToken success
	token := "test-token-123"
	tokenHash := []byte("hash")
	testHashProvider.On("OneTimeToken").Return(token, tokenHash, nil)

	// Mock Create token repository returning error
	appErr := applicationerror.NewApplicationError(
		status.InternalError,
		messages.MessageKeysInstance.SOMETHING_WENT_WRONG,
		"Failed to create token in repository",
	)
	testTokenRepository.On("Create", mock.AnythingOfType("dtos.OneTimeTokenCreate")).Return((*models.OneTimeToken)(nil), appErr)

	uc := NewCreateUserSendEmailUseCase(
		testHashProvider,
		testTokenRepository,
	)

	result := uc.Execute(ctx, locales.EN_US, testUser)

	assert.NotNil(result)
	assert.True(result.HasError())
	assert.False(result.IsSuccess())
	assert.Equal(status.InternalError, result.StatusCode)
}

func TestCreateUserSendEmailUseCase_Execute_EmailSendError(t *testing.T) {
	assert := assert.New(t)

	ctx := &app_context.AppContext{Context: context.Background()}

	testHashProvider := new(providersmocks.MockHashProvider)
	testTokenRepository := new(repositoriesmocks.MockOneTimeTokenRepository)

	mockRenderProvider := new(providersmocks.MockRenderProvider[emailmodels.NewUserEmailData])
	mockEmailProvider := new(providersmocks.MockEmailProvider)

	userStatus := models.UserStatusPending
	testUser := models.User{
		UserBase: models.UserBase{
			Name:   "Test User",
			Email:  "test@example.com",
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

	// Mock OneTimeToken success
	token := "test-token-123"
	tokenHash := []byte("hash")
	testHashProvider.On("OneTimeToken").Return(token, tokenHash, nil)

	// Mock Create token repository success
	testTokenRepository.On("Create", mock.AnythingOfType("dtos.OneTimeTokenCreate")).Return(&models.OneTimeToken{
		OneTimeTokenBase: models.OneTimeTokenBase{
			UserID:  1,
			Purpose: models.OneTimeTokenPurposeEmailVerify,
			Hash:    tokenHash,
			IsUsed:  false,
			Expires: time.Now().Add(24 * time.Hour),
		},
		DBBaseModel: models.DBBaseModel{
			ID: 1,
		},
	}, nil)

	// Mock email service returning error
	appErr := applicationerror.NewApplicationError(
		status.InternalError,
		messages.MessageKeysInstance.SOMETHING_WENT_WRONG,
		"Failed to send email",
	)
	mockRenderProvider.On("Render", mock.Anything, mock.Anything).Return("test-rendered-email", nil)
	mockEmailProvider.On("SendEmail", mock.Anything, mock.Anything, mock.Anything).Return(appErr)

	emailservice.RegisterUserEmailServiceInstance.SetUp(
		mockRenderProvider,
		mockEmailProvider,
	)

	uc := NewCreateUserSendEmailUseCase(
		testHashProvider,
		testTokenRepository,
	)

	result := uc.Execute(ctx, locales.EN_US, testUser)

	assert.NotNil(result)
	assert.True(result.HasError())
	assert.False(result.IsSuccess())
	assert.Equal(status.InternalError, result.StatusCode)
}

func TestCreateUserSendEmailUseCase_Execute_EmailRenderError(t *testing.T) {
	assert := assert.New(t)

	ctx := &app_context.AppContext{Context: context.Background()}

	testHashProvider := new(providersmocks.MockHashProvider)
	testTokenRepository := new(repositoriesmocks.MockOneTimeTokenRepository)

	mockRenderProvider := new(providersmocks.MockRenderProvider[emailmodels.NewUserEmailData])
	mockEmailProvider := new(providersmocks.MockEmailProvider)

	userStatus := models.UserStatusPending
	testUser := models.User{
		UserBase: models.UserBase{
			Name:   "Test User",
			Email:  "test@example.com",
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

	// Mock OneTimeToken success
	token := "test-token-123"
	tokenHash := []byte("hash")
	testHashProvider.On("OneTimeToken").Return(token, tokenHash, nil)

	// Mock Create token repository success
	testTokenRepository.On("Create", mock.AnythingOfType("dtos.OneTimeTokenCreate")).Return(&models.OneTimeToken{
		OneTimeTokenBase: models.OneTimeTokenBase{
			UserID:  1,
			Purpose: models.OneTimeTokenPurposeEmailVerify,
			Hash:    tokenHash,
			IsUsed:  false,
			Expires: time.Now().Add(24 * time.Hour),
		},
		DBBaseModel: models.DBBaseModel{
			ID: 1,
		},
	}, nil)

	// Mock email render returning error
	appErr := applicationerror.NewApplicationError(
		status.InternalError,
		messages.MessageKeysInstance.SOMETHING_WENT_WRONG,
		"Failed to render email template",
	)
	mockRenderProvider.On("Render", mock.Anything, mock.Anything).Return("", appErr)

	emailservice.RegisterUserEmailServiceInstance.SetUp(
		mockRenderProvider,
		mockEmailProvider,
	)

	uc := NewCreateUserSendEmailUseCase(
		testHashProvider,
		testTokenRepository,
	)

	result := uc.Execute(ctx, locales.EN_US, testUser)

	assert.NotNil(result)
	assert.True(result.HasError())
	assert.False(result.IsSuccess())
	assert.Equal(status.InternalError, result.StatusCode)
}

func TestCreateUserSendEmailUseCase_SetLocale(t *testing.T) {
	assert := assert.New(t)

	testHashProvider := new(providersmocks.MockHashProvider)
	testTokenRepository := new(repositoriesmocks.MockOneTimeTokenRepository)

	uc := NewCreateUserSendEmailUseCase(
		testHashProvider,
		testTokenRepository,
	)

	// Test setting locale
	uc.SetLocale(locales.ES_ES)
	assert.Equal(locales.ES_ES, uc.Locale)

	// Test setting empty locale (should not change)
	uc.SetLocale("")
	assert.Equal(locales.EN_US, uc.Locale)

	// Test setting another locale
	uc.SetLocale(locales.EN_US)
	assert.Equal(locales.EN_US, uc.Locale)
}
