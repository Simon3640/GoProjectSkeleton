package userusecases

import (
	"context"
	"testing"

	app_context "goprojectskeleton/src/application/shared/context"
	"goprojectskeleton/src/application/shared/locales"
	dtomocks "goprojectskeleton/src/application/shared/mocks/dtos"
	providersmocks "goprojectskeleton/src/application/shared/mocks/providers"
	repositoriesmocks "goprojectskeleton/src/application/shared/mocks/repositories"
	"goprojectskeleton/src/application/shared/status"

	"github.com/stretchr/testify/assert"
)

func TestDeleteUserUseCase(t *testing.T) {
	assert := assert.New(t)

	ctx := context.Background()

	actor := dtomocks.UserWithRole
	ctxWithUser := context.WithValue(ctx, app_context.UserKey, actor)

	testLogger := new(providersmocks.MockLoggerProvider)
	testUserRepository := new(repositoriesmocks.MockUserRepository)
	var testIDToDelete = actor.ID

	testUserRepository.On("SoftDelete", testIDToDelete).Return(nil)

	uc := NewDeleteUserUseCase(testLogger, testUserRepository)

	result := uc.Execute(ctxWithUser, locales.EN_US, testIDToDelete)

	assert.NotNil(result)
	assert.Equal(result.StatusCode, status.Success)
}

func TestDeleteUserUseCase_DifferentUser(t *testing.T) {
	assert := assert.New(t)

	ctx := context.Background()

	actor := dtomocks.UserWithRole
	ctxWithUser := context.WithValue(ctx, app_context.UserKey, actor)

	testLogger := new(providersmocks.MockLoggerProvider)
	testUserRepository := new(repositoriesmocks.MockUserRepository)
	var testIDToDelete = actor.ID + 1

	testUserRepository.On("Delete", testIDToDelete).Return(nil)

	uc := NewDeleteUserUseCase(testLogger, testUserRepository)

	result := uc.Execute(ctxWithUser, locales.EN_US, testIDToDelete)

	assert.NotNil(result)
	assert.Equal(result.StatusCode, status.Unauthorized)
}
