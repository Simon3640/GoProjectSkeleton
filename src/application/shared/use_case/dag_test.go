package usecase

import (
	"context"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/simon3640/goprojectskeleton/src/application/shared/locales"
	"github.com/simon3640/goprojectskeleton/src/application/shared/status"
	"github.com/simon3640/goprojectskeleton/src/application/shared/workers"

	"github.com/stretchr/testify/assert"
)

type UCStringToInt struct{}

func (uc *UCStringToInt) SetLocale(_ locales.LocaleTypeEnum) {}
func (uc *UCStringToInt) Execute(_ context.Context, _ locales.LocaleTypeEnum, input string) *UseCaseResult[int] {
	result := NewUseCaseResult[int]()
	intValue, err := strconv.Atoi(input)
	if err != nil {
		result.SetError(status.InternalError, err.Error())
		return result
	}
	result.SetData(status.Success, intValue, "Converted string to int")
	return result
}

var _ BaseUseCase[string, int] = (*UCStringToInt)(nil)

type UCIntExponent struct{}

func (uc *UCIntExponent) SetLocale(_ locales.LocaleTypeEnum) {}
func (uc *UCIntExponent) Execute(_ context.Context, _ locales.LocaleTypeEnum, input int) *UseCaseResult[int] {
	result := NewUseCaseResult[int]()
	result.SetData(status.Success, input*input, "Calculated exponent")
	return result
}

type UCIntToString struct{}

func (uc *UCIntToString) SetLocale(_ locales.LocaleTypeEnum) {}
func (uc *UCIntToString) Execute(_ context.Context, _ locales.LocaleTypeEnum, input int) *UseCaseResult[string] {
	result := NewUseCaseResult[string]()
	result.SetData(status.Success, strconv.Itoa(input), "Converted int to string")
	return result
}

// Use case para pruebas de background tasks con int
type UCLogBackgroundInt struct {
	mu       sync.Mutex
	logged   []int
	executed bool
}

func (uc *UCLogBackgroundInt) SetLocale(_ locales.LocaleTypeEnum) {}
func (uc *UCLogBackgroundInt) Execute(_ context.Context, _ locales.LocaleTypeEnum, input int) *UseCaseResult[int] {
	uc.mu.Lock()
	defer uc.mu.Unlock()
	uc.logged = append(uc.logged, input)
	uc.executed = true
	result := NewUseCaseResult[int]()
	result.SetData(status.Success, input, "Background task executed")
	return result
}

func (uc *UCLogBackgroundInt) GetLogged() []int {
	uc.mu.Lock()
	defer uc.mu.Unlock()
	return uc.logged
}

func (uc *UCLogBackgroundInt) WasExecuted() bool {
	uc.mu.Lock()
	defer uc.mu.Unlock()
	return uc.executed
}

func (uc *UCLogBackgroundInt) Reset() {
	uc.mu.Lock()
	defer uc.mu.Unlock()
	uc.logged = []int{}
	uc.executed = false
}

var _ BaseUseCase[int, int] = (*UCLogBackgroundInt)(nil)

// Use case para pruebas de background tasks con string
type UCLogBackgroundString struct {
	mu       sync.Mutex
	logged   []string
	executed bool
}

func (uc *UCLogBackgroundString) SetLocale(_ locales.LocaleTypeEnum) {}
func (uc *UCLogBackgroundString) Execute(_ context.Context, _ locales.LocaleTypeEnum, input string) *UseCaseResult[string] {
	uc.mu.Lock()
	defer uc.mu.Unlock()
	uc.logged = append(uc.logged, input)
	uc.executed = true
	result := NewUseCaseResult[string]()
	result.SetData(status.Success, "logged: "+input, "Background task executed")
	return result
}

func (uc *UCLogBackgroundString) GetLogged() []string {
	uc.mu.Lock()
	defer uc.mu.Unlock()
	return uc.logged
}

func (uc *UCLogBackgroundString) WasExecuted() bool {
	uc.mu.Lock()
	defer uc.mu.Unlock()
	return uc.executed
}

func (uc *UCLogBackgroundString) Reset() {
	uc.mu.Lock()
	defer uc.mu.Unlock()
	uc.logged = []string{}
	uc.executed = false
}

var _ BaseUseCase[string, string] = (*UCLogBackgroundString)(nil)

