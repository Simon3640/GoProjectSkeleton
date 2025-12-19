package passwordusecases

import (
	"testing"
	"time"

	dtos "github.com/simon3640/goprojectskeleton/src/application/modules/password/dtos"
	passwordmocks "github.com/simon3640/goprojectskeleton/src/application/modules/password/mocks"
	app_context "github.com/simon3640/goprojectskeleton/src/application/shared/context"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales"
	dtomocks "github.com/simon3640/goprojectskeleton/src/application/shared/mocks/dtos"
	providersmocks "github.com/simon3640/goprojectskeleton/src/application/shared/mocks/providers"
	"github.com/simon3640/goprojectskeleton/src/domain/models"

	"github.com/stretchr/testify/assert"
)

func TestCreatePasswordUseCase(t *testing.T) {
	assert := assert.New(t)

	actor := dtomocks.UserWithRole

	testPasswordRepository := new(passwordmocks.MockPasswordRepository)
	testHashProvider := new(providersmocks.MockHashProvider)
	testPassword := dtos.PasswordCreateNoHash{
		UserID:           actor.ID,
		NoHashedPassword: "TestPassword123!",
		ExpiresAt:        &time.Time{},
		IsActive:         true,
	}

	ctxWithUser := app_context.NewContextWithUser(&actor)

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

	uc := NewCreatePasswordUseCase(testPasswordRepository, testHashProvider)

	result := uc.Execute(ctxWithUser, locales.EN_US, testPassword)

	assert.NotNil(result)
	assert.True(result.IsSuccess())
	assert.Equal(*result.Data, true)
}
