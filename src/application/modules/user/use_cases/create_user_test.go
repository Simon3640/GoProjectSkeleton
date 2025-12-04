package userusecases

import (
	"context"
	"testing"
	"time"

	dtos "goprojectskeleton/src/application/shared/DTOs"
	"goprojectskeleton/src/application/shared/locales"
	dtomocks "goprojectskeleton/src/application/shared/mocks/dtos"
	providersmocks "goprojectskeleton/src/application/shared/mocks/providers"
	repositoriesmocks "goprojectskeleton/src/application/shared/mocks/repositories"
	"goprojectskeleton/src/application/shared/status"
	"goprojectskeleton/src/domain/models"

	"github.com/stretchr/testify/assert"
)

func TestCreateUserUseCase(t *testing.T) {
	assert := assert.New(t)

	ctx := context.Background()

	testLogger := new(providersmocks.MockLoggerProvider)
	testUserRepository := new(repositoriesmocks.MockUserRepository)
	testUser := dtomocks.UserCreate
	userStatus := models.UserStatusPending
	testUserRepository.On("Create", testUser).Return(&models.User{
		UserBase: models.UserBase{Name: "Test",
			Email:  "test@testing.com",
			Phone:  "1234567890",
			Status: &userStatus,
			RoleID: 2,
		},
		DBBaseModel: models.DBBaseModel{
			ID:        1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: time.Time{},
		},
	}, nil)

	uc := NewCreateUserUseCase(testLogger, testUserRepository)

	result := uc.Execute(ctx, locales.EN_US, testUser)

	assert.NotNil(result)
	assert.Equal(result.Data.ID == 1, true)
	assert.Equal(result.Data.Name == "Test", true)
}

func TestCreateUserUseCase_InvalidInput(t *testing.T) {
	assert := assert.New(t)

	ctx := context.Background()

	testLogger := new(providersmocks.MockLoggerProvider)
	testUserRepository := new(repositoriesmocks.MockUserRepository)
	userStatus := models.UserStatusPending
	testUserInvalidRoleID := dtos.UserCreate{
		UserBase: models.UserBase{Name: "Test",
			Email:  "invalid@gmail.com",
			Phone:  "1234567890",
			Status: &userStatus,
			RoleID: 1,
		},
	}

	uc := NewCreateUserUseCase(testLogger, testUserRepository)

	result := uc.Execute(ctx, locales.EN_US, testUserInvalidRoleID)

	assert.NotNil(result)

	assert.Equal(result.HasError(), true)
	assert.Equal(result.StatusCode, status.InvalidInput)

}
