/**
 * Description：
 * FileName：post.go
 * Author：CJiaの用心
 * Create：2025/6/13 17:22:26
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
	ErrPostNotFound             = gorm.ErrRecordNotFound
	ErrPostDuplicate            = errors.New("岗位已存在")
	ErrPostVersionInconsistency = errors.New("数据已被修改，请刷新后重试")
)

type PostDAO interface {
	Insert(ctx context.Context, model system.Post) error
	Delete(ctx context.Context, id string) (int64, error)
	BatchDelete(ctx context.Context, ids []string) error
	Update(ctx context.Context, model system.Post) error

	FindById(ctx context.Context, id string) (*system.Post, error)
	FindListPage(ctx context.Context, filter domainSystem.PostFilter) ([]*system.Post, int64, error)
	FindListAll(ctx context.Context, filter domainSystem.PostFilter) ([]*system.Post, error)

	CheckExistByNameAndCode(ctx context.Context, name, code, excludeId string) (bool, error)
}

type GORMPostDAO struct {
	db *gorm.DB
}

func NewGORMPostDAO(db *gorm.DB) PostDAO {
	return &GORMPostDAO{
		db: db,
	}
}

// Insert 新增
func (dao *GORMPostDAO) Insert(ctx context.Context, model system.Post) error {
	return dao.db.WithContext(ctx).Create(&model).Error
}

// Delete 删除
func (dao *GORMPostDAO) Delete(ctx context.Context, id string) (int64, error) {
	result := dao.db.WithContext(ctx).Where("id = ?", id).Delete(&system.Post{})
	return result.RowsAffected, result.Error
}

// BatchDelete 批量删除
func (dao *GORMPostDAO) BatchDelete(ctx context.Context, ids []string) error {
	return dao.db.WithContext(ctx).Where("id IN ?", ids).Delete(&system.Post{}).Error
}

// Update 更新
func (dao *GORMPostDAO) Update(ctx context.Context, model system.Post) error {
	result := dao.db.WithContext(ctx).Model(&model).
		Where("id = ? AND version = ?", model.Id, model.Version).
		Updates(map[string]any{
			"name":     model.Name,
			"code":     model.Code,
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
			Model(&system.Post{}).
			Select("1").
			Where("id = ?", model.Id).
			Limit(1).
			Find(&exists)

		if !exists {
			return ErrPostNotFound
		}
		return ErrPostVersionInconsistency
	}

	return result.Error
}

// FindById 根据id获取详情
func (dao *GORMPostDAO) FindById(ctx context.Context, id string) (*system.Post, error) {
	var model system.Post
	err := dao.db.WithContext(ctx).Where("id = ?", id).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &model, ErrPostNotFound
		}
		return &model, err
	}
	return &model, err
}

// FindListPage 分页查询
func (dao *GORMPostDAO) FindListPage(ctx context.Context, filter domainSystem.PostFilter) ([]*system.Post, int64, error) {
	var total int64
	var models []*system.Post

	query := dao.buildQuery(ctx, filter)

	err := query.Count(&total).
		Offset((filter.Page - 1) * filter.PageSize).
		Limit(filter.PageSize).
		Find(&models).Error

	return models, total, err
}

// FindListAll 获取所有列表
func (dao *GORMPostDAO) FindListAll(ctx context.Context, filter domainSystem.PostFilter) ([]*system.Post, error) {
	var models []*system.Post

	query := dao.buildQuery(ctx, filter)

	// 查询
	if err := query.Find(&models).Error; err != nil {
		return nil, err
	}

	return models, nil
}

// buildQuery 构建查询条件
func (dao *GORMPostDAO) buildQuery(ctx context.Context, filter domainSystem.PostFilter) *gorm.DB {
	builder := &domainSystem.PostFilter{
		Filters: filters.Filters{
			Creator:    filter.Creator,
			Modifier:   filter.Modifier,
			BelongDept: filter.BelongDept,
		},
		Status: filter.Status,
		Name:   filter.Name,
		Code:   filter.Code,
	}
	return builder.Apply(dao.db.WithContext(ctx).Model(&system.Post{}))
}

// CheckExistByNameAndCode 检查name、code是否同时存在
func (dao *GORMPostDAO) CheckExistByNameAndCode(ctx context.Context, name, code, excludeId string) (bool, error) {
	var model system.Post
	query := dao.db.WithContext(ctx).Model(&system.Post{}).
		Select("id"). // 只查询必要的字段
		Where("name = ? AND code = ?", name, code)

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
