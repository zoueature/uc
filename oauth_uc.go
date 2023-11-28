package ucgo

import (
	"fmt"
	"github.com/jiebutech/uc/model"
	"github.com/jiebutech/uc/oauth"
	"github.com/jiebutech/uc/types"
	"sync"
)

type OauthClient struct {
	oauthCliMap sync.Map
	userRepo    model.UserResource
	jwt         JwtEncoder
}

var oauthCli *OauthClient
var onceNewOauth = sync.Once{}

// NewOauthClient 实例化第三方登录客户端
func NewOauthClient(userRepo model.UserResource, jwtClient JwtEncoder) *OauthClient {
	onceNewOauth.Do(func() {
		oauthCli = &OauthClient{
			userRepo:    userRepo,
			oauthCliMap: sync.Map{},
			jwt:         jwtClient,
		}
	})
	return oauthCli
}

type OauthOption struct {
	App       string
	LoginType types.OauthLoginType
	Cfg       oauth.Config
}

func storeOauthCliKey(app string, loginType types.OauthLoginType) string {
	return app + "-" + string(loginType.LoginType())
}

// WithLoginType 注入登录方式
func (o *OauthClient) WithLoginType(option OauthOption, cover ...bool) *OauthClient {
	k := storeOauthCliKey(option.App, option.LoginType)
	_, ok := o.oauthCliMap.Load(k)
	if !ok || (len(cover) > 0 && cover[0]) {
		o.oauthCliMap.Store(k, option.LoginType.New(option.Cfg))
	}
	return o
}

func (o *OauthClient) oauthCli(app string, loginType types.OauthLoginType) (oauth.Oauth, error) {
	oc, ok := o.oauthCliMap.Load(storeOauthCliKey(app, loginType))
	if !ok {
		return nil, fmt.Errorf("login type not configure")
	}
	return oc.(oauth.Oauth), nil
}

// AuthURL 生成web端的授权地址
func (o *OauthClient) AuthURL(app string, loginType types.OauthLoginType) (string, error) {
	cli, err := o.oauthCli(app, loginType)
	if err != nil {
		return "", err
	}
	return cli.GenAuthLoginURL(), nil
}

// Login 登录
func (o *OauthClient) Login(app string, loginType types.OauthLoginType, code string) (string, model.UserEntity, error) {
	oauthCli, err := o.oauthCli(app, loginType)
	if err != nil {
		return "", nil, err
	}
	// 获取access token
	accessToken, err := oauthCli.GetAccessToken(code)
	if err != nil {
		return "", nil, err
	}
	// 获取用户信息
	oauthUser, err := oauthCli.GetOauthUserInfo(accessToken)
	if err != nil {
		return "", nil, err
	}
	oauthUserEntity := o.userRepo.GenOauthUser()
	oauthUserEntity.SetOpenid(oauthUser.Openid)
	oauthUserEntity.SetLoginType(loginType.LoginType())
	oauthUserEntity.SetApp(app)
	err = o.userRepo.GetOauthByOpenid(oauthUserEntity)
	user := o.userRepo.GenUser()
	if err != nil {
		if !o.userRepo.IsUserNotFound(err) {
			return "", nil, err
		}
		// 注册逻辑
		user, err = o.register(app, loginType, oauthUser)
		if err != nil {
			return "", nil, err
		}
	} else {
		// 已经注册过， 查询对应数据信息
		err = o.userRepo.GetOauthByOpenid(oauthUserEntity)
		if err != nil {
			return "", nil, err
		}
		user.SetId(oauthUserEntity.GetBindUserId())
		err = o.userRepo.GetUserById(user)
		if err != nil {
			return "", nil, err
		}
	}
	token, err := o.jwt.encodeJwt(user)
	if err != nil {
		return "", nil, err
	}
	return token, user, nil

}

func (o *OauthClient) register(app string, loginType types.OauthLoginType, oauthUser *oauth.UserInfo) (model.UserEntity, error) {
	user := o.userRepo.GenUser()
	user.SetLoginType(loginType.LoginType())
	user.SetUsername(fmt.Sprintf("%s-%s", string(loginType.LoginType()), oauthUser.Openid))
	user.SetIdentify(oauthUser.Openid)
	user.SetNickname(oauthUser.Nickname)
	user.SetAvatar(oauthUser.Avatar)
	user.SetApp(app)

	oauthUserEntity := o.userRepo.GenOauthUser()
	oauthUserEntity.SetBindUserId(user.GetID())
	oauthUserEntity.SetOpenid(oauthUser.Openid)
	oauthUserEntity.SetLoginType(loginType.LoginType())
	oauthUserEntity.SetApp(app)

	err := o.userRepo.TransactionCreate(map[model.Entity]func(){
		user: func() {
			oauthUserEntity.SetBindUserId(user.GetID())
		},
	})
	if err != nil {
		return nil, err
	}
	return user, nil
}
