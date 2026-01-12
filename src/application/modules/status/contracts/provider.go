// Package statuscontracts contains the contracts for the status module.
package statuscontracts

import (
	"time"

	statusmodels "github.com/simon3640/goprojectskeleton/src/domain/status/models"
)

// IApiStatusProvider is the contract for the API status provider.
type IApiStatusProvider interface {
	Get(date time.Time) statusmodels.Status
}
