package infrastructure

import (
	email_service "gormgoskeleton/src/application/shared/services/emails"
	settings "gormgoskeleton/src/application/shared/settings"
	config "gormgoskeleton/src/infrastructure/config"
	database "gormgoskeleton/src/infrastructure/database/gormgoskeleton"
	providers "gormgoskeleton/src/infrastructure/providers"
)

func Initialize() {
	if err := settings.AppSettingsInstance.Initialize(config.ConfigInstance.ToMap()); err != nil {
		providers.Logger.Error("Failed to initialize app settings", err)
		panic("Failed to initialize app settings: " + err.Error())
	}
	providers.Logger.Setup(
		settings.AppSettingsInstance.EnableLog,
		settings.AppSettingsInstance.DebugLog,
	)
	database.Gormgoskeletondb.SetUp(
		settings.AppSettingsInstance.DBHost,
		settings.AppSettingsInstance.DBPort,
		settings.AppSettingsInstance.DBUser,
		settings.AppSettingsInstance.DBPassword,
		settings.AppSettingsInstance.DBName,
		&settings.AppSettingsInstance.DBSSL,
		providers.Logger,
	)

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
}
