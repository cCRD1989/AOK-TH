package unit

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"math/rand"
	"os"
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

var charset = []byte(os.Getenv("MY_SECRET_KEY"))

func RandStr(n int) string {
	b := make([]byte, n)
	for i := range b {
		// randomly select 1 character from given charset
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
