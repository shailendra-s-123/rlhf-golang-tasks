package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	shutdownTimeout = 5 * time.Second // Timeout to wait for workers to finish during shutdown
)

func worker(ctx context.Context, id int) {
	defer log.Printf("Worker %d exited\n", id)

	for {
		select {
		case <-ctx.Done():
			log.Printf("Worker %d received shutdown signal, exiting\n", id)
			return
		default:
			log.Printf("Worker %d doing work...\n", id)
			duration := time.Duration(rand.Intn(5)) * time.Second
			select {
			case <-ctx.Done():
				log.Printf("Worker %d interrupted during work\n", id)
				return
			case <-time.After(duration):
			}
		}
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	numWorkers := 3
	workers := make([]context.Context, numWorkers)

	for i := 0; i < numWorkers; i++ {
		ctx, cancel := context.WithCancel(ctx)
		workers[i] = ctx
		go worker(ctx, i)
	}

	// Simulate some work being done
	time.Sleep(2 * time.Second)

	// Graceful shutdown: handle interrupt signals (SIGINT, SIGTERM)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Received interrupt signal, shutting down...")
		for _, ctx := range workers {
			cancel(ctx)
		}

		select {
		case <-time.After(shutdownTimeout):
			log.Printf("Shutdown timed out, forcing exit\n")
		default:
			log.Println("All workers terminated gracefully")
		}
	}()

	// Serve HTTP requests until shutdown
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	})
	log.Println("Listening on port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("Error starting HTTP server: %v", err)
	}

	log.Println("Shutdown complete")
}