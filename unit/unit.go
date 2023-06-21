package unit

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"math/rand"
)

func HashMD5(pass string) string {
	h := md5.New()
	io.WriteString(h, pass)
	return hex.EncodeToString(h.Sum(nil))
}

func index(slice []int, item int) int {
	for i := range slice {
		if slice[i] == item {
			return i
		}
	}
	return -1
}

func GenerateSecureToken(length int) string {
	b := make([]byte, (length / 2))
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}
