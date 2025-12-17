package otel

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"

	observabilitycontracts "github.com/simon3640/goprojectskeleton/src/application/contracts/observability"
)

// OtelTracePropagator is an implementation of TracePropagator using OpenTelemetry
type OtelTracePropagator struct {
	propagator propagation.TextMapPropagator
}

var _ observabilitycontracts.TracePropagator = (*OtelTracePropagator)(nil)

// NewOtelTracePropagator creates a new OtelTracePropagator
func NewOtelTracePropagator() *OtelTracePropagator {
	return &OtelTracePropagator{
		propagator: otel.GetTextMapPropagator(),
	}
}

// Extract extracts TraceContext from HTTP headers
func (o *OtelTracePropagator) Extract(headers map[string]string) (observabilitycontracts.TraceContext, bool) {
	// Create a carrier from the headers
	carrier := make(propagation.HeaderCarrier)
	for k, v := range headers {
		carrier.Set(k, v)
	}

	// Extract the context
	ctx := o.propagator.Extract(context.Background(), carrier)
	spanCtx := trace.SpanContextFromContext(ctx)

	if !spanCtx.IsValid() {
		return nil, false
	}

	return NewOtelTraceContext(spanCtx), true
}

// Inject injects TraceContext into HTTP headers
func (o *OtelTracePropagator) Inject(tc observabilitycontracts.TraceContext, headers map[string]string) {
	if tc == nil || !tc.IsValid() {
		return
	}

	otelTraceCtx, ok := tc.(*OtelTraceContext)
	if !ok {
		return
	}

	spanCtx := otelTraceCtx.SpanContext()
	if !spanCtx.IsValid() {
		return
	}

	// Create a context with the span context
	ctx := trace.ContextWithRemoteSpanContext(context.Background(), spanCtx)

	// Create a carrier and inject
	carrier := make(propagation.HeaderCarrier)
	o.propagator.Inject(ctx, carrier)

	// Copy the headers to the provided map
	for k, v := range carrier {
		if len(v) > 0 {
			headers[k] = v[0]
		}
	}
}
