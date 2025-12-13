package usecase

import (
	"context"
	"log"
	"runtime/debug"
	"sync"
	"time"

	app_context "github.com/simon3640/goprojectskeleton/src/application/shared/context"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales"
	"github.com/simon3640/goprojectskeleton/src/application/shared/status"
	"github.com/simon3640/goprojectskeleton/src/application/shared/workers"
)

// DagStep represents a single executable step in a DAG.
// It is a thin wrapper around a BaseUseCase and defines the
// input and output types for that step.
type DagStep[I any, O any] struct {
	uc BaseUseCase[I, O]
}

// NewStep creates a DagStep from the given use case.
// The returned step is immutable and can be safely reused
// across multiple DAGs.
func NewStep[I any, O any](uc BaseUseCase[I, O]) DagStep[I, O] { return DagStep[I, O]{uc: uc} }

// backgroundEntry represents a background task attached to a DAG node.
// Each entry captures a function to execute and an optional name
// used for logging or tracing purposes.
type backgroundEntry[O any] struct {
	name string
	fn   func(ctx *app_context.AppContext, out O)
}

// DAG represents an immutable, typed directed acyclic graph of use cases.
//
// A DAG is built by chaining steps using Then and ThenBackground.
// Each builder function returns a NEW DAG instance, making DAGs
// safe to reuse and share across goroutines.
//
// The DAG executes synchronously by default, with optional background
// tasks that run asynchronously after successful execution.
type DAG[I any, O any] struct {
	run         func(appContext *app_context.AppContext, input I) *UseCaseResult[O]
	ctx         *app_context.AppContext
	locale      locales.LocaleTypeEnum
	_background []backgroundEntry[O] // internal; treat as immutable — copy on write
	executor    *workers.BackgroundExecutor
}

// NewDag creates a new DAG starting with the provided initial step.
//
// The returned DAG:
//   - is immutable and safe for reuse,
//   - executes the first use case synchronously,
//   - may enqueue background tasks if a BackgroundExecutor is provided.
//
// Parameters:
//   - ctx: Application context shared by all steps in the DAG.
//   - first: The initial use case step of the DAG.
//   - locale: Locale passed to all use cases during execution.
//   - executor: Optional executor used to run background tasks asynchronously.
func NewDag[I any, O any](ctx *app_context.AppContext, first DagStep[I, O], locale locales.LocaleTypeEnum, executor *workers.BackgroundExecutor) *DAG[I, O] {
	run := func(appContext *app_context.AppContext, input I) *UseCaseResult[O] {
		return first.uc.Execute(appContext, locale, input)
	}
	return &DAG[I, O]{
		run:         run,
		ctx:         ctx,
		locale:      locale,
		_background: nil,
		executor:    executor,
	}
}

// Then chains a synchronous step to the DAG and returns a NEW DAG
// with a transformed output type.
//
// Execution semantics:
//   - The previous DAG is executed first.
//   - If the previous step fails, the error is propagated and the next
//     step is not executed.
//   - On success, the output of the previous step becomes the input
//     of the next step.
//
// Note:
//   - Background steps are not carried over when the output type changes,
//     since they are bound to the previous output type.
func Then[I any, O any, P any](d *DAG[I, O], next DagStep[O, P]) *DAG[I, P] {
	// capture previous background entries immutably by copying the slice
	prevBackground := make([]backgroundEntry[P], 0)
	// intentionally empty — background entries from previous output type cannot be reused

	run := func(ctx *app_context.AppContext, input I) *UseCaseResult[P] {
		r1 := d.run(ctx, input)
		if r1 == nil {
			return nil
		}
		if r1.HasError() {
			return &UseCaseResult[P]{
				StatusCode: r1.StatusCode,
				Success:    false,
				Error:      r1.Error,
				Details:    r1.Details,
			}
		}
		return next.uc.Execute(ctx, d.locale, *r1.Data)
	}
	return &DAG[I, P]{
		run:         run,
		ctx:         d.ctx,
		locale:      d.locale,
		_background: prevBackground,
		executor:    d.executor,
	}
}

// ThenBackground attaches a background step that consumes the current
// output type of the DAG.
//
// Background steps:
//   - Are executed only if the synchronous DAG execution succeeds.
//   - Do not affect the main DAG result.
//   - Are executed asynchronously using the DAG's BackgroundExecutor,
//     or synchronously if no executor is configured.
//
// The returned DAG:
//   - Preserves immutability by copying existing background steps.
//   - Keeps the same input and output types.
//
// Warnings:
//   - Do NOT use background steps when their result is required for
//     subsequent logic.
//   - Use background steps only for short, non-critical tasks.
//     For long-running or critical work, prefer queues or message brokers.
func ThenBackground[I any, O any, P any](d *DAG[I, O], next DagStep[O, P], name string) *DAG[I, O] {
	// We create a wrapper that executes 'next' but ignores its output/result. Any errors are logged.
	entr := backgroundEntry[O]{
		name: name,
		fn: func(ctx *app_context.AppContext, out O) {
			res := next.uc.Execute(ctx, d.locale, out)
			if res == nil {
				log.Printf("background %s returned nil result", name)
				return
			}
			if res.HasError() {
				log.Printf("background %s returned error: %v", name, res.Error)
			}
		},
	}

	// copy existing background slice to keep immutability
	newBg := make([]backgroundEntry[O], len(d._background), len(d._background)+1)
	copy(newBg, d._background)
	newBg = append(newBg, entr)

	// return new DAG with copied background list
	return &DAG[I, O]{
		run:         d.run,
		ctx:         d.ctx,
		locale:      d.locale,
		_background: newBg,
		executor:    d.executor,
	}
}

