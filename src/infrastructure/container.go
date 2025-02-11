package infrastructure

import (
	settings "gormgoskeleton/src/application/shared/settings"
	config "gormgoskeleton/src/infrastructure/config"
	providers "gormgoskeleton/src/infrastructure/providers"
)

func Initialize() {
	settings.AppSettingsInstance.Initialize(config.ConfigInstance.ToMap())
	providers.Logger.Setup(
		settings.AppSettingsInstance.EnableLog,
		settings.AppSettingsInstance.DebugLog,
	)
}
