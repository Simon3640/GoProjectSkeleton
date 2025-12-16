package otel

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	contractsobservability "github.com/simon3640/goprojectskeleton/src/application/contracts/observability"
)

// OtelMetricsCollector es una implementación de MetricsCollector usando OpenTelemetry
type OtelMetricsCollector struct {
	meter metric.Meter
}

var _ contractsobservability.MetricsCollector = (*OtelMetricsCollector)(nil)

// NewOtelMetricsCollector crea un nuevo OtelMetricsCollector
func NewOtelMetricsCollector(meterName string) *OtelMetricsCollector {
	meter := otel.Meter(meterName)
	return &OtelMetricsCollector{
		meter: meter,
	}
}

// RecordLatency registra latencia de una operación
func (o *OtelMetricsCollector) RecordLatency(
	operation string,
	duration time.Duration,
	tags map[string]string,
) {
	attrs := convertTagsToAttributes(tags)
	attrs = append(attrs, attribute.String("operation", operation))

	histogram, err := o.meter.Float64Histogram(
		operation+".duration",
		metric.WithDescription("Duration of "+operation),
		metric.WithUnit("ms"),
	)
	if err != nil {
		return
	}

	histogram.Record(context.Background(), float64(duration.Milliseconds()), metric.WithAttributes(attrs...))
}

// IncrementCounter incrementa un contador
func (o *OtelMetricsCollector) IncrementCounter(
	name string,
	tags map[string]string,
) {
	attrs := convertTagsToAttributes(tags)

	counter, err := o.meter.Int64Counter(
		name,
		metric.WithDescription("Counter for "+name),
	)
	if err != nil {
		return
	}

	counter.Add(context.Background(), 1, metric.WithAttributes(attrs...))
}

// convertTagsToAttributes convierte un mapa de tags a atributos de OpenTelemetry
func convertTagsToAttributes(tags map[string]string) []attribute.KeyValue {
	if tags == nil {
		return nil
	}

	attrs := make([]attribute.KeyValue, 0, len(tags))
	for k, v := range tags {
		attrs = append(attrs, attribute.String(k, v))
	}
	return attrs
}
