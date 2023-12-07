package uc

import (
	"errors"
	"github.com/zoueature/uc/cache"
	"github.com/zoueature/uc/model"
	"github.com/zoueature/uc/sender"
	"github.com/zoueature/uc/types"
	"sync"
)

type UserClient struct {
	cache      cache.Cache
	userRepo   model.UserResource
	jwtClients sync.Map
	senders    sync.Map
}

var ucCli *UserClient
var onceNewUc = sync.Once{}

// NewUserClient 实例化用户操作客户端
func NewUserClient(cache cache.Cache, repo model.UserResource) *UserClient {
	onceNewUc.Do(func() {
		ucCli = &UserClient{
			cache:      cache,
			userRepo:   repo,
			jwtClients: sync.Map{},
			senders:    sync.Map{},
		}
	})
	return ucCli
}

func (c *UserClient) WithEncoder(app string, jwtClient JwtEncoder) {
	c.jwtClients.Store(app, jwtClient)
}

func (c *UserClient) WithSender(app string, sender sender.SmsCodeSender) {
	c.senders.Store(app, sender)
}

func (c *UserClient) sender(app string) (sender.SmsCodeSender, error) {
	s, ok := c.senders.Load(app)
	if !ok {
		return nil, errors.New(app + " sender not config")
	}
	return s.(sender.SmsCodeSender), nil
}

func (c *UserClient) jwt(app string) (JwtEncoder, error) {
	s, ok := c.jwtClients.Load(app)
	if !ok {
		return nil, errors.New(app + " jwt encoder not config")
	}
	return s.(JwtEncoder), nil
}

// Login 用户登录并返回token
func (c *UserClient) Login(id UserIdentify, password Password) (string, model.UserEntity, error) {
	user := c.userRepo.GenUser()
	user.SetLoginType(id.Type)
	user.SetIdentify(id.Identify)
	user.SetApp(id.App)
	err := c.userRepo.GetUserByIdentify(user)
	if err != nil {
		return "", nil, err
	}
	if user.GetPassword() != password.marshalPassword() {
		return "", nil, types.PasswordNotMathErr
	}
	jwt, err := c.jwt(id.App)
	if err != nil {
		return "", nil, err
	}
	token, err := jwt.encodeJwt(user)
	if err != nil {
		return "", nil, err
	}
	return token, user, nil
}

// LoginByUsername 用户名密码登录并返回token
func (c *UserClient) LoginByUsername(app string, t types.LoginType, username string, password Password) (string, model.UserEntity, error) {
	user := c.userRepo.GenUser()
	user.SetApp(app)
	user.SetLoginType(t)
	user.SetUsername(username)
	err := c.userRepo.GetUserByUsername(user)
	if err != nil {
		return "", nil, err
	}
	if user.GetPassword() != password.marshalPassword() {
		return "", nil, types.PasswordNotMathErr
	}
	jwt, err := c.jwt(app)
	if err != nil {
		return "", nil, err
	}
	token, err := jwt.encodeJwt(user)
	if err != nil {
		return "", nil, err
	}
	return token, user, nil
}

// Register 注册
func (c *UserClient) Register(code string, info UserInfo) (string, model.UserEntity, error) {
	ok, err := c.checkCode(types.RegisterCodeType, info.App, info.Identify, code)
	if err != nil {
		return "", nil, err
	}
	if !ok {
		return "", nil, types.CodeNotMathErr
	}
	return c.register(info)
}

// RegisterWithNoCode 无验证码注册
func (c *UserClient) RegisterWithNoCode(info UserInfo) (string, model.UserEntity, error) {
	return c.register(info)
}

// GetUserInfoById  根据id获取用户信息
func (c *UserClient) GetUserInfoById(id int64) (model.UserEntity, error) {
	user := c.userRepo.GenUser()
	user.SetId(id)
	err := c.userRepo.GetUserById(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// SaveUserProfile  保存用户信息
func (c *UserClient) SaveUserProfile(id int64, userInfo SupportModifyUserInfo) error {
	user, err := c.GetUserInfoById(id)
	if err != nil {
		return err
	}
	user.SetNickname(userInfo.Name)
	user.SetAvatar(userInfo.Avatar)
	return c.userRepo.SaveUser(user)
}

func (c *UserClient) register(info UserInfo) (string, model.UserEntity, error) {
	user := c.userRepo.GenUser()
	user.SetApp(info.App)
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
	jwt, err := c.jwt(info.App)
	if err != nil {
		return "", nil, err
	}
	token, err := jwt.encodeJwt(user)
	if err != nil {
		return "", nil, err
	}
	return token, user, nil
}

func (c *UserClient) checkCode(t types.VerifyCodeType, app, identify string, code string) (bool, error) {
	ckey, err := t.CacheKey()
	if err != nil {
		return false, err
	}
	cacheCode := c.cache.Get(ckey.CacheKey(app, identify))
	return cacheCode == code, nil
}

// SendSmsCode 验证码发送
func (c *UserClient) SendSmsCode(t types.VerifyCodeType, identify UserIdentify) error {
	ckey, err := t.CacheKey()
	if err != nil {
		return err
	}
	code := generateSmsCode()
	// 缓存验证码
	err = c.cache.Set(ckey.CacheKey(identify.App, identify.Identify), code, codeCacheTTL)
	if err != nil {
		return err
	}
	// 调用发送器发送验证码
	codeSender, err := c.sender(identify.App)
	if err != nil {
		return err
	}
	return codeSender.Send(code, identify.Identify)
}

// ChangePasswordByCode 根据验证码修改密码
func (c *UserClient) ChangePasswordByCode(identify UserIdentify, code string, newPassword Password) error {
	ok, err := c.checkCode(types.PasswordCodeType, identify.App, identify.Identify, code)
	if err != nil {
		return err
	}
	if !ok {
		return types.CodeNotMathErr
	}
	user := c.userRepo.GenUser()
	user.SetApp(identify.App)
	user.SetIdentify(identify.Identify)
	user.SetLoginType(identify.Type)
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
func (c *UserClient) ChangePasswordByOld(identify UserIdentify, oldPassword, newPassword Password) error {
	user := c.userRepo.GenUser()
	user.SetIdentify(identify.Identify)
	user.SetLoginType(identify.Type)
	user.SetApp(identify.App)
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
