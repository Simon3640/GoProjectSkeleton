package usecases_user

import (
	"context"
	"testing"

	"gormgoskeleton/src/application/shared/locales"
	"gormgoskeleton/src/application/use_cases/mocks"

	"github.com/stretchr/testify/assert"
)

func TestDeleteUserUseCase(t *testing.T) {
	assert := assert.New(t)

	ctx := context.Background()

	testLogger := new(mocks.MockLoggerProvider)
	testUserRepository := new(mocks.MockUserRespository)
	testIDToDelete := 1

	testUserRepository.On("Delete", testIDToDelete).Return(nil)

	uc := NewDeleteUserUseCase(testLogger, testUserRepository)

	result := uc.Execute(ctx, locales.EN_US, testIDToDelete)

	assert.NotNil(result)
}
