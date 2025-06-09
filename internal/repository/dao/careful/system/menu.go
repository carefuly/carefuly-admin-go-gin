/**
 * Description：
 * FileName：menu.go
 * Author：CJiaの用心
 * Create：2025/5/13 16:29:32
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
	ErrMenuNotFound             = gorm.ErrRecordNotFound
	ErrMenuNameDuplicate        = errors.New("菜单名称已存在")
	ErrMenuDuplicate            = errors.New("菜单已存在")
	ErrMenuVersionInconsistency = errors.New("数据已被修改，请刷新后重试")
)

type MenuDAO interface {
	Insert(ctx context.Context, model system.Menu) error
	Delete(ctx context.Context, id string) (int64, error)
	BatchDelete(ctx context.Context, ids []string) error
	Update(ctx context.Context, model system.Menu) error

	FindById(ctx context.Context, id string) (*system.Menu, error)
	FindListAll(ctx context.Context, filter domainSystem.MenuFilter) ([]*system.Menu, error)

	CheckExistByTypeAndTitleAndParentId(ctx context.Context, menuType int, title, parentId, excludeId string) (bool, error)
}

type GORMMenuDAO struct {
	db *gorm.DB
}

func NewGORMMenuDAO(db *gorm.DB) MenuDAO {
	return &GORMMenuDAO{
		db: db,
	}
}

// Insert 新增
func (dao *GORMMenuDAO) Insert(ctx context.Context, model system.Menu) error {
	return dao.db.WithContext(ctx).Create(&model).Error
}

// Delete 删除
func (dao *GORMMenuDAO) Delete(ctx context.Context, id string) (int64, error) {
	result := dao.db.WithContext(ctx).Where("id = ?", id).Delete(&system.Menu{})
	return result.RowsAffected, result.Error
}

// BatchDelete 批量删除
func (dao *GORMMenuDAO) BatchDelete(ctx context.Context, ids []string) error {
	return dao.db.WithContext(ctx).Where("id IN ?", ids).Delete(&system.Menu{}).Error
}

// Update 更新
func (dao *GORMMenuDAO) Update(ctx context.Context, model system.Menu) error {
	result := dao.db.WithContext(ctx).Model(&model).
		Where("id = ? AND version = ?", model.Id, model.Version).
		Updates(map[string]any{
			"type":        model.Type,
			"icon":        model.Icon,
			"title":       model.Title,
			"name":        model.Name,
			"component":   model.Component,
			"path":        model.Path,
			"redirect":    model.Redirect,
			"isHide":      model.IsHide,
			"isLink":      model.IsLink,
			"isKeepAlive": model.IsKeepAlive,
			"isFull":      model.IsFull,
			"parent_id":   model.ParentID,
			"sort":        model.Sort,
			"version":     gorm.Expr("version + 1"),
			"modifier":    model.Modifier,
			"remark":      model.Remark,
		})
	// 处理行影响数为0的情况
	if result.RowsAffected == 0 {
		// 先检查记录是否存在
		var exists bool
		dao.db.WithContext(ctx).
			Model(&system.Menu{}).
			Select("1").
			Where("id = ?", model.Id).
			Limit(1).
			Find(&exists)

		if !exists {
			return ErrMenuNotFound
		}
		return ErrMenuVersionInconsistency
	}
	return result.Error
}

// FindById 根据id获取详情
func (dao *GORMMenuDAO) FindById(ctx context.Context, id string) (*system.Menu, error) {
	var model system.Menu
	err := dao.db.WithContext(ctx).Where("id = ?", id).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &model, ErrMenuNotFound
		}
		return &model, err
	}
	return &model, err
}

// FindListAll 获取所有列表
func (dao *GORMMenuDAO) FindListAll(ctx context.Context, filter domainSystem.MenuFilter) ([]*system.Menu, error) {
	var models []*system.Menu

	query := dao.buildQuery(ctx, filter)

	// 查询
	if err := query.Find(&models).Error; err != nil {
		return nil, err
	}

	return models, nil
}

// buildQuery 构建查询条件
func (dao *GORMMenuDAO) buildQuery(ctx context.Context, filter domainSystem.MenuFilter) *gorm.DB {
	builder := &domainSystem.MenuFilter{
		Filters: filters.Filters{
			Creator:  filter.Creator,
			Modifier: filter.Modifier,
		},
		Title: filter.Title,
	}
	return builder.Apply(dao.db.WithContext(ctx).Model(&system.Menu{}))
}

// CheckExistByTypeAndTitleAndParentId 检查type、title和parentId是否同时存在
func (dao *GORMMenuDAO) CheckExistByTypeAndTitleAndParentId(ctx context.Context, menuType int, title, parentId, excludeId string) (bool, error) {
	var model system.Menu
	query := dao.db.WithContext(ctx).Model(&system.Menu{}).
		Select("id"). // 只查询必要的字段
		Where("type = ? AND title = ? AND parent_id = ?", menuType, title, parentId)

	if excludeId != "" {
		query = query.Where("id != ?", excludeId)
	}

	// 使用 LIMIT 1 快速判断是否存在
	err := query.Limit(1).First(&model).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil // 不存在
	}
	return err == nil, err // 存在或查询出错
}
