package authusecases

import (
	"context"
	"testing"
	"time"

	dtos "github.com/simon3640/goprojectskeleton/src/application/modules/auth/dtos"
	authmocks "github.com/simon3640/goprojectskeleton/src/application/modules/auth/mocks"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales"
	dtomocks "github.com/simon3640/goprojectskeleton/src/application/shared/mocks/dtos"
	providersmocks "github.com/simon3640/goprojectskeleton/src/application/shared/mocks/providers"
	repositoriesmocks "github.com/simon3640/goprojectskeleton/src/application/shared/mocks/repositories"
	"github.com/simon3640/goprojectskeleton/src/application/shared/settings"
	"github.com/simon3640/goprojectskeleton/src/application/shared/status"
	"github.com/simon3640/goprojectskeleton/src/domain/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuthenticationUseCase(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()

	testLogger := new(providersmocks.MockLoggerProvider)
	testJWTProvider := new(authmocks.MockJWTProvider)
	testHashProvider := new(providersmocks.MockHashProvider)
	testPasswordRepository := new(authmocks.MockPasswordRepository)
	testUserRepository := new(authmocks.MockUserRepository)
	testOTPRepository := new(repositoriesmocks.MockOneTimePasswordRepository)

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

//TODO: Add test for OTP when enabled

func TestAuthenticationUseCase_InvalidCredentials(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()

	testLogger := new(providersmocks.MockLoggerProvider)
	testJWTProvider := new(authmocks.MockJWTProvider)
	testHashProvider := new(providersmocks.MockHashProvider)
	testPasswordRepository := new(repositoriesmocks.MockPasswordRepository)
	testUserRepository := new(repositoriesmocks.MockUserRepository)
	testOTPRepository := new(repositoriesmocks.MockOneTimePasswordRepository)

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
	ctx := context.Background()

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
	testOTPRepository := new(repositoriesmocks.MockOneTimePasswordRepository)
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
	ctx := context.Background()

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
	testOTPRepository := new(repositoriesmocks.MockOneTimePasswordRepository)
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
	ctx := context.Background()

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
	testOTPRepository := new(repositoriesmocks.MockOneTimePasswordRepository)
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
	ctx := context.Background()

	testLogger := new(providersmocks.MockLoggerProvider)
	testJWTProvider := new(authmocks.MockJWTProvider)
	testHashProvider := new(providersmocks.MockHashProvider)
	testPasswordRepository := new(authmocks.MockPasswordRepository)
	testUserRepository := new(authmocks.MockUserRepository)
	testOTPRepository := new(repositoriesmocks.MockOneTimePasswordRepository)

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
	ctx := context.Background()

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
	testOTPRepository := new(repositoriesmocks.MockOneTimePasswordRepository)
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
	ctx := context.Background()

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
	testOTPRepository := new(repositoriesmocks.MockOneTimePasswordRepository)
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
