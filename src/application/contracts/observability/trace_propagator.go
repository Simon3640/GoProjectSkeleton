package contractsobservability

// TracePropagator handles HTTP propagation (separate from Tracer)
type TracePropagator interface {
	// Extract extracts TraceContext from HTTP headers
	Extract(headers map[string]string) (TraceContext, bool)

	// Inject injects TraceContext into HTTP headers
	Inject(tc TraceContext, headers map[string]string)
}
