package utils

import (
	"crypto/rand"
	"encoding/hex"
)

func RandHex(n int) (string, error) {
	bytes := make([]byte, n)
	_, err := rand.Read(bytes)

	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
