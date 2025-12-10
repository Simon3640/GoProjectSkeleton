package aws

import (
	"fmt"
	"strings"

	contractsProviders "github.com/simon3640/goprojectskeleton/src/application/contracts/providers"
	application_errors "github.com/simon3640/goprojectskeleton/src/application/shared/errors"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales/messages"
	email_service "github.com/simon3640/goprojectskeleton/src/application/shared/services/emails"
	email_models "github.com/simon3640/goprojectskeleton/src/application/shared/services/emails/models"
	settings "github.com/simon3640/goprojectskeleton/src/application/shared/settings"
	"github.com/simon3640/goprojectskeleton/src/application/shared/status"
	"github.com/simon3640/goprojectskeleton/src/infrastructure/config"
	database "github.com/simon3640/goprojectskeleton/src/infrastructure/database/goprojectskeleton"
	"github.com/simon3640/goprojectskeleton/src/infrastructure/providers"
)

var initialized bool

func InitializeInfrastructure() *application_errors.ApplicationError {
	config, err := config.NewConfig(NewSecretsManagerConfigLoader())
	if err != nil {
		return err
	}
	if err := settings.AppSettingsInstance.Initialize(config.ToMap()); err != nil {
		return err
	}
	providers.Logger.Setup(
		settings.AppSettingsInstance.EnableLog,
		settings.AppSettingsInstance.DebugLog,
	)
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

	// Initialize JWT Provider
	providers.JWTProviderInstance.Setup(
		settings.AppSettingsInstance.JWTSecretKey,
		settings.AppSettingsInstance.JWTIssuer,
		settings.AppSettingsInstance.JWTAudience,
		settings.AppSettingsInstance.JWTAccessTTL,
		settings.AppSettingsInstance.JWTRefreshTTL,
		settings.AppSettingsInstance.JWTClockSkew,
	)

	// Initialize Email Provider
	providers.EmailProviderInstance.Setup(
		settings.AppSettingsInstance.MailHost,
		settings.AppSettingsInstance.MailPort,
		settings.AppSettingsInstance.MailFrom,
		settings.AppSettingsInstance.MailPassword,
	)

	// Initialize Cache Provider
	providers.CacheProviderInstance = providers.NewRedisCacheProviderTLS(
		settings.AppSettingsInstance.RedisHost,
		settings.AppSettingsInstance.RedisPassword,
		settings.AppSettingsInstance.RedisDB,
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
	return nil
}
