package noop

import (
	"time"

	contractsobservability "github.com/simon3640/goprojectskeleton/src/application/contracts/observability"
)

// NoOpClock is a no-op implementation of Clock using time.Now()
// It's essentially a wrapper around the standard time package
type NoOpClock struct{}

var _ contractsobservability.Clock = (*NoOpClock)(nil)

// NewNoOpClock creates a new NoOpClock
func NewNoOpClock() *NoOpClock {
	return &NoOpClock{}
}

// Now returns the current time
func (n *NoOpClock) Now() time.Time {
	return time.Now()
}
