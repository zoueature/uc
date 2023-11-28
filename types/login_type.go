package types

import "github.com/jiebutech/uc/oauth"

type OauthLoginType interface {
	New(config oauth.Config) oauth.Oauth
	LoginType() LoginType
}

const (
	facebookIdentify = "facebook"
	googleIdentify   = "google"
)

var (
	// FacebookLoginType facebook登录方式
	FacebookLoginType OauthLoginType = facebookLoginType(facebookIdentify)
	// GoogleLoginType google登录方式
	GoogleLoginType OauthLoginType = googleLoginType(googleIdentify)
)

// facebook登录方式相关实现
type facebookLoginType string

func (f facebookLoginType) LoginType() LoginType {
	return LoginType(f)
}

func (f facebookLoginType) New(config oauth.Config) oauth.Oauth {
	return oauth.NewFacebookOauth(config)
}

// google登录方式相关实现
type googleLoginType string

func (f googleLoginType) LoginType() LoginType {
	return LoginType(f)
}

func (f googleLoginType) New(config oauth.Config) oauth.Oauth {
	return oauth.NewGoogleOauth(config)
}
