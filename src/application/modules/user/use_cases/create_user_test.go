package userusecases

import (
	"context"
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

func TestCreateUserUseCase(t *testing.T) {
	assert := assert.New(t)

	ctx := &app_context.AppContext{Context: context.Background()}

	testUserRepository := new(usermocks.MockUserRepository)
	testUser := dtomocks.UserCreate
	userStatus := usermodels.UserStatusPending
	testUserRepository.On("Create", testUser).Return(&usermodels.User{
		UserBase: usermodels.UserBase{Name: "Test",
			Email:  "test@testing.com",
			Phone:  "1234567890",
			Status: &userStatus,
			RoleID: 2,
		},
		DBBaseModel: sharedmodels.DBBaseModel{
			ID:        1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: time.Time{},
		},
	}, nil)

	uc := NewCreateUserUseCase(testUserRepository)

	result := uc.Execute(ctx, locales.EN_US, testUser)

	assert.NotNil(result)
	assert.Equal(result.Data.ID == 1, true)
	assert.Equal(result.Data.Name == "Test", true)
}

func TestCreateUserUseCase_InvalidInput(t *testing.T) {
	assert := assert.New(t)

	ctx := &app_context.AppContext{Context: context.Background()}

	testUserRepository := new(usermocks.MockUserRepository)
	userStatus := usermodels.UserStatusPending
	testUserInvalidRoleID := userdtos.UserCreate{
		UserBase: usermodels.UserBase{Name: "Test",
			Email:  "invalid@gmail.com",
			Phone:  "1234567890",
			Status: &userStatus,
			RoleID: 1,
		},
	}

	uc := NewCreateUserUseCase(testUserRepository)

	result := uc.Execute(ctx, locales.EN_US, testUserInvalidRoleID)

	assert.NotNil(result)

	assert.Equal(result.HasError(), true)
	assert.Equal(result.StatusCode, status.InvalidInput)

}
