package authusecases

import (
	"context"
	"sync"
	"testing"
	"time"

	dtos "github.com/simon3640/goprojectskeleton/src/application/modules/auth/dtos"
	authmocks "github.com/simon3640/goprojectskeleton/src/application/modules/auth/mocks"
	app_context "github.com/simon3640/goprojectskeleton/src/application/shared/context"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales"
	dtomocks "github.com/simon3640/goprojectskeleton/src/application/shared/mocks/dtos"
	providersmocks "github.com/simon3640/goprojectskeleton/src/application/shared/mocks/providers"
	services "github.com/simon3640/goprojectskeleton/src/application/shared/services"
	"github.com/simon3640/goprojectskeleton/src/application/shared/settings"
	"github.com/simon3640/goprojectskeleton/src/application/shared/status"
	usecase "github.com/simon3640/goprojectskeleton/src/application/shared/use_case"
	"github.com/simon3640/goprojectskeleton/src/application/shared/workers"
	"github.com/simon3640/goprojectskeleton/src/domain/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuthenticationUseCase(t *testing.T) {
	assert := assert.New(t)
	ctx := &app_context.AppContext{Context: context.Background()}

	testLogger := new(providersmocks.MockLoggerProvider)
	testJWTProvider := new(authmocks.MockJWTProvider)
	testHashProvider := new(providersmocks.MockHashProvider)
	testPasswordRepository := new(authmocks.MockPasswordRepository)
	testUserRepository := new(authmocks.MockUserRepository)
	testOTPRepository := new(authmocks.MockOneTimePasswordRepository)

	uc := NewAuthenticateUseCase(testLogger, testPasswordRepository, testUserRepository, testOTPRepository, testHashProvider, testJWTProvider, nil)

	// Valid User Authentication
	userCredentials := dtos.UserCredentials{
		Email:    "user@example.com",
		Password: "plainPassword",
	}
	passwordExpiresAt := time.Now().Add(30 * 24 * time.Hour)
	passwordBase := models.PasswordBase{
		UserID:    uint(1),
		ExpiresAt: &passwordExpiresAt,
		IsActive:  true,
		Hash:      "hashedPassword123",
	}
	testPasswordRepository.On("GetActivePassword", "user@example.com").Return(&models.Password{
		PasswordBase: passwordBase,
		ID:           uint(1),
	}, nil)
	testHashProvider.On("VerifyPassword", passwordBase.Hash, userCredentials.Password).Return(true, nil)
	testJWTProvider.On("GenerateAccessToken", ctx, "1", mock.Anything).Return("accessToken", time.Now().Add(1*time.Hour), nil)
	testJWTProvider.On("GenerateRefreshToken", ctx, "1").Return("refreshToken", time.Now().Add(24*time.Hour), nil)
	testUserRepository.On("GetUserWithRole", uint(1)).Return(&dtomocks.UserWithRole, nil)
	result := uc.Execute(ctx, locales.EN_US, userCredentials)
	assert.NotNil(result)
	assert.True(result.IsSuccess())
	assert.Equal("accessToken", result.Data.AccessToken)
	assert.Equal("refreshToken", result.Data.RefreshToken)
}

