package cachestorage

import "time"

type CacheStorage interface {
	Set(key, value any, ttl time.Time) error
}
