package otel

import (
	contractsobservability "github.com/simon3640/goprojectskeleton/src/application/contracts/observability"
	"go.opentelemetry.io/otel/trace"
)

// OtelTraceContext is an implementation of TraceContext using OpenTelemetry
type OtelTraceContext struct {
	spanContext trace.SpanContext
}

var _ contractsobservability.TraceContext = (*OtelTraceContext)(nil)

// NewOtelTraceContext creates a new OtelTraceContext from a trace.SpanContext
func NewOtelTraceContext(spanContext trace.SpanContext) *OtelTraceContext {
	return &OtelTraceContext{
		spanContext: spanContext,
	}
}

// TraceID returns the trace ID in hexadecimal format
func (o *OtelTraceContext) TraceID() string {
	return o.spanContext.TraceID().String()
}

// SpanID returns the span ID in hexadecimal format
func (o *OtelTraceContext) SpanID() string {
	return o.spanContext.SpanID().String()
}

// IsValid indicates if the context is valid
func (o *OtelTraceContext) IsValid() bool {
	return o.spanContext.IsValid()
}

// SpanContext returns the span context of OpenTelemetry (for internal use)
func (o *OtelTraceContext) SpanContext() trace.SpanContext {
	return o.spanContext
}
