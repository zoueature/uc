package types

import "github.com/jiebutech/uc/oauth"

type OauthLoginType interface {
	New(config oauth.Config) oauth.Oauth
	LoginType() LoginType
}

const (
	facebookIdentify = "facebook"
)

var (
	// facebook登录方式
	FacebookLoginType OauthLoginType = facebookLoginType(facebookIdentify)
)

type facebookLoginType string

func (f facebookLoginType) LoginType() LoginType {
	return LoginType(f)
}

func (f facebookLoginType) New(config oauth.Config) oauth.Oauth {
	return oauth.NewFacebookOauth(config)
}
