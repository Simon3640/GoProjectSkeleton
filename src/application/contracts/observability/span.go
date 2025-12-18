package contractsobservability

// SpanOption is an option to configure a span
type SpanOption func(Span)

// SpanStatus represents the status of a span
type SpanStatus int

const (
	// SpanStatusUnset indicates that the status is not set
	SpanStatusUnset SpanStatus = iota
	// SpanStatusOK indicates that the operation was successful
	SpanStatusOK
	// SpanStatusError indicates that the operation failed
	SpanStatusError
)

// Span represents an instrumented operation
type Span interface {
	// SetAttribute adds attributes to the span
	SetAttribute(key string, value interface{})

	// SetStatus marks the span as successful or error
	SetStatus(status SpanStatus, description string)

	// End finalizes the span (must be called always)
	End()

	// UpdateAppContext updates the AppContext with the TraceContext of the span
	// This allows propagating the trace to child operations
	UpdateAppContext(appCtx interface{})
}
