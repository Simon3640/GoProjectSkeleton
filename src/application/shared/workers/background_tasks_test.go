package workers

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBackgroundExecutorSingleton(t *testing.T) {
	assert := assert.New(t)

	// Resetear el singleton antes de empezar
	ResetBackgroundExecutorSingleton()

	ctx := context.Background()

	// Inicializar el singleton
	InitializeBackgroundExecutor(ctx, 4, 100)

	// Obtener el singleton dos veces
	executor1 := GetBackgroundExecutor()
	executor2 := GetBackgroundExecutor()

	// Deberían ser la misma instancia
	assert.Same(executor1, executor2, "Singleton should return the same instance")
	assert.NotNil(executor1, "Executor should not be nil")

	// Limpiar
	ResetBackgroundExecutorSingleton()
}

func TestBackgroundExecutorSingletonThreadSafe(t *testing.T) {
	assert := assert.New(t)

	// Resetear el singleton antes de empezar
	ResetBackgroundExecutorSingleton()

	ctx := context.Background()

	// Inicializar el singleton primero
	InitializeBackgroundExecutor(ctx, 4, 100)

	var executors []*BackgroundExecutor
	var mu sync.Mutex
	var wg sync.WaitGroup

	// Crear múltiples goroutines que intentan obtener el singleton simultáneamente
	numGoroutines := 100
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			executor := GetBackgroundExecutor()
			mu.Lock()
			executors = append(executors, executor)
			mu.Unlock()
		}()
	}

	wg.Wait()

	// Todas las instancias deberían ser la misma
	assert.Equal(numGoroutines, len(executors))
	firstExecutor := executors[0]
	for _, executor := range executors {
		assert.Same(firstExecutor, executor, "All goroutines should get the same singleton instance")
	}

	// Limpiar
	ResetBackgroundExecutorSingleton()
}

func TestBackgroundExecutorSingletonIgnoresSubsequentParams(t *testing.T) {
	assert := assert.New(t)

	// Resetear el singleton antes de empezar
	ResetBackgroundExecutorSingleton()

	ctx := context.Background()

	// Primera inicialización con workers=4, queueSize=100
	InitializeBackgroundExecutor(ctx, 4, 100)
	executor1 := GetBackgroundExecutor()

	// Intentar inicializar de nuevo con diferentes parámetros - debería ignorarse
	InitializeBackgroundExecutor(ctx, 8, 200)
	executor2 := GetBackgroundExecutor()

	// Deberían ser la misma instancia
	assert.Same(executor1, executor2, "Singleton should ignore subsequent initialization parameters")

	// Limpiar
	ResetBackgroundExecutorSingleton()
}

func TestBackgroundExecutorSingletonReset(t *testing.T) {
	assert := assert.New(t)

	ctx := context.Background()

	// Inicializar el singleton
	InitializeBackgroundExecutor(ctx, 4, 100)
	executor1 := GetBackgroundExecutor()

	// Resetear
	ResetBackgroundExecutorSingleton()

	// Inicializar de nuevo - debería ser una nueva instancia
	InitializeBackgroundExecutor(ctx, 4, 100)
	executor2 := GetBackgroundExecutor()

	// Deberían ser diferentes instancias después del reset
	assert.NotSame(executor1, executor2, "After reset, should get a new instance")

	// Limpiar
	ResetBackgroundExecutorSingleton()
}

func TestBackgroundExecutorSingletonUsage(t *testing.T) {
	assert := assert.New(t)

	// Resetear el singleton antes de empezar
	ResetBackgroundExecutorSingleton()

	ctx := context.Background()

	// Inicializar el singleton
	InitializeBackgroundExecutor(ctx, 2, 50)
	executor := GetBackgroundExecutor()
	defer executor.Stop()

	// Ejecutar algunas tareas
	var wg sync.WaitGroup
	var executedTasks int
	var mu sync.Mutex

	numTasks := 10
	wg.Add(numTasks)

	for i := 0; i < numTasks; i++ {
		taskID := i
		err := executor.Submit(func(_ context.Context) {
			defer wg.Done()
			time.Sleep(10 * time.Millisecond)
			mu.Lock()
			executedTasks++
			mu.Unlock()
		})
		assert.NoError(err, "Should be able to submit task %d", taskID)
	}

	wg.Wait()

	// Verificar que todas las tareas se ejecutaron
	assert.Equal(numTasks, executedTasks, "All tasks should have been executed")

	// Limpiar
	ResetBackgroundExecutorSingleton()
}
