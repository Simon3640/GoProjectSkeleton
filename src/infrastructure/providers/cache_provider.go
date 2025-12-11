package providers

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"time"

	contractsProviders "github.com/simon3640/goprojectskeleton/src/application/contracts/providers"
	application_errors "github.com/simon3640/goprojectskeleton/src/application/shared/errors"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales/messages"
	"github.com/simon3640/goprojectskeleton/src/application/shared/status"

	"github.com/redis/go-redis/v9"
)

type RedisCacheProvider struct {
	client *redis.Client
}

var _ contractsProviders.ICacheProvider = &RedisCacheProvider{}

// NewRedisCacheProvider creates a new Redis cache provider
func NewRedisCacheProvider(addr string, password string, db int) *RedisCacheProvider {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &RedisCacheProvider{client: rdb}
}

// NewRedisCacheProviderTLS creates a new Redis cache provider with TLS enabled
func NewRedisCacheProviderTLS(addr, password string, db int) *RedisCacheProvider {
	rdb := NewRedisCacheProvider(addr, password, db)

	rdb.client.Options().TLSConfig = &tls.Config{
		MinVersion:         tls.VersionTLS12,
		InsecureSkipVerify: true,
	}
	return rdb
}

func (r *RedisCacheProvider) Set(key string, value any, ttl time.Duration) *application_errors.ApplicationError {
	ctx := context.Background()

	data, err := json.Marshal(value)
	if err != nil {
		return application_errors.NewApplicationError(
			status.ProviderError,
			messages.MessageKeysInstance.SOMETHING_WENT_WRONG,
			err.Error(),
		)
	}

	if err := r.client.Set(ctx, key, data, ttl).Err(); err != nil {
		return application_errors.NewApplicationError(
			status.ProviderError,
			messages.MessageKeysInstance.SOMETHING_WENT_WRONG,
			err.Error(),
		)
	}
	return nil
}

func (r *RedisCacheProvider) Get(key string, dest any) (bool, *application_errors.ApplicationError) {
	ctx := context.Background()

	data, err := r.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return false, nil
	} else if err != nil {
		appError := application_errors.NewApplicationError(
			status.ProviderError,
			messages.MessageKeysInstance.SOMETHING_WENT_WRONG,
			err.Error(),
		)
		return false, appError
	}

	if err := json.Unmarshal(data, dest); err != nil {
		return false, application_errors.NewApplicationError(
			status.ProviderError,
			messages.MessageKeysInstance.SOMETHING_WENT_WRONG,
			err.Error(),
		)
	}

	return true, nil
}

func (r *RedisCacheProvider) Delete(key string) *application_errors.ApplicationError {
	ctx := context.Background()
	if err := r.client.Del(ctx, key).Err(); err != nil {
		return application_errors.NewApplicationError(
			status.ProviderError,
			messages.MessageKeysInstance.SOMETHING_WENT_WRONG,
			err.Error(),
		)
	}
	return nil
}

func (r *RedisCacheProvider) Flush() *application_errors.ApplicationError {
	ctx := context.Background()
	if err := r.client.FlushDB(ctx).Err(); err != nil {
		return application_errors.NewApplicationError(
			status.ProviderError,
			messages.MessageKeysInstance.SOMETHING_WENT_WRONG,
			err.Error(),
		)
	}
	return nil
}

// Increment atomically increments the value of a key by 1 if the key does not exist creates it with the value 1 and sets the TTL
// returns the incremented value and an error if any
func (r *RedisCacheProvider) Increment(key string, ttl time.Duration) (int64, *application_errors.ApplicationError) {
	ctx := context.Background()
	cmd := r.client.Incr(ctx, key)
	if err := cmd.Err(); err != nil {
		return 0, application_errors.NewApplicationError(
			status.ProviderError,
			messages.MessageKeysInstance.SOMETHING_WENT_WRONG,
			err.Error(),
		)
	}
	if cmd.Val() == 1 {
		if err := r.client.Expire(ctx, key, ttl).Err(); err != nil {
			return 0, application_errors.NewApplicationError(
				status.ProviderError,
				messages.MessageKeysInstance.SOMETHING_WENT_WRONG,
				err.Error(),
			)
		}
	}
	return cmd.Val(), nil
}

// IncrementBy atomically increments the value of a key by a given amount
// returns the incremented value and an error if any
func (r *RedisCacheProvider) IncrementBy(key string, increment int64, ttl time.Duration) (int64, *application_errors.ApplicationError) {
	ctx := context.Background()
	cmd := r.client.IncrBy(ctx, key, increment)
	if err := cmd.Err(); err != nil {
		return 0, application_errors.NewApplicationError(
			status.ProviderError,
			messages.MessageKeysInstance.SOMETHING_WENT_WRONG,
			err.Error(),
		)
	}
	if cmd.Val() == increment {
		if err := r.client.Expire(ctx, key, ttl).Err(); err != nil {
			return 0, application_errors.NewApplicationError(
				status.ProviderError,
				messages.MessageKeysInstance.SOMETHING_WENT_WRONG,
				err.Error(),
			)
		}
	}
	return cmd.Val(), nil
}

// GetInt64 gets the value of a key as an int64
// returns the value and an error if any
func (r *RedisCacheProvider) GetInt64(key string) (int64, *application_errors.ApplicationError) {
	ctx := context.Background()
	data, err := r.client.Get(ctx, key).Int64()
	if err == redis.Nil {
		return 0, nil
	} else if err != nil {
		return 0, application_errors.NewApplicationError(
			status.ProviderError,
			messages.MessageKeysInstance.SOMETHING_WENT_WRONG,
			err.Error(),
		)
	}
	return data, nil
}

var CacheProviderInstance contractsProviders.ICacheProvider
