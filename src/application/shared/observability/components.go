package observability

import (
	contractsobservability "github.com/simon3640/goprojectskeleton/src/application/contracts/observability"
	"github.com/simon3640/goprojectskeleton/src/application/shared/observability/noop"
)

type ObservabilityComponents struct {
	Tracer     contractsobservability.Tracer
	Propagator contractsobservability.TracePropagator
	Metrics    contractsobservability.MetricsCollector
	Clock      contractsobservability.Clock
	Logger     contractsobservability.Logger
}

var ObservabilityComponentsInstance *ObservabilityComponents

// NewDefaultObservabilityComponents crea un conjunto completo de componentes de observabilidad no-op
// listos para usar. Todos los componentes son implementaciones no-op que solo loguean operaciones.
func NewDefaultObservabilityComponents() *ObservabilityComponents {
	return &ObservabilityComponents{
		Tracer:     noop.NewNoOpTracer(),
		Propagator: noop.NewNoOpTracePropagator(),
		Metrics:    noop.NewNoOpMetricsCollector(),
		Clock:      noop.NewNoOpClock(),
		Logger:     noop.NewNoOpLogger(),
	}
}

// init inicializa siempre con componentes no-op por defecto
func init() {
	ObservabilityComponentsInstance = NewDefaultObservabilityComponents()
}

// GetObservabilityComponents garantiza que siempre haya componentes de observabilidad disponibles.
// Si ObservabilityComponentsInstance es nil, lo inicializa con componentes no-op.
func GetObservabilityComponents() *ObservabilityComponents {
	if ObservabilityComponentsInstance == nil {
		ObservabilityComponentsInstance = NewDefaultObservabilityComponents()
	}
	return ObservabilityComponentsInstance
}
