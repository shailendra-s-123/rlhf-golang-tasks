package main

import (
	"fmt"
	"math/big"
)

func main() {
	// Reusable big.Int structures
	var x, y, result big.Int
	x.SetInt64(1234567890)
	y.SetInt64(987654321)

	// Perform multiple calculations with the same structures
	for i := 0; i < 100000; i++ {
		result.Mul(&x, &y)
		// Do something with the result
		// ...
	}

	fmt.Println("Result:", result.String())
}