func TestDagExecution(t *testing.T) {
	assert := assert.New(t)

	ctx := context.Background()
	executor := workers.NewBackgroundExecutor(ctx, 4, 100)
	executor.Start()
	defer executor.Stop()

	UC1 := &UCStringToInt{}
	UC2 := &UCIntExponent{}
	UC3 := &UCIntToString{}

	dag := NewDag(ctx, NewStep(UC1), locales.EN_US, executor)
	dag2 := Then(dag, NewStep(UC2))
	dag3 := Then(dag2, NewStep(UC3))

	input := "5"
	result := dag3.Execute(input)
	assert.NotNil(result)
	assert.True(result.IsSuccess())
	assert.NotNil(result.Data)
	assert.Equal("84", *result.Data)
}

func TestDagConcurrentExecution(t *testing.T) {
	assert := assert.New(t)

	ctx := context.Background()
	executor := workers.NewBackgroundExecutor(ctx, 4, 100)
	executor.Start()
	defer executor.Stop()

	UC1 := &UCStringToInt{}
	UC2 := &UCIntExponent{}
	UC3 := &UCIntToString{}
	ParallelUC := NewUseCaseParallelDag[string, int]()
	ParallelUC.Usecases = []BaseUseCase[string, int]{UC1, UC1, UC1, UC1, UC1}
	dag := NewDag(ctx, NewStep(UC1), locales.EN_US, executor)
	dag2 := Then(dag, NewStep(UC2))
	dag3 := Then(dag2, NewStep(UC3))
	dagParallel := Then(dag3, NewStep(ParallelUC))

	input := "5"
	result := dagParallel.Execute(input)
	assert.NotNil(result)
	assert.True(result.IsSuccess())
	assert.NotNil(result.Data)
	assert.Equal(5, len(*result.Data))
	for _, val := range *result.Data {
		assert.Equal(42, val)
	}
}

func TestDagWithBackgroundTask(t *testing.T) {
	assert := assert.New(t)

	ctx := context.Background()
	executor := workers.NewBackgroundExecutor(ctx, 4, 100)
	executor.Start()
	defer executor.Stop()

	UC1 := &UCStringToInt{}
	UC2 := &UCIntExponent{}
	backgroundUC := &UCLogBackgroundInt{}

	dag := NewDag(ctx, NewStep(UC1), locales.EN_US, executor)
	dag2 := Then(dag, NewStep(UC2))
	dagWithBackground := ThenBackground(dag2, NewStep(backgroundUC), "log-background")

	input := "5"
	result := dagWithBackground.Execute(input)

	assert.NotNil(result)
	assert.True(result.IsSuccess())
	assert.NotNil(result.Data)
	assert.Equal(84, *result.Data)

	// Esperar un poco para que la tarea en background se ejecute
	time.Sleep(100 * time.Millisecond)

	// Verificar que la tarea en background se ejecutó
	assert.True(backgroundUC.WasExecuted())
	logged := backgroundUC.GetLogged()
	assert.Equal(1, len(logged))
	assert.Equal(84, logged[0])
}

func TestDagWithMultipleBackgroundTasks(t *testing.T) {
	assert := assert.New(t)

	ctx := context.Background()
	executor := workers.NewBackgroundExecutor(ctx, 4, 100)
	executor.Start()
	defer executor.Stop()

	UC1 := &UCStringToInt{}
	UC2 := &UCIntExponent{}
	backgroundUC1 := &UCLogBackgroundInt{}
	backgroundUC2 := &UCLogBackgroundInt{}

	dag := NewDag(ctx, NewStep(UC1), locales.EN_US, executor)
	dag2 := Then(dag, NewStep(UC2))
	dagWithBg1 := ThenBackground(dag2, NewStep(backgroundUC1), "log-background-1")
	dagWithBg2 := ThenBackground(dagWithBg1, NewStep(backgroundUC2), "log-background-2")

	input := "5"
	result := dagWithBg2.Execute(input)

	assert.NotNil(result)
	assert.True(result.IsSuccess())
	assert.NotNil(result.Data)
	assert.Equal(84, *result.Data)

	// Esperar un poco para que las tareas en background se ejecuten
	time.Sleep(200 * time.Millisecond)

	// Verificar que ambas tareas en background se ejecutaron
	assert.True(backgroundUC1.WasExecuted())
	assert.True(backgroundUC2.WasExecuted())

	logged1 := backgroundUC1.GetLogged()
	logged2 := backgroundUC2.GetLogged()

	assert.Equal(1, len(logged1))
	assert.Equal(1, len(logged2))
	assert.Equal(84, logged1[0])
	assert.Equal(84, logged2[0])
}

