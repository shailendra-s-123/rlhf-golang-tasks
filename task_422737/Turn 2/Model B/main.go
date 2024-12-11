package main

import (
	"fmt"
	"math/big"
	"sync"
)

// Object pool for big.Int
type intPool struct {
	sync.Pool
}

func (p *intPool) Get() *big.Int {
	if i := p.Get(); i != nil {
		return i.(*big.Int)
	}
	return new(big.Int)
}

func (p *intPool) Put(*big.Int) {
	p.Put(&big.Int{}) // Clear the value before putting back into the pool
}

func main() {
	var intPool intPool

	// Reusable big.Int structures
	var x, y big.Int
	x.SetInt64(1234567890)
	y.SetInt64(987654321)

	// Function to perform calculations with pooled objects
	calculate := func() {
		result := intPool.Get()
		defer intPool.Put(result)
		result.Mul(&x, &y)
		// Do something with the result
		// ...
	}

	// Perform multiple calculations
	for i := 0; i < 100000; i++ {
		calculate()
	}

	fmt.Println("Result:", x.Mul(&x, &y).String())
}