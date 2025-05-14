/**
 * Description：
 * FileName：user.go
 * Author：CJiaの用心
 * Create：2025/5/12 15:07:23
 * Remark：
 */

package system

import (
	"context"
	"errors"
	domainSystem "github.com/carefuly/carefuly-admin-go-gin/internal/domain/careful/system"
	"github.com/carefuly/carefuly-admin-go-gin/internal/model/careful/system"
	"github.com/carefuly/carefuly-admin-go-gin/pkg/ginx/filters"
	"gorm.io/gorm"
)

var (
	ErrUserNotFound             = gorm.ErrRecordNotFound
	ErrUsernameDuplicate        = errors.New("用户名已存在")
	ErrUserDuplicate            = errors.New("用户信息已存在")
	ErrUserVersionInconsistency = errors.New("数据已被修改，请刷新后重试")
)

type UserDAO interface {
	Insert(ctx context.Context, model system.User) error
	Delete(ctx context.Context, id string) (int64, error)
	Update(ctx context.Context, model system.User) error
	UpdatePassword(ctx context.Context, userId string, hashedPassword string) error

	FindById(ctx context.Context, id string) (*system.User, error)
	FindByUsername(ctx context.Context, username string) (*system.User, error)
	FindListPage(ctx context.Context, filter domainSystem.UserFilter) ([]*system.User, int64, error)
	FindListAll(ctx context.Context, filter domainSystem.UserFilter) ([]*system.User, error)

	CheckExistByUsername(ctx context.Context, username, excludeId string) (bool, error)
}

type GORMUserDAO struct {
	db *gorm.DB
}

func NewGORMUserDAO(db *gorm.DB) UserDAO {
	return &GORMUserDAO{
		db: db,
	}
}

// Insert 新增
func (dao *GORMUserDAO) Insert(ctx context.Context, model system.User) error {
	return dao.db.WithContext(ctx).Create(&model).Error
}

// Delete 删除
func (dao *GORMUserDAO) Delete(ctx context.Context, id string) (int64, error) {
	result := dao.db.WithContext(ctx).Where("id = ?", id).Delete(&system.User{})
	return result.RowsAffected, result.Error
}

// Update 更新
func (dao *GORMUserDAO) Update(ctx context.Context, model system.User) error {
	result := dao.db.WithContext(ctx).Model(&model).
		Where("id = ? AND version = ?", model.Id, model.Version).
		Updates(map[string]any{
			"name":     model.Name,
			"gender":   model.Gender,
			"email":    model.Email,
			"mobile":   model.Mobile,
			"avatar":   model.Avatar,
			"version":  gorm.Expr("version + 1"),
			"modifier": model.Modifier,
			"remark":   model.Remark,
		})
	// 处理行影响数为0的情况
	if result.RowsAffected == 0 {
		// 先检查记录是否存在
		var exists bool
		dao.db.WithContext(ctx).
			Model(&system.User{}).
			Select("1").
			Where("id = ?", model.Id).
			Limit(1).
			Find(&exists)

		if !exists {
			return ErrUserNotFound
		}
		return ErrUserVersionInconsistency
	}

	return result.Error
}

// UpdatePassword 更新密码
func (dao *GORMUserDAO) UpdatePassword(ctx context.Context, userId string, hashedPassword string) error {
	result := dao.db.WithContext(ctx).Model(&system.User{}).
		Where("id = ?", userId).
		Update("password", hashedPassword)
	if result.RowsAffected == 0 {
		return ErrUserNotFound
	}
	return result.Error
}

// FindById 根据id获取详情
func (dao *GORMUserDAO) FindById(ctx context.Context, id string) (*system.User, error) {
	var user system.User
	err := dao.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &user, ErrUserNotFound
		}
		return &user, err
	}
	return &user, nil
}

// FindByUsername 根据用户名获取详情
func (dao *GORMUserDAO) FindByUsername(ctx context.Context, username string) (*system.User, error) {
	var user system.User
	err := dao.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &user, ErrUserNotFound
		}
		return &user, err
	}
	return &user, nil
}

// FindListPage 分页查询
func (dao *GORMUserDAO) FindListPage(ctx context.Context, filter domainSystem.UserFilter) ([]*system.User, int64, error) {
	var total int64
	var models []*system.User

	query := dao.db.WithContext(ctx).Model(&system.User{})

	if filter.Username != "" {
		query = query.Where("username LIKE ?", "%"+filter.Username+"%")
	}

	err := query.Count(&total).
		Offset((filter.Page - 1) * filter.PageSize).
		Limit(filter.PageSize).
		Find(&models).Error

	return models, total, err
}

// FindListAll 获取所有列表
func (dao *GORMUserDAO) FindListAll(ctx context.Context, filter domainSystem.UserFilter) ([]*system.User, error) {
	var models []*system.User

	query := dao.db.WithContext(ctx).Model(&system.User{})

	if filter.Username != "" {
		query = query.Where("username LIKE ?", "%"+filter.Username+"%")
	}

	// 查询
	if err := query.Find(&models).Error; err != nil {
		return nil, err
	}

	return models, nil
}

// buildQuery 构建查询条件
func (dao *GORMUserDAO) buildQuery(ctx context.Context, filter domainSystem.UserFilter) *gorm.DB {
	builder := &domainSystem.UserFilter{
		Filters: filters.Filters{
			Creator:  filter.Creator,
			Modifier: filter.Modifier,
		},
		Username: filter.Username,
	}
	return builder.Apply(dao.db.WithContext(ctx).Model(&system.User{}))
}

// CheckExistByUsername 检查用户名是否存在
func (dao *GORMUserDAO) CheckExistByUsername(ctx context.Context, username, excludeId string) (bool, error) {
	var user system.User
	query := dao.db.WithContext(ctx).Model(&system.User{}).
		Select("id"). // 只查询必要的字段
		Where("username = ?", username)

	if excludeId != "" {
		query = query.Where("id != ?", excludeId)
	}

	// 使用 LIMIT 1 快速判断是否存在
	err := query.Limit(1).First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil // 不存在
	}
	return err == nil, err // 存在或查询出错
}
