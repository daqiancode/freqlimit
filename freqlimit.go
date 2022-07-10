package freqlimit

import (
	"context"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	DefaultExpireDelay = 10
)

type Limit struct {
	Window int64
	Max    int64
}

type FreqLimit struct {
	red    *redis.Client
	limits []Limit
	key    string
	ctx    context.Context
}

func NewFreqLimit(client *redis.Client, key string) *FreqLimit {
	return &FreqLimit{
		red: client,
		key: key,
		ctx: context.Background(),
	}
}

//AddLimit windows: window length in seconds , max: max request count
func (s *FreqLimit) AddLimit(window, max int64) {
	s.limits = append(s.limits, Limit{Window: window, Max: max})
}

//SetCtx set context for redis API
func (s *FreqLimit) SetCtx(ctx context.Context) {
	s.ctx = ctx
}

//RedisKey the way set key in redis
func (s *FreqLimit) redisKey(window, now int64) string {
	timeKey := now - now%window
	return s.key + "/" + strconv.FormatInt(window, 10) + "/" + strconv.FormatInt(timeKey, 10)
}

//Incr increase the number of calls, throw error if exceed the max call number.
// return true if under frequency limit,return false if exceed the frequency limit
func (s *FreqLimit) Incr() (bool, error) {
	t, err := s.red.Time(s.ctx).Result()
	if err != nil {
		return false, err
	}
	for _, limit := range s.limits {
		redisKey := s.redisKey(limit.Window, t.Unix())
		count, err := s.red.Incr(s.ctx, redisKey).Result()
		if count == 1 {
			s.red.Expire(s.ctx, redisKey, time.Duration(limit.Window+int64(DefaultExpireDelay))*time.Second)
		}
		if err != nil {
			return false, err
		}
		if count > limit.Max {
			return false, nil
		}
	}
	return true, nil
}

//GetLeft return all window rest request count, window:left
func (s *FreqLimit) GetLeft() (map[int64]int64, error) {
	if len(s.limits) == 0 {
		return nil, nil
	}
	t, err := s.red.Time(s.ctx).Result()
	if err != nil {
		return nil, err
	}

	r := make(map[int64]int64, len(s.limits))
	keys := make([]string, len(s.limits))
	for i, v := range s.limits {
		keys[i] = s.redisKey(v.Window, t.Unix())
	}
	countStrs, err := s.red.MGet(s.ctx, keys...).Result()
	if err != nil {
		return nil, err
	}
	for i, v := range countStrs {
		w := s.limits[i]
		if v == nil {
			r[w.Window] = w.Max
			continue
		}
		count, err := strconv.ParseInt(v.(string), 10, 64)
		if err != nil {
			return nil, err
		}
		r[w.Window] = w.Max - count
	}
	return r, nil
}
