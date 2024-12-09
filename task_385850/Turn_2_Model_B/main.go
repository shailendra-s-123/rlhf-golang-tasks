package main

import (
	"fmt"
	"git/example.com/seamless-communications/consumer"
	"git/example.com/seamless-communications/producer"
)

func main() {
	// Create a channel for data transmission
	dataCh := make(chan int)

	go producer.StartProducing(dataCh)
	go consumer.StartConsuming(dataCh)

	fmt.Println("Main application started. Producer and consumer running...")

	// The main loop to keep the program running
	for {
		time.Sleep(10 * time.Second)
	}
}