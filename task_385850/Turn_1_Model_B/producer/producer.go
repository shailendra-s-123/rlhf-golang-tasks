package producer

import (
	"math/rand"
	"time"
)

const (
	channelBufferSize = 10
)

func Produce(numbers chan int) {
	defer close(numbers)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 10; i++ {
		num := rand.Intn(100)
		select {
		case numbers <- num:
			time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
		default:
		}
	}
}