package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

func main() {
	// Generate a large prime number using high precision
	const primeSize = 2048 // Example prime size in bits
	prime := new(big.Int)

	// Generate a random number in the range [2^(primeSize-1), 2^primeSize)
	if err := rand.Read(prime.Bytes((primeSize + 7) >> 3)); err != nil {
		panic("Error generating random number:", err)
	}

	// Set the least significant bit to 1 to ensure oddness
	prime.SetBit(prime, 0, 1)

	// Check if the number is prime
	if !prime.ProbablyPrime(100) {
		// If it's not prime, try again
		// (For a production application, you would handle this in a loop)
		panic("Generated number is not prime")
	}

	fmt.Printf("Generated %d-bit prime number: %v\n", primeSize, prime)
}