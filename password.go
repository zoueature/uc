package ucgo

import (
	"crypto/sha1"
	"encoding/hex"
)

var salt = "KASd9U)ud0a8ud0hIY*&^^daugsdsa"

func SetSalt(s string) {
	salt = s
}

func marshalPassword(password string) string {
	sh := sha1.New()
	sum := sh.Sum([]byte(salt + password))
	return hex.EncodeToString(sum)
}
