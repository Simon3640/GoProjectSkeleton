package main

import (
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	contractsobservability "github.com/simon3640/goprojectskeleton/src/application/contracts/observability"
	app_context "github.com/simon3640/goprojectskeleton/src/application/shared/context"
)

// OtelSpan is an implementation of Span using OpenTelemetry
type OtelSpan struct {
	span trace.Span
}

var _ contractsobservability.Span = (*OtelSpan)(nil)

// NewOtelSpan creates a new OtelSpan from a trace.Span
func NewOtelSpan(span trace.Span) *OtelSpan {
	return &OtelSpan{
		span: span,
	}
}

// SetAttribute adds attributes to the span
func (o *OtelSpan) SetAttribute(key string, value interface{}) {
	attr := convertToAttribute(key, value)
	if attr.Key != "" {
		o.span.SetAttributes(attr)
	}
}

// SetStatus marks the span as successful or error
func (o *OtelSpan) SetStatus(status contractsobservability.SpanStatus, description string) {
	switch status {
	case contractsobservability.SpanStatusOK:
		o.span.SetStatus(codes.Ok, description)
	case contractsobservability.SpanStatusError:
		o.span.SetStatus(codes.Error, description)
	default:
		o.span.SetStatus(codes.Unset, description)
	}
}

// End finalizes the span
func (o *OtelSpan) End() {
	o.span.End()
}

// UpdateAppContext updates the AppContext with the TraceContext of the span
func (o *OtelSpan) UpdateAppContext(appCtx interface{}) {
	if ctx, ok := appCtx.(*app_context.AppContext); ok {
		spanCtx := o.span.SpanContext()
		traceCtx := NewOtelTraceContext(spanCtx)
		ctx.WithTraceContext(traceCtx)

		// CRITICAL: Update the internal context with the context that contains the span
		// This allows the logger to get the span from the context
		ctxWithSpan := trace.ContextWithSpan(ctx.Context, o.span)
		ctx.Context = ctxWithSpan
	}
}

// convertToAttribute converts a value to an OpenTelemetry attribute
func convertToAttribute(key string, value interface{}) attribute.KeyValue {
	switch v := value.(type) {
	case string:
		return attribute.String(key, v)
	case int:
		return attribute.Int(key, v)
	case int64:
		return attribute.Int64(key, v)
	case float64:
		return attribute.Float64(key, v)
	case bool:
		return attribute.Bool(key, v)
	default:
		// If it can't be converted, return an empty attribute
		return attribute.String(key, "")
	}
}
