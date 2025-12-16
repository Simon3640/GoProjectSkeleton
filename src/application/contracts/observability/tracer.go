package contractsobservability

// TraceContextCarrier is an interface that allows obtaining TraceContext without creating circular dependency
type TraceContextCarrier interface {
	TraceContext() TraceContext
	WithTraceContext(tc TraceContext) interface{}
	HasTrace() bool
}

// Tracer creates spans from a TraceContextCarrier (typically AppContext)
// If the carrier does not have a TraceContext, creates the Trace root automatically
type Tracer interface {
	// StartSpan creates a span from a TraceContextCarrier
	// If there is no TraceContext in the carrier, creates the Trace root automatically
	StartSpan(
		carrier TraceContextCarrier,
		name string,
		opts ...SpanOption,
	) Span
}
