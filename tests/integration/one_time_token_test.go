package integrationtest

import (
	"testing"

	dtomocks "gormgoskeleton/src/application/shared/mocks/dtos"
	database "gormgoskeleton/src/infrastructure/database/gormgoskeleton"
	"gormgoskeleton/src/infrastructure/providers"
	"gormgoskeleton/src/infrastructure/repositories"

	"github.com/stretchr/testify/assert"
)

func TestOneTimeTokenGetByTokenHash(t *testing.T) {
	assert := assert.New(t)
	oneTimeTokenRepository := repositories.NewOneTimeTokenRepository(database.DB, providers.Logger)
	userRepository := repositories.NewUserRepository(database.DB, providers.Logger)

	// Create user to link the one-time token
	userCreated, _ := userRepository.Create(dtomocks.UserCreate)

	defer userRepository.Delete(userCreated.ID)

	oneTimeTokenCreate := dtomocks.OneTimeTokenCreate
	oneTimeTokenCreate.UserID = userCreated.ID

	// Create one-time token to test GetByTokenHash

	oneTimeTokenCreated, _ := oneTimeTokenRepository.Create(oneTimeTokenCreate)

	defer oneTimeTokenRepository.Delete(oneTimeTokenCreated.ID)

	// Test GetByTokenHash
	oneTimeTokenGotten, appErr := oneTimeTokenRepository.GetByTokenHash(oneTimeTokenCreate.Hash)

	assert.Nil(appErr)
	assert.NotNil(oneTimeTokenGotten)
	assert.Equal(oneTimeTokenCreate.UserID, oneTimeTokenGotten.UserID)
	assert.Equal(oneTimeTokenCreate.Hash, oneTimeTokenGotten.Hash)
	assert.Equal(oneTimeTokenCreate.Purpose, oneTimeTokenGotten.Purpose)
	assert.Equal(false, oneTimeTokenGotten.IsUsed)
}
