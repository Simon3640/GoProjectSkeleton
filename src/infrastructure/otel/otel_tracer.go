package main

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	contractsobservability "github.com/simon3640/goprojectskeleton/src/application/contracts/observability"
	app_context "github.com/simon3640/goprojectskeleton/src/application/shared/context"
)

// OtelTracer is an implementation of Tracer using OpenTelemetry
type OtelTracer struct {
	tracer trace.Tracer
}

var _ contractsobservability.Tracer = (*OtelTracer)(nil)

// NewOtelTracer creates a new OtelTracer
func NewOtelTracer(tracerName string) *OtelTracer {
	tracer := otel.Tracer(tracerName)
	return &OtelTracer{
		tracer: tracer,
	}
}

// StartSpan creates a span from a TraceContextCarrier
func (o *OtelTracer) StartSpan(
	carrier contractsobservability.TraceContextCarrier,
	name string,
	opts ...contractsobservability.SpanOption,
) contractsobservability.Span {
	var ctx context.Context
	var span trace.Span

	// Try to extract the context from the carrier
	if carrier != nil && carrier.HasTrace() {
		traceCtx := carrier.TraceContext()
		if otelTraceCtx, ok := traceCtx.(*OtelTraceContext); ok {
			// Create a context with the span context of OpenTelemetry
			spanCtx := otelTraceCtx.SpanContext()
			ctx = trace.ContextWithRemoteSpanContext(context.Background(), spanCtx)
		} else {
			ctx = context.Background()
		}
	} else {
		ctx = context.Background()
	}

	// If the carrier is an AppContext, use its internal context
	if appCtx, ok := carrier.(*app_context.AppContext); ok {
		ctx = appCtx.Context
	}

	// Create the span
	ctx, span = o.tracer.Start(ctx, name)

	// Apply options
	otelSpan := NewOtelSpan(span)
	for _, opt := range opts {
		opt(otelSpan)
	}

	// CRITICAL: If the carrier is an AppContext, update its Context with the context that contains the span
	// This allows the logger to get the span from the context
	if appCtx, ok := carrier.(*app_context.AppContext); ok {
		ctxWithSpan := trace.ContextWithSpan(ctx, span)
		appCtx.Context = ctxWithSpan
	}

	return otelSpan
}