func TestAuthenticationUseCase_OTPLoginEnabled(t *testing.T) {
	assert := assert.New(t)
	ctx := &app_context.AppContext{Context: context.Background()}

	// Setup background service factory for testing
	workers.ResetBackgroundExecutorSingleton()
	services.ResetBackgroundServiceFactory()

	// Initialize background executor and factory
	executorCtx := context.Background()
	workers.InitializeBackgroundExecutor(executorCtx, 2, 10)
	defer workers.ResetBackgroundExecutorSingleton()

	services.InitializeBackgroundServiceFactory()
	defer services.ResetBackgroundServiceFactory()

	testLogger := new(providersmocks.MockLoggerProvider)
	testJWTProvider := new(authmocks.MockJWTProvider)
	testHashProvider := new(providersmocks.MockHashProvider)
	testPasswordRepository := new(authmocks.MockPasswordRepository)
	testUserRepository := new(authmocks.MockUserRepository)
	testOTPRepository := new(authmocks.MockOneTimePasswordRepository)

	uc := NewAuthenticateUseCase(testLogger, testPasswordRepository, testUserRepository, testOTPRepository, testHashProvider, testJWTProvider, nil)

	// User with OTP login enabled
	userCredentials := dtos.UserCredentials{
		Email:    "user@example.com",
		Password: "plainPassword",
	}

	passwordExpiresAt := time.Now().Add(30 * 24 * time.Hour)
	passwordBase := models.PasswordBase{
		UserID:    uint(1),
		ExpiresAt: &passwordExpiresAt,
		IsActive:  true,
		Hash:      "hashedPassword123",
	}

	// Create user with OTP login enabled
	userWithOTP := dtomocks.UserWithRole
	userWithOTP.OTPLogin = true

	testPasswordRepository.On("GetActivePassword", "user@example.com").Return(&models.Password{
		PasswordBase: passwordBase,
		ID:           uint(1),
	}, nil)
	testHashProvider.On("VerifyPassword", passwordBase.Hash, userCredentials.Password).Return(true, nil)
	testUserRepository.On("GetUserWithRole", uint(1)).Return(&userWithOTP, nil)

	// Track execution time to verify it doesn't block
	startTime := time.Now()
	result := uc.Execute(ctx, locales.EN_US, userCredentials)
	executionTime := time.Since(startTime)

	// Verify result is returned immediately (should be very fast, < 50ms)
	// This is the key test: the execution should not wait for the background service
	assert.True(executionTime < 50*time.Millisecond, "Execution should not block, took %v", executionTime)
	assert.NotNil(result)
	assert.True(result.IsSuccess())
	assert.Nil(result.Data, "Result should not contain token when OTP login is enabled")
	assert.NotEmpty(result.Details, "Result should have details indicating OTP login is enabled")

	// Wait a bit for background service to execute (even if it fails, it shouldn't affect the test)
	// The background service may fail due to uninitialized email service, but that's OK for this test
	time.Sleep(200 * time.Millisecond)

	// Verify expectations for the synchronous part
	testPasswordRepository.AssertExpectations(t)
	testHashProvider.AssertExpectations(t)
	testUserRepository.AssertExpectations(t)
	// Note: We don't assert OTP repository expectations because the background service
	// may fail due to uninitialized email service, and that's acceptable for this test
}

