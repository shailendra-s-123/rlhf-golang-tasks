package main

import (
	"myapp/consumer"
	"myapp/producer"
	_"os"
)

func main() {
	// Create a channel to communicate between producer and consumer
	numCh := make(chan int)

	// Start the producer and consumer goroutines
	go producer.Produce(numCh)
	go consumer.Consume(numCh)
	

	// Keep the main program running
	select {}
}