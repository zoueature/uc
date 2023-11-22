package model

import (
	"errors"
	"fmt"
)
import "gorm.io/gorm"

type UserGenerator func() UserEntity

func DefaultUserRepo(orm *gorm.DB) UserResource {
	return &UserRepo{
		db: orm,
		userGen: func() UserEntity {
			return new(User)
		},
	}
}

type UserRepo struct {
	db      *gorm.DB
	userGen UserGenerator
}

func (u UserRepo) GenUser() UserEntity {
	return u.userGen()
}

// GetUserByIdentify 根据标识符获取用户信息
func (u *UserRepo) GetUserByIdentify(dest UserEntity) error {
	whereStr := fmt.Sprintf("%s = ? AND %s = ?", dest.IdentifyKey(), dest.LoginTypeKey())
	return u.db.Table(dest.TableName()).
		Where(whereStr, dest.GetIdentify(), dest.GetLoginType()).
		Find(dest).Error
}

// GetUserByUsername 根据用户名获取用户信息
func (u *UserRepo) GetUserByUsername(dest UserEntity) error {
	whereStr := fmt.Sprintf("%s = ? AND %s = ?", dest.UsernameKey(), dest.LoginTypeKey())
	return u.db.Table(dest.TableName()).
		Where(whereStr, dest.GetUserName(), dest.GetLoginType()).
		Find(dest).Error
}

// UpdatePassword 修改密码
func (u *UserRepo) UpdatePassword(dest UserEntity) error {
	if dest.GetID() == 0 || dest.GetPassword() == "" {
		return errors.New("empty user")
	}
	return u.db.Table(dest.TableName()).
		Where(dest.IdKey()+" = ?", dest.GetID()).
		Update(dest.PasswordKey(), dest.GetPassword()).Error
}
