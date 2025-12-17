package observability

import contractsobservability "github.com/simon3640/goprojectskeleton/src/application/contracts/observability"

// WithOperation sets the name of the operation
func WithOperation(operation string) contractsobservability.SpanOption {
	return func(s contractsobservability.Span) {
		s.SetAttribute("operation", operation)
	}
}

// WithFollowsFrom sets that this span follows from another (for background tasks)
func WithFollowsFrom(parentTraceContext contractsobservability.TraceContext) contractsobservability.SpanOption {
	return func(s contractsobservability.Span) {
		s.SetAttribute("follows_from.trace_id", parentTraceContext.TraceID())
		s.SetAttribute("follows_from.span_id", parentTraceContext.SpanID())
	}
}
