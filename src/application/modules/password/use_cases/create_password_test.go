package usecases_password

import (
	"context"
	"testing"
	"time"

	dtos "github.com/simon3640/goprojectskeleton/src/application/shared/DTOs"
	app_context "github.com/simon3640/goprojectskeleton/src/application/shared/context"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales"
	dtomocks "github.com/simon3640/goprojectskeleton/src/application/shared/mocks/dtos"
	providersmocks "github.com/simon3640/goprojectskeleton/src/application/shared/mocks/providers"
	repositoriesmocks "github.com/simon3640/goprojectskeleton/src/application/shared/mocks/repositories"
	"github.com/simon3640/goprojectskeleton/src/domain/models"

	"github.com/stretchr/testify/assert"
)

func TestCreatePasswordUseCase(t *testing.T) {
	assert := assert.New(t)

	ctx := context.Background()

	actor := dtomocks.UserWithRole

	testLogger := new(providersmocks.MockLoggerProvider)
	testPasswordRepository := new(repositoriesmocks.MockPasswordRepository)
	testHashProvider := new(providersmocks.MockHashProvider)
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
