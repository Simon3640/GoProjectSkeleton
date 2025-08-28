package usecases_user

import (
	"context"
	dtos "gormgoskeleton/src/application/shared/DTOs"
	app_context "gormgoskeleton/src/application/shared/context"
	"gormgoskeleton/src/application/shared/locales"
	"gormgoskeleton/src/application/shared/mocks"
	dto_mocks "gormgoskeleton/src/application/shared/mocks/dtos"
	"gormgoskeleton/src/application/shared/status"
	"gormgoskeleton/src/domain/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUpdateUserUseCase(t *testing.T) {
	assert := assert.New(t)

	ctx := context.Background()

	actor := dto_mocks.UserWithRole
	ctxWithUser := context.WithValue(ctx, app_context.UserKey, actor)

	testLogger := new(mocks.MockLoggerProvider)
	testUserRepository := new(mocks.MockUserRepository)
	name := "Update"
	testUser := dtos.UserUpdate{
		UserUpdateBase: dtos.UserUpdateBase{Name: &name},
		ID:             actor.ID,
	}

	testUserRepository.On("Update", testUser.ID, testUser).Return(&models.User{
		UserBase: models.UserBase{Name: "Update",
			Email:  "test@testing.com",
			Phone:  "1234567890",
			Status: "active"},
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

	actor := dto_mocks.UserWithRole
	ctxWithUser := context.WithValue(ctx, app_context.UserKey, actor)

	testLogger := new(mocks.MockLoggerProvider)
	testUserRepository := new(mocks.MockUserRepository)
	name := "Update"
	testUser := dtos.UserUpdate{
		UserUpdateBase: dtos.UserUpdateBase{Name: &name},
		ID:             actor.ID + 1,
	}

	testUserRepository.On("Update", testUser.ID, testUser).Return(nil)

	uc := NewUpdateUserUseCase(testLogger, testUserRepository)

	result := uc.Execute(ctxWithUser, locales.EN_US, testUser)

	assert.NotNil(result)
	assert.Equal(result.StatusCode, status.Unauthorized)
}