func TestAuthenticationUseCase_OTPLoginEnabled_NonBlocking(t *testing.T) {
	assert := assert.New(t)
	ctx := &app_context.AppContext{Context: context.Background()}

	// Setup background service factory for testing
	workers.ResetBackgroundExecutorSingleton()
	services.ResetBackgroundServiceFactory()

	executorCtx := context.Background()
	workers.InitializeBackgroundExecutor(executorCtx, 2, 10)
	defer workers.ResetBackgroundExecutorSingleton()

	services.InitializeBackgroundServiceFactory()
	defer services.ResetBackgroundServiceFactory()

	testLogger := new(providersmocks.MockLoggerProvider)
	testJWTProvider := new(authmocks.MockJWTProvider)
	testHashProvider := new(providersmocks.MockHashProvider)
	testPasswordRepository := new(authmocks.MockPasswordRepository)
	testUserRepository := new(authmocks.MockUserRepository)
	testOTPRepository := new(authmocks.MockOneTimePasswordRepository)

	uc := NewAuthenticateUseCase(testLogger, testPasswordRepository, testUserRepository, testOTPRepository, testHashProvider, testJWTProvider, nil)

	userCredentials := dtos.UserCredentials{
		Email:    "user@example.com",
		Password: "plainPassword",
	}

	passwordExpiresAt := time.Now().Add(30 * 24 * time.Hour)
	passwordBase := models.PasswordBase{
		UserID:    uint(1),
		ExpiresAt: &passwordExpiresAt,
		IsActive:  true,
		Hash:      "hashedPassword123",
	}

	userWithOTP := dtomocks.UserWithRole
	userWithOTP.OTPLogin = true

	testPasswordRepository.On("GetActivePassword", "user@example.com").Return(&models.Password{
		PasswordBase: passwordBase,
		ID:           uint(1),
	}, nil).Maybe()
	testHashProvider.On("VerifyPassword", passwordBase.Hash, userCredentials.Password).Return(true, nil).Maybe()
	testUserRepository.On("GetUserWithRole", uint(1)).Return(&userWithOTP, nil).Maybe()

	// Mock OTP generation for background service (will be called asynchronously)
	testHashProvider.On("GenerateOTP").Return("123456", []byte("hashedOTP"), nil).Maybe()
	testOTPRepository.On("Create", mock.Anything).Return(&models.OneTimePassword{}, nil).Maybe()

	// Execute multiple times to verify non-blocking behavior
	var wg sync.WaitGroup
	numExecutions := 5
	results := make([]*usecase.UseCaseResult[dtos.Token], numExecutions)
	executionTimes := make([]time.Duration, numExecutions)

	for i := 0; i < numExecutions; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			startTime := time.Now()
			results[index] = uc.Execute(ctx, locales.EN_US, userCredentials)
			executionTimes[index] = time.Since(startTime)
		}(i)
	}

	wg.Wait()

	// All executions should complete quickly (non-blocking)
	for i, execTime := range executionTimes {
		assert.True(execTime < 50*time.Millisecond, "Execution %d should not block, took %v", i, execTime)
		assert.NotNil(results[i], "Result %d should not be nil", i)
		assert.True(results[i].IsSuccess(), "Result %d should be successful", i)
	}

	// Wait for background services to complete
	time.Sleep(200 * time.Millisecond)
}

func TestAuthenticationUseCase_InvalidCredentials(t *testing.T) {
	assert := assert.New(t)
	ctx := &app_context.AppContext{Context: context.Background()}

	testLogger := new(providersmocks.MockLoggerProvider)
	testJWTProvider := new(authmocks.MockJWTProvider)
	testHashProvider := new(providersmocks.MockHashProvider)
	testPasswordRepository := new(authmocks.MockPasswordRepository)
	testUserRepository := new(authmocks.MockUserRepository)
	testOTPRepository := new(authmocks.MockOneTimePasswordRepository)

	uc := NewAuthenticateUseCase(testLogger, testPasswordRepository, testUserRepository, testOTPRepository, testHashProvider, testJWTProvider, nil)

	// Invalid User Authentication
	userCredentials := dtos.UserCredentials{
		Email:    "user@example.com",
		Password: "wrongPassword",
	}
	passwordExpiresAt := time.Now().Add(30 * 24 * time.Hour)
	passwordBase := models.PasswordBase{
		UserID:    uint(1),
		ExpiresAt: &passwordExpiresAt,
		IsActive:  true,
		Hash:      "hashedPassword123",
	}
	testPasswordRepository.On("GetActivePassword", "user@example.com").Return(&models.Password{
		PasswordBase: passwordBase,
		ID:           uint(1),
	}, nil)
	testHashProvider.On("VerifyPassword", passwordBase.Hash, userCredentials.Password).Return(false, nil)
	testUserRepository.On("GetUserWithRole", uint(1)).Return(&dtomocks.UserWithRole, nil)
	result := uc.Execute(ctx, locales.EN_US, userCredentials)
	assert.NotNil(result)
	assert.False(result.IsSuccess())
	assert.Equal(status.NotFound, result.StatusCode)
	assert.True(result.HasError())
	testPasswordRepository.AssertExpectations(t)
	testHashProvider.AssertExpectations(t)
	testJWTProvider.AssertExpectations(t)
}

