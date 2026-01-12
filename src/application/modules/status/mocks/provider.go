// Package providersmocks contains mock implementations of the contracts/providers/ interfaces
package statusmocks

import (
	"time"

	statuscontracts "github.com/simon3640/goprojectskeleton/src/application/modules/status/contracts"
	statusmodels "github.com/simon3640/goprojectskeleton/src/domain/status/models"

	"github.com/stretchr/testify/mock"
)

type MockStatusProvider struct {
	mock.Mock
}

var _ statuscontracts.IApiStatusProvider = (*MockStatusProvider)(nil)

// Get returns the status for the given date
func (msp *MockStatusProvider) Get(date time.Time) statusmodels.Status {
	args := msp.Called(date)
	return args.Get(0).(statusmodels.Status)
}
