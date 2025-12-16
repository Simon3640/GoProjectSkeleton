package api

import (
	"os"

	"github.com/simon3640/goprojectskeleton/src/application/shared/observability"
	"github.com/simon3640/goprojectskeleton/src/application/shared/settings"
	"github.com/simon3640/goprojectskeleton/src/infrastructure"
	"github.com/simon3640/goprojectskeleton/src/infrastructure/providers"

	otel "github.com/simon3640/goprojectskeleton/src/infrastructure/otel"
)

func Initialize() {
	providers.Logger.Info("Initializing infraestructure...")
	if err := infrastructure.Initialize(); err != nil {
		providers.Logger.Error("Error initializing infrastructure", err.ToError())
		os.Exit(1)
	}

	if settings.AppSettingsInstance.ObservabilityEnabled && settings.AppSettingsInstance.ObservabilityBackend == "opentelemetry" {
		providers.Logger.Info("Initializing OpenTelemetry...")
		otel.InitializeOtelSDK(otel.OtelConfig{
			ServiceName:    settings.AppSettingsInstance.AppName,
			ServiceVersion: settings.AppSettingsInstance.AppVersion,
			OTLPEndpoint:   settings.AppSettingsInstance.OTLPEndpoint,
			Environment:    settings.AppSettingsInstance.AppEnv,
		})

		observability.ObservabilityComponentsInstance = otel.NewOtelObservabilityComponents(otel.OtelConfig{
			ServiceName:    settings.AppSettingsInstance.AppName,
			ServiceVersion: settings.AppSettingsInstance.AppVersion,
			OTLPEndpoint:   settings.AppSettingsInstance.OTLPEndpoint,
			Environment:    settings.AppSettingsInstance.AppEnv,
		})
	}
}
