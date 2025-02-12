package usecases_user

import (
	"context"
	"testing"

	"gormgoskeleton/src/application/shared/locales"
	"gormgoskeleton/src/application/use_cases/mocks"
	"gormgoskeleton/src/domain/models"

	"github.com/stretchr/testify/assert"
)

func TestCreateUserUseCase(t *testing.T) {
	assert := assert.New(t)

	ctx := context.Background()

	testLogger := new(mocks.MockLoggerProvider)
	testUserRepository := new(mocks.MockUserRespository)
	testUser := models.UserCreate{
		UserBase: models.UserBase{Name: "Test",
			Email:  "test@testing.com",
			Phone:  "1234567890",
			Status: "active"},
	}

	testUserRepository.On("Create", testUser).Return(&models.User{
		UserBase: models.UserBase{Name: "Test",
			Email:  "test@testing.com",
			Phone:  "1234567890",
			Status: "active"},
		ID: 1,
	}, nil)

	uc := NewCreateUserUseCase(testLogger, testUserRepository)

	result := uc.Execute(ctx, locales.EN_US, testUser)

	assert.NotNil(result)
	assert.Equal(result.Data.ID == 1, true)
	assert.Equal(result.Data.Name == "Test", true)
}
