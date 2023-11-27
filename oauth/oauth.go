package oauth

type UserInfo struct {
	Openid   string
	Nickname string
	Avatar   string
	Email    string
}

type Oauth interface {
	GetAccessToken(code string) (string, error)
	GetOauthUserInfo(token string) (*UserInfo, error)
	GenAuthLoginURL() string
}

// Config 第三方登录配置
type Config interface {
	ClientId() string
	ClientSecret() string
	RedirectURI() string
}