func TestAuthenticationUseCase_RateLimitExceeded(t *testing.T) {
	assert := assert.New(t)
	ctx := &app_context.AppContext{Context: context.Background()}

	// Guardar valores originales
	originalMaxAttempts := settings.AppSettingsInstance.LoginMaxAttempts
	defer func() {
		settings.AppSettingsInstance.LoginMaxAttempts = originalMaxAttempts
	}()

	// Configurar max attempts para el test
	settings.AppSettingsInstance.LoginMaxAttempts = 5

	testLogger := new(providersmocks.MockLoggerProvider)
	testJWTProvider := new(authmocks.MockJWTProvider)
	testHashProvider := new(providersmocks.MockHashProvider)
	testPasswordRepository := new(authmocks.MockPasswordRepository)
	testUserRepository := new(authmocks.MockUserRepository)
	testOTPRepository := new(authmocks.MockOneTimePasswordRepository)
	cacheProvider := new(providersmocks.MockCacheProvider)

	uc := NewAuthenticateUseCase(testLogger, testPasswordRepository, testUserRepository, testOTPRepository, testHashProvider, testJWTProvider, cacheProvider)

	// Rate Limit Exceeded - usuario ha intentado 5 veces (igual al límite)
	userCredentials := dtos.UserCredentials{
		Email:    "user@example.com",
		Password: "plainPassword",
	}

	attempts := 5
	cacheProvider.On("GetInt64", "login_attempts:user@example.com").Return(int64(attempts), nil)

	result := uc.Execute(ctx, locales.EN_US, userCredentials)
	assert.NotNil(result)
	assert.False(result.IsSuccess())
	assert.Equal(status.TooManyRequests, result.StatusCode)
	assert.True(result.HasError())
	cacheProvider.AssertExpectations(t)
}

func TestAuthenticationUseCase_RateLimitNotExceeded(t *testing.T) {
	assert := assert.New(t)
	ctx := &app_context.AppContext{Context: context.Background()}

	// Guardar valores originales
	originalMaxAttempts := settings.AppSettingsInstance.LoginMaxAttempts
	defer func() {
		settings.AppSettingsInstance.LoginMaxAttempts = originalMaxAttempts
	}()

	// Configurar max attempts para el test
	settings.AppSettingsInstance.LoginMaxAttempts = 5

	testLogger := new(providersmocks.MockLoggerProvider)
	testJWTProvider := new(authmocks.MockJWTProvider)
	testHashProvider := new(providersmocks.MockHashProvider)
	testPasswordRepository := new(authmocks.MockPasswordRepository)
	testUserRepository := new(authmocks.MockUserRepository)
	testOTPRepository := new(authmocks.MockOneTimePasswordRepository)
	cacheProvider := new(providersmocks.MockCacheProvider)

	uc := NewAuthenticateUseCase(testLogger, testPasswordRepository, testUserRepository, testOTPRepository, testHashProvider, testJWTProvider, cacheProvider)

	// Rate Limit Not Exceeded - usuario ha intentado 3 veces (menos que el límite)
	userCredentials := dtos.UserCredentials{
		Email:    "user@example.com",
		Password: "plainPassword",
	}

	attempts := 3
	cacheProvider.On("GetInt64", "login_attempts:user@example.com").Return(int64(attempts), nil)

	passwordExpiresAt := time.Now().Add(30 * 24 * time.Hour)
	passwordBase := models.PasswordBase{
		UserID:    uint(1),
		ExpiresAt: &passwordExpiresAt,
		IsActive:  true,
		Hash:      "hashedPassword123",
	}
	testPasswordRepository.On("GetActivePassword", "user@example.com").Return(&models.Password{
		PasswordBase: passwordBase,
		ID:           uint(1),
	}, nil)
	testHashProvider.On("VerifyPassword", passwordBase.Hash, userCredentials.Password).Return(true, nil)
	testJWTProvider.On("GenerateAccessToken", ctx, "1", mock.Anything).Return("accessToken", time.Now().Add(1*time.Hour), nil)
	testJWTProvider.On("GenerateRefreshToken", ctx, "1").Return("refreshToken", time.Now().Add(24*time.Hour), nil)
	testUserRepository.On("GetUserWithRole", uint(1)).Return(&dtomocks.UserWithRole, nil)
	cacheProvider.On("Delete", "login_attempts:user@example.com").Return(nil)

	result := uc.Execute(ctx, locales.EN_US, userCredentials)
	assert.NotNil(result)
	assert.True(result.IsSuccess())
	assert.Equal("accessToken", result.Data.AccessToken)
	cacheProvider.AssertExpectations(t)
	testPasswordRepository.AssertExpectations(t)
	testHashProvider.AssertExpectations(t)
	testJWTProvider.AssertExpectations(t)
}

