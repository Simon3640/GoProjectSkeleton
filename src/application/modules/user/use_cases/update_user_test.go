package userusecases

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
	"github.com/simon3640/goprojectskeleton/src/application/shared/status"
	"github.com/simon3640/goprojectskeleton/src/domain/models"

	"github.com/stretchr/testify/assert"
)

func TestUpdateUserUseCase(t *testing.T) {
	assert := assert.New(t)

	ctx := context.Background()

	actor := dtomocks.UserWithRole
	ctxWithUser := context.WithValue(ctx, app_context.UserKey, actor)

	testLogger := new(providersmocks.MockLoggerProvider)
	testUserRepository := new(repositoriesmocks.MockUserRepository)
	name := "Update"
	testUser := dtos.UserUpdate{
		UserUpdateBase: models.UserUpdateBase{Name: &name},
		ID:             actor.ID,
	}
	userStatus := models.UserStatusActive
	testUserRepository.On("Update", testUser.ID, testUser).Return(&models.User{
		UserBase: models.UserBase{Name: "Update",
			Email:  "test@testing.com",
			Phone:  "1234567890",
			Status: &userStatus},
		DBBaseModel: models.DBBaseModel{
			ID:        actor.ID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: time.Time{},
		},
	}, nil)

	uc := NewUpdateUserUseCase(testLogger, testUserRepository)

	result := uc.Execute(ctxWithUser, locales.EN_US, testUser)

	assert.NotNil(result)
	assert.Equal(result.Data.ID == 1, true)
	assert.Equal(result.Data.Name == "Update", true)

}

func TestUpdateUserUseCase_DifferentUser(t *testing.T) {
	assert := assert.New(t)

	ctx := context.Background()

	actor := dtomocks.UserWithRole
	ctxWithUser := context.WithValue(ctx, app_context.UserKey, actor)

	testLogger := new(providersmocks.MockLoggerProvider)
	testUserRepository := new(repositoriesmocks.MockUserRepository)
	name := "Update"
	testUser := dtos.UserUpdate{
		UserUpdateBase: models.UserUpdateBase{Name: &name},
		ID:             actor.ID + 1,
	}

	testUserRepository.On("Update", testUser.ID, testUser).Return(nil)

	uc := NewUpdateUserUseCase(testLogger, testUserRepository)

	result := uc.Execute(ctxWithUser, locales.EN_US, testUser)

	assert.NotNil(result)
	assert.Equal(result.StatusCode, status.Unauthorized)
}
