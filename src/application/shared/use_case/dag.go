package usecase

import (
	"context"
	"sync"

	"github.com/simon3640/goprojectskeleton/src/application/shared/locales"
	"github.com/simon3640/goprojectskeleton/src/application/shared/status"
)

type DagStep[I any, O any] struct {
	uc BaseUseCase[I, O]
}

func NewStep[I any, O any](uc BaseUseCase[I, O]) DagStep[I, O] {
	return DagStep[I, O]{uc: uc}
}

type DAG[I any, O any] struct {
	run            func(I) *UseCaseResult[O]
	ctx            context.Context
	locale         locales.LocaleTypeEnum
	backgroundRuns []func(O) // Funciones que se ejecutan en background
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
		ctx:            ctx,
		locale:         locale,
		backgroundRuns: []func(O){},
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
		ctx:            d.ctx,
		locale:         d.locale,
		backgroundRuns: []func(P){}, // Reset background runs cuando cambia el tipo de salida
	}
}

// ThenBackground agrega un paso que se ejecutará en background después de que se retorne la respuesta.
// El resultado del paso anterior se retorna inmediatamente, y el siguiente paso se ejecuta en background.
// Útil para operaciones como envío de emails, notificaciones, etc.
//
// Ejemplo:
//
//	dag := NewDag(NewStep(createUserUC), locale, ctx)
//	dag = ThenBackground(dag, NewStep(sendEmailUC))
//	result := dag.Execute(input) // Retorna inmediatamente, sendEmailUC se ejecuta en background
func ThenBackground[I any, O any, P any](d *DAG[I, O], next DagStep[O, P]) *DAG[I, O] {
	// Guardar la función de ejecución en background
	backgroundRun := func(output O) {
		// Ejecutar el siguiente use case en background
		next.uc.Execute(
			d.ctx,
			d.locale,
			output,
		)
		// Nota: Los errores en background se pueden loggear pero no afectan la respuesta
	}

	return &DAG[I, O]{
		run: func(input I) *UseCaseResult[O] {
			// Ejecutar hasta el breakpoint
			result := d.run(input)

			// Si hay error, no ejecutar background
			if result.HasError() {
				return result
			}

			// Ejecutar tareas en background de forma asíncrona
			go func() {
				// Ejecutar todas las tareas en background acumuladas
				for _, bgRun := range d.backgroundRuns {
					bgRun(*result.Data)
				}
				// Ejecutar la nueva tarea en background
				backgroundRun(*result.Data)
			}()

			// Retornar inmediatamente sin esperar las tareas en background
			return result
		},
		ctx:            d.ctx,
		locale:         d.locale,
		backgroundRuns: append(d.backgroundRuns, backgroundRun),
	}
}

// Execute ejecuta el DAG y retorna el resultado inmediatamente.
// Si hay tareas en background, se ejecutan de forma asíncrona.
func (d *DAG[I, O]) Execute(input I) *UseCaseResult[O] {
	if d.run == nil {
		return nil
	}
	return d.run(input)
}

// ExecuteWithBackground ejecuta el DAG y espera a que todas las tareas en background terminen.
// Útil para testing o cuando necesitas asegurar que las tareas en background se completaron.
func (d *DAG[I, O]) ExecuteWithBackground(input I) *UseCaseResult[O] {
	if d.run == nil {
		return nil
	}

	result := d.run(input)

	// Si hay tareas en background, esperarlas
	if len(d.backgroundRuns) > 0 && result.Data != nil && !result.HasError() {
		var wg sync.WaitGroup
		for _, bgRun := range d.backgroundRuns {
			wg.Add(1)
			go func(bg func(O)) {
				defer wg.Done()
				bg(*result.Data)
			}(bgRun)
		}
		wg.Wait()
	}

	return result
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