func TestAuthenticationUseCase_IncrementFailedAttempts(t *testing.T) {
	assert := assert.New(t)
	ctx := &app_context.AppContext{Context: context.Background()}

	// Guardar valores originales
	originalMaxAttempts := settings.AppSettingsInstance.LoginMaxAttempts
	originalWindowMinutes := settings.AppSettingsInstance.LoginAttemptsWindowMinutes
	defer func() {
		settings.AppSettingsInstance.LoginMaxAttempts = originalMaxAttempts
		settings.AppSettingsInstance.LoginAttemptsWindowMinutes = originalWindowMinutes
	}()

	// Configurar para el test
	settings.AppSettingsInstance.LoginMaxAttempts = 5
	settings.AppSettingsInstance.LoginAttemptsWindowMinutes = 15

	testLogger := new(providersmocks.MockLoggerProvider)
	testJWTProvider := new(authmocks.MockJWTProvider)
	testHashProvider := new(providersmocks.MockHashProvider)
	testPasswordRepository := new(authmocks.MockPasswordRepository)
	testUserRepository := new(authmocks.MockUserRepository)
	testOTPRepository := new(authmocks.MockOneTimePasswordRepository)
	cacheProvider := new(providersmocks.MockCacheProvider)

	uc := NewAuthenticateUseCase(testLogger, testPasswordRepository, testUserRepository, testOTPRepository, testHashProvider, testJWTProvider, cacheProvider)

	// Invalid credentials - debe incrementar el contador
	userCredentials := dtos.UserCredentials{
		Email:    "user@example.com",
		Password: "wrongPassword",
	}

	// Primero verifica rate limit (no excedido)
	existingAttempts := 2
	cacheProvider.On("GetInt64", "login_attempts:user@example.com").Return(int64(existingAttempts), nil)

	passwordExpiresAt := time.Now().Add(30 * 24 * time.Hour)
	passwordBase := models.PasswordBase{
		UserID:    uint(1),
		ExpiresAt: &passwordExpiresAt,
		IsActive:  true,
		Hash:      "hashedPassword123",
	}
	testPasswordRepository.On("GetActivePassword", "user@example.com").Return(&models.Password{
		PasswordBase: passwordBase,
		ID:           uint(1),
	}, nil)
	testHashProvider.On("VerifyPassword", passwordBase.Hash, userCredentials.Password).Return(false, nil)
	testUserRepository.On("GetUserWithRole", uint(1)).Return(&dtomocks.UserWithRole, nil)

	// Debe incrementar el contador (2 -> 3)
	cacheProvider.On("GetInt64", "login_attempts:user@example.com").Return(int64(existingAttempts), nil)
	cacheProvider.On("Increment", "login_attempts:user@example.com", time.Duration(settings.AppSettingsInstance.LoginAttemptsWindowMinutes)*time.Minute).Return(int64(3), nil)

	result := uc.Execute(ctx, locales.EN_US, userCredentials)
	assert.NotNil(result)
	assert.False(result.IsSuccess())
	assert.Equal(status.NotFound, result.StatusCode)
	cacheProvider.AssertExpectations(t)
}

