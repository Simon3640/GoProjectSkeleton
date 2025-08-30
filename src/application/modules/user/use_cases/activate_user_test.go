package usecases_user

import (
	"context"
	"encoding/hex"
	"testing"
	"time"

	dtos "gormgoskeleton/src/application/shared/DTOs"
	"gormgoskeleton/src/application/shared/locales"
	"gormgoskeleton/src/application/shared/mocks"
	"gormgoskeleton/src/domain/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestActivateUserUseCase(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()

	testLogger := new(mocks.MockLoggerProvider)
	testUserRepository := new(mocks.MockUserRepository)
	testOneTimeTokenRepository := new(mocks.MockOneTimeTokenRepository)
	testHashProvider := new(mocks.MockHashProvider)

	// Entity Mocks

	hash := "hashed_token"
	tokenHash := []byte(hex.EncodeToString([]byte(hash)))
	oneTimeToken := models.OneTimeToken{
		OneTimeTokenBase: models.OneTimeTokenBase{
			UserID:  1,
			Purpose: models.OneTimeTokenPurposeEmailVerify,
			Hash:    []byte(hash),
			IsUsed:  false,
			Expires: time.Now().Add(1 * time.Hour),
		},
	}

	testHashProvider.On("HashOneTimeToken", "valid_token").Return(tokenHash)
	testOneTimeTokenRepository.On("GetByTokenHash", tokenHash).Return(&oneTimeToken, nil)
	testUserRepository.On(
		"Update",
		oneTimeToken.UserID,
		mock.AnythingOfType("dtos.UserUpdate"),
	).Return(&models.User{
		UserBase: models.UserBase{
			Name:   "Test User",
			Email:  "test@mail.com",
			Status: "active",
			RoleID: 2,
		},
		DBBaseModel: models.DBBaseModel{
			ID:        1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: time.Time{},
		},
	}, nil)

	// Create the use case
	useCase := NewActivateUserUseCase(
		testLogger,
		testUserRepository,
		testOneTimeTokenRepository,
		testHashProvider,
	)

	userActivate := dtos.UserActivate{
		Token: "valid_token",
	}

	result := useCase.Execute(ctx, locales.EN_US, userActivate)

	assert.NotNil(result)
	assert.True(result.IsSuccess())
	assert.NotNil(result.Data)
	assert.Equal(true, *result.Data)

}
