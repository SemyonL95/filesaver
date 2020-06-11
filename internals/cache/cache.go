package cache

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

type RedisCache struct {
	redis *redis.Client
}

func NewRedisCache(r *redis.Client) *RedisCache {
	return &RedisCache{
		redis: r,
	}
}

func (c *RedisCache) Set(ctx context.Context, key string, value string) (bool, error) {
	res := c.redis.SAdd(ctx, key, value)
	if res == nil {
		return false, fmt.Errorf("error during seting value to set")
	}

	return intToBool(res.Val()), nil
}

func (c *RedisCache) Exists(ctx context.Context, key string, value string) (bool, error) {
	res := c.redis.SIsMember(ctx, key, value)
	if res == nil {
		return false, fmt.Errorf("error during check value of set")
	}

	return res.Val(), nil
}

func (c *RedisCache) Remove(ctx context.Context, key string) (bool, error) {
	res := c.redis.Del(ctx, key)
	if res == nil {
		return false, fmt.Errorf("error during remove of set")
	}

	if res.Val() > 0 {
		return true, nil
	}

	return intToBool(res.Val()), nil
}

func intToBool(n int64) bool {
	if n > 0 {
		return true
	}

	return false
}
