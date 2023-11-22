package model

import "github.com/jiebutech/uc/types"

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

	IdKey() string
	LoginTypeKey() string
	IdentifyKey() string
	UsernameKey() string
	PasswordKey() string

	SetId(id int64)
	SetIdentify(identify string)
	SetLoginType(loginType types.LoginType)
	SetUsername(username string)
}

// UserResource 用户资源, 集成用户的相关操作
type UserResource interface {
	GenUser() UserEntity
	GetUserByIdentify(dest UserEntity) error
	GetUserByUsername(dest UserEntity) error
	UpdatePassword(dest UserEntity) error
}
