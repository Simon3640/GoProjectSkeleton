package contractsobservability

import "time"

// MetricsCollector records basic metrics
type MetricsCollector interface {
	// RecordLatency records the latency of an operation
	RecordLatency(operation string, duration time.Duration, tags map[string]string)

	// IncrementCounter increments a counter (primarily for errors)
	IncrementCounter(name string, tags map[string]string)
}
