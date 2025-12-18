package contractsproviders

import app_context "github.com/simon3640/goprojectskeleton/src/application/shared/context"

type ILoggerProvider interface {
	Error(message string, err error)
	Panic(message string, err error)
	ErrorMsg(message string)
	Info(message string)
	Warning(message string)
	Debug(message string, data any)
	ErrorWithContext(message string, err error, appCtx *app_context.AppContext)
	InfoWithContext(message string, appCtx *app_context.AppContext)
	WarningWithContext(message string, appCtx *app_context.AppContext)
	DebugWithContext(message string, data any, appCtx *app_context.AppContext)
}
