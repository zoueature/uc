package model

import (
	"github.com/zoueature/uc/types"
	"gorm.io/gorm/schema"
)

type Entity interface {
	schema.Tabler
}

// UserEntity 用户实体， 表示具体的用户
type UserEntity interface {
	TableName() string
	ToMap() map[string]interface{}

	GetID() int64
	GetIdentify() string
	GetUserName() string
	GetNickname() string
	GetAvatar() string
	GetPassword() string
	GetLoginType() types.LoginType
	GetApp() string

	AppKey() string
	IdKey() string
	LoginTypeKey() string
	IdentifyKey() string
	UsernameKey() string
	PasswordKey() string
	GetChannel() string

	SetId(id int64)
	SetIdentify(identify string)
	SetLoginType(loginType types.LoginType)
	SetUsername(username string)
	SetNickname(nickname string)
	SetAvatar(avatar string)
	SetPassword(password string)
	SetApp(string2 string)
	SetChannel(channel string)
}

type OauthUserEntity interface {
	TableName() string
	GetBindUserId() int64
	SetBindUserId(uuserid int64)
	GetOpenid() string
	SetOpenid(openid string)
	SetLoginType(loginType types.LoginType)
	LoginTypeKey() string
	GetLoginType() types.LoginType
	OpenidKey() string
	AppKey() string
	GetApp() string
	SetApp(string2 string)
}

// UserResource 用户资源, 集成用户的相关操作
type UserResource interface {
	IsUserNotFound(err error) bool

	GenUser() UserEntity
	GenOauthUser() OauthUserEntity
	GetUserByIdentify(dest UserEntity) error
	GetUserById(dest UserEntity) error
	GetUserByUsername(dest UserEntity) error
	UpdatePassword(dest UserEntity) error
	SaveUser(dest Entity) error
	CreateUser(dest Entity) error

	// oauth user
	GetOauthByOpenid(dest OauthUserEntity) error

	TransactionCreate(tablers map[Entity]func()) error
	TransactionSave(tablers map[Entity]func()) error
}
