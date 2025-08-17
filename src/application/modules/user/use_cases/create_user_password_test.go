package usecases_user

import (
	"context"
	"gormgoskeleton/src/application/shared/locales"
	"gormgoskeleton/src/application/shared/mocks"
	"gormgoskeleton/src/domain/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateUserAndPassword(t *testing.T) {
	assert := assert.New(t)

	ctx := context.Background()

	testLogger := new(mocks.MockLoggerProvider)
	testUserRepository := new(mocks.MockUserRespository)
	testHashProvider := new(mocks.MockHashProvider)

	userBase := models.UserBase{
		Name:   "Test User",
		Email:  "test@example.com",
		Phone:  "1234567890",
		Status: "active",
		RoleID: 2,
	}

	testUserAndPassword := models.UserAndPasswordCreate{
		UserCreate: models.UserCreate{
			UserBase: userBase,
		},
		Password: "P@ssw0rd",
	}

	testUserAndPasswordHash := testUserAndPassword
	testUserAndPasswordHash.Password = "hashed_password"

	testUserRepository.On("CreateWithPassword", testUserAndPasswordHash).Return(&models.User{
		UserBase: userBase,
		ID:       1,
	}, nil)

	testHashProvider.On("HashPassword", testUserAndPassword.Password).Return("hashed_password", nil)

	// Create the use case
	useCase := NewCreateUserAndPasswordUseCase(
		testLogger,
		testUserRepository,
		testHashProvider,
	)

	// Execute the use case
	result := useCase.Execute(ctx, locales.EN_US, testUserAndPassword)

	// Assert the result
	assert.True(result.IsSuccess())
	assert.NotNil(result.Data)
	assert.Equal(1, result.Data.ID)

}
