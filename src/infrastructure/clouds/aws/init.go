package aws

import (
	"fmt"
	"log"
	"strings"
	"sync"

	contractsProviders "github.com/simon3640/goprojectskeleton/src/application/contracts/providers"
	application_errors "github.com/simon3640/goprojectskeleton/src/application/shared/errors"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales/messages"
	email_service "github.com/simon3640/goprojectskeleton/src/application/shared/services/emails"
	email_models "github.com/simon3640/goprojectskeleton/src/application/shared/services/emails/models"
	settings "github.com/simon3640/goprojectskeleton/src/application/shared/settings"
	"github.com/simon3640/goprojectskeleton/src/application/shared/status"
	"github.com/simon3640/goprojectskeleton/src/infrastructure/config"
	database "github.com/simon3640/goprojectskeleton/src/infrastructure/databases/goprojectskeleton"
	"github.com/simon3640/goprojectskeleton/src/infrastructure/providers"
)

var (
	initializedBase     bool
	initializedDatabase bool
	initializedJWT      bool
	initializedCache    bool
	initializedEmail    bool
	initMutex           sync.Mutex
)

// InitializeBase initializes the base infrastructure: Config, Settings, and Logger.
// This should be called first by all Lambda functions.
func InitializeBase() *application_errors.ApplicationError {
	initMutex.Lock()
	defer initMutex.Unlock()

	if initializedBase {
		return nil
	}

	cfg, err := config.NewConfig(NewSecretsManagerConfigLoader())
	if err != nil {
		return err
	}
	if err := settings.AppSettingsInstance.Initialize(cfg.ToMap()); err != nil {
		return err
	}
	providers.Logger.Setup(
		settings.AppSettingsInstance.EnableLog,
		settings.AppSettingsInstance.DebugLog,
	)

	initializedBase = true
	log.Println("Base infrastructure initialized successfully")
	return nil
}

// InitializeDatabase initializes the database connection.
func InitializeDatabase() *application_errors.ApplicationError {
	initMutex.Lock()
	defer initMutex.Unlock()

	if initializedDatabase {
		return nil
	}

	if !initializedBase {
		if err := InitializeBase(); err != nil {
			return err
		}
	}

	if err := database.GoProjectSkeletondb.SetUp(
		settings.AppSettingsInstance.DBHost,
		settings.AppSettingsInstance.DBPort,
		settings.AppSettingsInstance.DBUser,
		settings.AppSettingsInstance.DBPassword,
		settings.AppSettingsInstance.DBName,
		&settings.AppSettingsInstance.DBSSL,
		providers.Logger,
	); err != nil {
		return err
	}

	initializedDatabase = true
	log.Println("Database initialized successfully")
	return nil
}

// InitializeJWT initializes the JWT provider.
func InitializeJWT() *application_errors.ApplicationError {
	initMutex.Lock()
	defer initMutex.Unlock()

	if initializedJWT {
		return nil
	}

	if !initializedBase {
		if err := InitializeBase(); err != nil {
			return err
		}
	}

	providers.JWTProviderInstance.Setup(
		settings.AppSettingsInstance.JWTSecretKey,
		settings.AppSettingsInstance.JWTIssuer,
		settings.AppSettingsInstance.JWTAudience,
		settings.AppSettingsInstance.JWTAccessTTL,
		settings.AppSettingsInstance.JWTRefreshTTL,
		settings.AppSettingsInstance.JWTClockSkew,
	)

	initializedJWT = true
	log.Println("JWT initialized successfully")
	return nil
}

// InitializeCache initializes the cache provider (Redis).
func InitializeCache() *application_errors.ApplicationError {
	initMutex.Lock()
	defer initMutex.Unlock()

	if initializedCache {
		return nil
	}

	if !initializedBase {
		if err := InitializeBase(); err != nil {
			return err
		}
	}

	providers.CacheProviderInstance = providers.NewRedisCacheProviderTLS(
		settings.AppSettingsInstance.RedisHost,
		settings.AppSettingsInstance.RedisPassword,
		settings.AppSettingsInstance.RedisDB,
	)

	initializedCache = true
	log.Println("Cache initialized successfully")
	return nil
}

// InitializeEmail initializes email provider, render providers, and email services.
func InitializeEmail() *application_errors.ApplicationError {
	initMutex.Lock()
	defer initMutex.Unlock()

	if initializedEmail {
		return nil
	}

	if !initializedBase {
		if err := InitializeBase(); err != nil {
			return err
		}
	}

	// Initialize Email Provider
	providers.EmailProviderInstance.Setup(
		settings.AppSettingsInstance.MailHost,
		settings.AppSettingsInstance.MailPort,
		settings.AppSettingsInstance.MailFrom,
		settings.AppSettingsInstance.MailPassword,
	)

	// Initialize Render Providers (S3 or Local)
	var renderNewUser contractsProviders.IRendererProvider[email_models.NewUserEmailData]
	var renderResetPassword contractsProviders.IRendererProvider[email_models.ResetPasswordEmailData]
	var renderOTP contractsProviders.IRendererProvider[email_models.OneTimePasswordEmailData]

	// Check if templates are stored in S3
	templatesPath := settings.AppSettingsInstance.TemplatesPath
	if strings.HasPrefix(templatesPath, "s3://") {
		// Extract bucket from S3 path: s3://bucket/path/
		path := strings.TrimPrefix(templatesPath, "s3://")
		parts := strings.SplitN(path, "/", 2)
		if parts[0] == "" {
			return application_errors.NewApplicationError(
				status.ProviderInitializationError,
				messages.MessageKeysInstance.SOMETHING_WENT_WRONG,
				fmt.Sprintf("Invalid S3 templates path, expected format: (s3://bucket/path/), got: %s", templatesPath),
			)
		}
		bucket := parts[0]

		s3RenderNewUser, s3RenderResetPassword, s3RenderOTP, err := NewS3RenderProviders(bucket)
		if err != nil {
			return application_errors.NewApplicationError(
				status.ProviderInitializationError,
				messages.MessageKeysInstance.SOMETHING_WENT_WRONG,
				fmt.Sprintf("Failed to initialize S3 render providers: %v", err),
			)
		}

		renderNewUser = s3RenderNewUser
		renderResetPassword = s3RenderResetPassword
		renderOTP = s3RenderOTP

		providers.Logger.Info(fmt.Sprintf("Using S3 render providers with bucket: %s", bucket))
	} else {
		return application_errors.NewApplicationError(
			status.ProviderInitializationError,
			messages.MessageKeysInstance.SOMETHING_WENT_WRONG,
			fmt.Sprintf("Not s3 templates path, expected format: (s3://bucket/path/), got: %s", templatesPath),
		)
	}

	// Services
	email_service.RegisterUserEmailServiceInstance.SetUp(
		renderNewUser,
		providers.EmailProviderInstance,
	)

	email_service.ResetPasswordEmailServiceInstance.SetUp(
		renderResetPassword,
		providers.EmailProviderInstance,
	)

	email_service.OneTimePasswordEmailServiceInstance.SetUp(
		renderOTP,
		providers.EmailProviderInstance,
	)

	initializedEmail = true
	log.Println("Email initialized successfully")
	return nil
}

