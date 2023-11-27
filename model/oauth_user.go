package model

import "github.com/jiebutech/uc/types"

type OauthUser struct {
	Id         int64           `json:"id" gorm:"id"`
	App        string          `json:"app" gorm:"app"`
	UserId     string          `json:"user_id" gorm:"user_id"` // openid， 跟php版本保持一致
	BindUserId int64           `json:"bind_user_id" gorm:"bind_user_id"`
	LoginType  types.LoginType `json:"login_type" gorm:"login_type"`
}

func (u *OauthUser) LoginTypeKey() string {
	return "login_type"
}

func (u *OauthUser) TableName() string {
	return "oauth_user"
}

func (u *OauthUser) GetBindUserId() int64 {
	return u.BindUserId
}

func (u *OauthUser) SetBindUserId(userid int64) {
	u.BindUserId = userid
}

func (u *OauthUser) SetLoginType(loginType types.LoginType) {
	u.LoginType = loginType
}

func (u *OauthUser) GetLoginType() types.LoginType {
	return u.LoginType
}

func (u *OauthUser) GetOpenid() string {
	return u.UserId
}

func (u *OauthUser) SetOpenid(openid string) {
	u.UserId = openid
}

func (u *OauthUser) OpenidKey() string {
	return "user_id"
}

func (u *OauthUser) GetApp() string {
	return u.App
}

func (u *OauthUser) SetApp(app string) {
	u.App = app
}

func (u *OauthUser) AppKey() string {
	return "app"
}
