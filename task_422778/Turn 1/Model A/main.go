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

// worker represents a goroutine that processes tasks
type worker struct {
	ctx context.Context
	quit chan struct{}
	id   int
}

// run starts the worker goroutine
func (w *worker) run() {
	defer log.Printf("Worker %d exited\n", w.id)

	for {
		select {
		case <-w.ctx.Done():
			return
		case <-w.quit:
			return
		default:
			log.Printf("Worker %d doing work...\n", w.id)
			time.Sleep(time.Duration(rand.Intn(5)) * time.Second)
		}
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	numWorkers := 3
	workers := make([]*worker, numWorkers)
	quit := make(chan struct{})

	for i := 0; i < numWorkers; i++ {
		w := &worker{
			ctx:  ctx,
			quit: quit,
			id:   i,
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
		cancel() // Cancel all worker contexts
		close(quit)
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