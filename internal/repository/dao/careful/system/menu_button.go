/**
 * Description：
 * FileName：menu_button.go
 * Author：CJiaの用心
 * Create：2025/6/9 13:09:44
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
	ErrMenuButtonNotFound             = gorm.ErrRecordNotFound
	ErrMenuButtonVersionInconsistency = errors.New("数据已被修改，请刷新后重试")
)

type MenuButtonDAO interface {
	Insert(ctx context.Context, model system.MenuButton) error
	Delete(ctx context.Context, id string) (int64, error)
	BatchDelete(ctx context.Context, ids []string) error
	Update(ctx context.Context, model system.MenuButton) error

	FindById(ctx context.Context, id string) (*system.MenuButton, error)
	FindListPage(ctx context.Context, filter domainSystem.MenuButtonFilter) ([]*system.MenuButton, int64, error)
	FindListAll(ctx context.Context, filter domainSystem.MenuButtonFilter) ([]*system.MenuButton, error)
}

type GORMMenuButtonDAO struct {
	db *gorm.DB
}

func NewGORMMenuButtonDAO(db *gorm.DB) MenuButtonDAO {
	return &GORMMenuButtonDAO{
		db: db,
	}
}

// Insert 新增
func (dao *GORMMenuButtonDAO) Insert(ctx context.Context, model system.MenuButton) error {
	return dao.db.WithContext(ctx).Create(&model).Error
}

// Delete 删除
func (dao *GORMMenuButtonDAO) Delete(ctx context.Context, id string) (int64, error) {
	result := dao.db.WithContext(ctx).Where("id = ?", id).Delete(&system.MenuButton{})
	return result.RowsAffected, result.Error
}

// BatchDelete 批量删除
func (dao *GORMMenuButtonDAO) BatchDelete(ctx context.Context, ids []string) error {
	return dao.db.WithContext(ctx).Where("id IN ?", ids).Delete(&system.MenuButton{}).Error
}

// Update 更新
func (dao *GORMMenuButtonDAO) Update(ctx context.Context, model system.MenuButton) error {
	result := dao.db.WithContext(ctx).Model(&model).
		Where("id = ? AND version = ?", model.Id, model.Version).
		Updates(map[string]any{
			"name":     model.Name,
			"code":     model.Code,
			"api":      model.Api,
			"method":   model.Method,
			"sort":     model.Sort,
			"status":   model.Status,
			"version":  gorm.Expr("version + 1"),
			"modifier": model.Modifier,
			"remark":   model.Remark,
		})
	// 处理行影响数为0的情况
	if result.RowsAffected == 0 {
		// 先检查记录是否存在
		var exists bool
		dao.db.WithContext(ctx).
			Model(&system.MenuButton{}).
			Select("1").
			Where("id = ?", model.Id).
			Limit(1).
			Find(&exists)

		if !exists {
			return ErrMenuButtonNotFound
		}
		return ErrMenuButtonVersionInconsistency
	}
	return result.Error
}

// FindById 根据id获取详情
func (dao *GORMMenuButtonDAO) FindById(ctx context.Context, id string) (*system.MenuButton, error) {
	var model system.MenuButton
	err := dao.db.WithContext(ctx).Where("id = ?", id).
		Preload("Menu").
		First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &model, ErrMenuButtonNotFound
		}
		return &model, err
	}
	return &model, err
}

// FindListPage 分页查询
func (dao *GORMMenuButtonDAO) FindListPage(ctx context.Context, filter domainSystem.MenuButtonFilter) ([]*system.MenuButton, int64, error) {
	var total int64
	var models []*system.MenuButton

	query := dao.buildQuery(ctx, filter)

	err := query.Count(&total).
		Offset((filter.Page - 1) * filter.PageSize).
		Limit(filter.PageSize).
		Find(&models).Error

	return models, total, err
}

// FindListAll 获取所有列表
func (dao *GORMMenuButtonDAO) FindListAll(ctx context.Context, filter domainSystem.MenuButtonFilter) ([]*system.MenuButton, error) {
	var models []*system.MenuButton

	query := dao.buildQuery(ctx, filter)

	// 查询
	if err := query.Find(&models).Error; err != nil {
		return nil, err
	}

	return models, nil
}

// buildQuery 构建查询条件
func (dao *GORMMenuButtonDAO) buildQuery(ctx context.Context, filter domainSystem.MenuButtonFilter) *gorm.DB {
	builder := &domainSystem.MenuButtonFilter{
		Filters: filters.Filters{
			Creator:    filter.Creator,
			Modifier:   filter.Modifier,
			BelongDept: filter.BelongDept,
		},
		Status: filter.Status,
		Name:   filter.Name,
		Code:   filter.Code,
		MenuId: filter.MenuId,
	}
	return builder.Apply(dao.db.WithContext(ctx).Model(&system.MenuButton{}))
}
