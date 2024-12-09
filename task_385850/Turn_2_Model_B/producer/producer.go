package producer

import (
	"math/rand"
	"time"
)

// StartProducing is the main function of the producer service
func StartProducing(dataChannel chan int) {
	fmt.Println("Producer service started...")
	for {
		randomData := rand.Intn(100)
		dataChannel <- randomData
		time.Sleep(2 * time.Second) // Generate data every 2 seconds
	}
}