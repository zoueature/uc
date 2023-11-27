package model

import "github.com/jiebutech/uc/types"

type User struct {
	Id        int64           `json:"id" gorm:"id"`
	App       string          `json:"app" gorm:"app"`
	LoginType types.LoginType `json:"loginType" gorm:"login_type"`
	Identify  string          `json:"identify" gorm:"identify"`
	Username  string          `json:"username" gorm:"username"`
	Nickname  string          `json:"nickname" gorm:"nickname"`
	Avatar    string          `json:"avatar" gorm:"avatar"`
	Password  string          `json:"password" gorm:"password"`
}

func (u *User) SetId(id int64) {
	u.Id = id
}

func (u *User) SetIdentify(identify string) {
	u.Identify = identify
}

func (u *User) SetLoginType(loginType types.LoginType) {
	u.LoginType = loginType
}

func (u *User) SetUsername(username string) {
	u.Username = username
}

func (u *User) IdKey() string {
	return "id"
}

func (u *User) LoginTypeKey() string {
	return "login_type"
}

func (u *User) IdentifyKey() string {
	return "identify"
}

func (u *User) UsernameKey() string {
	return "username"
}

func (u *User) PasswordKey() string {
	return "password"
}

func (u *User) TableName() string {
	return "user"
}

func (u *User) GetID() int64 {
	return u.Id
}

func (u *User) GetIdentify() string {
	return u.Identify
}

func (u *User) GetLoginType() types.LoginType {
	return u.LoginType
}

func (u *User) GetUserName() string {
	return u.Username
}

func (u *User) GetNickname() string {
	return u.Nickname
}

func (u *User) GetAvatar() string {
	return u.Avatar
}

func (u *User) GetPassword() string {
	return u.Password
}

func (u *User) SetNickname(nickname string) {
	u.Nickname = nickname
}

func (u *User) SetAvatar(avatar string) {
	u.Avatar = avatar
}

func (u *User) SetPassword(password string) {
	u.Password = password
}

func (u *User) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"id":        u.Id,
		"loginType": u.LoginType,
		"identify":  u.Identify,
		"username":  u.Username,
		"nickname":  u.Nickname,
		"avatar":    u.Avatar,
		"app":       u.App,
	}
}

func (u *User) GetApp() string {
	return u.App
}

func (u *User) SetApp(app string) {
	u.App = app
}

func (u *User) AppKey() string {
	return "app"
}
