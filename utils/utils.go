package utils

import (
	"crypto/rand"
	"encoding/hex"
)

// RandHex simply creates random hex string
// To use to query session data
func RandHex(n int) (string, error) {
	bytes := make([]byte, n)
	_, err := rand.Read(bytes)

	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
