package ucgo

// OauthConfig 第三方登录配置
type OauthConfig interface {
	ClientId() string
	ClientSecret() string
	RedirectURI() string
}
