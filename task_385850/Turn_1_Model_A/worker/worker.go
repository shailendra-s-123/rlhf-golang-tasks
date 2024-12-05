package worker

import (
	"fmt"
	"time"
)

// WorkerProcessor is a function that processes messages received from a channel
func WorkerProcessor(dataChannel chan int) {
	for {
		select {
		case data := <-dataChannel:
			processedData := data * data
			fmt.Printf("Worker: Received %d, Processed: %d\n", data, processedData)
		case <-time.After(5 * time.Second):
			fmt.Println("Worker: Shutting down...")
			return
		}
	}
}