package providers

import (
	"time"

	contractstatus "github.com/simon3640/goprojectskeleton/src/application/modules/status/contracts"
	"github.com/simon3640/goprojectskeleton/src/application/shared/settings"
	statusmodels "github.com/simon3640/goprojectskeleton/src/domain/status/models"
)

type ApiStatusProvider struct{}

var _ contractstatus.IApiStatusProvider = (*ApiStatusProvider)(nil)

// Get returns the API status for the given date
func (asp *ApiStatusProvider) Get(date time.Time) statusmodels.Status {
	return statusmodels.Status{
		AppName: settings.AppSettingsInstance.AppName,
		Version: settings.AppSettingsInstance.AppVersion,
		Status:  "OK",
		Date:    date.Format("2006-01-02 15:04:05"),
	}
}

func NewApiStatusProvider() *ApiStatusProvider {
	return &ApiStatusProvider{}
}