// Execute runs the DAG synchronously and schedules any attached background
// tasks using fire-and-forget semantics.
//
// Execution flow:
//   - The DAG is executed synchronously via the configured run function.
//   - If execution returns nil, an error, or no data, background tasks are
//     not scheduled.
//   - If execution succeeds and background steps are present, each background
//     task is scheduled for asynchronous execution.
//
// Background execution semantics:
//   - If a BackgroundExecutor is configured, background tasks are submitted
//     to the executor and executed under its managed context.
//   - If submission to the executor fails (for example, due to a full queue),
//     the task falls back to a standalone goroutine.
//   - If no executor is configured, background tasks are executed in
//     dedicated goroutines spawned by this method.
//
// Safety and cancellation:
//   - Background tasks respect context cancellation before starting execution.
//   - Panics in background tasks are recovered and logged to prevent crashing
//     the caller.
//   - Background task failures do not affect the synchronous DAG result.
//
// Concurrency and ordering:
//   - Background tasks are not ordered and may execute concurrently.
//   - This method does not wait for background tasks to complete.
//
// Returns:
//   - The UseCaseResult produced by the synchronous DAG execution.
func (d *DAG[I, O]) Execute(input I) *UseCaseResult[O] {
	if d.run == nil {
		return nil
	}
	res := d.run(d.ctx, input)
	if res == nil {
		return nil
	}
	if res.HasError() || res.Data == nil {
		return res
	}
	// schedule background tasks
	if len(d._background) == 0 {
		return res
	}

	out := *res.Data
	// submit tasks to executor if present
	if d.executor != nil {
		for _, be := range d._background {
			// capture be locally and the AppContext
			entry := be
			appCtx := d.ctx // Capture the AppContext to use in background task
			err := d.executor.Submit(func(ctx context.Context) {
				// respect ctx cancellation (using the executor's context)
				select {
				case <-ctx.Done():
					log.Printf("background %s cancelled before start\n", entry.name)
					return
				default:
				}
				// Use the captured AppContext for the background task
				// This allows the background task to access AppContext-specific data
				entry.fn(appCtx, out)
			})
			if err != nil {
				// fallback to fire-and-forget goroutine if queue is full
				log.Printf("executor submit failed for %s: %v — running in standalone goroutine", be.name, err)
				go func(entry backgroundEntry[O], o O) {
					defer func() {
						if r := recover(); r != nil {
							log.Printf("panic recovered in standalone background: %v\n%s", r, debug.Stack())
						}
					}()
					entry.fn(d.ctx, o)
				}(be, out)
			}
		}
		return res
	}

	// no executor — spawn goroutines per task but respecting ctx
	for _, be := range d._background {
		entry := be
		go func(entry backgroundEntry[O], o O, ctx *app_context.AppContext) {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("panic recovered in background goroutine: %v\n%s", r, debug.Stack())
				}
			}()
			// respect context
			select {
			case <-ctx.Done():
				log.Printf("background %s cancelled before start\n", entry.name)
				return
			default:
			}
			entry.fn(ctx, o)
		}(entry, out, d.ctx)
	}
	return res
}

