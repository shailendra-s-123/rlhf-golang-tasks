package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

const (
	maxRequestsPerSecond = 5
	rateLimitWindow      = time.Second
)

var (
	// rateLimiter is a global variable to track the rate limit for all clients.
	// In a more sophisticated implementation, this could be per-client or per-IP.
	rateLimiter = NewRateLimiter(maxRequestsPerSecond, rateLimitWindow)
)

// NewRateLimiter creates a new rate limiter with the specified maximum requests per second and window duration.
func NewRateLimiter(maxRequests int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		maxRequests: maxRequests,
		window:      window,
		requests:    0,
		mu:          &sync.Mutex{},
		queue:       make(chan struct{}, maxRequests),
	}
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

func handler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	err := rateLimiter.Acquire(ctx)
	if err != nil {
		log.Printf("Rate limit exceeded: %v", err)
		http.Error(w, err.Error(), http.StatusTooManyRequests)
		return
	}

	// Simulate some work
	select {
	case <-time.After(time.Duration(rand.Intn(1000)) * time.Millisecond):
	case <-ctx.Done():
		return
	}

	fmt.Fprintln(w, "Hello, World!")
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("Server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}