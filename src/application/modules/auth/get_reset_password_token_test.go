package auth

import (
	"context"
	"encoding/hex"
	"testing"
	"time"

	"gormgoskeleton/src/application/shared/locales"
	"gormgoskeleton/src/application/shared/mocks"
	"gormgoskeleton/src/domain/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetResetPasswordTokenUseCase(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()

	testLogger := new(mocks.MockLoggerProvider)
	testMockHashProvider := new(mocks.MockHashProvider)
	testOneTimeTokenRepo := new(mocks.MockOneTimeTokenRepository)
	testUserRepo := new(mocks.MockUserRepository)

	uc := NewGetResetPasswordTokenUseCase(testLogger,
		testOneTimeTokenRepo,
		testUserRepo,
		testMockHashProvider,
	)

	// mock models
	user := models.User{
		UserBase: models.UserBase{Name: "Test",
			Email:  "test@testing.com",
			Phone:  "1234567890",
			Status: "active"},
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
	assert.Equal(token, *result.Data)

}