// ExecuteWithBackground executes the DAG synchronously and then runs any
// registered background tasks associated with the DAG output.
//
// The method first invokes Execute(input) and immediately returns its result
// if one of the following conditions is met:
//   - the execution result is nil,
//   - the execution result contains an error,
//   - the execution result has no data,
//   - no background tasks are registered.
//
// Background execution semantics:
//   - If a background executor is configured, background tasks are assumed to
//     have been enqueued during Execute. In this case, this method optionally
//     waits for the executor queue to drain before returning.
//   - If waitTimeout is greater than zero, the wait is bounded by the provided
//     timeout. On timeout, the method returns the original execution result
//     without failing it.
//   - If waitTimeout is zero or negative, the method waits indefinitely for
//     the executor to finish processing queued background tasks.
//
// Fallback behavior (no executor):
//   - If no executor is configured, all background tasks are executed directly
//     in goroutines spawned by this method.
//   - The method blocks until all background tasks have completed.
//
// Safety and fault tolerance:
//   - Background task panics are recovered and logged to prevent crashing the
//     caller.
//   - Background task failures do not affect the main DAG execution result.
//
// Parameters:
//   - input: Input passed to the DAG execution.
//   - waitTimeout: Maximum duration to wait for background task completion when
//     using an executor. A zero or negative value means wait indefinitely.
//
// Returns:
//   - The UseCaseResult produced by the synchronous DAG execution. Background
//     task execution does not modify the returned result.
func (d *DAG[I, O]) ExecuteWithBackground(input I, waitTimeout time.Duration) *UseCaseResult[O] {
	res := d.Execute(input)
	if res == nil || res.HasError() || res.Data == nil || len(d._background) == 0 {
		return res
	}
	// if executor present, wait for its queued items to drain. Caller can pass a timeout.
	if d.executor != nil {
		if waitTimeout > 0 {
			// wait with timeout
			c := make(chan struct{})
			go func() {
				d.executor.Wait()
				close(c)
			}()
			select {
			case <-c:
				return res
			case <-time.After(waitTimeout):
				log.Printf("timeout waiting for background executor to finish")
				return res
			}
		}
		d.executor.Wait()
		return res
	}

	// no executor — run tasks and wait
	var wg sync.WaitGroup
	out := *res.Data
	for _, be := range d._background {
		entry := be
		wg.Add(1)
		go func(entry backgroundEntry[O], o O) {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					log.Printf("panic recovered in ExecuteWithBackground: %v\n%s", r, debug.Stack())
				}
			}()
			entry.fn(d.ctx, o)
		}(entry, out)
	}
	wg.Wait()
	return res
}

// UseCaseParallelDag is a parallel DAG. It contains the use cases to execute in parallel.
type UseCaseParallelDag[I any, O any] struct {
	Usecases []BaseUseCase[I, O]
}

var _ BaseUseCase[any, []any] = (*UseCaseParallelDag[any, any])(nil)

// Execute runs all configured use cases in parallel using the same input and
// aggregates their results into a single response.
//
// Each use case in d.Usecases is executed concurrently in its own goroutine.
// The order of the resulting outputs is preserved and matches the order of
// the use cases as defined in d.Usecases.
//
// Error handling semantics:
//   - If a use case returns a nil result, execution continues for the remaining
//     use cases, but the final result is marked as an internal error.
//   - If a use case returns an error result, execution continues for the
//     remaining use cases, but the first encountered error is propagated to the
//     final result.
//   - Subsequent errors are ignored once an error has been set in the final
//     result.
//
// Concurrency guarantees:
//   - A sync.WaitGroup is used to ensure all use cases complete before
//     returning.
//   - A mutex protects shared state (the aggregated result and output slice)
//     to prevent race conditions.
//   - The method blocks until all use cases have finished executing.
//
// On successful execution (no errors reported by any use case), the returned
// UseCaseResult contains a slice of outputs ([]O) with one element per use case
// and a success status.
//
// Parameters:
//   - ctx: Application context shared across all use cases (request-scoped data,
//     tracing, cancellation, etc.).
//   - locale: Locale used for localized behavior or messages inside each use case.
//   - input: Input passed unchanged to every use case.
//
// Returns:
//   - A UseCaseResult containing either the aggregated outputs on success or
//     the first error encountered during parallel execution.
func (d *UseCaseParallelDag[I, O]) Execute(
	ctx *app_context.AppContext,
	locale locales.LocaleTypeEnum,
	input I,
) *UseCaseResult[[]O] {
	result := NewUseCaseResult[[]O]()
	outputs := make([]O, len(d.Usecases))

	var wg sync.WaitGroup
	var mu sync.Mutex

	for i, uc := range d.Usecases {
		wg.Add(1)
		// capture i, uc
		go func(i int, uc BaseUseCase[I, O]) {
			defer wg.Done()
			res := uc.Execute(ctx, locale, input)
			if res == nil {
				mu.Lock()
				if !result.HasError() {
					result.SetError(status.InternalError, "nil result from usecase")
				}
				mu.Unlock()
				return
			}
			if res.HasError() {
				mu.Lock()
				if !result.HasError() {
					result.SetError(res.StatusCode, res.GetError().Error())
				}
				mu.Unlock()
				return
			}
			if res.Data != nil {
				mu.Lock()
				outputs[i] = *res.Data
				mu.Unlock()
			}
		}(i, uc)
	}

	wg.Wait()

	if !result.HasError() {
		result.SetData(
			status.Success,
			outputs,
			"All use cases executed successfully",
		)
	}
	return result
}

// SetLocale sets the locale for the UseCaseParallelDag.
func (d *UseCaseParallelDag[I, O]) SetLocale(locale locales.LocaleTypeEnum) {
	for _, uc := range d.Usecases {
		uc.SetLocale(locale)
	}
}

// SetAppContext sets the app context for the UseCaseParallelDag
func (d *UseCaseParallelDag[I, O]) SetAppContext(appContext *app_context.AppContext) {
	for _, uc := range d.Usecases {
		uc.SetAppContext(appContext)
	}
}

// NewUseCaseParallelDag creates a new UseCaseParallelDag.
func NewUseCaseParallelDag[I any, O any]() *UseCaseParallelDag[I, O] {
	return &UseCaseParallelDag[I, O]{}
}
