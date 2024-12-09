package main

import (
	"os"

	"github.com/example/task/consumer"
	"github.com/example/task/producer"
)

func main() {
	// Create a channel for communication between producer and consumer
	numCh := make(chan int)

	// Start the producer and consumer goroutines
	go producer.Produce(numCh)
	go consumer.Consume(numCh)

	// Keep the main program running
	select {}
}