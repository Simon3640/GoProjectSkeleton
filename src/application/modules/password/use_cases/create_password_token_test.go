package passwordusecases

import (
	"context"
	"encoding/hex"
	"testing"
	"time"

	dtos "github.com/simon3640/goprojectskeleton/src/application/modules/password/dtos"
	passwordmocks "github.com/simon3640/goprojectskeleton/src/application/modules/password/mocks"
	shareddtos "github.com/simon3640/goprojectskeleton/src/application/shared/DTOs"
	app_context "github.com/simon3640/goprojectskeleton/src/application/shared/context"
	applicationerrors "github.com/simon3640/goprojectskeleton/src/application/shared/errors"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales/messages"
	providersmocks "github.com/simon3640/goprojectskeleton/src/application/shared/mocks/providers"
	repositoriesmocks "github.com/simon3640/goprojectskeleton/src/application/shared/mocks/repositories"
	"github.com/simon3640/goprojectskeleton/src/application/shared/status"
	"github.com/simon3640/goprojectskeleton/src/domain/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreatePasswordTokenUseCase_Execute_Success(t *testing.T) {
	assert := assert.New(t)

	ctx := &app_context.AppContext{Context: context.Background()}
	testPasswordRepository := new(passwordmocks.MockPasswordRepository)
	testHashProvider := new(providersmocks.MockHashProvider)
	testOneTimeTokenRepository := new(repositoriesmocks.MockOneTimeTokenRepository)

	token := "test-token-123"
	tokenHash := []byte(hex.EncodeToString([]byte("hashed_token")))
	userID := uint(1)
	noHashedPassword := "NewPassword123!"
	hashedPassword := "HashedNewPassword123!"

	// Create a valid token
	validToken := &models.OneTimeToken{
		OneTimeTokenBase: models.OneTimeTokenBase{
			UserID:  userID,
			Purpose: models.OneTimeTokenPurposePasswordReset,
			Hash:    tokenHash,
			IsUsed:  false,
			Expires: time.Now().Add(1 * time.Hour),
		},
		DBBaseModel: models.DBBaseModel{
			ID:        1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	input := dtos.PasswordTokenCreate{
		Token:            token,
		NoHashedPassword: noHashedPassword,
	}

	// Configure mocks
	testHashProvider.On("HashOneTimeToken", token).Return(tokenHash)
	testOneTimeTokenRepository.On("GetByTokenHash", tokenHash).Return(validToken, nil)
	testHashProvider.On("HashPassword", noHashedPassword).Return(hashedPassword, nil)

	testPasswordRepository.On("Create", mock.MatchedBy(func(pc dtos.PasswordCreate) bool {
		return pc.UserID == userID &&
			pc.Hash == hashedPassword &&
			pc.IsActive == true &&
			pc.ExpiresAt != nil
	})).Return(&models.Password{
		PasswordBase: models.PasswordBase{
			UserID:   userID,
			Hash:     hashedPassword,
			IsActive: true,
		},
		ID: 1,
	}, nil)

	tokenUpdate := shareddtos.OneTimeTokenUpdate{IsUsed: true, ID: validToken.ID}
	updatedToken := &models.OneTimeToken{
		OneTimeTokenBase: models.OneTimeTokenBase{
			UserID:  userID,
			Purpose: models.OneTimeTokenPurposePasswordReset,
			Hash:    tokenHash,
			IsUsed:  true,
			Expires: validToken.Expires,
		},
		DBBaseModel: models.DBBaseModel{
			ID:        validToken.ID,
			CreatedAt: validToken.CreatedAt,
			UpdatedAt: time.Now(),
		},
	}
	testOneTimeTokenRepository.On("Update", validToken.ID, tokenUpdate).Return(updatedToken, nil)

	uc := NewCreatePasswordTokenUseCase(
		testPasswordRepository,
		testHashProvider,
		testOneTimeTokenRepository,
	)

	result := uc.Execute(ctx, locales.EN_US, input)

	assert.NotNil(result)
	assert.True(result.IsSuccess())
	assert.Equal(true, *result.Data)
	assert.Equal(status.Success, result.StatusCode)

	testHashProvider.AssertExpectations(t)
	testOneTimeTokenRepository.AssertExpectations(t)
	testPasswordRepository.AssertExpectations(t)
}

func TestCreatePasswordTokenUseCase_Execute_ErrorGettingToken(t *testing.T) {
	assert := assert.New(t)

	ctx := &app_context.AppContext{Context: context.Background()}
	testPasswordRepository := new(passwordmocks.MockPasswordRepository)
	testHashProvider := new(providersmocks.MockHashProvider)
	testOneTimeTokenRepository := new(repositoriesmocks.MockOneTimeTokenRepository)

	token := "test-token-123"
	tokenHash := []byte(hex.EncodeToString([]byte("hashed_token")))

	input := dtos.PasswordTokenCreate{
		Token:            token,
		NoHashedPassword: "NewPassword123!",
	}

	appError := applicationerrors.NewApplicationError(
		status.NotFound,
		messages.MessageKeysInstance.RESOURCE_NOT_FOUND,
		"Token not found",
	)

	// Configure mocks
	testHashProvider.On("HashOneTimeToken", token).Return(tokenHash)
	testOneTimeTokenRepository.On("GetByTokenHash", tokenHash).Return(nil, appError)

	uc := NewCreatePasswordTokenUseCase(
		testPasswordRepository,
		testHashProvider,
		testOneTimeTokenRepository,
	)

	result := uc.Execute(ctx, locales.EN_US, input)

	assert.NotNil(result)
	assert.True(result.HasError())
	assert.Equal(status.NotFound, result.StatusCode)

	testHashProvider.AssertExpectations(t)
	testOneTimeTokenRepository.AssertExpectations(t)
}

func TestCreatePasswordTokenUseCase_Execute_TokenIsNil(t *testing.T) {
	assert := assert.New(t)

	ctx := &app_context.AppContext{Context: context.Background()}
	testPasswordRepository := new(passwordmocks.MockPasswordRepository)
	testHashProvider := new(providersmocks.MockHashProvider)
	testOneTimeTokenRepository := new(repositoriesmocks.MockOneTimeTokenRepository)

	token := "test-token-123"
	tokenHash := []byte(hex.EncodeToString([]byte("hashed_token")))

	input := dtos.PasswordTokenCreate{
		Token:            token,
		NoHashedPassword: "NewPassword123!",
	}

	// Configure mocks
	testHashProvider.On("HashOneTimeToken", token).Return(tokenHash)
	testOneTimeTokenRepository.On("GetByTokenHash", tokenHash).Return(nil, nil)

	uc := NewCreatePasswordTokenUseCase(
		testPasswordRepository,
		testHashProvider,
		testOneTimeTokenRepository,
	)

	result := uc.Execute(ctx, locales.EN_US, input)

	assert.NotNil(result)
	assert.True(result.HasError())
	assert.Equal(status.Conflict, result.StatusCode)

	testHashProvider.AssertExpectations(t)
	testOneTimeTokenRepository.AssertExpectations(t)
}

func TestCreatePasswordTokenUseCase_Execute_TokenIsUsed(t *testing.T) {
	assert := assert.New(t)

	ctx := &app_context.AppContext{Context: context.Background()}
	testPasswordRepository := new(passwordmocks.MockPasswordRepository)
	testHashProvider := new(providersmocks.MockHashProvider)
	testOneTimeTokenRepository := new(repositoriesmocks.MockOneTimeTokenRepository)

	token := "test-token-123"
	tokenHash := []byte(hex.EncodeToString([]byte("hashed_token")))

	// Create a used token
	usedToken := &models.OneTimeToken{
		OneTimeTokenBase: models.OneTimeTokenBase{
			UserID:  1,
			Purpose: models.OneTimeTokenPurposePasswordReset,
			Hash:    tokenHash,
			IsUsed:  true,
			Expires: time.Now().Add(1 * time.Hour),
		},
		DBBaseModel: models.DBBaseModel{
			ID:        1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	input := dtos.PasswordTokenCreate{
		Token:            token,
		NoHashedPassword: "NewPassword123!",
	}

	// Configure mocks
	testHashProvider.On("HashOneTimeToken", token).Return(tokenHash)
	testOneTimeTokenRepository.On("GetByTokenHash", tokenHash).Return(usedToken, nil)

	uc := NewCreatePasswordTokenUseCase(
		testPasswordRepository,
		testHashProvider,
		testOneTimeTokenRepository,
	)

	result := uc.Execute(ctx, locales.EN_US, input)

	assert.NotNil(result)
	assert.True(result.HasError())
	assert.Equal(status.Conflict, result.StatusCode)

	testHashProvider.AssertExpectations(t)
	testOneTimeTokenRepository.AssertExpectations(t)
}

func TestCreatePasswordTokenUseCase_Execute_TokenExpired(t *testing.T) {
	assert := assert.New(t)

	ctx := &app_context.AppContext{Context: context.Background()}
	testPasswordRepository := new(passwordmocks.MockPasswordRepository)
	testHashProvider := new(providersmocks.MockHashProvider)
	testOneTimeTokenRepository := new(repositoriesmocks.MockOneTimeTokenRepository)

	token := "test-token-123"
	tokenHash := []byte(hex.EncodeToString([]byte("hashed_token")))

	// Create an expired token
	expiredToken := &models.OneTimeToken{
		OneTimeTokenBase: models.OneTimeTokenBase{
			UserID:  1,
			Purpose: models.OneTimeTokenPurposePasswordReset,
			Hash:    tokenHash,
			IsUsed:  false,
			Expires: time.Now().Add(-1 * time.Hour),
		},
		DBBaseModel: models.DBBaseModel{
			ID:        1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	input := dtos.PasswordTokenCreate{
		Token:            token,
		NoHashedPassword: "NewPassword123!",
	}

	// Configure mocks
	testHashProvider.On("HashOneTimeToken", token).Return(tokenHash)
	testOneTimeTokenRepository.On("GetByTokenHash", tokenHash).Return(expiredToken, nil)

	uc := NewCreatePasswordTokenUseCase(
		testPasswordRepository,
		testHashProvider,
		testOneTimeTokenRepository,
	)

	result := uc.Execute(ctx, locales.EN_US, input)

	assert.NotNil(result)
	assert.True(result.HasError())
	assert.Equal(status.Conflict, result.StatusCode)

	testHashProvider.AssertExpectations(t)
	testOneTimeTokenRepository.AssertExpectations(t)
}

func TestCreatePasswordTokenUseCase_Execute_ErrorCreatingPassword(t *testing.T) {
	assert := assert.New(t)

	ctx := &app_context.AppContext{Context: context.Background()}
	testPasswordRepository := new(passwordmocks.MockPasswordRepository)
	testHashProvider := new(providersmocks.MockHashProvider)
	testOneTimeTokenRepository := new(repositoriesmocks.MockOneTimeTokenRepository)

	token := "test-token-123"
	tokenHash := []byte(hex.EncodeToString([]byte("hashed_token")))
	userID := uint(1)
	noHashedPassword := "NewPassword123!"

	// Create a valid token
	validToken := &models.OneTimeToken{
		OneTimeTokenBase: models.OneTimeTokenBase{
			UserID:  userID,
			Purpose: models.OneTimeTokenPurposePasswordReset,
			Hash:    tokenHash,
			IsUsed:  false,
			Expires: time.Now().Add(1 * time.Hour),
		},
		DBBaseModel: models.DBBaseModel{
			ID:        1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	input := dtos.PasswordTokenCreate{
		Token:            token,
		NoHashedPassword: noHashedPassword,
	}

	appError := applicationerrors.NewApplicationError(
		status.InternalError,
		messages.MessageKeysInstance.SOMETHING_WENT_WRONG,
		"Error creating password",
	)

	// Configure mocks
	testHashProvider.On("HashOneTimeToken", token).Return(tokenHash)
	testOneTimeTokenRepository.On("GetByTokenHash", tokenHash).Return(validToken, nil)
	testHashProvider.On("HashPassword", noHashedPassword).Return("", appError)

	uc := NewCreatePasswordTokenUseCase(
		testPasswordRepository,
		testHashProvider,
		testOneTimeTokenRepository,
	)

	result := uc.Execute(ctx, locales.EN_US, input)

	assert.NotNil(result)
	assert.True(result.HasError())
	assert.Equal(status.InternalError, result.StatusCode)

	testHashProvider.AssertExpectations(t)
	testOneTimeTokenRepository.AssertExpectations(t)
}

func TestCreatePasswordTokenUseCase_Execute_ErrorMarkingTokenAsUsed(t *testing.T) {
	assert := assert.New(t)

	ctx := &app_context.AppContext{Context: context.Background()}
	testPasswordRepository := new(passwordmocks.MockPasswordRepository)
	testHashProvider := new(providersmocks.MockHashProvider)
	testOneTimeTokenRepository := new(repositoriesmocks.MockOneTimeTokenRepository)

	token := "test-token-123"
	tokenHash := []byte(hex.EncodeToString([]byte("hashed_token")))
	userID := uint(1)
	noHashedPassword := "NewPassword123!"
	hashedPassword := "HashedNewPassword123!"

	// Create a valid token
	validToken := &models.OneTimeToken{
		OneTimeTokenBase: models.OneTimeTokenBase{
			UserID:  userID,
			Purpose: models.OneTimeTokenPurposePasswordReset,
			Hash:    tokenHash,
			IsUsed:  false,
			Expires: time.Now().Add(1 * time.Hour),
		},
		DBBaseModel: models.DBBaseModel{
			ID:        1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	input := dtos.PasswordTokenCreate{
		Token:            token,
		NoHashedPassword: noHashedPassword,
	}

	appError := applicationerrors.NewApplicationError(
		status.InternalError,
		messages.MessageKeysInstance.SOMETHING_WENT_WRONG,
		"Error updating token",
	)

	// Configure mocks
	testHashProvider.On("HashOneTimeToken", token).Return(tokenHash)
	testOneTimeTokenRepository.On("GetByTokenHash", tokenHash).Return(validToken, nil)
	testHashProvider.On("HashPassword", noHashedPassword).Return(hashedPassword, nil)

	testPasswordRepository.On("Create", mock.MatchedBy(func(pc dtos.PasswordCreate) bool {
		return pc.UserID == userID &&
			pc.Hash == hashedPassword &&
			pc.IsActive == true &&
			pc.ExpiresAt != nil
	})).Return(&models.Password{
		PasswordBase: models.PasswordBase{
			UserID:   userID,
			Hash:     hashedPassword,
			IsActive: true,
		},
		ID: 1,
	}, nil)

	tokenUpdate := shareddtos.OneTimeTokenUpdate{IsUsed: true, ID: validToken.ID}
	testOneTimeTokenRepository.On("Update", validToken.ID, tokenUpdate).Return(nil, appError)

	uc := NewCreatePasswordTokenUseCase(
		testPasswordRepository,
		testHashProvider,
		testOneTimeTokenRepository,
	)

	result := uc.Execute(ctx, locales.EN_US, input)

	assert.NotNil(result)
	assert.True(result.HasError())
	assert.Equal(status.InternalError, result.StatusCode)

	testHashProvider.AssertExpectations(t)
	testOneTimeTokenRepository.AssertExpectations(t)
	testPasswordRepository.AssertExpectations(t)
}

func TestCreatePasswordTokenUseCase_Execute_TokenPurposeIsNotPasswordReset(t *testing.T) {
	assert := assert.New(t)

	ctx := &app_context.AppContext{Context: context.Background()}
	testPasswordRepository := new(passwordmocks.MockPasswordRepository)
	testHashProvider := new(providersmocks.MockHashProvider)
	testOneTimeTokenRepository := new(repositoriesmocks.MockOneTimeTokenRepository)

	token := "test-token-123"
	tokenHash := []byte(hex.EncodeToString([]byte("hashed_token")))
	userID := uint(1)
	noHashedPassword := "NewPassword123!"

	// Create a valid token
	validToken := &models.OneTimeToken{
		OneTimeTokenBase: models.OneTimeTokenBase{
			UserID:  userID,
			Purpose: models.OneTimeTokenPurposeEmailVerify,
			Hash:    tokenHash,
			IsUsed:  false,
			Expires: time.Now().Add(1 * time.Hour),
		},
		DBBaseModel: models.DBBaseModel{
			ID:        1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	input := dtos.PasswordTokenCreate{
		Token:            token,
		NoHashedPassword: noHashedPassword,
	}

	// Configure mocks
	testHashProvider.On("HashOneTimeToken", token).Return(tokenHash)
	testOneTimeTokenRepository.On("GetByTokenHash", tokenHash).Return(validToken, nil)

	uc := NewCreatePasswordTokenUseCase(
		testPasswordRepository,
		testHashProvider,
		testOneTimeTokenRepository,
	)

	result := uc.Execute(ctx, locales.EN_US, input)

	assert.NotNil(result)
	assert.True(result.HasError())
	assert.Equal(status.Conflict, result.StatusCode)

	testHashProvider.AssertExpectations(t)
	testOneTimeTokenRepository.AssertExpectations(t)
}
