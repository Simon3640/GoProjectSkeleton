// Package providersmocks contains mock implementations of the contracts/providers/ interfaces
package providersmocks

import (
	"time"

	contractsProviders "goprojectskeleton/src/application/contracts/providers"
	"goprojectskeleton/src/domain/models"

	"github.com/stretchr/testify/mock"
)

type MockStatusProvider struct {
	mock.Mock
}

var _ contractsProviders.IApiStatusProvider = (*MockStatusProvider)(nil)

func (msp *MockStatusProvider) Get(date time.Time) models.Status {
	args := msp.Called(date)
	return args.Get(0).(models.Status)
}
