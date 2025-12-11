package integrationtest

import (
	"testing"

	dtomocks "github.com/simon3640/goprojectskeleton/src/application/shared/mocks/dtos"
	database "github.com/simon3640/goprojectskeleton/src/infrastructure/databases/goprojectskeleton"
	repositories "github.com/simon3640/goprojectskeleton/src/infrastructure/databases/goprojectskeleton/repositories"
	"github.com/simon3640/goprojectskeleton/src/infrastructure/providers"

	"github.com/stretchr/testify/assert"
)

func TestUserRepository_CreateWithPassword(t *testing.T) {
	assert := assert.New(t)
	userRepository := repositories.NewUserRepository(database.GoProjectSkeletondb.DB, providers.Logger)

	// Test Create With password
	createdUser, appErr := userRepository.CreateWithPassword(dtomocks.UserAndPasswordCreate)

	assert.Nil(appErr)
	assert.NotNil(createdUser)
	assert.Equal(dtomocks.UserBase.Name, createdUser.Name)
	assert.Equal(dtomocks.UserBase.Email, createdUser.Email)
	assert.Equal(dtomocks.UserBase.Phone, createdUser.Phone)
	assert.Equal(dtomocks.UserBase.Status, createdUser.Status)
	assert.Equal(dtomocks.UserBase.RoleID, createdUser.RoleID)
	assert.NotEmpty(createdUser.ID)
	assert.NotEmpty(createdUser.CreatedAt)
	assert.NotEmpty(createdUser.UpdatedAt)
	assert.True(createdUser.DeletedAt.IsZero())

	userRepository.Delete(createdUser.ID)
}

func TestUserRepository_GetUserWithRole(t *testing.T) {
	assert := assert.New(t)
	userRepository := repositories.NewUserRepository(database.GoProjectSkeletondb.DB, providers.Logger)
	// Create user to test Get User With Role
	createdUser, appErr := userRepository.CreateWithPassword(dtomocks.UserAndPasswordCreate)

	// Test Get User With Role
	userWithRole, appErr := userRepository.GetUserWithRole(createdUser.ID)
	assert.Nil(appErr)
	assert.NotNil(userWithRole)
	assert.Equal(createdUser.ID, userWithRole.ID)
	assert.Equal(createdUser.Name, userWithRole.Name)
	assert.Equal(createdUser.Email, userWithRole.Email)
	assert.Equal(createdUser.Phone, userWithRole.Phone)
	assert.Equal(createdUser.Status, userWithRole.Status)
	assert.Equal(createdUser.RoleID, userWithRole.RoleID)
	assert.NotNil(userWithRole.GetRoleKey())

	// Delete created user
	userRepository.Delete(createdUser.ID)
}

func TestUserRepository_GetByEmailOrPhone(t *testing.T) {
	assert := assert.New(t)
	userRepository := repositories.NewUserRepository(database.GoProjectSkeletondb.DB, providers.Logger)

	createdUser, _ := userRepository.CreateWithPassword(dtomocks.UserAndPasswordCreate)

	// Test Get By Email
	userByEmail, appErr := userRepository.GetByEmailOrPhone(createdUser.Email)
	assert.Nil(appErr)
	assert.NotNil(userByEmail)
	assert.Equal(createdUser.Name, userByEmail.Name)
	assert.Equal(createdUser.Email, userByEmail.Email)
	assert.Equal(createdUser.Phone, userByEmail.Phone)
	assert.Equal(createdUser.Status, userByEmail.Status)
	assert.Equal(createdUser.RoleID, userByEmail.RoleID)
	assert.NotEmpty(userByEmail.ID)
	assert.NotEmpty(userByEmail.CreatedAt)
	assert.NotEmpty(userByEmail.UpdatedAt)
	assert.True(userByEmail.DeletedAt.IsZero())

	// Test Get By Phone
	userByPhone, appErr := userRepository.GetByEmailOrPhone(createdUser.Phone)
	assert.Nil(appErr)
	assert.NotNil(userByPhone)
	assert.Equal(createdUser.Name, userByPhone.Name)
	assert.Equal(createdUser.Email, userByPhone.Email)
	assert.Equal(createdUser.Phone, userByPhone.Phone)
	assert.Equal(createdUser.Status, userByPhone.Status)
	assert.Equal(createdUser.RoleID, userByPhone.RoleID)
	assert.NotEmpty(userByPhone.ID)
	assert.NotEmpty(userByPhone.CreatedAt)
	assert.NotEmpty(userByPhone.UpdatedAt)
	assert.True(userByPhone.DeletedAt.IsZero())

	userRepository.Delete(createdUser.ID)
}