func TestDagExecuteWithBackground(t *testing.T) {
	assert := assert.New(t)

	ctx := context.Background()
	executor := workers.NewBackgroundExecutor(ctx, 4, 100)
	executor.Start()
	defer executor.Stop()

	UC1 := &UCStringToInt{}
	UC2 := &UCIntExponent{}
	backgroundUC := &UCLogBackgroundInt{}

	dag := NewDag(ctx, NewStep(UC1), locales.EN_US, executor)
	dag2 := Then(dag, NewStep(UC2))
	dagWithBackground := ThenBackground(dag2, NewStep(backgroundUC), "log-background")

	input := "5"
	result := dagWithBackground.ExecuteWithBackground(input, 1*time.Second)

	assert.NotNil(result)
	assert.True(result.IsSuccess())
	assert.NotNil(result.Data)
	assert.Equal(84, *result.Data)

	// Con ExecuteWithBackground, las tareas deberían haber terminado
	assert.True(backgroundUC.WasExecuted())
	logged := backgroundUC.GetLogged()
	assert.Equal(1, len(logged))
	assert.Equal(84, logged[0])
}

func TestDagExecuteWithBackgroundTimeout(t *testing.T) {
	assert := assert.New(t)

	ctx := context.Background()
	executor := workers.NewBackgroundExecutor(ctx, 2, 10)
	executor.Start()
	defer executor.Stop()

	UC1 := &UCStringToInt{}
	UC2 := &UCIntExponent{}
	backgroundUC := &UCLogBackgroundInt{}

	dag := NewDag(ctx, NewStep(UC1), locales.EN_US, executor)
	dag2 := Then(dag, NewStep(UC2))
	dagWithBackground := ThenBackground(dag2, NewStep(backgroundUC), "log-background")

	input := "5"
	// Usar un timeout muy corto para probar el comportamiento
	result := dagWithBackground.ExecuteWithBackground(input, 1*time.Millisecond)

	assert.NotNil(result)
	assert.True(result.IsSuccess())
	assert.NotNil(result.Data)
	assert.Equal(84, *result.Data)
}

func TestDagWithoutExecutor(t *testing.T) {
	assert := assert.New(t)

	ctx := context.Background()
	backgroundUC := &UCLogBackgroundInt{}

	// Crear DAG sin executor (nil)
	UC1 := &UCStringToInt{}
	dag := NewDag(ctx, NewStep(UC1), locales.EN_US, nil)
	dagWithBackground := ThenBackground(dag, NewStep(backgroundUC), "log-background")

	input := "5"
	result := dagWithBackground.Execute(input)

	assert.NotNil(result)
	assert.True(result.IsSuccess())
	assert.NotNil(result.Data)
	assert.Equal(42, *result.Data)

	// Esperar un poco para que la goroutine se ejecute
	time.Sleep(100 * time.Millisecond)

	// Verificar que la tarea en background se ejecutó (aunque sin executor)
	assert.True(backgroundUC.WasExecuted())
	logged := backgroundUC.GetLogged()
	assert.Equal(1, len(logged))
	assert.Equal(42, logged[0])
}

func TestDagWithBackgroundChain(t *testing.T) {
	assert := assert.New(t)

	ctx := context.Background()
	executor := workers.NewBackgroundExecutor(ctx, 4, 100)
	executor.Start()
	defer executor.Stop()

	UC1 := &UCStringToInt{}
	UC2 := &UCIntExponent{}
	UC3 := &UCIntToString{}
	backgroundUC := &UCLogBackgroundString{}

	// Crear una cadena: UC1 -> UC2 -> UC3, y luego agregar background task
	dag := NewDag(ctx, NewStep(UC1), locales.EN_US, executor)
	dag2 := Then(dag, NewStep(UC2))
	dag3 := Then(dag2, NewStep(UC3))
	dagWithBackground := ThenBackground(dag3, NewStep(backgroundUC), "log-background")

	input := "5"
	result := dagWithBackground.Execute(input)

	assert.NotNil(result)
	assert.True(result.IsSuccess())
	assert.NotNil(result.Data)
	assert.Equal("84", *result.Data)

	// Esperar un poco para que la tarea en background se ejecute
	time.Sleep(100 * time.Millisecond)

	// Verificar que la tarea en background recibió el string final
	assert.True(backgroundUC.WasExecuted())
	logged := backgroundUC.GetLogged()
	assert.Equal(1, len(logged))
	assert.Equal("84", logged[0])
}
