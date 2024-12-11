package main

import (
	"fmt"
	"math/big"
)

// bigIntPool is a pool of pre-allocated big.Int objects
var bigIntPool = make(chan *big.Int, 1000)

func init() {
	for i := 0; i < 1000; i++ {
		bigIntPool <- new(big.Int)
	}
}

// getBigInt retrieves an object from the pool
func getBigInt() *big.Int {
	select {
	case b := <-bigIntPool:
		return b
	default:
		return new(big.Int)
	}
}

// putBigInt returns an object to the pool
func putBigInt(b *big.Int) {
	select {
	case bigIntPool <- b:
	default:
	}
}

func main() {
	// Reusable big.Int structures
	var x, y big.Int
	x.SetInt64(1234567890)
	y.SetInt64(987654321)

	// Perform multiple calculations with reused and pooled objects
	for i := 0; i < 100000; i++ {
		result := getBigInt() // Retrieve a reused or new big.Int from the pool
		result.Mul(&x, &y)

		// Do something with the result
		// ...

		putBigInt(result) // Return the result to the pool for reuse
	}

	fmt.Println("Result:", result.String())
}