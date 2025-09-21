package util

import (
	"crypto/sha256"
	"encoding/hex"
)

func HashPlainTextWithSHA256(plaintext string) string {
	hash := sha256.Sum256([]byte(plaintext))
	return hex.EncodeToString(hash[:])
}
