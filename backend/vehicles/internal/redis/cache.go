package redis

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

var ErrCacheMiss = errors.New("cache miss")

type Cache struct {
	Client *redis.Client
}

func NewCache(client *redis.Client) *Cache {
	return &Cache{
		Client: client,
	}
}

func (c *Cache) Set(key string, value any, ttl time.Duration) error {
	ctx, cancel := newContext()
	defer cancel()

	return c.Client.Set(ctx, key, value, ttl).Err()
}

func (c *Cache) Get(key string) (string, error) {
	ctx, cancel := newContext()
	defer cancel()

	value, err := c.Client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", ErrCacheMiss
		}
		return "", err
	}

	return value, nil
}

func (c *Cache) Delete(key string) error {
	ctx, cancel := newContext()
	defer cancel()

	return c.Client.Del(ctx, key).Err()
}

func newContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}

func (c *Cache) Exists(key string) (bool, error) {
	ctx, cancel := newContext()
	defer cancel()

	count, err := c.Client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (c *Cache) SetJSON(key string, value any, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	ctx, cancel := newContext()
	defer cancel()

	return c.Client.Set(ctx, key, data, ttl).Err()
}

func (c *Cache) GetJSON(key string, dest any) error {
	ctx, cancel := newContext()
	defer cancel()

	value, err := c.Client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return ErrCacheMiss
		}
		return err
	}

	return json.Unmarshal([]byte(value), dest)
}

func (c *Cache) DeletePatterns(pattern string) error {
	ctx, cancel := newContext()
	defer cancel()

	keys, err := c.Client.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}

	if len(keys) == 0 {
		return nil
	}

	return c.Client.Del(ctx, keys...).Err()
}
