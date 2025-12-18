package contractsobservability

// Logger is the interface for the logger for observability and avoids cyclic dependencies
// Uses interface{} for AppContext to avoid importing shared/context
// The implementers must make type assertion to *app_context.AppContext
// Note: This interface is compatible with contractsproviders.ILoggerProvider
// when the implementer has methods that accept *app_context.AppContext
type Logger interface {
	Error(message string, err error)
	Panic(message string, err error)
	ErrorMsg(message string)
	Info(message string)
	Warning(message string)
	Debug(message string, data any)
	// Los métodos WithContext aceptan interface{} pero deben recibir *app_context.AppContext
	ErrorWithContext(message string, err error, appCtx interface{})
	InfoWithContext(message string, appCtx interface{})
	WarningWithContext(message string, appCtx interface{})
	DebugWithContext(message string, data any, appCtx interface{})
}
