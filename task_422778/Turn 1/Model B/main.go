package main

import (
	"context"
	_"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"
)

func taskWorker(ctx context.Context, wg *sync.WaitGroup, workerID int) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			log.Printf("Worker %d received shutdown signal, exiting.", workerID)
			return
		default:
			log.Printf("Worker %d is working...", workerID)
			time.Sleep(2 * time.Second)
		}
	}
}

func main() {
	numWorkers := 3
	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(context.Background())

	defer func() {
		cancel()
		wg.Wait()
		log.Println("All workers terminated. Application shutting down.")
	}()

	// Start worker goroutines
	for workerID := 1; workerID <= numWorkers; workerID++ {
		wg.Add(1)
		go taskWorker(ctx, &wg, workerID)
	}

	// Wait for os.Interrupt (Ctrl+C)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	<-sigs

	log.Println("Received interrupt signal. Shutting down...")

	// Gracefully shutdown workers
	ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
}
