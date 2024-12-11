package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

const (
	numWorkers = 3  // Number of Goroutines per stage
	bufferSize = 10 // Buffer size for channels
)

func main() {
	rand.Seed(time.Now().UnixNano())

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize channels
	input := make(chan int, bufferSize)
	processed := make(chan int, bufferSize)
	output := make(chan int, bufferSize)

	// Start worker Goroutines
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go stage1(ctx, &wg, input, processed)
		wg.Add(1)
		go stage2(ctx, &wg, processed, output)
		wg.Add(1)
		go stage3(ctx, &wg, output)
	}

	// Generate input data
	for i := 0; i < 100; i++ {
		select {
		case input <- rand.Intn(100):
		case <-ctx.Done():
			return
		}
	}

	// Close input channel to signal end of work
	close(input)

	// Wait for all workers to finish
	wg.Wait()

	fmt.Println("Pipeline processing completed.")
}

func stage1(ctx context.Context, wg *sync.WaitGroup, in <-chan int, out chan<- int) {
	defer wg.Done()
	for num := range in {
		select {
		case <-ctx.Done():
			return
		default:
			processedNum := num * 2
			out <- processedNum
			fmt.Printf("Stage 1: Processed %d -> %d\n", num, processedNum)
		}
	}
}

func stage2(ctx context.Context, wg *sync.WaitGroup, in <-chan int, out chan<- int) {
	defer wg.Done()
	for num := range in {
		select {
		case <-ctx.Done():
			return
		default:
			time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond) // Simulate work
			processedNum := num + 10
			out <- processedNum
			fmt.Printf("Stage 2: Processed %d -> %d\n", num, processedNum)
		}
	}
}

func stage3(ctx context.Context, wg *sync.WaitGroup, in <-chan int) {
	defer wg.Done()
	for num := range in {
		select {
		case <-ctx.Done():
			return
		default:
			fmt.Printf("Stage 3: Final result: %d\n", num)
		}
	}
}