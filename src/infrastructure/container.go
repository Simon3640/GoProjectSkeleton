package infrastructure

import (
	settings "gormgoskeleton/src/application/shared/settings"
	config "gormgoskeleton/src/infrastructure/config"
	database "gormgoskeleton/src/infrastructure/database/gormgoskeleton"
	providers "gormgoskeleton/src/infrastructure/providers"
)

func Initialize() {
	settings.AppSettingsInstance.Initialize(config.ConfigInstance.ToMap())
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
}
