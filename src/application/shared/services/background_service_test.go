package services

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	app_context "github.com/simon3640/goprojectskeleton/src/application/shared/context"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales"
	"github.com/simon3640/goprojectskeleton/src/application/shared/workers"

	"github.com/stretchr/testify/assert"
)

// mockBackgroundService is a mock implementation of BackgroundService for testing
type mockBackgroundService struct {
	executeFunc func(ctx *app_context.AppContext, locale locales.LocaleTypeEnum, input string) error
	name        string
	executed    bool
	mu          sync.Mutex
}

func (m *mockBackgroundService) Execute(ctx *app_context.AppContext, locale locales.LocaleTypeEnum, input string) error {
	m.mu.Lock()
	m.executed = true
	m.mu.Unlock()
	if m.executeFunc != nil {
		return m.executeFunc(ctx, locale, input)
	}
	return nil
}

func (m *mockBackgroundService) Name() string {
	return m.name
}

func (m *mockBackgroundService) WasExecuted() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.executed
}

func (m *mockBackgroundService) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.executed = false
}

// mockExecutor is a mock implementation of BackgroundExecutorInterface for testing
type mockExecutor struct {
	submitFunc func(task func(ctx context.Context)) error
	submitted  bool
	mu         sync.Mutex
}

func (m *mockExecutor) Submit(task func(ctx context.Context)) error {
	m.mu.Lock()
	m.submitted = true
	m.mu.Unlock()
	if m.submitFunc != nil {
		return m.submitFunc(task)
	}
	// Execute the task immediately for testing
	go task(context.Background())
	return nil
}

func (m *mockExecutor) WasSubmitted() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.submitted
}

func (m *mockExecutor) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.submitted = false
}

func TestNewBackgroundExecutorAdapter(t *testing.T) {
	assert := assert.New(t)

	t.Run("returns nil when executor is nil", func(_ *testing.T) {
		adapter := NewBackgroundExecutorAdapter(nil)
		assert.Nil(adapter, "Adapter should be nil when executor is nil")
	})

	t.Run("returns adapter when executor is not nil", func(_ *testing.T) {
		ctx := context.Background()
		executor := workers.NewBackgroundExecutor(ctx, 2, 10)
		defer executor.Stop()

		adapter := NewBackgroundExecutorAdapter(executor)
		assert.NotNil(adapter, "Adapter should not be nil when executor is not nil")
	})
}

func TestBackgroundExecutorAdapter_Submit(t *testing.T) {
	assert := assert.New(t)

	ctx := context.Background()
	executor := workers.NewBackgroundExecutor(ctx, 2, 10)
	executor.Start()
	defer executor.Stop()

	adapter := NewBackgroundExecutorAdapter(executor)
	assert.NotNil(adapter)

	t.Run("submits task successfully", func(_ *testing.T) {
		var executed bool
		var wg sync.WaitGroup
		wg.Add(1)

		task := func(_ context.Context) {
			defer wg.Done()
			executed = true
		}

		err := adapter.Submit(task)
		assert.NoError(err, "Should submit task without error")

		wg.Wait()
		assert.True(executed, "Task should have been executed")
	})

	t.Run("handles multiple tasks", func(_ *testing.T) {
		var executedCount int
		var mu sync.Mutex
		var wg sync.WaitGroup
		numTasks := 5
		wg.Add(numTasks)

		for i := 0; i < numTasks; i++ {
			taskID := i
			task := func(_ context.Context) {
				defer wg.Done()
				mu.Lock()
				executedCount++
				mu.Unlock()
			}
			err := adapter.Submit(task)
			assert.NoError(err, "Should submit task %d without error", taskID)
		}

		wg.Wait()
		assert.Equal(numTasks, executedCount, "All tasks should have been executed")
	})
}

func TestNewBackgroundServiceFactory(t *testing.T) {
	assert := assert.New(t)

	t.Run("creates factory with executor", func(_ *testing.T) {
		mockExec := &mockExecutor{}
		factory := NewBackgroundServiceFactory(mockExec)

		assert.NotNil(factory, "Factory should not be nil")
		assert.Equal(mockExec, factory.executor, "Factory should have the provided executor")
	})

	t.Run("creates factory with nil executor", func(_ *testing.T) {
		factory := NewBackgroundServiceFactory(nil)

		assert.NotNil(factory, "Factory should not be nil even with nil executor")
		assert.Nil(factory.executor, "Factory executor should be nil")
	})
}