func TestAuthenticationUseCase_RateLimitWithNoCacheProvider(t *testing.T) {
	assert := assert.New(t)
	ctx := &app_context.AppContext{Context: context.Background()}

	testLogger := new(providersmocks.MockLoggerProvider)
	testJWTProvider := new(authmocks.MockJWTProvider)
	testHashProvider := new(providersmocks.MockHashProvider)
	testPasswordRepository := new(authmocks.MockPasswordRepository)
	testUserRepository := new(authmocks.MockUserRepository)
	testOTPRepository := new(authmocks.MockOneTimePasswordRepository)

	// Cache provider es nil - no debe aplicar rate limiting
	uc := NewAuthenticateUseCase(testLogger, testPasswordRepository, testUserRepository, testOTPRepository, testHashProvider, testJWTProvider, nil)

	userCredentials := dtos.UserCredentials{
		Email:    "user@example.com",
		Password: "plainPassword",
	}

	passwordExpiresAt := time.Now().Add(30 * 24 * time.Hour)
	passwordBase := models.PasswordBase{
		UserID:    uint(1),
		ExpiresAt: &passwordExpiresAt,
		IsActive:  true,
		Hash:      "hashedPassword123",
	}
	testPasswordRepository.On("GetActivePassword", "user@example.com").Return(&models.Password{
		PasswordBase: passwordBase,
		ID:           uint(1),
	}, nil)
	testHashProvider.On("VerifyPassword", passwordBase.Hash, userCredentials.Password).Return(true, nil)
	testJWTProvider.On("GenerateAccessToken", ctx, "1", mock.Anything).Return("accessToken", time.Now().Add(1*time.Hour), nil)
	testJWTProvider.On("GenerateRefreshToken", ctx, "1").Return("refreshToken", time.Now().Add(24*time.Hour), nil)
	testUserRepository.On("GetUserWithRole", uint(1)).Return(&dtomocks.UserWithRole, nil)

	result := uc.Execute(ctx, locales.EN_US, userCredentials)
	assert.NotNil(result)
	assert.True(result.IsSuccess())
	testPasswordRepository.AssertExpectations(t)
	testHashProvider.AssertExpectations(t)
	testJWTProvider.AssertExpectations(t)
}

func TestAuthenticationUseCase_IncrementFailedAttemptsWhenNoPreviousAttempts(t *testing.T) {
	assert := assert.New(t)
	ctx := &app_context.AppContext{Context: context.Background()}

	// Guardar valores originales
	originalMaxAttempts := settings.AppSettingsInstance.LoginMaxAttempts
	originalWindowMinutes := settings.AppSettingsInstance.LoginAttemptsWindowMinutes
	defer func() {
		settings.AppSettingsInstance.LoginMaxAttempts = originalMaxAttempts
		settings.AppSettingsInstance.LoginAttemptsWindowMinutes = originalWindowMinutes
	}()

	// Configurar para el test
	settings.AppSettingsInstance.LoginMaxAttempts = 5
	settings.AppSettingsInstance.LoginAttemptsWindowMinutes = 15

	testLogger := new(providersmocks.MockLoggerProvider)
	testJWTProvider := new(authmocks.MockJWTProvider)
	testHashProvider := new(providersmocks.MockHashProvider)
	testPasswordRepository := new(authmocks.MockPasswordRepository)
	testUserRepository := new(authmocks.MockUserRepository)
	testOTPRepository := new(authmocks.MockOneTimePasswordRepository)
	cacheProvider := new(providersmocks.MockCacheProvider)

	uc := NewAuthenticateUseCase(testLogger, testPasswordRepository, testUserRepository, testOTPRepository, testHashProvider, testJWTProvider, cacheProvider)

	// Invalid credentials - debe crear e incrementar el contador desde 0
	userCredentials := dtos.UserCredentials{
		Email:    "user@example.com",
		Password: "wrongPassword",
	}

	// Primero verifica rate limit (no existe contador previo)
	cacheProvider.On("GetInt64", "login_attempts:user@example.com").Return(int64(0), nil)

	passwordExpiresAt := time.Now().Add(30 * 24 * time.Hour)
	passwordBase := models.PasswordBase{
		UserID:    uint(1),
		ExpiresAt: &passwordExpiresAt,
		IsActive:  true,
		Hash:      "hashedPassword123",
	}
	testPasswordRepository.On("GetActivePassword", "user@example.com").Return(&models.Password{
		PasswordBase: passwordBase,
		ID:           uint(1),
	}, nil)
	testHashProvider.On("VerifyPassword", passwordBase.Hash, userCredentials.Password).Return(false, nil)
	testUserRepository.On("GetUserWithRole", uint(1)).Return(&dtomocks.UserWithRole, nil)

	// Debe crear el contador con valor 1 (desde 0)
	cacheProvider.On("GetInt64", "login_attempts:user@example.com").Return(int64(0), nil)
	cacheProvider.On("Increment", "login_attempts:user@example.com", time.Duration(settings.AppSettingsInstance.LoginAttemptsWindowMinutes)*time.Minute).Return(int64(1), nil)

	result := uc.Execute(ctx, locales.EN_US, userCredentials)
	assert.NotNil(result)
	assert.False(result.IsSuccess())
	assert.Equal(status.NotFound, result.StatusCode)
	cacheProvider.AssertExpectations(t)
}

