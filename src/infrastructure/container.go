package infrastructure

import (
	application_errors "github.com/simon3640/goprojectskeleton/src/application/shared/errors"
	email_service "github.com/simon3640/goprojectskeleton/src/application/shared/services/emails"
	settings "github.com/simon3640/goprojectskeleton/src/application/shared/settings"
	config "github.com/simon3640/goprojectskeleton/src/infrastructure/config"
	database "github.com/simon3640/goprojectskeleton/src/infrastructure/database/goprojectskeleton"
	providers "github.com/simon3640/goprojectskeleton/src/infrastructure/providers"
)

// Initialize initializes the infrastructure and returns an application error if it fails
func Initialize() *application_errors.ApplicationError {
	if err := settings.AppSettingsInstance.Initialize(config.NewConfig(nil).ToMap()); err != nil {
		providers.Logger.Error("Failed to initialize app settings", err)
		panic("Failed to initialize app settings: " + err.Error())
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
	providers.CacheProviderInstance = providers.NewRedisCacheProvider(
		settings.AppSettingsInstance.RedisHost,
		settings.AppSettingsInstance.RedisPassword,
		settings.AppSettingsInstance.RedisDB,
	)

	// Services
	email_service.RegisterUserEmailServiceInstance.SetUp(
		providers.RenderNewUserEmailInstance,
		providers.EmailProviderInstance,
	)

	email_service.ResetPasswordEmailServiceInstance.SetUp(
		providers.RenderResetPasswordEmailInstance,
		providers.EmailProviderInstance,
	)

	email_service.OneTimePasswordEmailServiceInstance.SetUp(
		providers.RenderOTPEmailInstance,
		providers.EmailProviderInstance,
	)
	return nil
}
