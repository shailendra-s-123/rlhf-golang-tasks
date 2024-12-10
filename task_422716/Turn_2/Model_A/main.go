package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// RequestValidator validates the request parameters
type RequestValidator func(*http.Request) error

// ValidateStringParam validates a string parameter is not empty
func ValidateStringParam(name string) RequestValidator {
	return func(r *http.Request) error {
		val := r.URL.Query().Get(name)
		if val == "" {
			return fmt.Errorf("param %s is required", name)
		}
		return nil
	}
}

// ValidateIntParam validates an int parameter can be parsed
func ValidateIntParam(name string) RequestValidator {
	return func(r *http.Request) error {
		val := r.URL.Query().Get(name)
		_, err := strconv.Atoi(val)
		return err
	}
}

// RateLimiter limits the number of requests per second
type RateLimiter struct {
	mu     sync.Mutex
	buckets map[string]*tokenBucket
}

type tokenBucket struct {
	rate   int // number of tokens per second
	capacity int // bucket size
	tokens  int // current number of tokens
	nextFill time.Time
}

func newRateLimiter(rates map[string]int) *RateLimiter {
	limiter := &RateLimiter{buckets: make(map[string]*tokenBucket)}
	for key, rate := range rates {
		limiter.buckets[key] = newTokenBucket(rate)
	}
	return limiter
}

func newTokenBucket(rate int) *tokenBucket {
	if rate <= 0 {
		panic("rate must be positive")
	}
	capacity := rate
	return &tokenBucket{rate: rate, capacity: capacity, tokens: capacity}
}

func (b *tokenBucket) refill() {
	if time.Now().After(b.nextFill) {
		b.tokens = b.capacity
		b.nextFill = time.Now().Add(time.Second / time.Duration(b.rate))
	}
}

func (b *tokenBucket) consume() bool {
	b.refill()
	if b.tokens == 0 {
		return false
	}
	b.tokens--
	return true
}

// requestValidationMiddleware validates the request using the provided validators
func requestValidationMiddleware(validators ...RequestValidator) http.HandlerFunc {
	return func(next http.Handler) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			for _, validator := range validators {
				if err := validator(r); err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
			}
			next.ServeHTTP(w, r)
		}
	}
}

// rateLimitingMiddleware limits the request rate based on the configured rate
func rateLimitingMiddleware(limiter *RateLimiter) http.HandlerFunc {
	return func(next http.Handler) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			path := r.URL.Path
			bucket, ok := limiter.buckets[path]
			if !ok {
				http.Error(w, "rate limit not configured for this path", http.StatusInternalServerError)
				return
			}

			if !bucket.consume() {
				http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		}
	}
}

func main() {
	fmt.Println("Server starting on :8080")
	// Create an empty router
	mux := http.NewServeMux()

	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/greet", greetHandler)

	// Custom middleware to handle errors and log requests
	handler := errorMiddleware(loggingMiddleware(mux))

	// Request validation
	validators := []RequestValidator{
		ValidateStringParam("name"),
	}

	// Rate limiting
	limiter := newRateLimiter(map[string]int{
		"/greet": 5, // 5 requests per second for /greet path
	})

	handler = rateLimitingMiddleware(limiter)(handler)
	handler = requestValidationMiddleware(validators...)(handler)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}