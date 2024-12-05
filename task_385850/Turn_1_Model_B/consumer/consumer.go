package consumer

import (
	"fmt"

)

func Consume(numbers chan int) {
	for num := range numbers {
		fmt.Println("Consumed:", num)
	}
	fmt.Println("Consumer finished")
}