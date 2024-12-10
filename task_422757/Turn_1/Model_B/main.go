package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

// TokenBucket implements a simple rate limiter.
type TokenBucket struct {
	capacity int
	tokens   int
	refill   time.Duration
	mu       sync.Mutex
	ticker   *time.Ticker
}

func NewTokenBucket(capacity int, refill time.Duration) *TokenBucket {
	tb := &TokenBucket{
		capacity: capacity,
		tokens:   capacity,
		refill:   refill,
	}
	tb.ticker = time.NewTicker(refill / time.Duration(capacity))
	go tb.refillTokens()
	return tb
}

func (tb *TokenBucket) refillTokens() {
	for range tb.ticker.C {
		tb.mu.Lock()
		defer tb.mu.Unlock()
		if tb.tokens < tb.capacity {
			tb.tokens++
		}
	}
}

func (tb *TokenBucket) TakeToken() (bool, error) {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	if tb.tokens == 0 {
		return false, fmt.Errorf("rate limit exceeded")
	}
	tb.tokens--
	return true, nil
}

// ApiRequest represents a request to be rate limited.
type ApiRequest struct {
	ctx   context.Context
	tb    *TokenBucket
	client string
}

func processRequest(req *ApiRequest) {
	select {
	case <-req.ctx.Done():
		log.Printf("Request from %s cancelled.\n", req.client)
		return
	default:
	}

	// Simulate API call delay
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)

	log.Printf("Request from %s completed.\n", req.client)
}

func main() {
	rand.Seed(time.Now().UnixNano())

	tb := NewTokenBucket(10, time.Second) // Allow 10 requests per second
	wg := &sync.WaitGroup{}

	// Generate random requests for 5 seconds
	for i := 0; i < 100; i++ {
		client := fmt.Sprintf("client-%d", i)
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)

		req := &ApiRequest{
			ctx:   ctx,
			tb:    tb,
			client: client,
		}

		go func() {
			if ok, _ := req.tb.TakeToken(); !ok {
				log.Printf("Request from %s rate limited.\n", client)
				cancel()
				return
			}

			wg.Add(1)
			processRequest(req)
			wg.Done()
		}()
	}

	// Wait for all requests to complete or timeout
	wg.Wait()
	log.Println("All requests completed.")
}
