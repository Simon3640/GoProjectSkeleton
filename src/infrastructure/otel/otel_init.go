// Package otel contains the implementation of the OtelInit
package otel

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/simon3640/goprojectskeleton/src/application/shared/observability"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.27.0"
)

var (
	tracerProvider *trace.TracerProvider
	metricProvider *metric.MeterProvider
	loggerProvider *OtelLogger
	initialized    bool
	shutdownFuncs  []func(context.Context) error
)

// OtelConfig configures OpenTelemetry
type OtelConfig struct {
	ServiceName    string
	ServiceVersion string
	OTLPEndpoint   string // Endpoint OTLP for traces and metrics (e.g. http://localhost:4318 or localhost:4318)
	Environment    string // e.g. "development", "production"
}

// normalizeEndpoint normalizes the endpoint by removing the scheme http:// or https://
// because the OpenTelemetry exporters add it automatically.
// Also resolves Docker host names to localhost if the app is running outside of Docker.
func normalizeEndpoint(endpoint string) string {
	if endpoint == "" {
		return endpoint
	}

	// If it's already just host:port, return it as is
	if !strings.Contains(endpoint, "://") {
		// If the host is a Docker name and we're outside of Docker, use localhost
		return resolveDockerHost(endpoint)
	}

	// Parse the URL to extract only host:port
	parsedURL, err := url.Parse(endpoint)
	if err != nil {
		// If parsing fails, try to remove manually http:// or https://
		endpoint = strings.TrimPrefix(endpoint, "http://")
		endpoint = strings.TrimPrefix(endpoint, "https://")
		return resolveDockerHost(endpoint)
	}

	// Return host:port
	var hostPort string
	if parsedURL.Port() != "" {
		hostPort = fmt.Sprintf("%s:%s", parsedURL.Hostname(), parsedURL.Port())
	} else {
		hostPort = parsedURL.Hostname()
	}
	return resolveDockerHost(hostPort)
}

// resolveDockerHost resolves Docker host names to localhost if the app is running outside of Docker
func resolveDockerHost(endpoint string) string {
	// Detect if we're inside Docker by checking if /.dockerenv exists
	_, err := os.Stat("/.dockerenv")
	isInsideDocker := err == nil

	// If we're inside Docker, return the endpoint as is
	if isInsideDocker {
		return endpoint
	}

	// If we're outside of Docker and the endpoint contains common Docker names,
	// try to resolve first, and if it fails, use localhost
	dockerHosts := []string{"otel-collector", "jaeger", "prometheus", "db", "redis", "mailhog"}

	parts := strings.Split(endpoint, ":")
	if len(parts) > 0 {
		host := parts[0]
		port := ""
		if len(parts) > 1 {
			port = ":" + strings.Join(parts[1:], ":")
		}

		// Check if it's a Docker host name
		for _, dockerHost := range dockerHosts {
			if host == dockerHost {
				// Try to resolve the name
				_, err := net.LookupHost(host)
				if err != nil {
					// If it can't be resolved, use localhost
					return "localhost" + port
				}
				// If it can be resolved, return as is
				return endpoint
			}
		}
	}

	return endpoint
}

// InitializeOtelSDK initializes the OpenTelemetry SDK
func InitializeOtelSDK(config OtelConfig) error {
	if initialized {
		return nil
	}

	ctx := context.Background()

	// Create resource with service information
	attrs := []attribute.KeyValue{
		semconv.ServiceNameKey.String(config.ServiceName),
		semconv.ServiceVersionKey.String(config.ServiceVersion),
	}
	if config.Environment != "" {
		attrs = append(attrs, attribute.String("deployment.environment", config.Environment))
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(attrs...),
		resource.WithFromEnv(),
		resource.WithProcess(),
	)
	if err != nil {
		return fmt.Errorf("failed to create resource: %w", err)
	}

	// Initialize Tracer Provider
	if err := initTracerProvider(ctx, res, config); err != nil {
		return fmt.Errorf("failed to initialize tracer provider: %w", err)
	}

	// Initialize Metric Provider
	if err := initMetricProvider(ctx, res, config); err != nil {
		return fmt.Errorf("failed to initialize metric provider: %w", err)
	}

	// Initialize Logger Provider
	if err := initLoggerProvider(*NewOtelTracer(config.ServiceName), config); err != nil {
		return fmt.Errorf("failed to initialize logger provider: %w", err)
	}

	// Configure propagation
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	initialized = true
	return nil
}

