package unit

import (
	"crypto/md5"
	"encoding/hex"
	"io"
)

func HashMD5(pass string) string {
	h := md5.New()
	io.WriteString(h, pass)
	return hex.EncodeToString(h.Sum(nil))
}
