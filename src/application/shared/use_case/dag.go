package usecase

import (
	"context"
	"sync"

	"gormgoskeleton/src/application/shared/locales"
	"gormgoskeleton/src/application/shared/status"
)

type DagStep[I any, O any] struct {
	uc BaseUseCase[I, O]
}

func NewStep[I any, O any](uc BaseUseCase[I, O]) DagStep[I, O] {
	return DagStep[I, O]{uc: uc}
}

type DAG[I any, O any] struct {
	run    func(I) *UseCaseResult[O]
	ctx    context.Context
	locale locales.LocaleTypeEnum
}

func NewDag[I any, O any](first DagStep[I, O], locale locales.LocaleTypeEnum, ctx context.Context) *DAG[I, O] {
	return &DAG[I, O]{
		run: func(input I) *UseCaseResult[O] {
			return first.uc.Execute(
				ctx,
				locale,
				input,
			)
		},
		ctx:    ctx,
		locale: locale,
	}
}

func Then[I any, O any, P any](d *DAG[I, O], next DagStep[O, P]) *DAG[I, P] {
	return &DAG[I, P]{
		run: func(input I) *UseCaseResult[P] {
			r1 := d.run(input)
			// TODO: Best error control must be done
			if r1.HasError() {
				return &UseCaseResult[P]{
					StatusCode: r1.StatusCode,
					Success:    false,
					Error:      r1.Error,
					Details:    r1.Details,
				}
			}
			return next.uc.Execute(
				d.ctx,
				d.locale,
				*r1.Data,
			)
		},
	}
}

func (d *DAG[I, O]) Execute(input I) *UseCaseResult[O] {
	if d.run == nil {
		return nil
	}
	return d.run(input)
}

type UseCaseParallelDag[I any, O any] struct {
	Usecases []BaseUseCase[I, O]
}

var _ BaseUseCase[any, []any] = (*UseCaseParallelDag[any, any])(nil)

// Execute for pipe purposes, input is an array of use
// cases to be executed in parallel and a list of outputs for next step
func (d *UseCaseParallelDag[I, O]) Execute(
	ctx context.Context,
	locale locales.LocaleTypeEnum,
	input I,
) *UseCaseResult[[]O] {
	result := NewUseCaseResult[[]O]()
	outputs := make([]O, len(d.Usecases))

	var wg sync.WaitGroup
	var mu sync.Mutex

	for i, uc := range d.Usecases {
		go func(i int, uc BaseUseCase[I, O]) {
			defer wg.Done()
			res := uc.Execute(ctx, locale, input)
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
		wg.Add(1)
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

func (d *UseCaseParallelDag[I, O]) SetLocale(locale locales.LocaleTypeEnum) {
}

func NewUseCaseParallelDag[I any, O any]() *UseCaseParallelDag[I, O] {
	return &UseCaseParallelDag[I, O]{}
}
