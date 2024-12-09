package producer

import (
	"math/rand"
	"time"
)

// Produce generates random numbers and sends them to the provided channel
func Produce(ch chan int) {
	rand.Seed(time.Now().UnixNano())
	for {
		num := rand.Intn(100)
		ch <- num
		time.Sleep(1 * time.Second)
	}
}