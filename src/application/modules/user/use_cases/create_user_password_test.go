package usecases_user

import (
	"context"
	"testing"
	"time"

	dtos "goprojectskeleton/src/application/shared/DTOs"
	"goprojectskeleton/src/application/shared/locales"
	"goprojectskeleton/src/application/shared/mocks"
	"goprojectskeleton/src/domain/models"

	"github.com/stretchr/testify/assert"
)

func TestCreateUserAndPassword(t *testing.T) {
	assert := assert.New(t)

	ctx := context.Background()

	testLogger := new(mocks.MockLoggerProvider)
	testUserRepository := new(mocks.MockUserRepository)
	testHashProvider := new(mocks.MockHashProvider)

	status := models.UserStatusPending
	userBase := models.UserBase{
		Name:   "Test User",
		Email:  "test@example.com",
		Phone:  "1234567890",
		Status: &status,
		RoleID: 2,
	}

	testUserAndPassword := dtos.UserAndPasswordCreate{
		UserCreate: dtos.UserCreate{
			UserBase: userBase,
		},
		Password: "P@ssw0rd",
	}

	testUserAndPasswordHash := testUserAndPassword
	testUserAndPasswordHash.Password = "hashed_password"

	testUserRepository.On("CreateWithPassword", testUserAndPasswordHash).Return(&models.User{
		UserBase: userBase,
		DBBaseModel: models.DBBaseModel{
			ID:        1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: time.Time{},
		},
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
	assert.Equal(uint(1), result.Data.ID)

}
