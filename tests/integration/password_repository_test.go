package integrationtest

import (
	"testing"

	passwordmocks "github.com/simon3640/goprojectskeleton/src/application/modules/password/mocks"
	dtomocks "github.com/simon3640/goprojectskeleton/src/application/shared/mocks/dtos"
	database "github.com/simon3640/goprojectskeleton/src/infrastructure/databases/goprojectskeleton"
	repositories "github.com/simon3640/goprojectskeleton/src/infrastructure/databases/goprojectskeleton/repositories"
	"github.com/simon3640/goprojectskeleton/src/infrastructure/providers"

	"github.com/stretchr/testify/assert"
)

func TestPasswordCreate(t *testing.T) {
	assert := assert.New(t)
	passwordRepository := repositories.NewPasswordRepository(database.GoProjectSkeletondb.DB, providers.Logger)

	passwordCreated, appErr := passwordRepository.Create(passwordmocks.PasswordCreate)

	assert.Nil(appErr)
	assert.NotNil(passwordCreated)
	assert.Equal(passwordmocks.PasswordCreate.UserID, passwordCreated.UserID)
	assert.Equal(passwordmocks.PasswordCreate.Hash, passwordCreated.Hash)
	assert.Equal(passwordmocks.PasswordCreate.ExpiresAt, passwordCreated.ExpiresAt)
	assert.Equal(passwordmocks.PasswordCreate.IsActive, passwordCreated.IsActive)

	passwordRepository.Delete(passwordCreated.ID)
}

func TestPasswordGetActivePassword(t *testing.T) {
	assert := assert.New(t)
	passwordRepository := repositories.NewPasswordRepository(database.GoProjectSkeletondb.DB, providers.Logger)
	userRepository := repositories.NewUserRepository(database.GoProjectSkeletondb.DB, providers.Logger)

	// Create user to link the password
	userCreated, appErr := userRepository.Create(dtomocks.UserCreate)
	assert.NotNil(userCreated)
	assert.Nil(appErr)

	defer userRepository.Delete(userCreated.ID)

	passwordCreate := passwordmocks.PasswordCreate
	passwordCreate.UserID = userCreated.ID

	// Create password to test Get Active Password
	passwordCreated, _ := passwordRepository.Create(passwordCreate)

	defer passwordRepository.Delete(passwordCreated.ID)

	// Test Get Active Password
	passwordGotten, appErr := passwordRepository.GetActivePassword(userCreated.Email)

	assert.Nil(appErr)
	assert.NotNil(passwordGotten)
	assert.Equal(passwordCreate.UserID, passwordGotten.UserID)
	assert.Equal(passwordCreate.Hash, passwordGotten.Hash)
	assert.Equal(passwordCreate.ExpiresAt, passwordGotten.ExpiresAt)
	assert.Equal(passwordCreate.IsActive, passwordGotten.IsActive)
	passwordRepository.Delete(passwordCreated.ID)
}
