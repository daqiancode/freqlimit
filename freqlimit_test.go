package freqlimit_test

import (
	"fmt"
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

	left, err := limit.GetLeft()
	assert.Nil(t, err)
	fmt.Println(left)
	ok, err := limit.Incr()
	assert.Nil(t, err)
	assert.True(t, ok)
	left, err = limit.GetLeft()
	assert.Nil(t, err)
	fmt.Println(left)
	ok, err = limit.Incr()
	assert.Nil(t, err)
	assert.True(t, ok)
	ok, err = limit.Incr()
	assert.Nil(t, err)
	assert.False(t, ok)
}
