package ucgo

import (
	"github.com/jiebutech/uc/cache"
	"github.com/jiebutech/uc/model"
	"github.com/jiebutech/uc/sender"
	"github.com/jiebutech/uc/types"
)

type InternalClient struct {
	cache    cache.Cache
	sender   sender.SmsCodeSender
	userRepo model.UserResource
}

func NewInternalClient(cache cache.Cache, sender sender.SmsCodeSender, repo model.UserResource) *InternalClient {
	return &InternalClient{
		cache:    cache,
		sender:   sender,
		userRepo: repo,
	}
}

// Login 用户登录并返回token
func (c *InternalClient) Login(t types.LoginType, identify, password string) (string, model.UserEntity, error) {
	user := c.userRepo.GenUser()
	user.SetLoginType(t)
	user.SetIdentify(identify)
	err := c.userRepo.GetUserByIdentify(user)
	if err != nil {
		return "", nil, err
	}
	if user.GetPassword() != marshalPassword(password) {
		return "", nil, types.PasswordNotMathErr
	}
	token, err := encodeJwt(user)
	if err != nil {
		return "", nil, err
	}
	return token, user, nil
}

func (c InternalClient) LoginByUsername(t types.LoginType, username, password string) (string, model.UserEntity, error) {
	user := c.userRepo.GenUser()
	user.SetLoginType(t)
	user.SetUsername(username)
	err := c.userRepo.GetUserByUsername(user)
	if err != nil {
		return "", nil, err
	}
	if user.GetPassword() != marshalPassword(password) {
		return "", nil, types.PasswordNotMathErr
	}
	token, err := encodeJwt(user)
	if err != nil {
		return "", nil, err
	}
	return token, user, nil
}

// SendSmsCode 验证码发送
func (c *InternalClient) SendSmsCode(t types.VerifyCodeType, identify string) error {
	ckey, err := t.CacheKey()
	if err != nil {
		return err
	}
	code := generateSmsCode()
	// 缓存验证码
	err = c.cache.Set(ckey.CacheKey(identify), code, codeCacheTTL)
	if err != nil {
		return err
	}
	// 调用发送器发送验证码
	return c.sender.Send(code, identify)
}
