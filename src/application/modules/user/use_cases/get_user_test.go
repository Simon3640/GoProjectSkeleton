package usecases_user

import (
	"context"
	"testing"
	"time"

	app_context "gormgoskeleton/src/application/shared/context"
	"gormgoskeleton/src/application/shared/locales"
	"gormgoskeleton/src/application/shared/mocks"
	dto_mocks "gormgoskeleton/src/application/shared/mocks/dtos"
	"gormgoskeleton/src/application/shared/status"
	"gormgoskeleton/src/domain/models"

	"github.com/stretchr/testify/assert"
)

func TestGetUserUseCase(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()

	actor := dto_mocks.UserWithRole
	ctxWithUser := context.WithValue(ctx, app_context.UserKey, actor)

	testLogger := new(mocks.MockLoggerProvider)
	testUserRepository := new(mocks.MockUserRepository)
	var testId uint = actor.ID

	testUserRepository.On("GetByID", testId).Return(&models.User{
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
	}, nil)

	uc := NewGetUserUseCase(testLogger, testUserRepository)

	result := uc.Execute(ctxWithUser, locales.EN_US, testId)

	assert.NotNil(result)
	assert.Equal(result.Data.ID == 1, true)
	assert.Equal(result.Data.Name == "Test", true)
	assert.Equal(result.Details == "", true)
}

func TestGetUserUseCase_DifferentUser(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()

	actor := dto_mocks.UserWithRole
	ctxWithUser := context.WithValue(ctx, app_context.UserKey, actor)

	testLogger := new(mocks.MockLoggerProvider)
	testUserRepository := new(mocks.MockUserRepository)
	var testId uint = actor.ID + 1 // Different user ID

	testUserRepository.On("GetByID", testId).Return(&models.User{
		UserBase: models.UserBase{Name: "Test",
			Email:  "test@testing.com",
			Phone:  "1234567890",
			Status: "active"},
		DBBaseModel: models.DBBaseModel{
			ID:        2,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: time.Time{},
		},
	}, nil)

	uc := NewGetUserUseCase(testLogger, testUserRepository)

	result := uc.Execute(ctxWithUser, locales.EN_US, testId)

	assert.NotNil(result)
	assert.Equal(result.HasError(), true)
	assert.Equal(result.StatusCode, status.Unauthorized)

}
