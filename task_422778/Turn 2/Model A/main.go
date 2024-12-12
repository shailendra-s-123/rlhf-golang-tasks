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
	workDuration = time.Duration(rand.Intn(5)) * time.Second
	shutdownTimeout = 5 * time.Second
)

func processTask(ctx context.Context, id int) error {
	log.Printf("Worker %d started task\n", id)

	select {
	case <-ctx.Done():
		return fmt.Errorf("task %d cancelled", id)
	default:
		time.Sleep(workDuration)
	}

	log.Printf("Worker %d completed task\n", id)
	return nil
}

func worker(ctx context.Context, id int) {
	defer log.Printf("Worker %d exited\n", id)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			if err := processTask(ctx, id); err != nil {
				log.Printf("Worker %d task error: %v\n", id, err)
				return
			}
		}
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	numWorkers := 3
	workers := make([]*worker, numWorkers)

	for i := 0; i < numWorkers; i++ {
		w := &worker{
			ctx: ctx,
			id:  i,
		}
		workers[i] = w
		go w.run()
	}

	// Simulate some work being done
	time.Sleep(2 * time.Second)

	// Graceful shutdown: handle interrupt signals (SIGINT, SIGTERM)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Received interrupt signal, shutting down...")
		// Cancel the context to signal workers to stop
		cancel()
		// Wait for workers to complete or timeout
		time.AfterFunc(shutdownTimeout, func() {
			log.Println("Shutdown timeout reached, forcing exit")
			os.Exit(1)
		})
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