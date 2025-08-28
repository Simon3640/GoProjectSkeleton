package usecases_password

import (
	"context"
	"testing"
	"time"

	dtos "gormgoskeleton/src/application/shared/DTOs"
	app_context "gormgoskeleton/src/application/shared/context"
	"gormgoskeleton/src/application/shared/locales"
	"gormgoskeleton/src/application/shared/mocks"
	dto_mocks "gormgoskeleton/src/application/shared/mocks/dtos"
	"gormgoskeleton/src/domain/models"

	"github.com/stretchr/testify/assert"
)

func TestCreatePasswordUseCase(t *testing.T) {
	assert := assert.New(t)

	ctx := context.Background()

	actor := dto_mocks.UserWithRole

	testLogger := new(mocks.MockLoggerProvider)
	testPasswordRepository := new(mocks.MockPasswordRepository)
	testHashProvider := new(mocks.MockHashProvider)
	testPassword := dtos.PasswordCreateNoHash{
		UserID:           actor.ID,
		NoHashedPassword: "TestPassword123!",
		ExpiresAt:        &time.Time{},
		IsActive:         true,
	}

	contextWithUser := context.WithValue(ctx, app_context.UserKey, actor)

	testPasswordCreate := dtos.NewPasswordCreate(
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

	uc := NewCreatePasswordUseCase(testLogger, testPasswordRepository, testHashProvider, false)

	result := uc.Execute(contextWithUser, locales.EN_US, testPassword)

	assert.NotNil(result)
	assert.True(result.IsSuccess())
	assert.Equal(*result.Data, true)
}
