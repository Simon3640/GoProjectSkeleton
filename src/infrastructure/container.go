package infrastructure

import (
	"context"

	application_errors "github.com/simon3640/goprojectskeleton/src/application/shared/errors"
	"github.com/simon3640/goprojectskeleton/src/application/shared/observability/noop"
	services "github.com/simon3640/goprojectskeleton/src/application/shared/services"
	email_service "github.com/simon3640/goprojectskeleton/src/application/shared/services/emails"
	settings "github.com/simon3640/goprojectskeleton/src/application/shared/settings"
	"github.com/simon3640/goprojectskeleton/src/application/shared/workers"
	config "github.com/simon3640/goprojectskeleton/src/infrastructure/config"
	database "github.com/simon3640/goprojectskeleton/src/infrastructure/databases/goprojectskeleton"
	initdb "github.com/simon3640/goprojectskeleton/src/infrastructure/databases/goprojectskeleton/init_db"
	providers "github.com/simon3640/goprojectskeleton/src/infrastructure/providers"
)

// Initialize initializes the infrastructure and returns an application error if it fails
func Initialize() *application_errors.ApplicationError {
	config, err := config.NewConfig(nil)
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
	noop.Logger.Setup(
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

	// Migrate database
	if err := initdb.InitMigrate(database.GoProjectSkeletondb.DB, providers.Logger); err != nil {
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

	// Initialize Background Executor
	ctx := context.Background()
	workers.InitializeBackgroundExecutor(
		ctx,
		settings.AppSettingsInstance.BackgroundWorkers,
		settings.AppSettingsInstance.BackgroundQueueSize,
	)

	services.InitializeBackgroundServiceFactory()
	return nil
}
