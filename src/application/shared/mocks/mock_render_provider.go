package mocks

import (
	contractsProviders "gormgoskeleton/src/application/contracts/providers"
	application_errors "gormgoskeleton/src/application/shared/errors"

	"github.com/stretchr/testify/mock"
)

type MockRenderProvider[D any] struct {
	mock.Mock
}

var _ contractsProviders.IRendererProvider[any] = (*MockRenderProvider[any])(nil)

func (m *MockRenderProvider[D]) Render(template string, data D) (string, *application_errors.ApplicationError) {
	args := m.Called(template, data)
	errorArg := args.Get(1)
	if errorArg != nil {
		return args.String(0), errorArg.(*application_errors.ApplicationError)
	}
	return args.String(0), nil
}
