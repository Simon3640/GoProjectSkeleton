package otel

import (
	"time"

	contractsobservability "github.com/simon3640/goprojectskeleton/src/application/contracts/observability"
)

// RealClock is a real implementation of Clock using time.Now()
type RealClock struct{}

var _ contractsobservability.Clock = (*RealClock)(nil)

// NewRealClock creates a new RealClock
func NewRealClock() *RealClock {
	return &RealClock{}
}

// Now returns the current time
func (r *RealClock) Now() time.Time {
	return time.Now()
}
