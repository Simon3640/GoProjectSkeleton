package app_context

import "time"

// Trace represents a trace in the application
type Trace struct {
	TraceID   string
	ParentID  *string
	StartedAt time.Time
	Metadata  map[string]string
}
