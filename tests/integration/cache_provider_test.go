package integrationtest

import (
	"sync"
	"testing"
	"time"

	contractsproviders "github.com/simon3640/goprojectskeleton/src/application/contracts/providers"
	"github.com/simon3640/goprojectskeleton/src/application/shared/settings"
	"github.com/simon3640/goprojectskeleton/src/infrastructure/providers"
	"github.com/stretchr/testify/assert"
)

func setupCacheProvider() contractsproviders.ICacheProvider {
	providers.CacheProviderInstance.Flush()
	return providers.CacheProviderInstance
}

func TestCacheProvider_Increment(t *testing.T) {
	assert := assert.New(t)

	cacheProvider := setupCacheProvider()

	increment, err := cacheProvider.Increment("testincrement", time.Duration(settings.AppSettingsInstance.RedisTTL)*time.Second)
	assert.Nil(err)
	assert.Equal(int64(1), increment)

	increment, err = cacheProvider.Increment("testincrement", time.Duration(settings.AppSettingsInstance.RedisTTL)*time.Second)
	assert.Nil(err)
	assert.Equal(int64(2), increment)

	increment, err = cacheProvider.Increment("testincrement", time.Duration(settings.AppSettingsInstance.RedisTTL)*time.Second)
	assert.Nil(err)
	assert.Equal(int64(3), increment)
}

func TestCacheProvider_IncrementBy(t *testing.T) {
	assert := assert.New(t)

	cacheProvider := setupCacheProvider()

	increment, err := cacheProvider.IncrementBy("testincrementby", 10, time.Duration(settings.AppSettingsInstance.RedisTTL)*time.Second)
	assert.Nil(err)
	assert.Equal(int64(10), increment)

	increment, err = cacheProvider.IncrementBy("testincrementby", 10, time.Duration(settings.AppSettingsInstance.RedisTTL)*time.Second)
	assert.Nil(err)
	assert.Equal(int64(20), increment)
}

func TestCacheProvider_IncrementByConcurrency(t *testing.T) {
	assert := assert.New(t)

	cacheProvider := setupCacheProvider()
	key := "testincrementbyconcurrency"

	var wg sync.WaitGroup
	wg.Add(10)

	cacheProvider.Set(key, 0, 10*time.Second)

	results := make(chan int64, 10)

	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()

			inc, err := cacheProvider.IncrementBy(key, 1, 10*time.Second)
			assert.Nil(err)

			results <- inc
		}()
	}

	wg.Wait()
	close(results)

	// Validate that 10 increments were performed
	assert.Equal(10, len(results))

	// Validate final value
	var finalValue int64
	exists, _ := cacheProvider.Get(key, &finalValue)
	assert.True(exists)
	assert.Equal(int64(10), finalValue)

	// Validate that the results are between 1 and 10 and without repetition
	seen := map[int64]bool{}
	for v := range results {
		assert.GreaterOrEqual(v, int64(1))
		assert.LessOrEqual(v, int64(10))
		seen[v] = true
	}

	assert.Equal(10, len(seen))
}

func TestCacheProvider_GetInt64(t *testing.T) {
	assert := assert.New(t)

	cacheProvider := setupCacheProvider()

	int64Value, err := cacheProvider.GetInt64("testint64")
	assert.Nil(err)
	assert.Equal(int64(0), int64Value)
}

func TestCacheProvider_Delete(t *testing.T) {
	assert := assert.New(t)

	cacheProvider := setupCacheProvider()

	cacheProvider.Set("testdelete", 1, time.Duration(settings.AppSettingsInstance.RedisTTL)*time.Second)

	err := cacheProvider.Delete("testdelete")
	assert.Nil(err)
}

func TestCacheProvider_Flush(t *testing.T) {
	assert := assert.New(t)

	cacheProvider := setupCacheProvider()

	cacheProvider.Set("testflush", 1, time.Duration(settings.AppSettingsInstance.RedisTTL)*time.Second)
	cacheProvider.Set("testflush2", 2, time.Duration(settings.AppSettingsInstance.RedisTTL)*time.Second)
	err := cacheProvider.Flush()
	assert.Nil(err)

	int64Value, err := cacheProvider.GetInt64("testflush")
	assert.Nil(err)
	assert.Equal(int64(0), int64Value)

	int64Value, err = cacheProvider.GetInt64("testflush2")
	assert.Nil(err)
	assert.Equal(int64(0), int64Value)
}

func TestCacheProvider_Get(t *testing.T) {
	assert := assert.New(t)

	cacheProvider := setupCacheProvider()

	value := 8451

	cacheProvider.Set("testget", value, time.Duration(settings.AppSettingsInstance.RedisTTL)*time.Second)
	intValue := 0
	exists, err := cacheProvider.Get("testget", &intValue)
	assert.Nil(err)
	assert.True(exists)
	assert.Equal(value, intValue)
}

func TestCacheProvider_Set(t *testing.T) {
	assert := assert.New(t)

	cacheProvider := setupCacheProvider()

	value := "testset"

	err := cacheProvider.Set("testset", value, time.Duration(settings.AppSettingsInstance.RedisTTL)*time.Second)
	assert.Nil(err)
}

func TestCacheProvider_TTL_Expiration(t *testing.T) {
	assert := assert.New(t)

	cacheProvider := setupCacheProvider()

	ttl := time.Duration(2) * time.Second

	cacheProvider.Set("testttl", 1, ttl)

	time.Sleep(ttl + 1*time.Second)

	var intValue int
	exists, err := cacheProvider.Get("testttl", &intValue)
	assert.Nil(err)
	assert.False(exists)
}
