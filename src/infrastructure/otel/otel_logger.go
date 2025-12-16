package otel

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	observabilitycontracts "github.com/simon3640/goprojectskeleton/src/application/contracts/observability"
	contractsproviders "github.com/simon3640/goprojectskeleton/src/application/contracts/providers"
	app_context "github.com/simon3640/goprojectskeleton/src/application/shared/context"
)

// OtelLogger implementa ILoggerProvider directamente con métodos que aceptan *AppContext
var _ contractsproviders.ILoggerProvider = (*OtelLogger)(nil)

// LoggerWrapper es un wrapper que adapta OtelLogger para implementar contractsobservability.Logger
type LoggerWrapper struct {
	*OtelLogger
}

var _ observabilitycontracts.Logger = (*LoggerWrapper)(nil)

// NewLoggerWrapper crea un wrapper para OtelLogger que implementa contractsobservability.Logger
func NewLoggerWrapper(logger *OtelLogger) *LoggerWrapper {
	return &LoggerWrapper{OtelLogger: logger}
}

// OtelLogger is a logger specific to OpenTelemetry that preserves the span context
type OtelLogger struct {
	tracer           observabilitycontracts.Tracer
	fallbackToStdout bool
}

// NewOtelLogger creates a new logger for OpenTelemetry
func NewOtelLogger(tracer observabilitycontracts.Tracer, fallbackToStdout bool) *OtelLogger {
	return &OtelLogger{
		tracer:           tracer,
		fallbackToStdout: fallbackToStdout,
	}
}

// logEvent adds a log event to the active span
func (l *OtelLogger) logEvent(ctx context.Context, level string, message string, attrs map[string]interface{}) {
	// Get the span from the OpenTelemetry context
	span := trace.SpanFromContext(ctx)

	// If there is no valid span, try to get it from the AppContext if available
	if !span.SpanContext().IsValid() {
		// Try to get AppContext from the context
		if appCtxValue := ctx.Value("app_context"); appCtxValue != nil {
			if appCtx, ok := appCtxValue.(*app_context.AppContext); ok && appCtx != nil {
				// Use directly the context of the AppContext that now contains the span
				span = trace.SpanFromContext(appCtx.Context)
			}
		}
	}

	// If there is a valid span, add the event
	if span.SpanContext().IsValid() {
		otelAttrs := []attribute.KeyValue{
			attribute.String("log.message", message),
			attribute.String("log.level", level),
			attribute.String("log.severity", level),
		}

		// Add additional attributes
		for key, value := range attrs {
			otelAttrs = append(otelAttrs, l.convertToAttribute(key, value))
		}

		// Add the event to the span
		eventName := fmt.Sprintf("log.%s", level)
		span.AddEvent(eventName, trace.WithAttributes(otelAttrs...))
	}

	// Write to stdout if enabled
	if l.fallbackToStdout {
		l.writeToStdout(level, message)
	}
}

// convertToAttribute converts a value to an OpenTelemetry attribute
func (l *OtelLogger) convertToAttribute(key string, value interface{}) attribute.KeyValue {
	attrKey := "log." + key
	switch v := value.(type) {
	case string:
		return attribute.String(attrKey, v)
	case int:
		return attribute.Int(attrKey, v)
	case int64:
		return attribute.Int64(attrKey, v)
	case uint64:
		return attribute.Int64(attrKey, int64(v))
	case float64:
		return attribute.Float64(attrKey, v)
	case bool:
		return attribute.Bool(attrKey, v)
	case time.Duration:
		return attribute.String(attrKey, v.String())
	case time.Time:
		return attribute.String(attrKey, v.Format(time.RFC3339))
	case error:
		return attribute.String(attrKey, v.Error())
	default:
		return attribute.String(attrKey, fmt.Sprintf("%v", v))
	}
}

// writeToStdout writes the log to stdout with colors
func (l *OtelLogger) writeToStdout(level, message string) {
	var color string
	switch level {
	case "ERROR", "error":
		color = "\033[31m"
	case "WARN", "warn":
		color = "\033[33m"
	case "INFO", "info":
		color = "\033[32m"
	case "DEBUG", "debug":
		color = "\033[36m"
	default:
		color = "\033[0m"
	}

	os.Stdout.WriteString(fmt.Sprintf("%s[%s] %s\033[0m\n", color, level, message))
}

// Implementation of the ILoggerProvider contract

// Error registers an error without context
func (l *OtelLogger) Error(message string, err error) {
	ctx := context.Background()
	errMsg := message
	if err != nil {
		errMsg = fmt.Sprintf("%s: %v", message, err)
	}
	l.logEvent(ctx, "ERROR", errMsg, map[string]interface{}{
		"error": err,
	})
}

