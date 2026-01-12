package userusecases

import (
	"testing"
	"time"

	userdtos "github.com/simon3640/goprojectskeleton/src/application/modules/user/dtos"
	usermocks "github.com/simon3640/goprojectskeleton/src/application/modules/user/mocks"
	app_context "github.com/simon3640/goprojectskeleton/src/application/shared/context"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales"
	dtomocks "github.com/simon3640/goprojectskeleton/src/application/shared/mocks/dtos"
	"github.com/simon3640/goprojectskeleton/src/application/shared/status"
	sharedmodels "github.com/simon3640/goprojectskeleton/src/domain/shared/models"
	usermodels "github.com/simon3640/goprojectskeleton/src/domain/user/models"

	"github.com/stretchr/testify/assert"
)

func TestUpdateUserUseCase(t *testing.T) {
	assert := assert.New(t)

	actor := dtomocks.UserWithRole
	ctxWithUser := app_context.NewContextWithUser(&actor)

	testUserRepository := new(usermocks.MockUserRepository)
	name := "Update"
	testUser := userdtos.UserUpdate{
		UserUpdateBase: usermodels.UserUpdateBase{Name: &name},
		ID:             actor.ID,
	}
	userStatus := usermodels.UserStatusActive
	testUserRepository.On("Update", testUser.ID, testUser).Return(&usermodels.User{
		UserBase: usermodels.UserBase{Name: "Update",
			Email:  "test@testing.com",
			Phone:  "1234567890",
			Status: &userStatus},
		DBBaseModel: sharedmodels.DBBaseModel{
			ID:        actor.ID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: time.Time{},
		},
	}, nil)

	uc := NewUpdateUserUseCase(testUserRepository)

	result := uc.Execute(ctxWithUser, locales.EN_US, testUser)

	assert.NotNil(result)
	assert.Equal(result.Data.ID == 1, true)
	assert.Equal(result.Data.Name == "Update", true)

}

func TestUpdateUserUseCase_DifferentUser(t *testing.T) {
	assert := assert.New(t)

	actor := dtomocks.UserWithRole
	ctxWithUser := app_context.NewContextWithUser(&actor)

	testUserRepository := new(usermocks.MockUserRepository)
	name := "Update"
	testUser := userdtos.UserUpdate{
		UserUpdateBase: usermodels.UserUpdateBase{Name: &name},
		ID:             actor.ID + 1,
	}

	testUserRepository.On("Update", testUser.ID, testUser).Return(nil)

	uc := NewUpdateUserUseCase(testUserRepository)

	result := uc.Execute(ctxWithUser, locales.EN_US, testUser)

	assert.NotNil(result)
	assert.Equal(result.StatusCode, status.Unauthorized)
}
