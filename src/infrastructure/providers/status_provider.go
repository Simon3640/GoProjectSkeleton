package providers

import (
	"time"

	contractsProviders "goprojectskeleton/src/application/contracts/providers"
	"goprojectskeleton/src/application/shared/settings"
	"goprojectskeleton/src/domain/models"
)

type ApiStatusProvider struct{}

var _ contractsProviders.IApiStatusProvider = (*ApiStatusProvider)(nil)

func (asp *ApiStatusProvider) Get(date time.Time) models.Status {
	return models.Status{
		AppName: settings.AppSettingsInstance.AppName,
		Version: settings.AppSettingsInstance.AppVersion,
		Status:  "OK",
		Date:    date.Format("2006-01-02 15:04:05"),
	}
}

func NewApiStatusProvider() *ApiStatusProvider {
	return &ApiStatusProvider{}
}
