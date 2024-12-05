package main

import (
	"fmt"

	"github.com/example/worker/consumer"
	"github.com/example/worker/producer"
)

func main() {
	// Create a buffer channel
	numbers := make(chan int, producer.channelBufferSize)

	go producer.Produce(numbers)
	go consumer.Consume(numbers)

	fmt.Println("Main waiting...")
	// Main waits for both producer and consumer to finish
}