// initTracerProvider initializes the Tracer Provider
func initTracerProvider(ctx context.Context, res *resource.Resource, config OtelConfig) error {
	var exporter trace.SpanExporter
	var err error

	if config.OTLPEndpoint != "" {
		endpoint := normalizeEndpoint(config.OTLPEndpoint)
		exporter, err = otlptracehttp.New(ctx,
			otlptracehttp.WithEndpoint(endpoint),
			otlptracehttp.WithInsecure(), // Use TLS in production
		)
		if err != nil {
			return fmt.Errorf("failed to create OTLP trace exporter: %w", err)
		}
	} else {
		// If there is no endpoint, use a no-op exporter
		exporter = &noOpSpanExporter{}
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(res),
		trace.WithSampler(trace.AlwaysSample()), // In production, use a more intelligent sampler
	)

	otel.SetTracerProvider(tp)
	tracerProvider = tp

	shutdownFuncs = append(shutdownFuncs, func(ctx context.Context) error {
		return tp.Shutdown(ctx)
	})

	return nil
}

// initMetricProvider initializes the Metric Provider
func initMetricProvider(ctx context.Context, res *resource.Resource, config OtelConfig) error {
	if config.OTLPEndpoint != "" {
		endpoint := normalizeEndpoint(config.OTLPEndpoint)
		exporter, err := otlpmetrichttp.New(ctx,
			otlpmetrichttp.WithEndpoint(endpoint),
			otlpmetrichttp.WithInsecure(), // Use TLS in production
		)
		if err != nil {
			return fmt.Errorf("failed to create OTLP metric exporter: %w", err)
		}

		mp := metric.NewMeterProvider(
			metric.WithResource(res),
			metric.WithReader(metric.NewPeriodicReader(exporter,
				metric.WithInterval(10*time.Second),
			)),
		)

		otel.SetMeterProvider(mp)
		metricProvider = mp

		shutdownFuncs = append(shutdownFuncs, func(ctx context.Context) error {
			return mp.Shutdown(ctx)
		})
	} else {
		// If there is no endpoint, create a basic provider without exporter
		mp := metric.NewMeterProvider(
			metric.WithResource(res),
		)

		otel.SetMeterProvider(mp)
		metricProvider = mp

		shutdownFuncs = append(shutdownFuncs, func(ctx context.Context) error {
			return mp.Shutdown(ctx)
		})
	}

	return nil
}

// initLoggerProvider initializes the Logger Provider
func initLoggerProvider(tracer OtelTracer, config OtelConfig) error {
	// Create the logger with fallback to stdout
	loggerProvider = NewOtelLogger(&tracer, true)
	return nil
}

// Shutdown closes all OpenTelemetry providers
func Shutdown(ctx context.Context) error {
	for _, shutdown := range shutdownFuncs {
		if err := shutdown(ctx); err != nil {
			fmt.Printf("Error shutting down OpenTelemetry provider: %v", err)
		}
	}
	return nil
}

// noOpSpanExporter is a no-op exporter for when no endpoint is configured
type noOpSpanExporter struct{}

// ExportSpans is a no-op implementation of ExportSpans
func (n *noOpSpanExporter) ExportSpans(ctx context.Context, spans []trace.ReadOnlySpan) error {
	return nil
}

// Shutdown is a no-op implementation of Shutdown
func (n *noOpSpanExporter) Shutdown(ctx context.Context) error {
	return nil
}

// NewOtelObservabilityComponents creates a new OtelObservabilityComponents
func NewOtelObservabilityComponents(config OtelConfig) *observability.ObservabilityComponents {
	tracer := NewOtelTracer(config.ServiceName)
	logger := NewOtelLogger(tracer, true)
	return &observability.ObservabilityComponents{
		Tracer:     tracer,
		Propagator: NewOtelTracePropagator(),
		Metrics:    NewOtelMetricsCollector(config.ServiceName),
		Clock:      NewRealClock(),
		Logger:     NewLoggerWrapper(logger),
	}
}
