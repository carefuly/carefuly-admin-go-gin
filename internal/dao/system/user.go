/**
 * Description：
 * FileName：user.go
 * Author：CJiaの用心
 * Create：2025/3/26 13:07:02
 * Remark：
 */

package dao

import (
	"context"
	"errors"
	model "github.com/carefuly/carefuly-admin-go-gin/internal/model/system"
	"gorm.io/gorm"
)

var (
	ErrDuplicateUsername = errors.New("用户账号")
	ErrUserNotFound   = gorm.ErrRecordNotFound
)

type UserDAO interface {
	Insert(ctx context.Context, u model.User) error

	ExistsByUserName(ctx context.Context, username string) (bool, error)
}

type GORMUserDAO struct {
	db *gorm.DB
}

func NewGORMUserDAO(db *gorm.DB) UserDAO {
	return &GORMUserDAO{
		db: db,
	}
}

func (dao *GORMUserDAO) Insert(ctx context.Context, u model.User) error {
	return dao.db.WithContext(ctx).Create(&u).Error
}

func (dao *GORMUserDAO) ExistsByUserName(ctx context.Context, username string) (bool, error) {
	var count int64
	err := dao.db.WithContext(ctx).Model(&model.User{}).
		Where("username = ?", username).
		Count(&count).Error
	return count > 0, err
}