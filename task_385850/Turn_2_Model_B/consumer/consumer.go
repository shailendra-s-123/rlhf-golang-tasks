package consumer

import (
	"fmt"
)

// StartConsuming  is the main function of the consumer service.
func StartConsuming(dataChannel chan int) {
	fmt.Println("Consumer service started...")
	for {
		receivedData := <-dataChannel
		fmt.Println("Consumer received data:", receivedData)
	}
}