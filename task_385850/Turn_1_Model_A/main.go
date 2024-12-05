package main

import (
	"github.com/example/worker/worker"// Replace with your actual import path
	"time"

	
)

func main() {
	dataChannel := make(chan int)

	go worker.WorkerProcessor(dataChannel)

	time.Sleep(1 * time.Second)
	dataChannel <- 10
	time.Sleep(1 * time.Second)
	dataChannel <- 20

	// Graceful shutdown after sending some data
	time.Sleep(2 * time.Second)
	close(dataChannel)
	time.Sleep(2 * time.Second)
}