package contracts_providers

import (
	"gormgoskeleton/src/domain/models"
	"time"
)

type IApiStatusProvider interface {
	Get(date time.Time) models.Status
}
