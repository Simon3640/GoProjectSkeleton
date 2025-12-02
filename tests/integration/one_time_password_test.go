package integrationtest

import (
	"testing"

	dtomocks "goprojectskeleton/src/application/shared/mocks/dtos"
	database "goprojectskeleton/src/infrastructure/database/goprojectskeleton"
	"goprojectskeleton/src/infrastructure/providers"
	"goprojectskeleton/src/infrastructure/repositories"

	"github.com/stretchr/testify/assert"
)

func TestOneTimePasswordGetByPasswordHash(t *testing.T) {
	assert := assert.New(t)
	oneTimePasswordRepository := repositories.NewOneTimePasswordRepository(database.GoProjectSkeletondb.DB, providers.Logger)
	userRepository := repositories.NewUserRepository(database.GoProjectSkeletondb.DB, providers.Logger)

	// Create user to link the one-time token
	userCreated, _ := userRepository.Create(dtomocks.UserCreate)

	defer userRepository.Delete(userCreated.ID)

	oneTimePasswordCreate := dtomocks.OneTimePasswordCreate
	oneTimePasswordCreate.UserID = userCreated.ID

	// Create one-time token to test GetByPasswordHash

	oneTimePasswordCreated, _ := oneTimePasswordRepository.Create(oneTimePasswordCreate)

	defer oneTimePasswordRepository.Delete(oneTimePasswordCreated.ID)

	// Test GetByPasswordHash
	oneTimePasswordGotten, appErr := oneTimePasswordRepository.GetByPasswordHash(oneTimePasswordCreate.Hash)

	assert.Nil(appErr)
	assert.NotNil(oneTimePasswordGotten)
	assert.Equal(oneTimePasswordCreate.UserID, oneTimePasswordGotten.UserID)
	assert.Equal(oneTimePasswordCreate.Hash, oneTimePasswordGotten.Hash)
	assert.Equal(oneTimePasswordCreate.Purpose, oneTimePasswordGotten.Purpose)
	assert.Equal(false, oneTimePasswordGotten.IsUsed)
}
