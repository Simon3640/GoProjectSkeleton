package userusecases

import (
	"context"
	"encoding/hex"
	"testing"
	"time"

	userdtos "github.com/simon3640/goprojectskeleton/src/application/modules/user/dtos"
	usermocks "github.com/simon3640/goprojectskeleton/src/application/modules/user/mocks"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales"
	providersmocks "github.com/simon3640/goprojectskeleton/src/application/shared/mocks/providers"
	repositoriesmocks "github.com/simon3640/goprojectskeleton/src/application/shared/mocks/repositories"
	"github.com/simon3640/goprojectskeleton/src/domain/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestActivateUserUseCase(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()

	testLogger := new(providersmocks.MockLoggerProvider)
	testUserRepository := new(usermocks.MockUserRepository)
	testOneTimeTokenRepository := new(repositoriesmocks.MockOneTimeTokenRepository)
	testHashProvider := new(providersmocks.MockHashProvider)

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
	userStatusActive := models.UserStatusActive
	testUserRepository.On(
		"Update",
		oneTimeToken.UserID,
		mock.AnythingOfType("userdtos.UserUpdate"),
	).Return(&models.User{
		UserBase: models.UserBase{
			Name:   "Test User",
			Email:  "test@mail.com",
			Status: &userStatusActive,
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

	userActivate := userdtos.UserActivate{
		Token: "valid_token",
	}

	result := useCase.Execute(ctx, locales.EN_US, userActivate)

	assert.NotNil(result)
	assert.True(result.IsSuccess())
	assert.NotNil(result.Data)
	assert.Equal(true, *result.Data)

}
