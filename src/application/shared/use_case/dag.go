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

// DagStep is a step in the DAG. It contains the use case to execute.
type DagStep[I any, O any] struct {
	uc BaseUseCase[I, O]
}

// NewStep creates a new DagStep.
func NewStep[I any, O any](uc BaseUseCase[I, O]) DagStep[I, O] { return DagStep[I, O]{uc: uc} }

// backgroundEntry captures the function and an optional name for tracing
type backgroundEntry[O any] struct {
	name string
	fn   func(ctx *app_context.AppContext, out O)
}

// DAG is immutable: builder functions return a new DAG instance.
type DAG[I any, O any] struct {
	run         func(appContext *app_context.AppContext, input I) *UseCaseResult[O]
	ctx         *app_context.AppContext
	locale      locales.LocaleTypeEnum
	_background []backgroundEntry[O] // internal; treat as immutable — copy on write
	executor    *workers.BackgroundExecutor
}

// NewDag creates a new DAG starting with `first` use case. The returned dag is safe to reuse.
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

// Then chains a synchronous step. It returns a NEW DAG with the output type changed.
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

// ThenBackground attaches a background step that consumes the current output type O.
// The background step will be executed asynchronously via the DAG's BackgroundExecutor.
// The returned DAG preserves previous background steps (copied) and appends the new one.
// Warning: Do not use background tasks if you need to wait for the result of the background task.
// Warning: Use background for small tasks, otherwise consider using other background solutions like a queue or a message broker.
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

// Execute runs the DAG synchronously and schedules background tasks (fire-and-forget).
// If an executor is configured tasks will be submitted there. Otherwise they run in dedicated goroutines
// that respect context cancellation and include panic recovery.
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

// ExecuteWithBackground runs the DAG and waits for all background tasks to finish.
// If an executor exists we call Wait() on it. Otherwise we create a temporary WaitGroup and run each task.
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

func (d *UseCaseParallelDag[I, O]) SetAppContext(appContext *app_context.AppContext) {
	for _, uc := range d.Usecases {
		uc.SetAppContext(appContext)
	}
}

// NewUseCaseParallelDag creates a new UseCaseParallelDag.
func NewUseCaseParallelDag[I any, O any]() *UseCaseParallelDag[I, O] {
	return &UseCaseParallelDag[I, O]{}
}
