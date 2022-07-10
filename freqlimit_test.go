package freqlimit_test

import (
	"testing"

	"github.com/daqiancode/freqlimit"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

func getRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})
}

func TestFreqLimit(t *testing.T) {
	limit := freqlimit.NewFreqLimit(getRedisClient(), "user/1")
	limit.AddLimit(1, 2)
	ok, err := limit.Incr()
	assert.Nil(t, err)
	assert.True(t, ok)
	ok, err = limit.Incr()
	assert.Nil(t, err)
	assert.True(t, ok)
	ok, err = limit.Incr()
	assert.Nil(t, err)
	assert.False(t, ok)
}
