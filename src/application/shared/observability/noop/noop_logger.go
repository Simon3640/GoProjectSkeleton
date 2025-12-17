package noop

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"

	contractsobservability "github.com/simon3640/goprojectskeleton/src/application/contracts/observability"
	app_context "github.com/simon3640/goprojectskeleton/src/application/shared/context"
)

var _ contractsobservability.Logger = (*NoOpLogger)(nil)

// NewNoOpLogger creates a new NoOpLogger
func NewNoOpLogger() *NoOpLogger {
	return &NoOpLogger{}
}

type NoOpLogger struct {
	enableLog bool
	debugLog  bool
	logger    *slog.Logger
}

type LogHandler struct {
	slog.Handler
}

func (h LogHandler) Enabled(_ context.Context, level slog.Level) bool {
	return true
}

func (h LogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

func (h LogHandler) WithGroup(name string) slog.Handler {
	return h
}

func (h LogHandler) Handle(ctx context.Context, record slog.Record) error {
	var color string
	switch record.Level {
	case slog.LevelDebug:
		color = "\033[36m"
	case slog.LevelInfo:
		color = "\033[32m"
	case slog.LevelWarn:
		color = "\033[33m"
	case slog.LevelError:
		color = "\033[31m"
	default:
		color = "\033[0m"
	}

	if _, err := os.Stdout.WriteString(color + record.Message + "\033[0m" + "\n"); err != nil {
		return err
	}
	return nil
}

func (lp *NoOpLogger) Setup(enableLog bool, debugLog bool) {
	lp.enableLog = enableLog
	lp.debugLog = debugLog
	lp.logger = slog.New(LogHandler{})
}

func (lp *NoOpLogger) Error(message string, err error) {
	if lp.enableLog {
		lp.logger.Error(fmt.Sprintf("Error: %s, %s", message, err))
	}
}

func (lp *NoOpLogger) Panic(message string, err error) {
	if lp.enableLog {
		log.Panicf("%s", fmt.Sprintf("Panic: %s, %s", message, err))
	}
}

func (lp *NoOpLogger) ErrorMsg(message string) {
	if lp.enableLog {
		lp.logger.Error(message)
	}
}

func (lp *NoOpLogger) Info(message string) {
	if lp.enableLog {
		lp.logger.Info(message)
	}
}

func (lp *NoOpLogger) Warning(message string) {
	if lp.enableLog {
		lp.logger.Warn(message)
	}
}

func (lp *NoOpLogger) Debug(message string, data any) {
	if lp.enableLog && lp.debugLog {
		if data != nil {
			lp.logger.Debug(fmt.Sprintf("%s: %v", message, data))
		} else {
			lp.logger.Debug(message)
		}
	}
}

// formatMessageWithTrace formatea el mensaje incluyendo trace_id y span_id si están disponibles
func (lp *NoOpLogger) formatMessageWithTrace(message string, appCtx interface{}) string {
	var appCtxTyped *app_context.AppContext
	if appCtx != nil {
		if typed, ok := appCtx.(*app_context.AppContext); ok {
			appCtxTyped = typed
		}
	}
	if appCtxTyped == nil || !appCtxTyped.HasTrace() {
		return message
	}

	traceCtx := appCtxTyped.TraceContext()
	if traceCtx == nil || !traceCtx.IsValid() {
		return message
	}

	return fmt.Sprintf("[trace_id=%s span_id=%s] %s", traceCtx.TraceID(), traceCtx.SpanID(), message)
}

// ErrorWithContext registra un error con contexto de trace
func (lp *NoOpLogger) ErrorWithContext(message string, err error, appCtx interface{}) {
	if lp.enableLog {
		formattedMsg := lp.formatMessageWithTrace(fmt.Sprintf("Error: %s, %s", message, err), appCtx)
		lp.logger.Error(formattedMsg)
	}
}

// InfoWithContext registra un mensaje info con contexto de trace
func (lp *NoOpLogger) InfoWithContext(message string, appCtx interface{}) {
	if lp.enableLog {
		formattedMsg := lp.formatMessageWithTrace(message, appCtx)
		lp.logger.Info(formattedMsg)
	}
}

// WarningWithContext registra un warning con contexto de trace
func (lp *NoOpLogger) WarningWithContext(message string, appCtx interface{}) {
	if lp.enableLog {
		formattedMsg := lp.formatMessageWithTrace(message, appCtx)
		lp.logger.Warn(formattedMsg)
	}
}

// DebugWithContext registra un debug con contexto de trace
func (lp *NoOpLogger) DebugWithContext(message string, data any, appCtx interface{}) {
	if lp.enableLog && lp.debugLog {
		formattedMsg := lp.formatMessageWithTrace(message, appCtx)
		if data != nil {
			lp.logger.Debug(fmt.Sprintf("%s: %v", formattedMsg, data))
		} else {
			lp.logger.Debug(formattedMsg)
		}
	}
}

var Logger *NoOpLogger

func init() {
	fmt.Println("NoOpLogger initialized")
	Logger = NewNoOpLogger()
}
