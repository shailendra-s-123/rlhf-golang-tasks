package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
)

const keySize = 32 // 256 bits

var errInvalidKey = errors.New("invalid key size")

func Encrypt(plaintext []byte, key []byte) ([]byte, error) {
	if len(key) != keySize {
		return nil, errInvalidKey
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return nonce + ciphertext, nil
}

func Decrypt(ciphertext []byte, key []byte) ([]byte, error) {
	if len(key) != keySize {
		return nil, errInvalidKey
	}

	if len(ciphertext) < aes.BlockSize+cipher.NonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:cipher.NonceSize], ciphertext[cipher.NonceSize:]
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}