package contractsproviders

import (
	application_errors "goprojectskeleton/src/application/shared/errors"
	"time"
)

type ICacheProvider interface {
	// Get retrieves the value for a given key from the cache.
	Get(key string, dest any) (bool, *application_errors.ApplicationError)
	// Set stores a value with the specified key and time-to-live (TTL) in the cache.
	Set(key string, value any, ttl time.Duration) *application_errors.ApplicationError
	// Delete removes the value associated with the specified key from the cache.
	Delete(key string) *application_errors.ApplicationError
	// Flush removes all values from the cache.
	Flush() *application_errors.ApplicationError
}
