package common

import (
	"crypto/rand"
	"math/big"
)

// lowercase letters only — avoids visual ambiguity (0/O, 1/l) and is URL-safe without encoding
const charset = "abcdefghijklmnopqrstuvwxyz"

// GenerateID returns a random 15-character lowercase alphabetic ID
func GenerateID() string {
	const length int = 15
	b := make([]byte, length)

	for i := range b {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		b[i] = charset[n.Int64()]
	}

	return string(b)
}
