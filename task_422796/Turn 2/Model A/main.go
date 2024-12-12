package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

const (
	targetCPUUsage   = 80   // Target CPU usage percentage
	targetMemUsage   = 80   // Target memory usage percentage in percentage
	checkInterval   = 2 * time.Second
	minConcurrency   = 1
	maxConcurrency   = 10
	adjustmentFactor = 0.2 // How much to adjust concurrency based on load
)

var (
	wg        sync.WaitGroup
	mu        sync.Mutex
	workCh    = make(chan struct{}, 0)
	cpuUsage  float64
	memUsage  float64
)

func doWork() {
	defer wg.Done()
	// Simulate work
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	rand.Seed(time.Now().UnixNano())

	go monitorLoad()

	wg.Add(10) // Start with some initial workers
	for i := 0; i < 10; i++ {
		go func() {
			workCh <- struct{}{}
			defer func() { <-workCh }()
			doWork()
		}()
	}

	wg.Wait()
}

func monitorLoad() {
	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()

	for range ticker.C {
		var err error
		cpuPercent, err = cpu.Percent(time.Second, false)
		if err != nil {
			fmt.Println("Error getting CPU usage:", err)
			continue
		}

		v, err := mem.VirtualMemory()
		if err != nil {
			fmt.Println("Error getting memory usage:", err)
			continue
		}
		memUsage = float64(v.Used) / float64(v.Total) * 100

		mu.Lock()
		defer mu.Unlock()

		cpuUsage = cpuPercent
		currentConcurrency := len(workCh)
		newConcurrency := adjustConcurrency(cpuUsage, memUsage, currentConcurrency)
		diff := newConcurrency - currentConcurrency

		if diff > 0 {
			for i := 0; i < diff; i++ {
				workCh <- struct{}{}
			}
			fmt.Printf("Increased concurrency to %d: CPU %.2f%% Mem %.2f%%\n", newConcurrency, cpuUsage, memUsage)
		} else if diff < 0 {
			for i := 0; i > diff; i-- {
				<-workCh
			}
			fmt.Printf("Decreased concurrency to %d: CPU %.2f%% Mem %.2f%%\n", newConcurrency, cpuUsage, memUsage)
		}
	}
}

func adjustConcurrency(cpuUsage, memUsage float64, currentConcurrency int) int {
	load := (100 - cpuUsage) * adjustmentFactor + (100 - memUsage) * adjustmentFactor

	if load < 0 {
		return max(minConcurrency, currentConcurrency/2)
	} else if load > 100 {
		return min(maxConcurrency, currentConcurrency*2)
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