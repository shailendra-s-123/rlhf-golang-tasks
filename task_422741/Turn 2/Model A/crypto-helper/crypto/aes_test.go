package crypto

import (
	"bytes"
	"testing"
)

const (
	sampleKey   = "12345678901234567890123456789012"
	sampleText = "Hello, this is a secret message."
)

func TestEncryptDecrypt(t *testing.T) {
	key := []byte(sampleKey)
	plaintext := []byte(sampleText)

	ciphertext, err := Encrypt(plaintext, key)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	decrypted, err := Decrypt(ciphertext, key)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}

	if string(decrypted) != sampleText {
		t.Errorf("Decrypted text does not match plaintext. Expected: %q, got: %q", sampleText, string(decrypted))
	}
}

func TestInvalidKeySize(t *testing.T) {
	key := []byte("123456789012345") // Invalid key size: 15 bytes instead of 32
	plaintext := []byte(sampleText)

	_, err := Encrypt(plaintext, key)
	if err != errInvalidKey {
		t.Errorf("Expected errInvalidKey, got %v", err)
	}
}

func TestCiphertextTooShort(t *testing.T) {
	key := []byte(sampleKey)
	plaintext := []byte(sampleText)

	ciphertext, err := Encrypt(plaintext, key)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	ciphertextTooShort := ciphertext[:1] // Truncate ciphertext

	_, err = Decrypt(ciphertextTooShort, key)
	if err == nil {
		t.Error("Expected error for ciphertext too short, but none was returned")
	}
}