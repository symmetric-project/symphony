package utils

import (
	"crypto/rand"
	"encoding/hex"
)

func NewOctid() string {
	// Bytes worth 8 characters
	bytes := make([]byte, 4)

	// Randomize bytes
	rand.Read(bytes)

	// Encode the randomized bytes as a string
	return hex.EncodeToString(bytes)
}
