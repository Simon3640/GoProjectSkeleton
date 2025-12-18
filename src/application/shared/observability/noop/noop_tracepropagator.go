package noop

import (
	"log"

	contractsobservability "github.com/simon3640/goprojectskeleton/src/application/contracts/observability"
)

// NoOpTracePropagator is a no-op implementation of TracePropagator that only logs operations
type NoOpTracePropagator struct{}

var _ contractsobservability.TracePropagator = (*NoOpTracePropagator)(nil)

// NewNoOpTracePropagator creates a new NoOpTracePropagator
func NewNoOpTracePropagator() *NoOpTracePropagator {
	return &NoOpTracePropagator{}
}

// Extract extracts TraceContext from HTTP headers (no-op, always returns nil, false)
func (n *NoOpTracePropagator) Extract(headers map[string]string) (contractsobservability.TraceContext, bool) {
	log.Printf("[TRACE] Extract called with %d headers (no-op: returning nil)", len(headers))
	return nil, false
}

// Inject injects TraceContext into HTTP headers (no-op, does nothing)
func (n *NoOpTracePropagator) Inject(tc contractsobservability.TraceContext, headers map[string]string) {
	if tc == nil {
		log.Printf("[TRACE] Inject called with nil TraceContext (no-op: doing nothing)")
		return
	}

	traceID := tc.TraceID()
	spanID := tc.SpanID()

	log.Printf("[TRACE] Inject called with traceID=%s, spanID=%s (no-op: not injecting into headers)", traceID, spanID)
	// No-op: do not modify headers
}
