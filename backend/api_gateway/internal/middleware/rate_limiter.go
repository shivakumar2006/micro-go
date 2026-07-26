package middleware

import (
	"sync"
	"time"
)

type Bucket struct {
	CurrentTokens  float64
	LastRefillTime time.Time
}

type RateLimiter struct {
	mu         sync.Mutex
	buckets    map[string]*Bucket
	Capacity   float64
	RefillRate float64
}

func NewRateLimiter(capacity float64, refillRate float64) *RateLimiter {
	return &RateLimiter{
		buckets:    make(map[string]*Bucket),
		Capacity:   float64(capacity),
		RefillRate: refillRate,
	}
}

func (r *RateLimiter) Refill(bucket *Bucket) {
	now := time.Now()

	elapsed := now.Sub(bucket.LastRefillTime).Seconds()

	// add token based on time elapsed since last refill
	newTokens := elapsed * r.RefillRate

	bucket.CurrentTokens += newTokens

	if bucket.CurrentTokens > r.Capacity {
		bucket.CurrentTokens = r.Capacity
	}

	bucket.LastRefillTime = now
}

func (r *RateLimiter) Allow(key string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	bucket, exists := r.buckets[key]
	if !exists {
		newBucket := &Bucket{
			CurrentTokens:  r.Capacity,
			LastRefillTime: time.Now(),
		}
		r.buckets[key] = newBucket
		bucket = newBucket
	}

	r.Refill(bucket)

	if bucket.CurrentTokens >= 1 {
		bucket.CurrentTokens--
		return true
	}
	return false
}
