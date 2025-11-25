package mocks

import (
	"time"

	contractsProviders "gormgoskeleton/src/application/contracts/providers"
	"gormgoskeleton/src/domain/models"

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
