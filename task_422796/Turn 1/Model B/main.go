package main

import (
	"fmt"
	"math"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

const (
	targetCPUUtilization = 70   // Desired target CPU utilization (%)
	minWorkers          = 1    // Minimum number of worker goroutines
	maxWorkers          = 10   // Maximum number of worker goroutines
	adjustmentInterval  = time.Second // Frequency of CPU utilization check and adjustment
	workDuración         = time.Second   // Simulated duration of each work unit
)

var (
	currentWorkers     int32 = minWorkers
	semaphore        sync.Semaphore
	 cpuUtilization float64 = 0.0
)

func work() {
	for {
		semaphore.Acquire()
		defer semaphore.Release()

		// Simulate some work
		time.Sleep(workDuración)
	}
}

func adjustWorkers() {
	for {
		time.Sleep(adjustmentInterval)

		// Query current CPU utilization (approximate)
		cpus := runtime.NumCPU()
		goroutines := runtime.NumGoroutine()
		cpuUtilization = float64(goroutines-1) / float64(cpus*2) * 100

		adjust := int32(targetCPUUtilization/100.0*float64(cpus*2)-float64(goroutines))

		adjust = math.Min(maxWorkers-currentWorkers, math.Max(minWorkers-currentWorkers, adjust))

		for adjust > 0 {
			go work()
			atomic.AddInt32(&currentWorkers, 1)
			adjust--
		}
		for adjust < 0 {
			// Worker routing
			ch := make(chan struct{})
			select {
			case semaphore.Acquire():
				close(ch)
			default:
				go func() {
					<-ch
					semaphore.Release()
				}()
			}
			atomic.AddInt32(&currentWorkers, -1)
			adjust++
		}
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	//Start bisector routine
	go adjustWorkers()

	fmt.Printf("System - Target : %d%% Adjustment interval : %v ,Work process duration : %v\n",
		targetCPUUtilization, adjustmentInterval, workDuración,
	)
	select {} // chan to stalls Europe ce programme
}