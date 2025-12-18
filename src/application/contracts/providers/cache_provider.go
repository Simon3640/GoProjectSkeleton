package contractsproviders

import (
	"time"

	application_errors "github.com/simon3640/goprojectskeleton/src/application/shared/errors"
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

	// Increment atomically increments the value of a key by 1 if the key does not exist creates it with the value 1 and sets the TTL
	// returns the incremented value and an error if any
	Increment(key string, ttl time.Duration) (int64, *application_errors.ApplicationError)
	// IncrementBy atomically increments the value of a key by a given amount if the key does not exist creates it with the value increment and sets the TTL
	// returns the incremented value and an error if any
	IncrementBy(key string, increment int64, ttl time.Duration) (int64, *application_errors.ApplicationError)
	// GetInt64 gets the value of a key as an int64
	// returns the value and an error if any
	GetInt64(key string) (int64, *application_errors.ApplicationError)
}
