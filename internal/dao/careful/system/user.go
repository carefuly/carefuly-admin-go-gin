/**
 * Description：
 * FileName：user.go
 * Author：CJiaの用心
 * Create：2025/3/28 11:20:54
 * Remark：
 */

package system

import (
	"context"
	"errors"
	"github.com/carefuly/carefuly-admin-go-gin/internal/model/careful/system"
	"gorm.io/gorm"
)

var (
	ErrDuplicateUsername = errors.New("用户账号冲突")
	ErrUserNotFound      = gorm.ErrRecordNotFound
)

type UserDAO interface {
	Insert(ctx context.Context, u system.User) error

	FindByUserName(ctx context.Context, username string) (*system.User, error)

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

func (dao *GORMUserDAO) Insert(ctx context.Context, u system.User) error {
	return dao.db.WithContext(ctx).Create(&u).Error
}

func (dao *GORMUserDAO) FindByUserName(ctx context.Context, username string) (*system.User, error) {
	var user system.User
	err := dao.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	return &user, err
}

func (dao *GORMUserDAO) ExistsByUserName(ctx context.Context, username string) (bool, error) {
	var count int64
	err := dao.db.WithContext(ctx).Model(&system.User{}).
		Where("username = ?", username).
		Count(&count).Error
	return count > 0, err
}
