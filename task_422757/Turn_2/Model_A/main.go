package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

const (
	defaultMaxRequestsPerSecond = 5
	rateLimitWindow             = time.Second
	dynamicRateAdjustmentInterval = time.Minute
)

var (
	rateLimiters = make(map[string]*RateLimiter) // Map of client/IP to RateLimiter
	mu           sync.RWMutex
)

func getRateLimiter(key string) *RateLimiter {
	mu.RLock()
	defer mu.RUnlock()
	rl, ok := rateLimiters[key]
	if !ok {
		rl = NewRateLimiter(defaultMaxRequestsPerSecond, rateLimitWindow)
		rateLimiters[key] = rl
	}
	return rl
}

// NewRateLimiter creates a new rate limiter with the specified maximum requests per second and window duration.
func NewRateLimiter(maxRequests int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		maxRequests: maxRequests,
		window:      window,
		requests:    0,
		mu:          &sync.Mutex{},
		queue:       make(chan struct{}, maxRequests),
		lastAdjust:  time.Now(),
	}
	go rl.dynamicRateAdjustment()
	go rl.cleanup()
	return rl
}

// RateLimiter struct to track the rate limit.
type RateLimiter struct {
	maxRequests int
	window      time.Duration
	requests    int
	mu          *sync.Mutex
	queue       chan struct{}
	lastAdjust  time.Time
}

// Acquire acquires a permit from the rate limiter.
func (rl *RateLimiter) Acquire(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case rl.queue <- struct{}{}:
		defer func() {
			<-rl.queue
		}()
		return nil
	}
}

func (rl *RateLimiter) dynamicRateAdjustment() {
	ticker := time.NewTicker(dynamicRateAdjustmentInterval)
	defer ticker.Stop()
	for range ticker.C {
		rl.mu.Lock()
		requestsPerSecond := float64(rl.requests) / (float64(rl.window) / time.Second)
		if requestsPerSecond > 1.5*float64(rl.maxRequests) {
			rl.maxRequests /= 2
		} else if requestsPerSecond < 0.5*float64(rl.maxRequests) {
			rl.maxRequests *= 2
		}
		rl.requests = 0
		rl.mu.Unlock()
	}
}

// cleanup removes expired requests from the count.
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(rl.window)
	defer ticker.Stop()
	for range ticker.C {
		rl.mu.Lock()
		rl.requests = 0
		rl.mu.Unlock()
	}
}

func rateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		key := getClientIP(r)
		rl := getRateLimiter(key)

		err := rl.Acquire(ctx)
		if err != nil {
			log.Printf("Rate limit exceeded for %s: %v", key, err)
			http.Error(w, err.Error(), http.StatusTooManyRequests)
			return
		}