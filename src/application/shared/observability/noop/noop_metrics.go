package noop

import (
	"fmt"
	"log"
	"strings"
	"time"

	contractsobservability "github.com/simon3640/goprojectskeleton/src/application/contracts/observability"
)

// NoOpMetricsCollector is a no-op implementation of MetricsCollector that only logs metrics
type NoOpMetricsCollector struct{}

var _ contractsobservability.MetricsCollector = (*NoOpMetricsCollector)(nil)

// NewNoOpMetricsCollector creates a new NoOpMetricsCollector
func NewNoOpMetricsCollector() *NoOpMetricsCollector {
	return &NoOpMetricsCollector{}
}

// RecordLatency logs latency metrics in a structured format
func (n *NoOpMetricsCollector) RecordLatency(
	operation string,
	duration time.Duration,
	tags map[string]string,
) {
	tagsStr := formatTags(tags)
	log.Printf("[METRICS] operation: %s, duration: %v, tags: %s", operation, duration, tagsStr)
}

// IncrementCounter logs counter increments in a structured format
func (n *NoOpMetricsCollector) IncrementCounter(
	name string,
	tags map[string]string,
) {
	tagsStr := formatTags(tags)
	log.Printf("[METRICS] counter: %s, value: +1, tags: %s", name, tagsStr)
}

// formatTags formats a map of tags into a string for logging
func formatTags(tags map[string]string) string {
	if len(tags) == 0 {
		return "{}"
	}

	var parts []string
	for k, v := range tags {
		parts = append(parts, fmt.Sprintf("%s=%s", k, v))
	}
	return "{" + strings.Join(parts, ", ") + "}"
}
