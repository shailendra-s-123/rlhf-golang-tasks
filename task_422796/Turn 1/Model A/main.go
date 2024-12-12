package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
	"github.com/shirou/gopsutil/cpu"
	
)

const (
	targetCPUUsage = 80 // Target CPU usage percentage
	checkInterval   = 1 * time.Second
	minConcurrency  = 1
	maxConcurrency  = 10
)

var (
	wg sync.WaitGroup
	mu sync.Mutex
	// We'll use a channel to control the number of concurrent goroutines.
	// Goroutines will block on this channel when concurrency limit is reached.
	workCh = make(chan struct{}, 0)
)

func doWork() {
	defer wg.Done()
	// Simulate work
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	rand.Seed(time.Now().UnixNano())

	go monitorCPU()

	wg.Add(100) // Start with some initial workers
	for i := 0; i < 100; i++ {
		go func() {
			workCh <- struct{}{}
			defer func() { <-workCh }()
			doWork()
		}()
	}

	wg.Wait()
}

func monitorCPU() {
	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()

	for range ticker.C {
		cpuPercent, err := cpu.Percent(time.Second, false)
		if err != nil {
			fmt.Println("Error getting CPU usage:", err)
			continue
		}

		mu.Lock()
		defer mu.Unlock()

		currentConcurrency := len(workCh)
		newConcurrency := adjustConcurrency(cpuPercent, currentConcurrency)
		diff := newConcurrency - currentConcurrency

		if diff > 0 {
			for i := 0; i < diff; i++ {
				workCh <- struct{}{}
			}
			fmt.Printf("Increased concurrency to %d\n", newConcurrency)
		} else if diff < 0 {
			for i := 0; i > diff; i-- {
				<-workCh
			}
			fmt.Printf("Decreased concurrency to %d\n", newConcurrency)
		}
	}
}

func adjustConcurrency(cpuPercent float64, currentConcurrency int) int {
	if cpuPercent < targetCPUUsage/2 {
		return min(maxConcurrency, currentConcurrency*2)
	} else if cpuPercent > targetCPUUsage {
		return max(minConcurrency, currentConcurrency/2)
	}
	return currentConcurrency
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}