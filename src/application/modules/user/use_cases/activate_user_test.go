package userusecases

import (
	"context"
	"encoding/hex"
	"testing"
	"time"

	userdtos "github.com/simon3640/goprojectskeleton/src/application/modules/user/dtos"
	usermocks "github.com/simon3640/goprojectskeleton/src/application/modules/user/mocks"
	app_context "github.com/simon3640/goprojectskeleton/src/application/shared/context"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales"
	providersmocks "github.com/simon3640/goprojectskeleton/src/application/shared/mocks/providers"
	repositoriesmocks "github.com/simon3640/goprojectskeleton/src/application/shared/mocks/repositories"
	sharedmodels "github.com/simon3640/goprojectskeleton/src/domain/shared/models"
	usermodels "github.com/simon3640/goprojectskeleton/src/domain/user/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestActivateUserUseCase(t *testing.T) {
	assert := assert.New(t)
	ctx := &app_context.AppContext{Context: context.Background()}

	testUserRepository := new(usermocks.MockUserRepository)
	testOneTimeTokenRepository := new(repositoriesmocks.MockOneTimeTokenRepository)
	testHashProvider := new(providersmocks.MockHashProvider)

	// Entity Mocks

	hash := "hashed_token"
	tokenHash := []byte(hex.EncodeToString([]byte(hash)))
	oneTimeToken := sharedmodels.OneTimeToken{
		OneTimeTokenBase: sharedmodels.OneTimeTokenBase{
			UserID:  1,
			Purpose: sharedmodels.OneTimeTokenPurposeEmailVerify,
			Hash:    []byte(hash),
			IsUsed:  false,
			Expires: time.Now().Add(1 * time.Hour),
		},
	}

	testHashProvider.On("HashOneTimeToken", "valid_token").Return(tokenHash)
	testOneTimeTokenRepository.On("GetByTokenHash", tokenHash).Return(&oneTimeToken, nil)
	userStatusActive := usermodels.UserStatusActive
	testUserRepository.On(
		"Update",
		oneTimeToken.UserID,
		mock.AnythingOfType("userdtos.UserUpdate"),
	).Return(&usermodels.User{
		UserBase: usermodels.UserBase{
			Name:   "Test User",
			Email:  "test@mail.com",
			Status: &userStatusActive,
			RoleID: 2,
		},
		DBBaseModel: sharedmodels.DBBaseModel{
			ID:        1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: time.Time{},
		},
	}, nil)

	// Create the use case
	useCase := NewActivateUserUseCase(
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
