package cache

import (
	contractsproviders "github.com/simon3640/goprojectskeleton/src/application/contracts/providers"
	app_context "github.com/simon3640/goprojectskeleton/src/application/shared/context"
	application_errors "github.com/simon3640/goprojectskeleton/src/application/shared/errors"
	"github.com/simon3640/goprojectskeleton/src/application/shared/observability"
)

// Executor is the cache executor for the cache policy.
type Executor[T any, TInput any] struct {
	cacheProvider contractsproviders.ICacheProvider
}

// NewExecutor crea un nuevo Executor con el proveedor de caché dado.
func NewExecutor[T any, TInput any](cacheProvider contractsproviders.ICacheProvider) *Executor[T, TInput] {
	return &Executor[T, TInput]{cacheProvider: cacheProvider}
}

// Execute cache executor
// - Build the key
// - Check if the data exists in the cache
// - If the data exists in the cache, deserialize the data
// - If the data does not exist in the cache, fetch the data from the source
// - Set the data in the cache
// - Return the data
// - If there is an error, return the error
func (e *Executor[T, TInput]) Execute(policy ICachePolicy[T, TInput], input TInput, context *app_context.AppContext) (*T, *application_errors.ApplicationError) {
	// Build the key
	key := policy.BuildKey(input, context)
	var data T

	// Check if the data exists in the cache
	exists, _ := e.cacheProvider.Get(key, &data)
	if exists {
		// Deserialize the data
		deserializedData, err := policy.Deserialize(data)
		if err != nil {
			observability.GetObservabilityComponents().Logger.ErrorWithContext("Failed to deserialize data from cache", err.ToError(), context)
			return nil, err
		}
		return policy.OnCacheHit(deserializedData), nil
	}

	// Fetch the data from the source
	data, err := policy.FetchData(input, context)
	if err != nil {
		return nil, err
	}
	e.cacheProvider.Set(key, data, policy.GetTTL())
	return &data, nil
}
