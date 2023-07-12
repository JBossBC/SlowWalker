package utils

import (
	"crypto/md5"
	"encoding/hex"
)

func Encrypt(name string) string {
	byt := md5.Sum([]byte(name))
	return hex.EncodeToString(byt[:])
}
