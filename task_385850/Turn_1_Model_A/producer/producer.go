package producer

import (
	"math/rand"
	"time"
)

// Produce sends random numbers to the given channel
func Produce(ch chan int) {
	for {
		num := rand.Intn(100)
		ch <- num
		time.Sleep(1 * time.Second)
	}
}