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
	logPrefix        = "concurrency-adjuster: "
)

var (
	wg       sync.WaitGroup
	mu       sync.Mutex
	workCh   = make(chan struct{}, 0)
	currentConcurrency int
)

func doWork() {
	defer wg.Done()
	// Simulate work using random CPU and memory usage
	time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
	useCPU := float64(rand.Intn(20))
	useMem := float64(rand.Intn(20))
	fmt.Printf("%sWorker finished: CPU %.1f%%, Mem %.1f%%\n", logPrefix, useCPU, useMem)
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	rand.Seed(time.Now().UnixNano())

	go monitorSystemLoad()

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

func monitorSystemLoad() {
	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()

	for range ticker.C {
		cpuPercent, err := cpu.Percent(time.Second, false)
		if err != nil {
			fmt.Println(logPrefix, "Error getting CPU usage:", err)
			continue
		}

		memStats, err := mem.VirtualMemory()
		if err != nil {
			fmt.Println(logPrefix, "Error getting memory usage:", err)
			continue
		}

		memUsage := (memStats.Used / memStats.Total) * 100

		mu.Lock()
		defer mu.Unlock()

		newConcurrency := adjustConcurrency(cpuPercent, memUsage, currentConcurrency)
		diff := newConcurrency - currentConcurrency

		if diff > 0 {
			for i := 0; i < diff; i++ {
				workCh <- struct{}{}
			}
			fmt.Printf(
				"%sAdjusted concurrency to %d: CPU %.1f%%, Mem %.1f%%\n",
				logPrefix,
				newConcurrency,
				cpuPercent,
				memUsage,
			)
		} else if diff < 0 {
			for i := 0; i > diff; i-- {
				<-workCh
			}
			fmt.Printf(
				"%sAdjusted concurrency to %d: CPU %.1f%%, Mem %.1f%%\n",
				logPrefix,
				newConcurrency,
				cpuPercent,
				memUsage,
			)
		}
		currentConcurrency = newConcurrency
	}
}

func adjustConcurrency(cpuPercent, memUsage float64, currentConcurrency int) int {
	if cpuPercent < targetCPUUsage/2 && memUsage < targetMemUsage/2 {
		return min(maxConcurrency, currentConcurrency*2)
	} else if cpuPercent > targetCPUUsage || memUsage > targetMemUsage {
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