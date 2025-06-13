/**
 * Description：
 * FileName：menu_column.go
 * Author：CJiaの用心
 * Create：2025/6/9 13:09:55
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
	ErrMenuColumnNotFound             = gorm.ErrRecordNotFound
	ErrMenuColumnVersionInconsistency = errors.New("数据已被修改，请刷新后重试")
)

type MenuColumnDAO interface {
	Insert(ctx context.Context, model system.MenuColumn) error
	Delete(ctx context.Context, id string) (int64, error)
	BatchDelete(ctx context.Context, ids []string) error
	Update(ctx context.Context, model system.MenuColumn) error

	FindById(ctx context.Context, id string) (*system.MenuColumn, error)
	FindListPage(ctx context.Context, filter domainSystem.MenuColumnFilter) ([]*system.MenuColumn, int64, error)
	FindListByMenuIds(ctx context.Context) ([]*system.MenuColumn, error)
	FindListAll(ctx context.Context, filter domainSystem.MenuColumnFilter) ([]*system.MenuColumn, error)
}

type GORMMenuColumnDAO struct {
	db *gorm.DB
}

func NewGORMMenuColumnDAO(db *gorm.DB) MenuColumnDAO {
	return &GORMMenuColumnDAO{
		db: db,
	}
}

// Insert 新增
func (dao *GORMMenuColumnDAO) Insert(ctx context.Context, model system.MenuColumn) error {
	return dao.db.WithContext(ctx).Create(&model).Error
}

// Delete 删除
func (dao *GORMMenuColumnDAO) Delete(ctx context.Context, id string) (int64, error) {
	result := dao.db.WithContext(ctx).Where("id = ?", id).Delete(&system.MenuColumn{})
	return result.RowsAffected, result.Error
}

// BatchDelete 批量删除
func (dao *GORMMenuColumnDAO) BatchDelete(ctx context.Context, ids []string) error {
	return dao.db.WithContext(ctx).Where("id IN ?", ids).Delete(&system.MenuColumn{}).Error
}

// Update 更新
func (dao *GORMMenuColumnDAO) Update(ctx context.Context, model system.MenuColumn) error {
	result := dao.db.WithContext(ctx).Model(&model).
		Where("id = ? AND version = ?", model.Id, model.Version).
		Updates(map[string]any{
			"title":    model.Title,
			"field":    model.Field,
			"width":    model.Width,
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
			Model(&system.MenuColumn{}).
			Select("1").
			Where("id = ?", model.Id).
			Limit(1).
			Find(&exists)

		if !exists {
			return ErrMenuColumnNotFound
		}
		return ErrMenuColumnVersionInconsistency
	}
	return result.Error
}

// FindById 根据id获取详情
func (dao *GORMMenuColumnDAO) FindById(ctx context.Context, id string) (*system.MenuColumn, error) {
	var model system.MenuColumn
	err := dao.db.WithContext(ctx).Where("id = ?", id).
		Preload("Menu").
		First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &model, ErrMenuColumnNotFound
		}
		return &model, err
	}
	return &model, err
}

// FindListPage 分页查询
func (dao *GORMMenuColumnDAO) FindListPage(ctx context.Context, filter domainSystem.MenuColumnFilter) ([]*system.MenuColumn, int64, error) {
	var total int64
	var models []*system.MenuColumn

	query := dao.buildQuery(ctx, filter)

	err := query.Count(&total).
		Offset((filter.Page - 1) * filter.PageSize).
		Limit(filter.PageSize).
		Find(&models).Error

	return models, total, err
}

// FindListByMenuIds 获取指定菜单下的所有列
func (dao *GORMMenuColumnDAO) FindListByMenuIds(ctx context.Context) ([]*system.MenuColumn, error) {
	var models []*system.MenuColumn

	query := dao.db.WithContext(ctx)

	// 查询
	if err := query.Find(&models).Error; err != nil {
		return nil, err
	}

	return models, nil
}

// FindListAll 获取所有列表
func (dao *GORMMenuColumnDAO) FindListAll(ctx context.Context, filter domainSystem.MenuColumnFilter) ([]*system.MenuColumn, error) {
	var models []*system.MenuColumn

	query := dao.buildQuery(ctx, filter)

	// 查询
	if err := query.Find(&models).Error; err != nil {
		return nil, err
	}

	return models, nil
}

// buildQuery 构建查询条件
func (dao *GORMMenuColumnDAO) buildQuery(ctx context.Context, filter domainSystem.MenuColumnFilter) *gorm.DB {
	builder := &domainSystem.MenuColumnFilter{
		Filters: filters.Filters{
			Creator:    filter.Creator,
			Modifier:   filter.Modifier,
			BelongDept: filter.BelongDept,
		},
		Status: filter.Status,
		Title:  filter.Title,
		Field:  filter.Field,
		MenuId: filter.MenuId,
	}
	return builder.Apply(dao.db.WithContext(ctx).Model(&system.MenuColumn{}))
}
