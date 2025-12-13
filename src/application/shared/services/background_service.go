// Package services provides background services for the application.
package services

import (
	"context"
	"log"

	app_context "github.com/simon3640/goprojectskeleton/src/application/shared/context"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales"
	"github.com/simon3640/goprojectskeleton/src/application/shared/workers"
)

// BackgroundService is the interface that background services must implement.
// It defines a service that can be executed in the background using the BackgroundExecutor.
type BackgroundService[Input any] interface {
	// Execute executes the service with the given context, locale, and input.
	// This method will be called asynchronously in a background worker.
	Execute(ctx *app_context.AppContext, locale locales.LocaleTypeEnum, input Input) error
	// Name returns a human-readable name for the service, useful for logging and tracing.
	Name() string
}

// BackgroundExecutorInterface is an interface that abstracts the background executor.
// This allows the factory to work with the workers.BackgroundExecutor.
type BackgroundExecutorInterface interface {
	Submit(task func(ctx context.Context)) error
}

// backgroundExecutorAdapter adapts workers.BackgroundExecutor to BackgroundExecutorInterface
type backgroundExecutorAdapter struct {
	executor *workers.BackgroundExecutor
}

func (a *backgroundExecutorAdapter) Submit(task func(ctx context.Context)) error {
	return a.executor.Submit(workers.BackgroundTask(task))
}

// NewBackgroundExecutorAdapter creates an adapter for workers.BackgroundExecutor
func NewBackgroundExecutorAdapter(executor *workers.BackgroundExecutor) BackgroundExecutorInterface {
	if executor == nil {
		return nil
	}
	return &backgroundExecutorAdapter{executor: executor}
}

// BackgroundServiceFactory provides methods to execute background services.
type BackgroundServiceFactory struct {
	executor BackgroundExecutorInterface
}

// NewBackgroundServiceFactory creates a new factory for executing background services.
// It uses the provided executor to submit background tasks.
func NewBackgroundServiceFactory(executor BackgroundExecutorInterface) *BackgroundServiceFactory {
	return &BackgroundServiceFactory{
		executor: executor,
	}
}

// ExecuteService submits a background service to be executed asynchronously.
// The service will be executed in a background worker using the configured executor.
// This method returns immediately without waiting for the service to complete.
func ExecuteService[Input any](
	factory *BackgroundServiceFactory,
	service BackgroundService[Input],
	appCtx *app_context.AppContext,
	locale locales.LocaleTypeEnum,
	input Input,
) error {
	if factory == nil || factory.executor == nil {
		// Fallback: execute in a goroutine if no executor is configured
		go func() {
			if err := service.Execute(appCtx, locale, input); err != nil {
				log.Printf("background service %s returned error: %v", service.Name(), err)
			}
		}()
		return nil
	}

	serviceInput := input
	serviceName := service.Name()

	// Submit the task to the executor
	return factory.executor.Submit(func(_ context.Context) {
		// Execute the service with the captured AppContext
		if err := service.Execute(appCtx, locale, serviceInput); err != nil {
			// Log error but don't propagate it (fire-and-forget)
			log.Printf("background service %s returned error: %v", serviceName, err)
		}
	})
}
