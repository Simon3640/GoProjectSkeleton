package userusecases

import (
	"context"
	"testing"
	"time"

	userdtos "github.com/simon3640/goprojectskeleton/src/application/modules/user/dtos"
	usermocks "github.com/simon3640/goprojectskeleton/src/application/modules/user/mocks"
	app_context "github.com/simon3640/goprojectskeleton/src/application/shared/context"
	applicationerror "github.com/simon3640/goprojectskeleton/src/application/shared/errors"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales/messages"
	providersmocks "github.com/simon3640/goprojectskeleton/src/application/shared/mocks/providers"
	appstatus "github.com/simon3640/goprojectskeleton/src/application/shared/status"
	sharedmodels "github.com/simon3640/goprojectskeleton/src/domain/shared/models"
	usermodels "github.com/simon3640/goprojectskeleton/src/domain/user/models"

	"github.com/stretchr/testify/assert"
)

func TestCreateUserAndPassword(t *testing.T) {
	assert := assert.New(t)

	ctx := &app_context.AppContext{Context: context.Background()}

	testUserRepository := new(usermocks.MockUserRepository)
	testHashProvider := new(providersmocks.MockHashProvider)

	status := usermodels.UserStatusPending
	userBase := usermodels.UserBase{
		Name:   "Test User",
		Email:  "test@example.com",
		Phone:  "1234567890",
		Status: &status,
		RoleID: 2,
	}

	testUserAndPassword := userdtos.UserAndPasswordCreate{
		UserCreate: userdtos.UserCreate{
			UserBase: userBase,
		},
		Password: "P@ssw0rd",
	}

	testUserAndPasswordHash := testUserAndPassword
	testUserAndPasswordHash.Password = "hashed_password"

	testUserRepository.On("CreateWithPassword", testUserAndPasswordHash).Return(&usermodels.User{
		UserBase: userBase,
		DBBaseModel: sharedmodels.DBBaseModel{
			ID:        1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: time.Time{},
		},
	}, nil)

	testHashProvider.On("HashPassword", testUserAndPassword.Password).Return("hashed_password", nil)

	// Create the use case
	useCase := NewCreateUserAndPasswordUseCase(
		testUserRepository,
		testHashProvider,
	)

	// Execute the use case
	result := useCase.Execute(ctx, locales.EN_US, testUserAndPassword)

	// Assert the result
	assert.True(result.IsSuccess())
	assert.NotNil(result.Data)
	assert.Equal(uint(1), result.Data.ID)
}

func TestCreateUserAndPassword_InvalidPassword(t *testing.T) {
	assert := assert.New(t)

	ctx := &app_context.AppContext{Context: context.Background()}

	testUserRepository := new(usermocks.MockUserRepository)
	testHashProvider := new(providersmocks.MockHashProvider)

	status := usermodels.UserStatusPending
	userBase := usermodels.UserBase{
		Name:   "Test User",
		Email:  "test@example.com",
		Phone:  "1234567890",
		Status: &status,
		RoleID: 2,
	}

	testUserAndPassword := userdtos.UserAndPasswordCreate{
		UserCreate: userdtos.UserCreate{
			UserBase: userBase,
		},
		Password: "short", // Invalid password (too short)
	}

	// Create the use case
	useCase := NewCreateUserAndPasswordUseCase(
		testUserRepository,
		testHashProvider,
	)

	// Execute the use case
	result := useCase.Execute(ctx, locales.EN_US, testUserAndPassword)

	// Assert the result - should fail validation
	assert.False(result.IsSuccess())
	assert.True(result.HasError())
	assert.Equal(appstatus.InvalidInput, result.StatusCode)
}

func TestCreateUserAndPassword_InvalidEmail(t *testing.T) {
	assert := assert.New(t)

	ctx := &app_context.AppContext{Context: context.Background()}

	testUserRepository := new(usermocks.MockUserRepository)
	testHashProvider := new(providersmocks.MockHashProvider)

	status := usermodels.UserStatusPending
	userBase := usermodels.UserBase{
		Name:   "Test User",
		Email:  "invalid-email", // Invalid email
		Phone:  "1234567890",
		Status: &status,
		RoleID: 2,
	}

	testUserAndPassword := userdtos.UserAndPasswordCreate{
		UserCreate: userdtos.UserCreate{
			UserBase: userBase,
		},
		Password: "P@ssw0rd123",
	}

	// Create the use case
	useCase := NewCreateUserAndPasswordUseCase(
		testUserRepository,
		testHashProvider,
	)

	// Execute the use case
	result := useCase.Execute(ctx, locales.EN_US, testUserAndPassword)

	// Assert the result - should fail validation
	assert.False(result.IsSuccess())
	assert.True(result.HasError())
	assert.Equal(appstatus.InvalidInput, result.StatusCode)
}

