package main

import (
	"fmt"
	"math/big"
)

// Function to perform high-precision calculations with reused big.Int instances
func highPrecisionCalculation(n int) *big.Int {
	// Reusable big.Int instances
	var result, one, temp big.Int

	// Initialise result, one, and temp
	result.SetInt64(0)
	one.SetInt64(1)
	temp.SetInt64(0)

	// Loop for multiplication from 1 to n
	for i := 1; i <= n; i++ {
		// Reuse temp for each iteration
		temp.Mul(&one, big.NewInt(int64(i)))
		result.Mul(&result, &temp)
	}

	return &result
}

func main() {
	// Example usage
	n := 10
	fmt.Println("Factorial of", n, ":", highPrecisionCalculation(n))
}