func TestExecuteService(t *testing.T) {
	assert := assert.New(t)

	ctx := &app_context.AppContext{Context: context.Background()}
	locale := locales.EN_US
	input := "test-input"

	t.Run("executes service when factory is nil", func(_ *testing.T) {
		service := &mockBackgroundService{name: "test-service"}
		service.Reset()

		err := ExecuteService(nil, service, ctx, locale, input)
		assert.NoError(err, "Should not return error even with nil factory")

		// Wait a bit for goroutine to execute
		time.Sleep(50 * time.Millisecond)
		assert.True(service.WasExecuted(), "Service should have been executed in goroutine")
	})

	t.Run("executes service when executor is nil", func(_ *testing.T) {
		factory := NewBackgroundServiceFactory(nil)
		service := &mockBackgroundService{name: "test-service"}
		service.Reset()

		err := ExecuteService(factory, service, ctx, locale, input)
		assert.NoError(err, "Should not return error even with nil executor")

		// Wait a bit for goroutine to execute
		time.Sleep(50 * time.Millisecond)
		assert.True(service.WasExecuted(), "Service should have been executed in goroutine")
	})

	t.Run("executes service through executor", func(_ *testing.T) {
		mockExec := &mockExecutor{}
		mockExec.Reset()
		factory := NewBackgroundServiceFactory(mockExec)
		service := &mockBackgroundService{name: "test-service"}
		service.Reset()

		err := ExecuteService(factory, service, ctx, locale, input)
		assert.NoError(err, "Should not return error")
		assert.True(mockExec.WasSubmitted(), "Task should have been submitted to executor")

		// Wait a bit for task to execute
		time.Sleep(50 * time.Millisecond)
		assert.True(service.WasExecuted(), "Service should have been executed")
	})

	t.Run("handles service execution error", func(_ *testing.T) {
		mockExec := &mockExecutor{}
		factory := NewBackgroundServiceFactory(mockExec)
		expectedError := errors.New("service error")
		service := &mockBackgroundService{
			name: "error-service",
			executeFunc: func(_ *app_context.AppContext, _ locales.LocaleTypeEnum, _ string) error {
				return expectedError
			},
		}

		err := ExecuteService(factory, service, ctx, locale, input)
		assert.NoError(err, "ExecuteService should not return error (fire-and-forget)")

		// Wait a bit for task to execute
		time.Sleep(50 * time.Millisecond)
		assert.True(service.WasExecuted(), "Service should have been executed even with error")
	})

	t.Run("preserves context and input", func(_ *testing.T) {
		mockExec := &mockExecutor{}
		factory := NewBackgroundServiceFactory(mockExec)
		var capturedCtx *app_context.AppContext
		var capturedInput string

		service := &mockBackgroundService{
			name: "capture-service",
			executeFunc: func(ctx *app_context.AppContext, _ locales.LocaleTypeEnum, input string) error {
				capturedCtx = ctx
				capturedInput = input
				return nil
			},
		}

		err := ExecuteService(factory, service, ctx, locale, input)
		assert.NoError(err)

		// Wait a bit for task to execute
		time.Sleep(50 * time.Millisecond)
		assert.Equal(ctx, capturedCtx, "Context should be preserved")
		assert.Equal(input, capturedInput, "Input should be preserved")
	})

	t.Run("handles executor submit error", func(_ *testing.T) {
		submitError := errors.New("queue full")
		mockExec := &mockExecutor{
			submitFunc: func(_ func(ctx context.Context)) error {
				return submitError
			},
		}
		factory := NewBackgroundServiceFactory(mockExec)
		service := &mockBackgroundService{name: "test-service"}

		err := ExecuteService(factory, service, ctx, locale, input)
		assert.Error(err, "Should return error when executor submit fails")
		assert.Equal(submitError, err, "Should return the executor submit error")
	})

	t.Run("executes multiple services concurrently", func(_ *testing.T) {
		ctx := context.Background()
		executor := workers.NewBackgroundExecutor(ctx, 4, 100)
		executor.Start()
		defer executor.Stop()

		adapter := NewBackgroundExecutorAdapter(executor)
		factory := NewBackgroundServiceFactory(adapter)

		appCtx := &app_context.AppContext{Context: context.Background()}
		var wg sync.WaitGroup
		numServices := 10
		executedCount := 0
		var mu sync.Mutex

		for i := 0; i < numServices; i++ {
			wg.Add(1)
			service := &mockBackgroundService{
				name: "concurrent-service",
				executeFunc: func(_ *app_context.AppContext, _ locales.LocaleTypeEnum, _ string) error {
					defer wg.Done()
					mu.Lock()
					executedCount++
					mu.Unlock()
					return nil
				},
			}

			err := ExecuteService(factory, service, appCtx, locale, input)
			assert.NoError(err, "Should submit service %d without error", i)
		}

		wg.Wait()
		assert.Equal(numServices, executedCount, "All services should have been executed")
	})
}

func TestBackgroundServiceFactory_Integration(t *testing.T) {
	assert := assert.New(t)

	ctx := context.Background()
	executor := workers.NewBackgroundExecutor(ctx, 2, 10)
	executor.Start()
	defer executor.Stop()

	adapter := NewBackgroundExecutorAdapter(executor)
	factory := NewBackgroundServiceFactory(adapter)

	appCtx := &app_context.AppContext{Context: context.Background()}
	locale := locales.EN_US

	t.Run("full integration test", func(_ *testing.T) {
		var executed bool
		var wg sync.WaitGroup
		wg.Add(1)

		service := &mockBackgroundService{
			name: "integration-service",
			executeFunc: func(_ *app_context.AppContext, _ locales.LocaleTypeEnum, _ string) error {
				defer wg.Done()
				executed = true
				return nil
			},
		}

		err := ExecuteService(factory, service, appCtx, locale, "test")
		assert.NoError(err)

		wg.Wait()
		assert.True(executed, "Service should have been executed")
	})
}
