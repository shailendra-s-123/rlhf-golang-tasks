package crypto_utils

import "crypto/aes"

// AESEncrypt encrypts data using AES-128 CBC mode.
func AESEncrypt(key []byte, plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// TODO: Implement CBC mode encryption
	return nil, nil
}