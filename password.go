package ucgo

import (
	"crypto/sha1"
	"encoding/hex"
)

type Password struct {
	Salt     string
	Password string
}

func (p Password) marshalPassword() string {
	sh := sha1.New()
	sum := sh.Sum([]byte(p.Salt + p.Password))
	return hex.EncodeToString(sum)
}
