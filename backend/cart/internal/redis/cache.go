package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var ErrCacheMiss = errors.New("cache miss")

type Cache struct {
	Redis *redis.Client
}

func NewCache(redis *redis.Client) *Cache {
	return &Cache{Redis: redis}
}

func (c *Cache) Set(key string, value any, ttl time.Duration) error {
	ctx, cancel := c.newContext()
	defer cancel()

	err := c.Redis.Set(ctx, key, value, ttl).Err()
	if err != nil {
		panic(err)
	}

	return nil
}

func (c *Cache) Get(key string) (string, error) {
	ctx, cancel := c.newContext()
	defer cancel()

	value, err := c.Redis.Get(ctx, key).Result()
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

	return c.Redis.Del(ctx, key).Err()
}

func (c *Cache) SetJSON(key string, value any, ttl time.Duration) error {
	byteValue, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON : %w", err)
	}

	ctx, cancel := c.newContext()
	defer cancel()

	err = c.Redis.Set(ctx, key, byteValue, ttl).Err()
	if err != nil {
		panic(err)
	}

	return nil
}

func (c *Cache) GetJSON(key string, dest any) error {
	ctx, cancel := c.newContext()
	defer cancel()

	value, err := c.Redis.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return ErrCacheMiss
		}
		panic(err)
	}

	return json.Unmarshal([]byte(value), dest)
}

func (c *Cache) Exists(key string) (bool, error) {
	ctx, cancel := c.newContext()
	defer cancel()

	count, err := c.Redis.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (c *Cache) DeletePatterns(patterns string) error {
	ctx, cancel := c.newContext()
	defer cancel()

	keys, err := c.Redis.Keys(ctx, patterns).Result()
	if err != nil {
		return err
	}

	if len(keys) == 0 {
		return nil
	}

	return c.Redis.Del(ctx, keys...).Err()
}

func (c *Cache) newContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}