func TestCreateUserAndPassword_InvalidRoleID(t *testing.T) {
	assert := assert.New(t)

	ctx := &app_context.AppContext{Context: context.Background()}

	testUserRepository := new(usermocks.MockUserRepository)
	testHashProvider := new(providersmocks.MockHashProvider)

	status := usermodels.UserStatusPending
	userBase := usermodels.UserBase{
		Name:   "Test User",
		Email:  "test@example.com",
		Phone:  "1234567890",
		Status: &status,
		RoleID: 1, // Invalid role ID (admin role cannot be created)
	}

	testUserAndPassword := userdtos.UserAndPasswordCreate{
		UserCreate: userdtos.UserCreate{
			UserBase: userBase,
		},
		Password: "P@ssw0rd123",
	}

	// Create the use case
	useCase := NewCreateUserAndPasswordUseCase(
		testUserRepository,
		testHashProvider,
	)

	// Execute the use case
	result := useCase.Execute(ctx, locales.EN_US, testUserAndPassword)

	// Assert the result - should fail validation
	assert.False(result.IsSuccess())
	assert.True(result.HasError())
	assert.Equal(appstatus.InvalidInput, result.StatusCode)
}

func TestCreateUserAndPassword_HashPasswordError(t *testing.T) {
	assert := assert.New(t)

	ctx := &app_context.AppContext{Context: context.Background()}

	testUserRepository := new(usermocks.MockUserRepository)
	testHashProvider := new(providersmocks.MockHashProvider)

	status := usermodels.UserStatusPending
	userBase := usermodels.UserBase{
		Name:   "Test User",
		Email:  "test@example.com",
		Phone:  "1234567890",
		Status: &status,
		RoleID: 2,
	}

	testUserAndPassword := userdtos.UserAndPasswordCreate{
		UserCreate: userdtos.UserCreate{
			UserBase: userBase,
		},
		Password: "P@ssw0rd123",
	}

	// Mock HashPassword returning error
	appErr := applicationerror.NewApplicationError(
		appstatus.InternalError,
		messages.MessageKeysInstance.SOMETHING_WENT_WRONG,
		"Failed to hash password",
	)
	testHashProvider.On("HashPassword", testUserAndPassword.Password).Return("", appErr)

	// Create the use case
	useCase := NewCreateUserAndPasswordUseCase(
		testUserRepository,
		testHashProvider,
	)

	// Execute the use case
	result := useCase.Execute(ctx, locales.EN_US, testUserAndPassword)

	// Assert the result - should fail with hash error
	assert.False(result.IsSuccess())
	assert.True(result.HasError())
	assert.Equal(appstatus.InternalError, result.StatusCode)
}

func TestCreateUserAndPassword_RepositoryError(t *testing.T) {
	assert := assert.New(t)

	ctx := &app_context.AppContext{Context: context.Background()}

	testUserRepository := new(usermocks.MockUserRepository)
	testHashProvider := new(providersmocks.MockHashProvider)

	status := usermodels.UserStatusPending
	userBase := usermodels.UserBase{
		Name:   "Test User",
		Email:  "test@example.com",
		Phone:  "1234567890",
		Status: &status,
		RoleID: 2,
	}

	testUserAndPassword := userdtos.UserAndPasswordCreate{
		UserCreate: userdtos.UserCreate{
			UserBase: userBase,
		},
		Password: "P@ssw0rd123",
	}

	testUserAndPasswordHash := testUserAndPassword
	testUserAndPasswordHash.Password = "hashed_password"

	// Mock HashPassword success
	testHashProvider.On("HashPassword", testUserAndPassword.Password).Return("hashed_password", nil)

	// Mock CreateWithPassword returning error
	appErr := applicationerror.NewApplicationError(
		appstatus.InternalError,
		messages.MessageKeysInstance.SOMETHING_WENT_WRONG,
		"Failed to create user",
	)
	var nilUser *usermodels.User
	testUserRepository.On("CreateWithPassword", testUserAndPasswordHash).Return(nilUser, appErr)

	// Create the use case
	useCase := NewCreateUserAndPasswordUseCase(
		testUserRepository,
		testHashProvider,
	)

	// Execute the use case
	result := useCase.Execute(ctx, locales.EN_US, testUserAndPassword)

	// Assert the result - should fail with repository error
	assert.False(result.IsSuccess())
	assert.True(result.HasError())
	assert.Equal(appstatus.InternalError, result.StatusCode)
}

func TestCreateUserAndPassword_SetLocale(t *testing.T) {
	assert := assert.New(t)

	testUserRepository := new(usermocks.MockUserRepository)
	testHashProvider := new(providersmocks.MockHashProvider)

	uc := NewCreateUserAndPasswordUseCase(
		testUserRepository,
		testHashProvider,
	)

	// Test setting locale
	uc.SetLocale(locales.ES_ES)
	assert.Equal(locales.ES_ES, uc.Locale)

	// Test setting empty locale (should not change)
	uc.SetLocale("")
	assert.Equal(locales.EN_US, uc.Locale)

	// Test setting another locale
	uc.SetLocale(locales.EN_US)
	assert.Equal(locales.EN_US, uc.Locale)
}
