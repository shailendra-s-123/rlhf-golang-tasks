package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"math/big"
)

func main() {
	// Generate an ECDSA key pair using the P-256 curve
	curve := elliptic.P256()
	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		panic("error generating ECDSA key:", err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		panic("public key is not of type *ecdsa.PublicKey")
	}

	fmt.Println("\nGenerated ECDSA Key Pair:")
	fmt.Println("Private Key:")
	fmt.Println(privateKey)
	fmt.Println("Public Key:")
	fmt.Println(publicKeyECDSA)

	// Example of secure computation: ECDSA signature generation
	message := []byte("Hello, this is a secure message!")
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, message)
	if err != nil {
		panic("error generating ECDSA signature:", err)
	}

	fmt.Println("\nGenerated ECDSA Signature:")
	fmt.Println("R:", r)
	fmt.Println("S:", s)

	// Example of secure computation: ECDSA signature verification
	if ecdsa.Verify(publicKeyECDSA, message, r, s) {
		fmt.Println("\nSignature verification: Success")
	} else {
		fmt.Println("\nSignature verification: Failure")
	}
}