// Panic registers a panic and throws it
func (l *OtelLogger) Panic(message string, err error) {
	ctx := context.Background()
	errMsg := fmt.Sprintf("%s: %v", message, err)
	l.logEvent(ctx, "ERROR", errMsg, map[string]interface{}{
		"error": err,
		"panic": true,
	})
	panic(errMsg)
}

// ErrorMsg registers an error message without error object
func (l *OtelLogger) ErrorMsg(message string) {
	ctx := context.Background()
	l.logEvent(ctx, "ERROR", message, nil)
}

// Info registers an INFO level message without context
func (l *OtelLogger) Info(message string) {
	ctx := context.Background()
	l.logEvent(ctx, "INFO", message, nil)
}

// Warning registers a WARNING level message without context
func (l *OtelLogger) Warning(message string) {
	ctx := context.Background()
	l.logEvent(ctx, "WARN", message, nil)
}

// Debug registers a DEBUG level message without context
func (l *OtelLogger) Debug(message string, data any) {
	ctx := context.Background()
	attrs := make(map[string]interface{})
	if data != nil {
		attrs["data"] = data
	}
	l.logEvent(ctx, "DEBUG", message, attrs)
}

// ErrorWithContext implementa contractsproviders.ILoggerProvider (acepta *AppContext)
func (l *OtelLogger) ErrorWithContext(message string, err error, appCtx *app_context.AppContext) {
	ctx := l.createContextFromAppContext(appCtx)
	errMsg := message
	if err != nil {
		errMsg = fmt.Sprintf("%s: %v", message, err)
	}
	l.logEvent(ctx, "ERROR", errMsg, map[string]interface{}{
		"error": err,
	})
}

// InfoWithContext implementa contractsproviders.ILoggerProvider (acepta *AppContext)
func (l *OtelLogger) InfoWithContext(message string, appCtx *app_context.AppContext) {
	ctx := l.createContextFromAppContext(appCtx)
	l.logEvent(ctx, "INFO", message, nil)
}

// WarningWithContext implementa contractsproviders.ILoggerProvider (acepta *AppContext)
func (l *OtelLogger) WarningWithContext(message string, appCtx *app_context.AppContext) {
	ctx := l.createContextFromAppContext(appCtx)
	l.logEvent(ctx, "WARN", message, nil)
}

// DebugWithContext implementa contractsproviders.ILoggerProvider (acepta *AppContext)
func (l *OtelLogger) DebugWithContext(message string, data any, appCtx *app_context.AppContext) {
	ctx := l.createContextFromAppContext(appCtx)
	attrs := make(map[string]interface{})
	if data != nil {
		attrs["data"] = data
	}
	l.logEvent(ctx, "DEBUG", message, attrs)
}

// Métodos del LoggerWrapper que implementan contractsobservability.Logger

func (w *LoggerWrapper) ErrorWithContext(message string, err error, appCtx interface{}) {
	var appCtxTyped *app_context.AppContext
	if appCtx != nil {
		if typed, ok := appCtx.(*app_context.AppContext); ok {
			appCtxTyped = typed
		}
	}
	w.OtelLogger.ErrorWithContext(message, err, appCtxTyped)
}

func (w *LoggerWrapper) InfoWithContext(message string, appCtx interface{}) {
	var appCtxTyped *app_context.AppContext
	if appCtx != nil {
		if typed, ok := appCtx.(*app_context.AppContext); ok {
			appCtxTyped = typed
		}
	}
	w.OtelLogger.InfoWithContext(message, appCtxTyped)
}

func (w *LoggerWrapper) WarningWithContext(message string, appCtx interface{}) {
	var appCtxTyped *app_context.AppContext
	if appCtx != nil {
		if typed, ok := appCtx.(*app_context.AppContext); ok {
			appCtxTyped = typed
		}
	}
	w.OtelLogger.WarningWithContext(message, appCtxTyped)
}

func (w *LoggerWrapper) DebugWithContext(message string, data any, appCtx interface{}) {
	var appCtxTyped *app_context.AppContext
	if appCtx != nil {
		if typed, ok := appCtx.(*app_context.AppContext); ok {
			appCtxTyped = typed
		}
	}
	w.OtelLogger.DebugWithContext(message, data, appCtxTyped)
}

// createContextFromAppContext creates a context that includes the AppContext
func (l *OtelLogger) createContextFromAppContext(appCtx interface{}) context.Context {
	if appCtx == nil {
		return context.Background()
	}
	// Type assertion para obtener el AppContext
	if appCtxTyped, ok := appCtx.(*app_context.AppContext); ok && appCtxTyped != nil {
		// Usar directamente el contexto del AppContext que ahora contiene el span
		// Agregar el AppContext al contexto para que pueda ser accedido si es necesario
		ctx := context.WithValue(appCtxTyped.Context, "app_context", appCtxTyped)
		return ctx
	}
	return context.Background()
}
