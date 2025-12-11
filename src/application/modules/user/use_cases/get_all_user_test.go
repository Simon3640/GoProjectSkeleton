package userusecases

import (
	"context"
	"testing"
	"time"

	app_context "github.com/simon3640/goprojectskeleton/src/application/shared/context"
	application_errors "github.com/simon3640/goprojectskeleton/src/application/shared/errors"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales/messages"
	dtomocks "github.com/simon3640/goprojectskeleton/src/application/shared/mocks/dtos"
	providersmocks "github.com/simon3640/goprojectskeleton/src/application/shared/mocks/providers"
	repositoriesmocks "github.com/simon3640/goprojectskeleton/src/application/shared/mocks/repositories"
	appstatus "github.com/simon3640/goprojectskeleton/src/application/shared/status"
	"github.com/simon3640/goprojectskeleton/src/domain/models"
	domain_utils "github.com/simon3640/goprojectskeleton/src/domain/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetAllUserUseCase_Execute_SuccessFromCache(t *testing.T) {
	assert := assert.New(t)

	ctx := context.Background()

	// Create admin user for context
	adminUser := models.UserWithRole{
		UserBase: models.UserBase{
			Name:   "Admin User",
			Email:  "admin@example.com",
			Phone:  "1234567890",
			RoleID: 1,
		},
		ID: 1,
	}
	adminUser.SetRole(dtomocks.AdminRole)
	ctxWithUser := context.WithValue(ctx, app_context.UserKey, adminUser)

	testLogger := new(providersmocks.MockLoggerProvider)
	testUserRepository := new(repositoriesmocks.MockUserRepository)
	testCacheProvider := new(providersmocks.MockCacheProvider)

	// Create query payload
	page := 1
	pageSize := 10
	queryPayload := domain_utils.NewQueryPayloadBuilder[models.User](nil, nil, &page, &pageSize)

	// Mock cache - data found
	testUsers := []models.User{
		{
			UserBase: models.UserBase{
				Name:   "User 1",
				Email:  "user1@example.com",
				Phone:  "1234567890",
				RoleID: 2,
			},
			DBBaseModel: models.DBBaseModel{
				ID:        1,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		},
	}
	var total int64 = 1

	cacheKey := "users:" + queryPayload.GetQueryKey()
	testCacheProvider.On("Get", cacheKey, mock.AnythingOfType("*[]models.User")).Return(true, nil, testUsers).Run(func(args mock.Arguments) {
		dest := args.Get(1).(*[]models.User)
		*dest = testUsers
	})
	testCacheProvider.On("Get", cacheKey+":total", mock.AnythingOfType("*int64")).Return(true, nil, total).Run(func(args mock.Arguments) {
		dest := args.Get(1).(*int64)
		*dest = total
	})
	testLogger.On("Debug", mock.Anything, mock.Anything).Return()

	uc := NewGetAllUserUseCase(
		testLogger,
		testUserRepository,
		testCacheProvider,
	)

	result := uc.Execute(ctxWithUser, locales.EN_US, queryPayload)

	assert.NotNil(result)
	assert.True(result.IsSuccess())
	assert.False(result.HasError())
	assert.NotNil(result.Data)
	assert.Equal(1, len(result.Data.Records))
	assert.Equal(int64(1), result.Data.Meta.Total)
	assert.True(result.Data.Meta.Cached)
}

func TestGetAllUserUseCase_Execute_SuccessFromRepository(t *testing.T) {
	assert := assert.New(t)

	ctx := context.Background()

	// Create admin user for context
	adminUser := models.UserWithRole{
		UserBase: models.UserBase{
			Name:   "Admin User",
			Email:  "admin@example.com",
			Phone:  "1234567890",
			RoleID: 1,
		},
		ID: 1,
	}
	adminUser.SetRole(dtomocks.AdminRole)
	ctxWithUser := context.WithValue(ctx, app_context.UserKey, adminUser)

	testLogger := new(providersmocks.MockLoggerProvider)
	testUserRepository := new(repositoriesmocks.MockUserRepository)
	testCacheProvider := new(providersmocks.MockCacheProvider)

	// Create query payload
	page := 1
	pageSize := 10
	queryPayload := domain_utils.NewQueryPayloadBuilder[models.User](nil, nil, &page, &pageSize)

	// Mock cache - not found
	cacheKey := "users:" + queryPayload.GetQueryKey()
	testCacheProvider.On("Get", cacheKey, mock.AnythingOfType("*[]models.User")).Return(false, nil)
	testCacheProvider.On("Get", cacheKey+":total", mock.AnythingOfType("*int64")).Return(false, nil)

	// Mock repository
	testUsers := []models.User{
		{
			UserBase: models.UserBase{
				Name:   "User 1",
				Email:  "user1@example.com",
				Phone:  "1234567890",
				RoleID: 2,
			},
			DBBaseModel: models.DBBaseModel{
				ID:        1,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		},
	}
	var total int64 = 1
	testUserRepository.On("GetAll", mock.Anything, 0, 10).Return(testUsers, total, nil)

	// Mock cache set
	testCacheProvider.On("Set", cacheKey, testUsers, mock.AnythingOfType("time.Duration")).Return(nil)
	testCacheProvider.On("Set", cacheKey+":total", total, mock.AnythingOfType("time.Duration")).Return(nil)

	uc := NewGetAllUserUseCase(
		testLogger,
		testUserRepository,
		testCacheProvider,
	)

	result := uc.Execute(ctxWithUser, locales.EN_US, queryPayload)

	assert.NotNil(result)
	assert.True(result.IsSuccess())
	assert.False(result.HasError())
	assert.NotNil(result.Data)
	assert.Equal(1, len(result.Data.Records))
	assert.Equal(int64(1), result.Data.Meta.Total)
	assert.False(result.Data.Meta.Cached)
}

func TestGetAllUserUseCase_Execute_RepositoryError(t *testing.T) {
	assert := assert.New(t)

	ctx := context.Background()

	// Create admin user for context
	adminUser := models.UserWithRole{
		UserBase: models.UserBase{
			Name:   "Admin User",
			Email:  "admin@example.com",
			Phone:  "1234567890",
			RoleID: 1,
		},
		ID: 1,
	}
	adminUser.SetRole(dtomocks.AdminRole)
	ctxWithUser := context.WithValue(ctx, app_context.UserKey, adminUser)

	testLogger := new(providersmocks.MockLoggerProvider)
	testUserRepository := new(repositoriesmocks.MockUserRepository)
	testCacheProvider := new(providersmocks.MockCacheProvider)

	// Create query payload
	page := 1
	pageSize := 10
	queryPayload := domain_utils.NewQueryPayloadBuilder[models.User](nil, nil, &page, &pageSize)

	// Mock cache - not found
	cacheKey := "users:" + queryPayload.GetQueryKey()
	testCacheProvider.On("Get", cacheKey, mock.AnythingOfType("*[]models.User")).Return(false, nil)
	testCacheProvider.On("Get", cacheKey+":total", mock.AnythingOfType("*int64")).Return(false, nil)

	// Mock repository error
	appErr := application_errors.NewApplicationError(
		appstatus.InternalError,
		messages.MessageKeysInstance.SOMETHING_WENT_WRONG,
		"Database error",
	)
	testUserRepository.On("GetAll", mock.Anything, 0, 10).Return(nil, int64(0), appErr)
	testLogger.On("Error", mock.Anything, mock.Anything).Return()

	uc := NewGetAllUserUseCase(
		testLogger,
		testUserRepository,
		testCacheProvider,
	)

	result := uc.Execute(ctxWithUser, locales.EN_US, queryPayload)

	assert.NotNil(result)
	assert.False(result.IsSuccess())
	assert.True(result.HasError())
	assert.Equal(appstatus.InternalError, result.StatusCode)
}

func TestGetAllUserUseCase_Execute_CacheGetError(t *testing.T) {
	assert := assert.New(t)

	ctx := context.Background()

	// Create admin user for context
	adminUser := models.UserWithRole{
		UserBase: models.UserBase{
			Name:   "Admin User",
			Email:  "admin@example.com",
			Phone:  "1234567890",
			RoleID: 1,
		},
		ID: 1,
	}
	adminUser.SetRole(dtomocks.AdminRole)
	ctxWithUser := context.WithValue(ctx, app_context.UserKey, adminUser)

	testLogger := new(providersmocks.MockLoggerProvider)
	testUserRepository := new(repositoriesmocks.MockUserRepository)
	testCacheProvider := new(providersmocks.MockCacheProvider)

	// Create query payload
	page := 1
	pageSize := 10
	queryPayload := domain_utils.NewQueryPayloadBuilder[models.User](nil, nil, &page, &pageSize)

	// Mock cache - error getting cache (should not fail, just log)
	cacheKey := "users:" + queryPayload.GetQueryKey()
	appErr := application_errors.NewApplicationError(
		appstatus.InternalError,
		messages.MessageKeysInstance.SOMETHING_WENT_WRONG,
		"Cache error",
	)
	testCacheProvider.On("Get", cacheKey, mock.AnythingOfType("*[]models.User")).Return(false, appErr)
	testLogger.On("Error", mock.Anything, mock.Anything).Return()

	// Mock repository
	testUsers := []models.User{
		{
			UserBase: models.UserBase{
				Name:   "User 1",
				Email:  "user1@example.com",
				Phone:  "1234567890",
				RoleID: 2,
			},
			DBBaseModel: models.DBBaseModel{
				ID:        1,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		},
	}
	var total int64 = 1
	testUserRepository.On("GetAll", mock.Anything, 0, 10).Return(testUsers, total, nil)

	// Mock cache set
	testCacheProvider.On("Set", cacheKey, testUsers, mock.AnythingOfType("time.Duration")).Return(nil)
	testCacheProvider.On("Set", cacheKey+":total", total, mock.AnythingOfType("time.Duration")).Return(nil)

	uc := NewGetAllUserUseCase(
		testLogger,
		testUserRepository,
		testCacheProvider,
	)

	result := uc.Execute(ctxWithUser, locales.EN_US, queryPayload)

	// Should still succeed even if cache has error
	assert.NotNil(result)
	assert.True(result.IsSuccess())
	assert.False(result.HasError())
	assert.NotNil(result.Data)
}

func TestGetAllUserUseCase_Execute_CacheSetError(t *testing.T) {
	assert := assert.New(t)

	ctx := context.Background()

	// Create admin user for context
	adminUser := models.UserWithRole{
		UserBase: models.UserBase{
			Name:   "Admin User",
			Email:  "admin@example.com",
			Phone:  "1234567890",
			RoleID: 1,
		},
		ID: 1,
	}
	adminUser.SetRole(dtomocks.AdminRole)
	ctxWithUser := context.WithValue(ctx, app_context.UserKey, adminUser)

	testLogger := new(providersmocks.MockLoggerProvider)
	testUserRepository := new(repositoriesmocks.MockUserRepository)
	testCacheProvider := new(providersmocks.MockCacheProvider)

	// Create query payload
	page := 1
	pageSize := 10
	queryPayload := domain_utils.NewQueryPayloadBuilder[models.User](nil, nil, &page, &pageSize)

	// Mock cache - not found
	cacheKey := "users:" + queryPayload.GetQueryKey()
	testCacheProvider.On("Get", cacheKey, mock.AnythingOfType("*[]models.User")).Return(false, nil)
	testCacheProvider.On("Get", cacheKey+":total", mock.AnythingOfType("*int64")).Return(false, nil)

	// Mock repository
	testUsers := []models.User{
		{
			UserBase: models.UserBase{
				Name:   "User 1",
				Email:  "user1@example.com",
				Phone:  "1234567890",
				RoleID: 2,
			},
			DBBaseModel: models.DBBaseModel{
				ID:        1,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		},
	}
	var total int64 = 1
	testUserRepository.On("GetAll", mock.Anything, 0, 10).Return(testUsers, total, nil)

	// Mock cache set error (should not fail, just log)
	appErr := application_errors.NewApplicationError(
		appstatus.InternalError,
		messages.MessageKeysInstance.SOMETHING_WENT_WRONG,
		"Cache set error",
	)
	testCacheProvider.On("Set", cacheKey, testUsers, mock.AnythingOfType("time.Duration")).Return(appErr)
	testCacheProvider.On("Set", cacheKey+":total", total, mock.AnythingOfType("time.Duration")).Return(appErr)
	testLogger.On("Error", mock.Anything, mock.Anything).Return()

	uc := NewGetAllUserUseCase(
		testLogger,
		testUserRepository,
		testCacheProvider,
	)

	result := uc.Execute(ctxWithUser, locales.EN_US, queryPayload)

	// Should still succeed even if cache set has error
	assert.NotNil(result)
	assert.True(result.IsSuccess())
	assert.False(result.HasError())
	assert.NotNil(result.Data)
}

func TestGetAllUserUseCase_Execute_Unauthorized(t *testing.T) {
	assert := assert.New(t)

	ctx := context.Background()

	// Create non-admin user for context
	regularUser := models.UserWithRole{
		UserBase: models.UserBase{
			Name:   "Regular User",
			Email:  "user@example.com",
			Phone:  "1234567890",
			RoleID: 2,
		},
		ID: 1,
	}
	regularUser.SetRole(dtomocks.UserRole)
	ctxWithUser := context.WithValue(ctx, app_context.UserKey, regularUser)

	testLogger := new(providersmocks.MockLoggerProvider)
	testUserRepository := new(repositoriesmocks.MockUserRepository)
	testCacheProvider := new(providersmocks.MockCacheProvider)

	// Create query payload
	page := 1
	pageSize := 10
	queryPayload := domain_utils.NewQueryPayloadBuilder[models.User](nil, nil, &page, &pageSize)

	uc := NewGetAllUserUseCase(
		testLogger,
		testUserRepository,
		testCacheProvider,
	)

	result := uc.Execute(ctxWithUser, locales.EN_US, queryPayload)

	assert.NotNil(result)
	assert.False(result.IsSuccess())
	assert.True(result.HasError())
	assert.Equal(appstatus.Unauthorized, result.StatusCode)
}

func TestGetAllUserUseCase_Execute_InvalidInput(t *testing.T) {
	assert := assert.New(t)

	ctx := context.Background()

	// Create admin user for context
	adminUser := models.UserWithRole{
		UserBase: models.UserBase{
			Name:   "Admin User",
			Email:  "admin@example.com",
			Phone:  "1234567890",
			RoleID: 1,
		},
		ID: 1,
	}
	adminUser.SetRole(dtomocks.AdminRole)
	ctxWithUser := context.WithValue(ctx, app_context.UserKey, adminUser)

	testLogger := new(providersmocks.MockLoggerProvider)
	testUserRepository := new(repositoriesmocks.MockUserRepository)
	testCacheProvider := new(providersmocks.MockCacheProvider)

	// Create invalid query payload (page = 0, which will be validated)
	queryPayload := domain_utils.QueryPayloadBuilder[models.User]{
		Pagination: domain_utils.Pagination{
			Page:     0, // Invalid
			PageSize: 10,
		},
	}

	uc := NewGetAllUserUseCase(
		testLogger,
		testUserRepository,
		testCacheProvider,
	)

	result := uc.Execute(ctxWithUser, locales.EN_US, queryPayload)

	assert.NotNil(result)
	assert.False(result.IsSuccess())
	assert.True(result.HasError())
	assert.Equal(appstatus.InvalidInput, result.StatusCode)
}

func TestGetAllUserUseCase_SetLocale(t *testing.T) {
	assert := assert.New(t)

	testLogger := new(providersmocks.MockLoggerProvider)
	testUserRepository := new(repositoriesmocks.MockUserRepository)
	testCacheProvider := new(providersmocks.MockCacheProvider)

	uc := NewGetAllUserUseCase(
		testLogger,
		testUserRepository,
		testCacheProvider,
	)

	// Test setting locale
	uc.SetLocale(locales.ES_ES)
	assert.Equal(locales.ES_ES, uc.Locale)

	// Test setting empty locale (should not change)
	uc.SetLocale("")
	assert.Equal(locales.ES_ES, uc.Locale)

	// Test setting another locale
	uc.SetLocale(locales.EN_US)
	assert.Equal(locales.EN_US, uc.Locale)
}