// InitializeInfrastructure initializes all infrastructure components.
// This is kept for backward compatibility but should be replaced with
// specific initialization functions for better tree-shaking.
func InitializeInfrastructure() *application_errors.ApplicationError {
	if err := InitializeBase(); err != nil {
		return err
	}
	if err := InitializeDatabase(); err != nil {
		return err
	}
	if err := InitializeJWT(); err != nil {
		return err
	}
	if err := InitializeCache(); err != nil {
		return err
	}
	if err := InitializeEmail(); err != nil {
		return err
	}
	return nil
}

// InitializeForStatus initializes infrastructure for status handlers (health check).
// Only initializes base (config, settings, logger).
func InitializeForStatus() *application_errors.ApplicationError {
	return InitializeBase()
}

// InitializeForAuthLogin initializes infrastructure for auth login handler.
// Requires: Base, Database, JWT, Cache.
func InitializeForAuthLogin() *application_errors.ApplicationError {
	if err := InitializeBase(); err != nil {
		return err
	}
	if err := InitializeDatabase(); err != nil {
		return err
	}
	if err := InitializeJWT(); err != nil {
		return err
	}
	if err := InitializeCache(); err != nil {
		return err
	}
	if err := InitializeEmail(); err != nil {
		return err
	}
	return nil
}

// InitializeForAuthRefresh initializes infrastructure for auth refresh handler.
// Requires: Base, JWT.
func InitializeForAuthRefresh() *application_errors.ApplicationError {
	if err := InitializeBase(); err != nil {
		return err
	}
	if err := InitializeJWT(); err != nil {
		return err
	}
	return nil
}

// InitializeForAuthLoginOTP initializes infrastructure for auth login OTP handler.
// Requires: Base, Database, JWT.
func InitializeForAuthLoginOTP() *application_errors.ApplicationError {
	if err := InitializeBase(); err != nil {
		return err
	}
	if err := InitializeDatabase(); err != nil {
		return err
	}
	if err := InitializeJWT(); err != nil {
		return err
	}
	return nil
}

// InitializeForAuthPasswordReset initializes infrastructure for auth password reset handler.
// Requires: Base, Database, Email.
func InitializeForAuthPasswordReset() *application_errors.ApplicationError {
	if err := InitializeBase(); err != nil {
		return err
	}
	if err := InitializeDatabase(); err != nil {
		return err
	}
	if err := InitializeEmail(); err != nil {
		return err
	}
	return nil
}

// InitializeForUser initializes infrastructure for basic user handlers.
// Requires: Base, Database.
func InitializeForUser() *application_errors.ApplicationError {
	if err := InitializeBase(); err != nil {
		return err
	}
	if err := InitializeDatabase(); err != nil {
		return err
	}
	return nil
}

// InitializeForUserWithCache initializes infrastructure for user handlers that need cache.
// Requires: Base, Database, Cache.
func InitializeForUserWithCache() *application_errors.ApplicationError {
	if err := InitializeBase(); err != nil {
		return err
	}
	if err := InitializeDatabase(); err != nil {
		return err
	}
	if err := InitializeCache(); err != nil {
		return err
	}
	return nil
}

// InitializeForUserWithEmail initializes infrastructure for user handlers that need email.
// Requires: Base, Database, Email.
func InitializeForUserWithEmail() *application_errors.ApplicationError {
	if err := InitializeBase(); err != nil {
		return err
	}
	if err := InitializeDatabase(); err != nil {
		return err
	}
	if err := InitializeEmail(); err != nil {
		return err
	}
	return nil
}

// InitializeForPassword initializes infrastructure for password handlers.
// Requires: Base, Database.
func InitializeForPassword() *application_errors.ApplicationError {
	if err := InitializeBase(); err != nil {
		return err
	}
	if err := InitializeDatabase(); err != nil {
		return err
	}
	return nil
}

// InitializeForPasswordWithEmail initializes infrastructure for password handlers that need email.
// Requires: Base, Database, Email.
func InitializeForPasswordWithEmail() *application_errors.ApplicationError {
	if err := InitializeBase(); err != nil {
		return err
	}
	if err := InitializeDatabase(); err != nil {
		return err
	}
	if err := InitializeEmail(); err != nil {
		return err
	}
	return nil
}
