package integrationtest

import (
	"testing"

	dtomocks "github.com/simon3640/goprojectskeleton/src/application/shared/mocks/dtos"
	database "github.com/simon3640/goprojectskeleton/src/infrastructure/databases/goprojectskeleton"
	authrepositories "github.com/simon3640/goprojectskeleton/src/infrastructure/databases/goprojectskeleton/repositories/auth"
	userrepositories "github.com/simon3640/goprojectskeleton/src/infrastructure/databases/goprojectskeleton/repositories/user"
	"github.com/simon3640/goprojectskeleton/src/infrastructure/providers"

	"github.com/stretchr/testify/assert"
)

func TestOneTimeTokenGetByTokenHash(t *testing.T) {
	assert := assert.New(t)
	oneTimeTokenRepository := authrepositories.NewOneTimeTokenRepository(database.GoProjectSkeletondb.DB, providers.Logger)
	userRepository := userrepositories.NewUserRepository(database.GoProjectSkeletondb.DB, providers.Logger)

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
