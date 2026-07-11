package redis

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

var ErrCacheMiss = errors.New("cache miss")

type Cache struct {
	Redis *Redis
}

func NewCache(redis *Redis) *Cache {
	return &Cache{Redis: redis}
}

func (c *Cache) Set(key string, value any, ttl time.Duration) error {
	ctx, cancel := c.newContext()
	defer cancel()

	err := c.Redis.Client.Set(ctx, key, value, ttl).Err()
	if err != nil {
		panic(err)
	}

	return nil
}

func (c *Cache) Get(key string) (string, error) {
	ctx, cancel := c.newContext()
	defer cancel()

	value, err := c.Redis.Client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(nil, redis.Nil) {
			return "", ErrCacheMiss
		}
	}

	return value, nil
}

func (c *Cache) Delete(key string) error {
	ctx, cancel := c.newContext()
	defer cancel()

	return c.Redis.Client.Del(ctx, key).Err()
}

func (c *Cache) newContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}
