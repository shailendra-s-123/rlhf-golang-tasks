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

var (
	rates       map[string]*rateLimiter
	rateLimiterLock sync.RWMutex
)

// rateLimiter is a rate limiting struct
type rateLimiter struct {
	limit int
	tokens int
	tick   <-chan time.Time
}

// newRateLimiter initializes a new rateLimiter
func newRateLimiter(limit int) (*rateLimiter, error) {
	if limit <= 0 {
		return nil, fmt.Errorf("rate limit must be positive: %d", limit)
	}

	interval := time.Duration(1*time.Second / time.Duration(limit))
	tick := time.Tick(interval)

	return &rateLimiter{
		limit: limit,
		tokens: limit,
		tick:   tick,
	}, nil
}

// validates the request parameter
func requestValidationMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Example validation: check if "id" parameter is present and an integer
		if _, ok := r.URL.Query()["id"]; !ok {
			http.Error(w, "Missing 'id' parameter", http.StatusBadRequest)
			return
		}

		idStr := r.URL.Query()["id"][0]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid 'id' parameter: not an integer", http.StatusBadRequest)
			return
		}

		// Proceed with the handler if validation passes
		handler.ServeHTTP(w, r)
	})
}

// implements the rate limiting function
func rateLimitingMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "id", r.URL.Query()["id"][0])
		currentID := context.Get(ctx, "id").(string)

		rateLimiterLock.RLock()
		r, exists := rates[currentID]
		rateLimiterLock.RUnlock()

		if !exists {
			r, err := newRateLimiter(5) // Assuming a default rate limit of 5 requests per second
			if err != nil {
				log.Printf("Error creating rate limiter: %v", err)
				http.Error(w, "Server error", http.StatusInternalServerError)
				return
			}

			rateLimiterLock.Lock()
			defer rateLimiterLock.Unlock()
			rates[currentID] = r
		}

		// Wait for new tokens if needed
		for r.tokens == 0 {
			select {
			case <-time.After(100 * time.Millisecond): // Wait for some time before checking again
			case <-r.tick:
				r.tokens = r.limit
			}
		}

		// Grant a token for the current request
		r.tokens--

		// Proceed with the handler
		handler.ServeHTTP(w, r.Request)
	})
}
func main() {
	fmt.Println("Server starting on :8080")

	// Create a new router
	mux := http.NewServeMux()
	mux.HandleFunc("/item", itemHandler)
	handler := requestValidationMiddleware(mux)
	handler = rateLimitingMiddleware(handler)

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
func itemHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
} 