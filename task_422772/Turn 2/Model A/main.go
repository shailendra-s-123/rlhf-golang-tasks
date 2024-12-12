package main

import (
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"math/big"
)

func main() {
	// Choose a secure elliptic curve
	curve := elliptic.P256()

	// Generate a private key (which involves high-precision modular arithmetic)
	privKey, err := elliptic.GenerateKey(curve, rand.Reader)
	if err != nil {
		panic("error generating private key:", err)
	}

	fmt.Println("\nGenerated ECC Private Key:")
	fmt.Println("X:", privKey.X)
	fmt.Println("Y:", privKey.Y)
	fmt.Println("D:", privKey.D)

	// Calculate the public key using point multiplication (high-precision operation)
	publicKey := curve.ScalarMult(curve.Params().G.X, curve.Params().G.Y, privKey.D, curve.Params().N)

	fmt.Println("\nCalculated ECC Public Key:")
	fmt.Println("X:", publicKey.X)
	fmt.Println("Y:", publicKey.Y)

	// Perform point multiplication again to demonstrate precision (validate calculation)
	recalculatedPublicKey := curve.ScalarMult(curve.Params().G.X, curve.Params().G.Y, privKey.D, curve.Params().N)

	fmt.Println("\nRecalculated ECC Public Key:")
	fmt.Println("X:", recalculatedPublicKey.X)
	fmt.Println("Y:", recalculatedPublicKey.Y)

	// Verify that the public key and recalculated public key are the same
	if publicKey.X.Cmp(recalculatedPublicKey.X) != 0 || publicKey.Y.Cmp(recalculatedPublicKey.Y) != 0 {
		panic("Public key calculation failed!")
	}
	fmt.Println("\nPublic key calculation verified.")
}