package usecase

import (
	contractsobservability "github.com/simon3640/goprojectskeleton/src/application/contracts/observability"
	app_context "github.com/simon3640/goprojectskeleton/src/application/shared/context"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales"
	"github.com/simon3640/goprojectskeleton/src/application/shared/status"
)

// InstrumentUseCase executes a UseCase with automatic instrumentation
func InstrumentUseCase[Input any, Output any](
	uc BaseUseCase[Input, Output],
	appCtx *app_context.AppContext,
	locale locales.LocaleTypeEnum,
	input Input,
	tracer contractsobservability.Tracer,
	metrics contractsobservability.MetricsCollector,
	clock contractsobservability.Clock,
	useCaseName string,
) *UseCaseResult[Output] {
	// Create span for the use case
	span := tracer.StartSpan(appCtx, "usecase."+useCaseName)
	defer span.End()

	// Update AppContext with the TraceContext of the span
	span.UpdateAppContext(appCtx)

	// Measure latency
	start := clock.Now()
	result := uc.Execute(appCtx, locale, input)
	duration := clock.Now().Sub(start)

	// Register metrics
	tags := map[string]string{
		"usecase": useCaseName,
	}
	metrics.RecordLatency("usecase.execute", duration, tags)

	// Mark span with status according to result
	if result != nil {
		if result.HasError() {
			span.SetStatus(contractsobservability.SpanStatusError, result.GetError().Error())
			tags["status"] = "error"
			tags["status_code"] = string(result.StatusCode)
			metrics.IncrementCounter("usecase.error", tags)
		} else {
			span.SetStatus(contractsobservability.SpanStatusOK, "")
			tags["status"] = "success"
			metrics.IncrementCounter("usecase.success", tags)
		}
	} else {
		span.SetStatus(contractsobservability.SpanStatusError, "nil result")
		tags["status"] = "error"
		tags["status_code"] = string(status.InternalError)
		metrics.IncrementCounter("usecase.error", tags)
	}

	return result
}

// InstrumentDAGStep executes a step of the DAG with instrumentation
func InstrumentDAGStep[I any, O any](
	step DagStep[I, O],
	appCtx *app_context.AppContext,
	locale locales.LocaleTypeEnum,
	input I,
	tracer contractsobservability.Tracer,
	metrics contractsobservability.MetricsCollector,
	clock contractsobservability.Clock,
	stepName string,
) *UseCaseResult[O] {
	// Create span for the step of the DAG
	span := tracer.StartSpan(appCtx, "dag.step."+stepName)
	defer span.End()

	// Update AppContext with the TraceContext of the span
	span.UpdateAppContext(appCtx)

	// Measure latency
	start := clock.Now()
	result := step.uc.Execute(appCtx, locale, input)
	duration := clock.Now().Sub(start)

	// Register metrics
	tags := map[string]string{
		"dag_step": stepName,
	}
	metrics.RecordLatency("dag.step.execute", duration, tags)

	// Mark span with status according to result
	if result != nil {
		if result.HasError() {
			span.SetStatus(contractsobservability.SpanStatusError, result.GetError().Error())
			tags["status"] = "error"
			metrics.IncrementCounter("dag.step.error", tags)
		} else {
			span.SetStatus(contractsobservability.SpanStatusOK, "")
			tags["status"] = "success"
			metrics.IncrementCounter("dag.step.success", tags)
		}
	} else {
		span.SetStatus(contractsobservability.SpanStatusError, "nil result")
		tags["status"] = "error"
		metrics.IncrementCounter("dag.step.error", tags)
	}

	return result
}
