package usecases_user

import (
	"context"
	"testing"
	"time"

	dtos "gormgoskeleton/src/application/shared/DTOs"
	"gormgoskeleton/src/application/shared/locales"
	"gormgoskeleton/src/application/shared/mocks"
	dto_mocks "gormgoskeleton/src/application/shared/mocks/dtos"
	"gormgoskeleton/src/application/shared/status"
	"gormgoskeleton/src/domain/models"

	"github.com/stretchr/testify/assert"
)

func TestCreateUserUseCase(t *testing.T) {
	assert := assert.New(t)

	ctx := context.Background()

	testLogger := new(mocks.MockLoggerProvider)
	testUserRepository := new(mocks.MockUserRepository)
	testUser := dto_mocks.UserCreate

	testUserRepository.On("Create", testUser).Return(&models.User{
		UserBase: models.UserBase{Name: "Test",
			Email:  "test@testing.com",
			Phone:  "1234567890",
			Status: "active",
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

	testLogger := new(mocks.MockLoggerProvider)
	testUserRepository := new(mocks.MockUserRepository)
	testUserInvalidRoleID := dtos.UserCreate{
		UserBase: models.UserBase{Name: "Test",
			Email:  "invalid@gmail.com",
			Phone:  "1234567890",
			Status: "active",
			RoleID: 1,
		},
	}

	uc := NewCreateUserUseCase(testLogger, testUserRepository)

	result := uc.Execute(ctx, locales.EN_US, testUserInvalidRoleID)

	assert.NotNil(result)

	assert.Equal(result.HasError(), true)
	assert.Equal(result.StatusCode, status.InvalidInput)

}
