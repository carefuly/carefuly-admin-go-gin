/**
 * Description：
 * FileName：userPassword.go
 * Author：CJiaの用心
 * Create：2025/3/26 13:09:26
 * Remark：
 */

package dao

import (
	"context"
	model "github.com/carefuly/carefuly-admin-go-gin/internal/model/system"
	"gorm.io/gorm"
)

type UserPasswordDAO interface {
	Insert(ctx context.Context, u model.UserPassword) error

	ExistsByUserName(ctx context.Context, username string) (bool, error)
}

type GORMUserPasswordDAO struct {
	db *gorm.DB
}

func NewUserPasswordDAO(db *gorm.DB) UserPasswordDAO {
	return &GORMUserPasswordDAO{
		db: db,
	}
}

func (dao *GORMUserPasswordDAO) Insert(ctx context.Context, u model.UserPassword) error {
	return dao.db.WithContext(ctx).Create(&u).Error
}

func (dao *GORMUserPasswordDAO) ExistsByUserName(ctx context.Context, username string) (bool, error) {
	var count int64
	err := dao.db.WithContext(ctx).Model(&model.UserPassword{}).
		Where("username = ?", username).
		Count(&count).Error
	return count > 0, err
}
