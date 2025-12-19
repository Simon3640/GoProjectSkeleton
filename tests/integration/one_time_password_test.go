package integrationtest

import (
	"testing"

	authmocks "github.com/simon3640/goprojectskeleton/src/application/modules/auth/mocks"
	dtomocks "github.com/simon3640/goprojectskeleton/src/application/shared/mocks/dtos"
	database "github.com/simon3640/goprojectskeleton/src/infrastructure/databases/goprojectskeleton"
	authrepositories "github.com/simon3640/goprojectskeleton/src/infrastructure/databases/goprojectskeleton/repositories/auth"
	userrepositories "github.com/simon3640/goprojectskeleton/src/infrastructure/databases/goprojectskeleton/repositories/user"
	"github.com/simon3640/goprojectskeleton/src/infrastructure/providers"

	"github.com/stretchr/testify/assert"
)

func TestOneTimePasswordGetByPasswordHash(t *testing.T) {
	assert := assert.New(t)
	oneTimePasswordRepository := authrepositories.NewOneTimePasswordRepository(database.GoProjectSkeletondb.DB, providers.Logger)
	userRepository := userrepositories.NewUserRepository(database.GoProjectSkeletondb.DB, providers.Logger)

	// Create user to link the one-time token
	userCreated, _ := userRepository.Create(dtomocks.UserCreate)

	defer userRepository.Delete(userCreated.ID)

	oneTimePasswordCreate := authmocks.OneTimePasswordCreate
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
