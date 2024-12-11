package crypto_utils

import "crypto/sha256"

// SHA256Sum calculates the SHA-256 hash of the data.
func SHA256Sum(data []byte) []byte {
	h := sha256.New()
	h.Write(data)
	return h.Sum(nil)
}