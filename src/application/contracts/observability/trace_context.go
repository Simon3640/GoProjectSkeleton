package contractsobservability

// TraceContext transports information about trace/span
type TraceContext interface {
	// TraceID returns the ID of the trace
	TraceID() string

	// SpanID returns the ID of the current span
	SpanID() string

	// IsValid indicates if the context is valid
	IsValid() bool
}
