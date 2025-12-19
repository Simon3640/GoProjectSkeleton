// Package statuscontracts contains the contracts for the status module.
package statuscontracts

import (
	"time"

	"github.com/simon3640/goprojectskeleton/src/domain/models"
)

// IApiStatusProvider is the contract for the API status provider.
type IApiStatusProvider interface {
	Get(date time.Time) models.Status
}
