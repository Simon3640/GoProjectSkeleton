package usecase

import (
	"context"

	"gormgoskeleton/src/application/shared/locales"
)

type DagStep[I any, O any] struct {
	uc BaseUseCase[I, O]
}

func NewStep[I any, O any](uc BaseUseCase[I, O]) DagStep[I, O] {
	return DagStep[I, O]{uc: uc}
}

// func Parallel(steps ...DagStep) DagStep {
// 	return DagStep{run: func(input any) (*any, error) {
// 		var wg sync.WaitGroup
// 		results := make([]*any, len(steps))
// 		errors := make(chan error, len(steps))
// 		for i, step := range steps {
// 			wg.Add(1)
// 			go func(i int, step DagStep) {
// 				defer wg.Done()
// 				result, err := step.run(input)
// 				if err != nil {
// 					errors <- err
// 					return
// 				}
// 				results[i] = result
// 			}(i, step)
// 		}
// 		wg.Wait()
// 		close(errors)
// 		for err := range errors {
// 			if err != nil {
// 				return nil, err
// 			}
// 		}
// 		// Convert []*any to []any, then take its address to match *any
// 		finalResults := make([]any, len(results))
// 		for i, r := range results {
// 			if r != nil {
// 				finalResults[i] = *r
// 			} else {
// 				finalResults[i] = nil
// 			}
// 		}
// 		return any(&finalResults).(*any), nil
// 	}}
// }

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
			return next.uc.Execute(
				d.ctx,
				d.locale,
				*r1.Data,
			)
		},
	}
}

func (d *DAG[I, O]) Execute(input I) (*UseCaseResult[O], error) {
	if d.run == nil {
		return nil, nil
	}
	result := d.run(input)
	if result.HasError() {
		return nil, result.GetError()
	}
	return result, nil
}
