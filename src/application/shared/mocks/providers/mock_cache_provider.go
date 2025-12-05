// Package mocks provides mock implementations of the contracts/providers/cache_provider.go interface
package providersmocks

import (
	contractsProviders "goprojectskeleton/src/application/contracts/providers"
	application_errors "goprojectskeleton/src/application/shared/errors"
	"time"

	"github.com/stretchr/testify/mock"
)

// MockCacheProvider is a mock implementation of the ICacheProvider interface
type MockCacheProvider struct {
	mock.Mock
}

var _ contractsProviders.ICacheProvider = (*MockCacheProvider)(nil)

// Get get the value for a given key from the cache
func (m *MockCacheProvider) Get(key string, dest any) (bool, *application_errors.ApplicationError) {
	args := m.Called(key, dest)

	// The first argument is the bool (exists)
	exists := args.Bool(0)

	// The second argument can be the error
	if len(args) > 1 {
		if errorArg := args.Get(1); errorArg != nil {
			if err, ok := errorArg.(*application_errors.ApplicationError); ok {
				return exists, err
			}
		}
	}

	// If exists and there is a value in the third argument, copy it to the dest
	if exists && dest != nil && len(args) > 2 {
		if val := args.Get(2); val != nil {
			if dest, ok := dest.(*int); ok {
				if v, ok := val.(int); ok {
					*dest = v
				}
			}
		}
	}

	return exists, nil
}

// Set set the value for a given key in the cache
func (m *MockCacheProvider) Set(key string, value any, ttl time.Duration) *application_errors.ApplicationError {
	args := m.Called(key, value, ttl)
	errorArg := args.Get(0)
	if errorArg != nil {
		return errorArg.(*application_errors.ApplicationError)
	}
	return nil
}

// Delete delete the value for a given key from the cache
func (m *MockCacheProvider) Delete(key string) *application_errors.ApplicationError {
	args := m.Called(key)
	errorArg := args.Get(0)
	if errorArg != nil {
		return errorArg.(*application_errors.ApplicationError)
	}
	return nil
}

// Flush flush the cache
func (m *MockCacheProvider) Flush() *application_errors.ApplicationError {
	args := m.Called()
	errorArg := args.Get(0)
	if errorArg != nil {
		return errorArg.(*application_errors.ApplicationError)
	}
	return nil
}
