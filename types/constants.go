package types

import (
	"errors"
	"fmt"
)

var (
	UndefinedCodeTypeErr = errors.New("undefined verify code type")
	PasswordNotMathErr   = errors.New("password not match")
	CodeNotMathErr       = errors.New("code not match")
	EqualOldPasswordErr  = errors.New("password equal the old one")
	UserExistsErr        = errors.New("user exists")
)

type LoginType string

const (
	EmailLogin  LoginType = "email"
	MobileLogin LoginType = "mobile"
)

type VerifyCodeType string

func (v VerifyCodeType) CacheKey() (CacheKey, error) {
	c, ok := codeTypeKeyMap[v]
	if !ok {
		return "", UndefinedCodeTypeErr
	}
	return c, nil
}

const (
	RegisterCodeType VerifyCodeType = "register"
	LoginCodeType    VerifyCodeType = "login"
	PasswordCodeType VerifyCodeType = "password"
)

var codeTypeKeyMap = map[VerifyCodeType]CacheKey{
	RegisterCodeType: registerCodeCacheKey,
	LoginCodeType:    loginCodeCacheKey,
	PasswordCodeType: passwordCodeCacheKey,
}

type CacheKey string

func (c CacheKey) CacheKey(args ...interface{}) string {
	return fmt.Sprintf(string(c), args...)
}

const (
	registerCodeCacheKey CacheKey = "register_code_%s_%s"
	loginCodeCacheKey    CacheKey = "login_code_%s_%s"
	passwordCodeCacheKey CacheKey = "password_code_%s_%s"
)
