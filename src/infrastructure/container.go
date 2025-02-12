package infrastructure

import (
	settings "gingoskeleton/src/application/shared/settings"
	config "gingoskeleton/src/infrastructure/config"
	providers "gingoskeleton/src/infrastructure/providers"
)

func Initialize() {
	settings.AppSettingsInstance.Initialize(config.ConfigInstance.ToMap())
	providers.Logger.Setup(
		settings.AppSettingsInstance.EnableLog,
		settings.AppSettingsInstance.DebugLog,
	)
}
