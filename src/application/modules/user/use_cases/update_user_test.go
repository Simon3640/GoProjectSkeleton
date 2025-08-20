package usecases_user

import (
	"context"
	"gormgoskeleton/src/application/shared/locales"
	"gormgoskeleton/src/application/shared/mocks"
	"gormgoskeleton/src/domain/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdateUserUseCase(t *testing.T) {
	assert := assert.New(t)

	ctx := context.Background()

	testLogger := new(mocks.MockLoggerProvider)
	testUserRepository := new(mocks.MockUserRepository)
	name := "Update"
	testUser := models.UserUpdate{
		UserUpdateBase: models.UserUpdateBase{Name: &name},
		ID:             1,
	}

	testUserRepository.On("Update", testUser.ID, testUser).Return(&models.User{
		UserBase: models.UserBase{Name: "Update",
			Email:  "test@testing.com",
			Phone:  "1234567890",
			Status: "active"},
		ID: 1,
	}, nil)

	uc := NewUpdateUserUseCase(testLogger, testUserRepository)

	result := uc.Execute(ctx, locales.EN_US, testUser)

	assert.NotNil(result)
	assert.Equal(result.Data.ID == 1, true)
	assert.Equal(result.Data.Name == "Update", true)

}
