package oauth

import (
	"github.com/jiebutech/uc/model"
	"github.com/jiebutech/uc/types"
)

type Oauth interface {
	GetAccessToken(code string) (string, error)
	GetOauthUserInfo(token string) (*model.OauthUserInfo, error)
	GenAuthLoginURL() string
}

// Config 第三方登录配置
type Config interface {
	ClientId() string
	ClientSecret() string
	RedirectURI() string
}

type OauthOption struct {
	LoginType types.OauthLoginType
	Cfg       Config
}
