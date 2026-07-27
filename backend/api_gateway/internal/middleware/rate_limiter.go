package middleware

import (
	"net"
	"net/http"
	"strings"
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
	r := &RateLimiter{
		buckets:    make(map[string]*Bucket),
		Capacity:   float64(capacity),
		RefillRate: refillRate,
	}

	r.Cleanup(10 * time.Minute)

	return r
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
		bucket = &Bucket{
			CurrentTokens:  r.Capacity,
			LastRefillTime: time.Now(),
		}

		r.buckets[key] = bucket
	} else {
		r.Refill(bucket)
	}

	if bucket.CurrentTokens >= 1 {
		bucket.CurrentTokens--
		return true
	}
	return false
}

func (r *RateLimiter) Cleanup(inactiveDuration time.Duration) {
	ticker := time.NewTicker(5 * time.Minute)

	go func() {
		defer ticker.Stop()

		for range ticker.C {
			r.mu.Lock()

			now := time.Now()

			for key, bucket := range r.buckets {
				if now.Sub(bucket.LastRefillTime) > inactiveDuration {
					delete(r.buckets, key)
				}
			}

			r.mu.Unlock()
		}
	}()
}

func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := rl.GetClientIP(r)

		if !rl.Allow(key) {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Retry-After", "1")

			w.WriteHeader(http.StatusTooManyRequests)

			w.Write([]byte(`{"success": false, 
			"message": "Rate limit exceeded. Please try again later."}`))
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (rl *RateLimiter) GetClientIP(r *http.Request) string {
	userID := r.Header.Get("X-User-ID")
	if userID != "" {
		return "user:" + userID
	}

	ip := r.Header.Get("X-Real-IP")
	if ip == "" {
		ip = r.Header.Get("X-Forwarded-For")
	}

	if ip == "" {
		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err == nil {
			ip = host
		} else {
			ip = r.RemoteAddr
		}
	}

	// X-Forwarded-For may contains multiple IPs
	if strings.Contains(ip, ",") {
		ip = strings.TrimSpace(strings.Split(ip, ",")[0])
	}

	return "ip:" + ip
}
