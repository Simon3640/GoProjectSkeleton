package integrationtest

import (
	"testing"

	dtomocks "goprojectskeleton/src/application/shared/mocks/dtos"
	database "goprojectskeleton/src/infrastructure/database/goprojectskeleton"
	"goprojectskeleton/src/infrastructure/providers"
	"goprojectskeleton/src/infrastructure/repositories"

	"github.com/stretchr/testify/assert"
)

func TestPasswordCreate(t *testing.T) {
	assert := assert.New(t)
	passwordRepository := repositories.NewPasswordRepository(database.GoProjectSkeletondb.DB, providers.Logger)

	passwordCreated, appErr := passwordRepository.Create(dtomocks.PasswordCreate)

	assert.Nil(appErr)
	assert.NotNil(passwordCreated)
	assert.Equal(dtomocks.PasswordCreate.UserID, passwordCreated.UserID)
	assert.Equal(dtomocks.PasswordCreate.Hash, passwordCreated.Hash)
	assert.Equal(dtomocks.PasswordCreate.ExpiresAt, passwordCreated.ExpiresAt)
	assert.Equal(dtomocks.PasswordCreate.IsActive, passwordCreated.IsActive)

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

	password_create := dtomocks.PasswordCreate
	password_create.UserID = userCreated.ID

	// Create password to test Get Active Password
	passwordCreated, appErr := passwordRepository.Create(password_create)

	defer passwordRepository.Delete(passwordCreated.ID)

	// Test Get Active Password
	passwordGotten, appErr := passwordRepository.GetActivePassword(userCreated.Email)

	assert.Nil(appErr)
	assert.NotNil(passwordGotten)
	assert.Equal(password_create.UserID, passwordGotten.UserID)
	assert.Equal(password_create.Hash, passwordGotten.Hash)
	assert.Equal(password_create.ExpiresAt, passwordGotten.ExpiresAt)
	assert.Equal(password_create.IsActive, passwordGotten.IsActive)
	passwordRepository.Delete(passwordCreated.ID)
}
