package helpers

import (
	"crypto/sha256"
	"encoding/hex"
)

func SHA256Sum(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}
