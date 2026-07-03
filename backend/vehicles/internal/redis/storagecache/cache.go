package storage

import (
	"time"
)

type StorageCache interface {
	Set(key string, value any, ttl time.Duration) error
	Get(key string) (string, error)
	Delete(key string) error
	Exists(key string) (bool, error)
	GetJSON(key string, dest any) error
	SetJSON(key string, value any, ttl time.Duration) error
	DeletePatterns(pattern string) error
}
