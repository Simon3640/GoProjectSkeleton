package providersmocks

import (
	"fmt"

	contractsProviders "github.com/simon3640/goprojectskeleton/src/application/contracts/providers"
	app_context "github.com/simon3640/goprojectskeleton/src/application/shared/context"

	"github.com/stretchr/testify/mock"
)

// MockLoggerProvider is a mock of LoggerProvider for testing
type MockLoggerProvider struct {
	mock.Mock
}

var _ contractsProviders.ILoggerProvider = (*MockLoggerProvider)(nil)

// Error logs an error message
func (mlp *MockLoggerProvider) Error(message string, err error) {
	fmt.Printf("Error: %s, %s", message, err)
}

// Panic logs a panic message
func (mlp *MockLoggerProvider) Panic(message string, err error) {
	fmt.Printf("Panic: %s, %s", message, err)
}

// ErrorMsg logs an error message
func (mlp *MockLoggerProvider) ErrorMsg(message string) {
	fmt.Print(message)
}

// Info logs an info message
func (mlp *MockLoggerProvider) Info(message string) {
	fmt.Print(message)
}

// Warning logs a warning message
func (mlp *MockLoggerProvider) Warning(message string) {
	fmt.Print(message)
}

// Debug logs a debug message
func (mlp *MockLoggerProvider) Debug(message string, data any) {
	fmt.Printf("%s: %v", message, data)
}

// ErrorWithContext logs an error message with context
func (mlp *MockLoggerProvider) ErrorWithContext(message string, err error, appCtx *app_context.AppContext) {
	fmt.Printf("Error: %s, %s", message, err)
}

// InfoWithContext logs an info message with context
func (mlp *MockLoggerProvider) InfoWithContext(message string, appCtx *app_context.AppContext) {
	fmt.Print(message)
}

// WarningWithContext logs a warning message with context
func (mlp *MockLoggerProvider) WarningWithContext(message string, appCtx *app_context.AppContext) {
	fmt.Print(message)
}

// DebugWithContext logs a debug message with context
func (mlp *MockLoggerProvider) DebugWithContext(message string, data any, appCtx *app_context.AppContext) {
	fmt.Printf("%s: %v", message, data)
}
