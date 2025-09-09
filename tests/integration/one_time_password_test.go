package integrationtest

import (
	"testing"

	dtomocks "gormgoskeleton/src/application/shared/mocks/dtos"
	database "gormgoskeleton/src/infrastructure/database/gormgoskeleton"
	"gormgoskeleton/src/infrastructure/providers"
	"gormgoskeleton/src/infrastructure/repositories"

	"github.com/stretchr/testify/assert"
)

func TestOneTimePasswordGetByPasswordHash(t *testing.T) {
	assert := assert.New(t)
	oneTimePasswordRepository := repositories.NewOneTimePasswordRepository(database.DB, providers.Logger)
	userRepository := repositories.NewUserRepository(database.DB, providers.Logger)

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
