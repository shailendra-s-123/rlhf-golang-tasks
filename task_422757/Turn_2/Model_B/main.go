package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/ip"
	"sync"
	"sync/atomic"
	"time"
)

const (
	maxRequestsPerSecond  = 10
	rateLimitWindow       = time.Second
	windowBuckets         = 10 // Determines precision of rate limit calculation (per ticker interval)
	bufferSize            = 1000
	pollingInterval       = 500 * time.Millisecond
	initialBackoff         = 100 * time.Millisecond
	maxBackoff            = 5 * time.Second
	maxAdjustmentPercent = 20
)

var (
	rateLimiters = make(map[string]*RateLimiter) // Map of client IPs or API keys to rate limiters
	mu           sync.RWMutex
	updateChan   = make(chan struct{}, 1)          // Channel for notifying watchers about updates
	stats        struct {
		totalLimited int32
	}
	watchers         = make([]chan *rateLimiterStats, 0)
	rateExponentialBackoff = func(attempt int) time.Duration {
		return time.Duration(rand.Intn(int(initialBackoff.Seconds())+1)) * time.Millisecond * time.Power(2, float64(attempt))
	}
)

func getClientKey(r *http.Request) string {
	remoteIP := r.RemoteAddr
	ipSplit := strings.Split(remoteIP, ":")
	return ipSplit[0]
	// In a real-world application, you could also use an API key or other identifier
}

// NewRateLimiter creates a new rate limiter with the specified parameters.
func NewRateLimiter(maxRequests int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		maxRequests:   maxRequests,
		window:        window,
		requests:      make([]int, windowBuckets),
		bucketInterval time.Duration(window / time.Duration(windowBuckets)),
		queue:         make(chan struct{}, maxRequests),
	}
	go rl.cleanup()
	return rl
}

// RateLimiter struct now has a channels for updating limits dynamically.
type RateLimiter struct {
	maxRequests   int
	window        time.Duration
	requests      []int
	bucketInterval time.Duration
	queue         chan struct{}
	limitChan     chan int // Channel for dynamic rate adjustments
}

// cleanup removes expired requests from the count.
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(rl.bucketInterval)
	defer ticker.Stop()
	for range ticker.C {
		rl.requests = append(rl.requests[1:], 0)
	}
}

// Acquire acquires a permit from the rate limiter, and returns the state of rate-limited
func (rl *RateLimiter) Acquire(ctx context.Context, requestStats *rateLimiterStats) (bool, error) {
	select {
	case <-ctx.Done():
		return false, ctx.Err()
	case rl.queue <- struct{}{}:
		defer func() {
			<-rl.queue
		}()
		return true, nil
	}
}

// GetRate stats
func GetRate(clientKey string) *rateLimiterStats {
	mu.RLock()
	rl, ok := rateLimiters[clientKey]
	mu.RUnlock()

	if !ok {
		return nil
	}

	stats := &rateLimiterStats{
		RPS:    rl.window / time.Duration(rl.bucketInterval) * float64(rl.sumRequests(5*time.Second)),
		MaxRPS: float64(rl.maxRequests),
	}

	return stats
}

func addRateLimiterWatcher(w chan *rateLimiterStats) {
	mu.Lock()
	watchers = append(watchers, w)
	mu.Unlock()
	go func() {
		for range updateChan {
			var statsList []*rateLimiterStats
			mu.RLock()
			for _, rl := range rateLimiters {
				stat := &rateLimiterStats{
					RPS:    rl.window / time.Duration(rl.bucketInterval) * float64(rl.sumRequests(5*time.Second)),
					MaxRPS: float64(rl.maxRequests),
				}
				statsList = append(statsList, stat)
			}
			mu.RUnlock()
			w <- &rateLimiterStats{
				RPSList: statsList,
			}
		}
	}()
}