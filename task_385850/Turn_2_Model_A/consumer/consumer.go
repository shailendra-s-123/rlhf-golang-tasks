package consumer

import (
	"fmt"
)

// Consume receives numbers from the channel and prints them
func Consume(ch chan int) {
	for num := range ch {
		fmt.Println("Consumed:", num)
	}
}