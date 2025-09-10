package providers

import (
	"context"
	"encoding/json"
	"time"

	contractsProviders "gormgoskeleton/src/application/contracts/providers"
	application_errors "gormgoskeleton/src/application/shared/errors"
	"gormgoskeleton/src/application/shared/locales/messages"
	"gormgoskeleton/src/application/shared/status"

	"github.com/redis/go-redis/v9"
)

type RedisCacheProvider struct {
	client *redis.Client
}

var _ contractsProviders.ICacheProvider = &RedisCacheProvider{}

func NewRedisCacheProvider(addr string, password string, db int) *RedisCacheProvider {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &RedisCacheProvider{client: rdb}
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

var CacheProviderInstance contractsProviders.ICacheProvider
