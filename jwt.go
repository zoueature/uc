package ucgo

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/jiebutech/uc/model"
	"time"
)

const loginAtKey = "loginAt"

var jwtKey = "dasjodas6899u0ouj9912bk891"

func SetJwtKey(key string) {
	jwtKey = key
}

// encodeJwt jwt token 签名
func encodeJwt(user model.UserEntity) (string, error) {
	um := user.ToMap()
	um[loginAtKey] = time.Now().Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(um))
	return token.SignedString([]byte(jwtKey))
}
