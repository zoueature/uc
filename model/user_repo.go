package model

import (
	"errors"
	"fmt"
)
import "gorm.io/gorm"

type userGenerator func() UserEntity
type oauthUserGenerator func() OauthUserEntity

func DefaultUserRepo(orm *gorm.DB) UserResource {
	return &UserRepo{
		db: orm,
		userGen: func() UserEntity {
			return new(User)
		},
		oauthUserGen: func() OauthUserEntity {
			return new(OauthUser)
		},
	}
}

type UserRepo struct {
	db           *gorm.DB
	userGen      userGenerator
	oauthUserGen oauthUserGenerator
}

func (repo *UserRepo) IsUserNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

func (repo *UserRepo) GetOauthByOpenid(dest OauthUserEntity) error {
	where := fmt.Sprintf("%s = ? AND %s = ? AND %s = ?", dest.AppKey(), dest.OpenidKey(), dest.LoginTypeKey())
	return repo.db.Table(dest.TableName()).Where(where, dest.GetApp(), dest.GetOpenid(), dest.GetLoginType()).First(dest).Error
}

func (repo *UserRepo) GetUserById(dest UserEntity) error {
	where := fmt.Sprintf("%s = ?", dest.IdKey())
	return repo.db.Table(dest.TableName()).Where(where, dest.GetID()).First(dest).Error
}

func (repo *UserRepo) GenOauthUser() OauthUserEntity {
	return repo.oauthUserGen()
}

func (repo UserRepo) GenUser() UserEntity {
	return repo.userGen()
}

// GetUserByIdentify 根据标识符获取用户信息
func (repo *UserRepo) GetUserByIdentify(dest UserEntity) error {
	whereStr := fmt.Sprintf("%s = ? AND %s = ? AND %s = ?", dest.AppKey(), dest.IdentifyKey(), dest.LoginTypeKey())
	return repo.db.Table(dest.TableName()).
		Where(whereStr, dest.GetApp(), dest.GetIdentify(), dest.GetLoginType()).
		First(dest).Error
}

// GetUserByUsername 根据用户名获取用户信息
func (repo *UserRepo) GetUserByUsername(dest UserEntity) error {
	whereStr := fmt.Sprintf("%s = ? AND %s = ? AND %s = ?", dest.AppKey(), dest.UsernameKey(), dest.LoginTypeKey())
	return repo.db.Table(dest.TableName()).
		Where(whereStr, dest.GetApp(), dest.GetUserName(), dest.GetLoginType()).
		First(dest).Error
}

// UpdatePassword 修改密码
func (repo *UserRepo) UpdatePassword(dest UserEntity) error {
	if dest.GetID() == 0 || dest.GetPassword() == "" {
		return errors.New("empty user")
	}
	return repo.db.Table(dest.TableName()).
		Where(dest.IdKey()+" = ?", dest.GetID()).
		Update(dest.PasswordKey(), dest.GetPassword()).Error
}

func (repo *UserRepo) SaveUser(dest Entity) error {
	return repo.db.Table(dest.TableName()).Save(dest).Error
}

func (repo *UserRepo) CreateUser(dest Entity) error {
	return repo.db.Table(dest.TableName()).Create(dest).Error
}

// TransactionCreate 事务创建
func (repo *UserRepo) TransactionCreate(tablers map[Entity]func()) error {
	return repo.db.Transaction(func(tx *gorm.DB) error {
		for tabler, do := range tablers {
			err := tx.Table(tabler.TableName()).Create(tabler).Error
			if err != nil {
				return err
			}
			do()
		}
		return nil
	})
}

// TransactionSave 事务保存
func (repo *UserRepo) TransactionSave(tablers map[Entity]func()) error {
	return repo.db.Transaction(func(tx *gorm.DB) error {
		for tabler, do := range tablers {
			err := tx.Table(tabler.TableName()).Save(tabler).Error
			if err != nil {
				return err
			}
			do()
		}
		return nil
	})
}
