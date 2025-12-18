// Package services provides a singleton factory for background services.
package services

import (
	"sync"

	app_context "github.com/simon3640/goprojectskeleton/src/application/shared/context"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales"
	"github.com/simon3640/goprojectskeleton/src/application/shared/workers"
)

var (
	backgroundServiceFactorySingleton *BackgroundServiceFactory
	factoryOnce                       sync.Once
	factoryMu                         sync.Mutex
)

// InitializeBackgroundServiceFactory initializes the singleton BackgroundServiceFactory instance.
// This should be called during application startup in the infrastructure layer.
// This function is thread-safe and can only be called once.
func InitializeBackgroundServiceFactory() {
	factoryOnce.Do(func() {
		executor := workers.GetBackgroundExecutor()
		adapter := NewBackgroundExecutorAdapter(executor)
		backgroundServiceFactorySingleton = NewBackgroundServiceFactory(adapter)
	})
}

// GetBackgroundServiceFactory returns the singleton instance of BackgroundServiceFactory.
// The factory must be initialized first by calling InitializeBackgroundServiceFactory.
// This function is thread-safe.
func GetBackgroundServiceFactory() *BackgroundServiceFactory {
	return backgroundServiceFactorySingleton
}

// ResetBackgroundServiceFactory resets the singleton instance.
// This is primarily useful for testing purposes.
// WARNING: Only call this when you're sure no other goroutines are using the factory.
func ResetBackgroundServiceFactory() {
	factoryMu.Lock()
	defer factoryMu.Unlock()
	backgroundServiceFactorySingleton = nil
	factoryOnce = sync.Once{}
}

// ExecuteBackgroundService is a convenience function that executes a background service
// using the singleton factory. This is the recommended way to execute background services.
func ExecuteBackgroundService[Input any](
	service BackgroundService[Input],
	ctx *app_context.AppContext,
	locale locales.LocaleTypeEnum,
	input Input,
) error {
	factory := GetBackgroundServiceFactory()
	return ExecuteService(factory, service, ctx, locale, input)
}
