// Package providersmocks contains mock implementations of the contracts/providers/ interfaces
package statusmocks

import (
	"time"

	statuscontracts "github.com/simon3640/goprojectskeleton/src/application/modules/status/contracts"
	"github.com/simon3640/goprojectskeleton/src/domain/models"

	"github.com/stretchr/testify/mock"
)

type MockStatusProvider struct {
	mock.Mock
}

var _ statuscontracts.IApiStatusProvider = (*MockStatusProvider)(nil)

func (msp *MockStatusProvider) Get(date time.Time) models.Status {
	args := msp.Called(date)
	return args.Get(0).(models.Status)
}
