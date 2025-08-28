package usecases_user

import (
	"context"
	"testing"

	app_context "gormgoskeleton/src/application/shared/context"
	"gormgoskeleton/src/application/shared/locales"
	"gormgoskeleton/src/application/shared/mocks"
	dto_mocks "gormgoskeleton/src/application/shared/mocks/dtos"
	"gormgoskeleton/src/application/shared/status"

	"github.com/stretchr/testify/assert"
)

func TestDeleteUserUseCase(t *testing.T) {
	assert := assert.New(t)

	ctx := context.Background()

	actor := dto_mocks.UserWithRole
	ctxWithUser := context.WithValue(ctx, app_context.UserKey, actor)

	testLogger := new(mocks.MockLoggerProvider)
	testUserRepository := new(mocks.MockUserRepository)
	var testIDToDelete uint = actor.ID

	testUserRepository.On("SoftDelete", testIDToDelete).Return(nil)

	uc := NewDeleteUserUseCase(testLogger, testUserRepository)

	result := uc.Execute(ctxWithUser, locales.EN_US, testIDToDelete)

	assert.NotNil(result)
	assert.Equal(result.StatusCode, status.Success)
}

func TestDeleteUserUseCase_DifferentUser(t *testing.T) {
	assert := assert.New(t)

	ctx := context.Background()

	actor := dto_mocks.UserWithRole
	ctxWithUser := context.WithValue(ctx, app_context.UserKey, actor)

	testLogger := new(mocks.MockLoggerProvider)
	testUserRepository := new(mocks.MockUserRepository)
	var testIDToDelete uint = actor.ID + 1

	testUserRepository.On("Delete", testIDToDelete).Return(nil)

	uc := NewDeleteUserUseCase(testLogger, testUserRepository)

	result := uc.Execute(ctxWithUser, locales.EN_US, testIDToDelete)

	assert.NotNil(result)
	assert.Equal(result.StatusCode, status.Unauthorized)
}
