package ucgo

import (
	"github.com/jiebutech/uc/cache"
	"github.com/jiebutech/uc/model"
	"github.com/jiebutech/uc/sender"
	"github.com/jiebutech/uc/types"
	"sync"
)

type UserClient struct {
	cache    cache.Cache
	sender   sender.SmsCodeSender
	userRepo model.UserResource
}

var ucCli *UserClient
var onceNewUc = sync.Once{}

// NewUserClient 实例化用户操作客户端
func NewUserClient(cache cache.Cache, sender sender.SmsCodeSender, repo model.UserResource) *UserClient {
	onceNewUc.Do(func() {
		ucCli = &UserClient{
			cache:    cache,
			sender:   sender,
			userRepo: repo,
		}
	})
	return ucCli
}

// Login 用户登录并返回token
func (c *UserClient) Login(t types.LoginType, identify string, password Password) (string, model.UserEntity, error) {
	user := c.userRepo.GenUser()
	user.SetLoginType(t)
	user.SetIdentify(identify)
	err := c.userRepo.GetUserByIdentify(user)
	if err != nil {
		return "", nil, err
	}
	if user.GetPassword() != password.marshalPassword() {
		return "", nil, types.PasswordNotMathErr
	}
	token, err := encodeJwt(user)
	if err != nil {
		return "", nil, err
	}
	return token, user, nil
}

// LoginByUsername 用户名密码登录并返回token
func (c UserClient) LoginByUsername(t types.LoginType, username string, password Password) (string, model.UserEntity, error) {
	user := c.userRepo.GenUser()
	user.SetLoginType(t)
	user.SetUsername(username)
	err := c.userRepo.GetUserByUsername(user)
	if err != nil {
		return "", nil, err
	}
	if user.GetPassword() != password.marshalPassword() {
		return "", nil, types.PasswordNotMathErr
	}
	token, err := encodeJwt(user)
	if err != nil {
		return "", nil, err
	}
	return token, user, nil
}

type UserInfo struct {
	Type     types.LoginType
	Identify string
	Password Password
	Avatar   string
	Nickname string
	Username string
}

// Register 注册
func (c UserClient) Register(code string, info UserInfo) (string, model.UserEntity, error) {
	ok, err := c.checkCode(types.RegisterCodeType, info.Identify, code)
	if err != nil {
		return "", nil, err
	}
	if !ok {
		return "", nil, types.CodeNotMathErr
	}
	return c.register(info)
}

// RegisterWithNoCode 无验证码注册
func (c UserClient) RegisterWithNoCode(info UserInfo) (string, model.UserEntity, error) {
	return c.register(info)
}

func (c UserClient) register(info UserInfo) (string, model.UserEntity, error) {
	user := c.userRepo.GenUser()
	user.SetLoginType(info.Type)
	user.SetIdentify(info.Identify)
	user.SetUsername(info.Username)
	user.SetPassword(info.Password.marshalPassword())
	user.SetAvatar(info.Avatar)
	user.SetNickname(info.Nickname)
	err := c.userRepo.GetUserByIdentify(user)
	if err == nil {
		return "", nil, types.UserExistsErr
	}
	err = c.userRepo.GetUserByUsername(user)
	if err == nil {
		return "", nil, types.UserExistsErr
	}
	err = c.userRepo.CreateUser(user)
	if err != nil {
		return "", nil, err
	}
	token, err := encodeJwt(user)
	if err != nil {
		return "", nil, err
	}
	return token, user, nil
}

func (c *UserClient) checkCode(t types.VerifyCodeType, identify string, code string) (bool, error) {
	ckey, err := t.CacheKey()
	if err != nil {
		return false, err
	}
	cacheCode := c.cache.Get(ckey.CacheKey(identify))
	return cacheCode == code, nil
}

// SendSmsCode 验证码发送
func (c *UserClient) SendSmsCode(t types.VerifyCodeType, identify string) error {
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

// ChangePasswordByCode 根据验证码修改密码
func (c *UserClient) ChangePasswordByCode(t types.LoginType, identify, code string, newPassword Password) error {
	ok, err := c.checkCode(types.PasswordCodeType, identify, code)
	if err != nil {
		return err
	}
	if !ok {
		return types.CodeNotMathErr
	}
	user := c.userRepo.GenUser()
	user.SetIdentify(identify)
	user.SetLoginType(t)
	err = c.userRepo.GetUserByIdentify(user)
	if err != nil {
		return err
	}
	newPswd := newPassword.marshalPassword()
	if newPswd == user.GetPassword() {
		return types.EqualOldPasswordErr
	}
	user.SetPassword(newPswd)
	return c.userRepo.SaveUser(user)
}

// ChangePasswordByOld 根据旧密码修改密码
func (c *UserClient) ChangePasswordByOld(t types.LoginType, identify string, oldPassword, newPassword Password) error {
	user := c.userRepo.GenUser()
	user.SetIdentify(identify)
	user.SetLoginType(t)
	err := c.userRepo.GetUserByIdentify(user)
	if err != nil {
		return err
	}
	oldPswd := oldPassword.marshalPassword()
	if oldPswd != user.GetPassword() {
		return types.PasswordNotMathErr
	}
	newPswd := newPassword.marshalPassword()
	if newPswd == user.GetPassword() {
		return types.EqualOldPasswordErr
	}
	user.SetPassword(newPswd)
	return c.userRepo.SaveUser(user)
}
