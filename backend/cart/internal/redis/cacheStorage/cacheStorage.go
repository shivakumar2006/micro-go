package cachestorage

import "time"

type CacheStorage interface {
	Set(key string, value any, ttl time.Duration) error
	Get(key string) (string, error)
	GetJSON(key string, dest any) error
	SetJSON(key string, value any, ttl time.Duration) error
	Delete(key string) error
	DeletePatterns(key string) error
	Exists(key string) (bool, error)
}
