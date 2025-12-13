package userusecases

import (
	"testing"

	usermocks "github.com/simon3640/goprojectskeleton/src/application/modules/user/mocks"
	app_context "github.com/simon3640/goprojectskeleton/src/application/shared/context"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales"
	dtomocks "github.com/simon3640/goprojectskeleton/src/application/shared/mocks/dtos"
	providersmocks "github.com/simon3640/goprojectskeleton/src/application/shared/mocks/providers"
	"github.com/simon3640/goprojectskeleton/src/application/shared/status"

	"github.com/stretchr/testify/assert"
)

func TestDeleteUserUseCase(t *testing.T) {
	assert := assert.New(t)

	actor := dtomocks.UserWithRole
	ctxWithUser := app_context.NewContextWithUser(&actor)

	testLogger := new(providersmocks.MockLoggerProvider)
	testUserRepository := new(usermocks.MockUserRepository)
	var testIDToDelete = actor.ID

	testUserRepository.On("SoftDelete", testIDToDelete).Return(nil)

	uc := NewDeleteUserUseCase(testLogger, testUserRepository)

	result := uc.Execute(ctxWithUser, locales.EN_US, testIDToDelete)

	assert.NotNil(result)
	assert.Equal(result.StatusCode, status.Success)
	assert.Equal(*result.Data, true)
}

func TestDeleteUserUseCase_DifferentUser(t *testing.T) {
	assert := assert.New(t)

	actor := dtomocks.UserWithRole
	ctxWithUser := app_context.NewContextWithUser(&actor)

	testLogger := new(providersmocks.MockLoggerProvider)
	testUserRepository := new(usermocks.MockUserRepository)
	var testIDToDelete = actor.ID + 1

	testUserRepository.On("SoftDelete", testIDToDelete).Return(nil)

	uc := NewDeleteUserUseCase(testLogger, testUserRepository)

	result := uc.Execute(ctxWithUser, locales.EN_US, testIDToDelete)

	assert.NotNil(result)
	assert.Equal(result.StatusCode, status.Unauthorized)
}
