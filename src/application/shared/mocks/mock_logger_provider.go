package mocks

import (
	"fmt"
	contractsProviders "gormgoskeleton/src/application/contracts/providers"

	"github.com/stretchr/testify/mock"
)

type MockLoggerProvider struct {
	mock.Mock
}

var _ contractsProviders.ILoggerProvider = (*MockLoggerProvider)(nil)

func (mlp *MockLoggerProvider) Error(message string, err error) {
	fmt.Printf("Error: %s, %s", message, err)
}

func (mlp *MockLoggerProvider) Panic(message string, err error) {
	fmt.Printf("Panic: %s, %s", message, err)
}

func (mlp *MockLoggerProvider) ErrorMsg(message string) {
	fmt.Print(message)
}

func (mlp *MockLoggerProvider) Info(message string) {
	fmt.Print(message)
}

func (mlp *MockLoggerProvider) Warning(message string) {
	fmt.Print(message)
}

func (mlp *MockLoggerProvider) Debug(message string, data any) {
	fmt.Printf("%s: %v", message, data)
}
