package usecases_user

import (
	"context"
	"testing"

	"gormgoskeleton/src/application/shared/locales"
	"gormgoskeleton/src/application/shared/mocks"
	"gormgoskeleton/src/domain/models"

	"github.com/stretchr/testify/assert"
)

func TestGetUserUseCase(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()

	testLogger := new(mocks.MockLoggerProvider)
	testUserRepository := new(mocks.MockUserRespository)
	var testId uint = 1

	testUserRepository.On("GetByID", testId).Return(&models.User{
		UserBase: models.UserBase{Name: "Test",
			Email:  "test@testing.com",
			Phone:  "1234567890",
			Status: "active"},
		ID: 1,
	}, nil)

	uc := NewGetUserUseCase(testLogger, testUserRepository)

	result := uc.Execute(ctx, locales.EN_US, testId)

	assert.NotNil(result)
	assert.Equal(result.Data.ID == 1, true)
	assert.Equal(result.Data.Name == "Test", true)
	assert.Equal(result.Details == "", true)
}
