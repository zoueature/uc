package uc

import (
	"crypto/sha1"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jiebutech/uc/model"
	"github.com/jiebutech/uc/types"
	"github.com/spf13/cast"
	"time"
)

type UserIdentify struct {
	App      string
	Type     types.LoginType
	Identify string
}

type UserInfo struct {
	UserIdentify
	Password Password
	Avatar   string
	Nickname string
	Username string
}
type BasicUserInfo struct {
	Id     int64
	Name   string
	Avatar string
}

type SupportModifyUserInfo struct {
	Name   string
	Avatar string
}

type Password struct {
	Salt     string
	Password string
}

func (p Password) marshalPassword() string {
	sh := sha1.New()
	sh.Write([]byte(p.Salt + p.Password))
	sum := sh.Sum(nil)
	return fmt.Sprintf("%x", sum)
}

type JwtEncoder interface {
	encodeJwt(user model.UserEntity) (string, error)
	decodeJwt(tokenStr string) (*BasicUserInfo, error)
}

func DefaultJwtEncoder(secret string, ttl int) JwtEncoder {
	return &jwtCli{
		secret: secret,
		ttl:    time.Duration(ttl),
	}
}

type jwtCli struct {
	secret string
	ttl    time.Duration
}

type myClaim struct {
	Id     int64  `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
	jwt.RegisteredClaims
}

func (c *jwtCli) encodeJwt(user model.UserEntity) (string, error) {
	claim := &myClaim{
		Id:     user.GetID(),
		Name:   user.GetNickname(),
		Avatar: user.GetAvatar(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(c.ttl * time.Second)),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	return token.SignedString([]byte(c.secret))
}

func (c *jwtCli) decodeJwt(tokenStr string) (*BasicUserInfo, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(c.secret), nil
	})
	if err != nil {
		return nil, err
	}
	claim, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("claim error")
	}
	return &BasicUserInfo{
		Id:     cast.ToInt64(claim["id"]),
		Name:   cast.ToString(claim["name"]),
		Avatar: cast.ToString(claim["avatar"]),
	}, nil
}
