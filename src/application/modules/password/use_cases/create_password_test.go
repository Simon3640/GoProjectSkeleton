package usecases_password

import (
	"context"
	"testing"
	"time"

	"gormgoskeleton/src/application/shared/locales"
	"gormgoskeleton/src/application/shared/mocks"
	"gormgoskeleton/src/domain/models"

	"github.com/stretchr/testify/assert"
)

func TestCreatePasswordUseCase(t *testing.T) {
	assert := assert.New(t)

	ctx := context.Background()

	testLogger := new(mocks.MockLoggerProvider)
	testPasswordRepository := new(mocks.MockPasswordRepository)
	testHashProvider := new(mocks.MockHashProvider)
	testPassword := models.PasswordCreateNoHash{
		UserID:           1,
		NoHashedPassword: "TestPassword123!",
		ExpiresAt:        &time.Time{},
		IsActive:         true,
	}

	testPasswordCreate := models.NewPasswordCreate(
		testPassword.UserID,
		"HashedPassword123!",
		testPassword.ExpiresAt,
		testPassword.IsActive,
	)

	testPasswordRepository.On("Create", testPasswordCreate).Return(&models.Password{
		PasswordBase: models.PasswordBase{
			UserID:    1,
			Hash:      "HashedPassword123!",
			ExpiresAt: &time.Time{},
			IsActive:  true,
		},
		ID: 1,
	}, nil)

	testHashProvider.On("HashPassword", testPassword.NoHashedPassword).Return("HashedPassword123!", nil)

	uc := NewCreatePasswordUseCase(testLogger, testPasswordRepository, testHashProvider)

	result := uc.Execute(ctx, locales.EN_US, testPassword)

	assert.NotNil(result)
	assert.True(result.IsSuccess())
	assert.Equal(*result.Data, true)
}
