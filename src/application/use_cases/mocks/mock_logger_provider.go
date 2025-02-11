package mocks

import (
	"fmt"

	"gingoskeleton/src/application/contracts"

	"github.com/stretchr/testify/mock"
)

type MockLoggerProvider struct {
	mock.Mock
}

var _ contracts.ILoggerProvider = (*MockLoggerProvider)(nil)

func (mlp *MockLoggerProvider) Error(message string, err error) {
	fmt.Printf("Error: %s, %s", message, err)
}

func (mlp *MockLoggerProvider) Panic(message string, err error) {
	fmt.Printf("Panic: %s, %s", message, err)
}

func (mlp *MockLoggerProvider) ErrorMsg(message string) {
	fmt.Printf(message)
}

func (mlp *MockLoggerProvider) Info(message string) {
	fmt.Printf(message)
}

func (mlp *MockLoggerProvider) Warning(message string) {
	fmt.Printf(message)
}

func (mlp *MockLoggerProvider) Debug(message string, data any) {
	fmt.Printf("%s: %v", message, data)
}
