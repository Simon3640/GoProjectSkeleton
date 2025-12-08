package contractsproviders

import (
	"time"

	"github.com/simon3640/goprojectskeleton/src/domain/models"
)

type IApiStatusProvider interface {
	Get(date time.Time) models.Status
}
