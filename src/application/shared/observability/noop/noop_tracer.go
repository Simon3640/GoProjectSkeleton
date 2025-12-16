package noop

import (
	"fmt"
	"log"

	contractsobservability "github.com/simon3640/goprojectskeleton/src/application/contracts/observability"
	app_context "github.com/simon3640/goprojectskeleton/src/application/shared/context"
)

// NoOpTracer is a no-op implementation of Tracer that only logs operations
type NoOpTracer struct{}

var _ contractsobservability.Tracer = (*NoOpTracer)(nil)

// NewNoOpTracer creates a new NoOpTracer
func NewNoOpTracer() *NoOpTracer {
	return &NoOpTracer{}
}

// StartSpan creates a span that only logs the operation
func (n *NoOpTracer) StartSpan(
	carrier contractsobservability.TraceContextCarrier,
	name string,
	opts ...contractsobservability.SpanOption,
) contractsobservability.Span {
	span := &NoOpSpan{
		name:    name,
		started: true,
		attrs:   make(map[string]interface{}),
	}

	// Apply options
	for _, opt := range opts {
		opt(span)
	}

	// Log span start
	log.Printf("[TRACE] Starting span: %s", name)

	// Update AppContext if it's an AppContext
	if appCtx, ok := carrier.(*app_context.AppContext); ok {
		// Create a no-op trace context for the span
		traceCtx := NewNoOpTraceContext()
		appCtx.WithTraceContext(traceCtx)
	}

	return span
}

// NoOpSpan is a no-op implementation of Span that only logs operations
type NoOpSpan struct {
	name      string
	started   bool
	ended     bool
	status    contractsobservability.SpanStatus
	statusMsg string
	attrs     map[string]interface{}
}

var _ contractsobservability.Span = (*NoOpSpan)(nil)

// SetAttribute adds an attribute to the span (logged)
func (n *NoOpSpan) SetAttribute(key string, value interface{}) {
	if n.attrs == nil {
		n.attrs = make(map[string]interface{})
	}
	n.attrs[key] = value
	log.Printf("[TRACE] Span '%s' attribute: %s = %v", n.name, key, value)
}

// SetStatus marks the span with a status (logged)
func (n *NoOpSpan) SetStatus(status contractsobservability.SpanStatus, description string) {
	n.status = status
	n.statusMsg = description
	statusStr := "UNSET"
	switch status {
	case contractsobservability.SpanStatusOK:
		statusStr = "OK"
	case contractsobservability.SpanStatusError:
		statusStr = "ERROR"
	}
	if description != "" {
		log.Printf("[TRACE] Span '%s' status: %s - %s", n.name, statusStr, description)
	} else {
		log.Printf("[TRACE] Span '%s' status: %s", n.name, statusStr)
	}
}

// End finalizes the span (logged)
func (n *NoOpSpan) End() {
	if n.ended {
		return
	}
	n.ended = true

	statusStr := "UNSET"
	switch n.status {
	case contractsobservability.SpanStatusOK:
		statusStr = "OK"
	case contractsobservability.SpanStatusError:
		statusStr = "ERROR"
	}

	if len(n.attrs) > 0 {
		log.Printf("[TRACE] Ending span: %s, status: %s, attributes: %v", n.name, statusStr, n.attrs)
	} else {
		log.Printf("[TRACE] Ending span: %s, status: %s", n.name, statusStr)
	}
}

// UpdateAppContext updates the AppContext with the TraceContext of the span
func (n *NoOpSpan) UpdateAppContext(appCtx interface{}) {
	if ctx, ok := appCtx.(*app_context.AppContext); ok {
		traceCtx := NewNoOpTraceContext()
		ctx.WithTraceContext(traceCtx)
	}
}

// NoOpTraceContext is a no-op implementation of TraceContext
type NoOpTraceContext struct {
	traceID string
	spanID  string
}

var _ contractsobservability.TraceContext = (*NoOpTraceContext)(nil)

// NewNoOpTraceContext creates a new NoOpTraceContext
func NewNoOpTraceContext() *NoOpTraceContext {
	// Generate simple IDs for logging purposes
	return &NoOpTraceContext{
		traceID: fmt.Sprintf("noop-trace-%d", 0),
		spanID:  fmt.Sprintf("noop-span-%d", 0),
	}
}

// TraceID returns a placeholder trace ID
func (n *NoOpTraceContext) TraceID() string {
	return n.traceID
}

// SpanID returns a placeholder span ID
func (n *NoOpTraceContext) SpanID() string {
	return n.spanID
}

// IsValid always returns false for no-op context
func (n *NoOpTraceContext) IsValid() bool {
	return true
}
