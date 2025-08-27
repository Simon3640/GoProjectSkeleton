package mocks

import (
	contracts_providers "gormgoskeleton/src/application/contracts/providers"
	application_errors "gormgoskeleton/src/application/shared/errors"

	"github.com/stretchr/testify/mock"
)

type MockEmailProvider struct {
	mock.Mock
}

var _ contracts_providers.IEmailProvider = (*MockEmailProvider)(nil)

func (m *MockEmailProvider) SendEmail(to string, subject string, body string) *application_errors.ApplicationError {
	args := m.Called(to, subject, body)
	errorArg := args.Get(0)
	if errorArg != nil {
		return errorArg.(*application_errors.ApplicationError)
	}
	return nil
}