func TestAuthenticationUseCase_RateLimitWithMaxAttemptsZero(t *testing.T) {
	assert := assert.New(t)
	ctx := &app_context.AppContext{Context: context.Background()}

	// Guardar valores originales
	originalMaxAttempts := settings.AppSettingsInstance.LoginMaxAttempts
	defer func() {
		settings.AppSettingsInstance.LoginMaxAttempts = originalMaxAttempts
	}()

	// Configurar max attempts a 0 - no debe aplicar rate limiting
	settings.AppSettingsInstance.LoginMaxAttempts = 0

	testLogger := new(providersmocks.MockLoggerProvider)
	testJWTProvider := new(authmocks.MockJWTProvider)
	testHashProvider := new(providersmocks.MockHashProvider)
	testPasswordRepository := new(authmocks.MockPasswordRepository)
	testUserRepository := new(authmocks.MockUserRepository)
	testOTPRepository := new(authmocks.MockOneTimePasswordRepository)
	cacheProvider := new(providersmocks.MockCacheProvider)

	uc := NewAuthenticateUseCase(testLogger, testPasswordRepository, testUserRepository, testOTPRepository, testHashProvider, testJWTProvider, cacheProvider)

	userCredentials := dtos.UserCredentials{
		Email:    "user@example.com",
		Password: "plainPassword",
	}

	// No debe llamar a Get del cache porque maxAttempts es 0
	passwordExpiresAt := time.Now().Add(30 * 24 * time.Hour)
	passwordBase := models.PasswordBase{
		UserID:    uint(1),
		ExpiresAt: &passwordExpiresAt,
		IsActive:  true,
		Hash:      "hashedPassword123",
	}
	testPasswordRepository.On("GetActivePassword", "user@example.com").Return(&models.Password{
		PasswordBase: passwordBase,
		ID:           uint(1),
	}, nil)
	testHashProvider.On("VerifyPassword", passwordBase.Hash, userCredentials.Password).Return(true, nil)
	testJWTProvider.On("GenerateAccessToken", ctx, "1", mock.Anything).Return("accessToken", time.Now().Add(1*time.Hour), nil)
	testJWTProvider.On("GenerateRefreshToken", ctx, "1").Return("refreshToken", time.Now().Add(24*time.Hour), nil)
	testUserRepository.On("GetUserWithRole", uint(1)).Return(&dtomocks.UserWithRole, nil)
	// Cuando la autenticación es exitosa, se limpia el contador
	cacheProvider.On("Delete", "login_attempts:user@example.com").Return(nil)

	result := uc.Execute(ctx, locales.EN_US, userCredentials)
	assert.NotNil(result)
	assert.True(result.IsSuccess())
	// No debe llamar a Get del cache porque maxAttempts es 0
	cacheProvider.AssertNotCalled(t, "GetInt64", "login_attempts:user@example.com")
	testPasswordRepository.AssertExpectations(t)
	testHashProvider.AssertExpectations(t)
	testJWTProvider.AssertExpectations(t)
	cacheProvider.AssertExpectations(t)
}
