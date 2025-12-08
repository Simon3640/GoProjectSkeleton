package authusecases

import (
	"context"
	"encoding/hex"
	"testing"
	"time"

	"github.com/simon3640/goprojectskeleton/src/application/shared/locales"
	providersmocks "github.com/simon3640/goprojectskeleton/src/application/shared/mocks/providers"
	repositoriesmocks "github.com/simon3640/goprojectskeleton/src/application/shared/mocks/repositories"
	"github.com/simon3640/goprojectskeleton/src/domain/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetResetPasswordTokenUseCase(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()

	testLogger := new(providersmocks.MockLoggerProvider)
	testMockHashProvider := new(providersmocks.MockHashProvider)
	testOneTimeTokenRepo := new(repositoriesmocks.MockOneTimeTokenRepository)
	testUserRepo := new(repositoriesmocks.MockUserRepository)

	uc := NewGetResetPasswordTokenUseCase(testLogger,
		testOneTimeTokenRepo,
		testUserRepo,
		testMockHashProvider,
	)

	// mock models
	userStatus := models.UserStatusActive
	user := models.User{
		UserBase: models.UserBase{Name: "Test",
			Email:  "test@testing.com",
			Phone:  "1234567890",
			Status: &userStatus},
		DBBaseModel: models.DBBaseModel{
			ID:        1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: time.Time{},
		},
	}

	token := "validResetPasswordToken.123"
	tokenHash := []byte(hex.EncodeToString([]byte(token)))

	oneTimeToken := models.OneTimeToken{
		OneTimeTokenBase: models.OneTimeTokenBase{
			UserID:  user.ID,
			Purpose: models.OneTimeTokenPurposePasswordReset,
			Hash:    tokenHash,
			IsUsed:  false,
			Expires: time.Now().Add(1 * time.Hour),
		},
		DBBaseModel: models.DBBaseModel{
			ID:        1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: time.Time{},
		},
	}

	// Mocking
	testUserRepo.On("GetByEmailOrPhone", user.Email).Return(&user, nil)
	testMockHashProvider.On("OneTimeToken").Return(token, tokenHash, nil)
	testOneTimeTokenRepo.On("Create", mock.AnythingOfType("dtos.OneTimeTokenCreate")).Return(&oneTimeToken, nil)

	result := uc.Execute(ctx, locales.EN_US, user.Email)

	assert.NotNil(result)
	assert.True(result.IsSuccess())
	assert.Equal(token, result.Data.Token)

}
