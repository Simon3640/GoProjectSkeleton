package userusecases

import (
	"testing"
	"time"

	usermocks "github.com/simon3640/goprojectskeleton/src/application/modules/user/mocks"
	app_context "github.com/simon3640/goprojectskeleton/src/application/shared/context"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales"
	dtomocks "github.com/simon3640/goprojectskeleton/src/application/shared/mocks/dtos"
	"github.com/simon3640/goprojectskeleton/src/application/shared/status"
	sharedmodels "github.com/simon3640/goprojectskeleton/src/domain/shared/models"
	usermodels "github.com/simon3640/goprojectskeleton/src/domain/user/models"

	"github.com/stretchr/testify/assert"
)

func TestGetUserUseCase(t *testing.T) {
	assert := assert.New(t)

	actor := dtomocks.UserWithRole
	ctxWithUser := app_context.NewContextWithUser(&actor)

	testUserRepository := new(usermocks.MockUserRepository)
	var testId = actor.ID
	userStatus := usermodels.UserStatusActive
	testUserRepository.On("GetByID", testId).Return(&usermodels.User{
		UserBase: usermodels.UserBase{Name: "Test",
			Email:  "test@testing.com",
			Phone:  "1234567890",
			Status: &userStatus},
		DBBaseModel: sharedmodels.DBBaseModel{
			ID:        1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: time.Time{},
		},
	}, nil)

	uc := NewGetUserUseCase(testUserRepository)

	result := uc.Execute(ctxWithUser, locales.EN_US, testId)

	assert.NotNil(result)
	assert.Equal(result.Data.ID == 1, true)
	assert.Equal(result.Data.Name == "Test", true)
	assert.Equal(result.Details == "", true)
}

func TestGetUserUseCase_DifferentUser(t *testing.T) {
	assert := assert.New(t)

	actor := dtomocks.UserWithRole
	ctxWithUser := app_context.NewContextWithUser(&actor)

	testUserRepository := new(usermocks.MockUserRepository)
	var testId = actor.ID + 1 // Different user ID
	userStatus := usermodels.UserStatusActive
	testUserRepository.On("GetByID", testId).Return(&usermodels.User{
		UserBase: usermodels.UserBase{Name: "Test",
			Email:  "test@testing.com",
			Phone:  "1234567890",
			Status: &userStatus},
		DBBaseModel: sharedmodels.DBBaseModel{
			ID:        2,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: time.Time{},
		},
	}, nil)

	uc := NewGetUserUseCase(testUserRepository)

	result := uc.Execute(ctxWithUser, locales.EN_US, testId)

	assert.NotNil(result)
	assert.Equal(result.HasError(), true)
	assert.Equal(result.StatusCode, status.Unauthorized)

}
