/**
 * Description：
 * FileName：bucket.go
 * Author：CJiaの用心
 * Create：2025/7/14 16:33:59
 * Remark：
 */

package tools

import (
	"context"
	"errors"
	domainTools "github.com/carefuly/carefuly-admin-go-gin/internal/domain/careful/tools"
	"github.com/carefuly/carefuly-admin-go-gin/internal/model/careful/tools"
	"github.com/carefuly/carefuly-admin-go-gin/pkg/ginx/filters"
	"gorm.io/gorm"
)

var (
	ErrBucketNotFound             = gorm.ErrRecordNotFound
	ErrBucketNameDuplicate        = errors.New("存储桶名称已存在")
	ErrBucketCodeDuplicate        = errors.New("存储桶编码已存在")
	ErrBucketDuplicate            = errors.New("存储桶已存在")
	ErrBucketVersionInconsistency = errors.New("数据已被修改，请刷新后重试")
)

type BucketDAO interface {
	Insert(ctx context.Context, model tools.Bucket) error
	Delete(ctx context.Context, id string) error
	BatchDelete(ctx context.Context, ids []string) error
	Update(ctx context.Context, model tools.Bucket) error

	FindById(ctx context.Context, id string) (*tools.Bucket, error)
	FindListPage(ctx context.Context, filter domainTools.BucketFilter) ([]*tools.Bucket, int64, error)
	FindListAll(ctx context.Context, filter domainTools.BucketFilter) ([]*tools.Bucket, error)

	CheckExistByCode(ctx context.Context, code, excludeId string) (bool, error)
	CheckExistByName(ctx context.Context, name, excludeId string) (bool, error)
}

type GORMBucketDAO struct {
	db *gorm.DB
}

func NewGORMBucketDAO(db *gorm.DB) BucketDAO {
	return &GORMBucketDAO{
		db: db,
	}
}

// Insert 新增
func (dao *GORMBucketDAO) Insert(ctx context.Context, model tools.Bucket) error {
	return dao.db.WithContext(ctx).Create(&model).Error
}

// Delete 删除
func (dao *GORMBucketDAO) Delete(ctx context.Context, id string) error {
	return dao.db.WithContext(ctx).Where("id = ?", id).Delete(&tools.Bucket{}).Error
}

// BatchDelete 批量删除
func (dao *GORMBucketDAO) BatchDelete(ctx context.Context, ids []string) error {
	return dao.db.WithContext(ctx).Where("id IN ?", ids).Delete(&tools.Bucket{}).Error
}

// Update 更新
func (dao *GORMBucketDAO) Update(ctx context.Context, model tools.Bucket) error {
	result := dao.db.WithContext(ctx).Model(&model).
		Where("id = ? AND version = ?", model.Id, model.Version).
		Updates(map[string]any{
			"name":     model.Name,
			"code":     model.Code,
			"size":     model.Size,
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
			Model(&tools.Bucket{}).
			Select("1").
			Where("id = ?", model.Id).
			Limit(1).
			Find(&exists)

		if !exists {
			return ErrBucketNotFound
		}
		return ErrBucketVersionInconsistency
	}
	return result.Error
}

// FindById 根据id获取详情
func (dao *GORMBucketDAO) FindById(ctx context.Context, id string) (*tools.Bucket, error) {
	var model tools.Bucket
	err := dao.db.WithContext(ctx).Where("id = ?", id).First(&model).Error
	return &model, err
}

// FindListPage 分页查询
func (dao *GORMBucketDAO) FindListPage(ctx context.Context, filter domainTools.BucketFilter) ([]*tools.Bucket, int64, error) {
	var total int64
	var models []*tools.Bucket

	query := dao.buildQuery(ctx, filter)

	err := query.Count(&total).
		Offset((filter.Page - 1) * filter.PageSize).
		Limit(filter.PageSize).
		Find(&models).Error

	return models, total, err
}

// FindListAll 获取所有列表
func (dao *GORMBucketDAO) FindListAll(ctx context.Context, filter domainTools.BucketFilter) ([]*tools.Bucket, error) {
	var models []*tools.Bucket
	err := dao.buildQuery(ctx, filter).Find(&models).Error
	return models, err
}

// buildQuery 构建查询条件
func (dao *GORMBucketDAO) buildQuery(ctx context.Context, filter domainTools.BucketFilter) *gorm.DB {
	builder := &domainTools.BucketFilter{
		Filters: filters.Filters{
			Creator:    filter.Creator,
			Modifier:   filter.Modifier,
			BelongDept: filter.BelongDept,
		},
		Status: filter.Status,
		Name:   filter.Name,
		Code:   filter.Code,
	}
	return builder.Apply(dao.db.WithContext(ctx).Model(&tools.Bucket{}))
}

// CheckExistByName 检查name是否存在
func (dao *GORMBucketDAO) CheckExistByName(ctx context.Context, name, excludeId string) (bool, error) {
	var model tools.Bucket
	query := dao.db.WithContext(ctx).Model(&tools.Bucket{}).
		Select("id"). // 只查询必要的字段
		Where("name = ?", name)

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

// CheckExistByCode 检查code是否存在
func (dao *GORMBucketDAO) CheckExistByCode(ctx context.Context, code, excludeId string) (bool, error) {
	var model tools.Bucket
	query := dao.db.WithContext(ctx).Model(&tools.Bucket{}).
		Select("id"). // 只查询必要的字段
		Where("code = ?", code)

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
