// Package cache provides the cache policy interface and cache executor.
package cache

import (
	"time"

	app_context "github.com/simon3640/goprojectskeleton/src/application/shared/context"
	application_errors "github.com/simon3640/goprojectskeleton/src/application/shared/errors"
)

// ICachePolicy define la política de caché para un caso de uso.
// Permite implementar buildKey, fetchData, serialize/deserialize y onCacheHit desde otros paquetes.
type ICachePolicy[T any, TInput any] interface {
	BuildKey(input TInput, context *app_context.AppContext) string
	FetchData(input TInput, context *app_context.AppContext) (T, *application_errors.ApplicationError)
	Serialize(data T) (any, *application_errors.ApplicationError)
	Deserialize(data any) (T, *application_errors.ApplicationError)
	GetTTL() time.Duration
	OnCacheHit(data T) *T
}
