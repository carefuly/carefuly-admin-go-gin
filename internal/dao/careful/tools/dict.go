/**
 * Description：
 * FileName：dict.go
 * Author：CJiaの用心
 * Create：2025/4/14 20:42:35
 * Remark：
 */

package tools

import (
	"context"
	"errors"
	domainTools "github.com/carefuly/carefuly-admin-go-gin/internal/domain/careful/tools"
	"github.com/carefuly/carefuly-admin-go-gin/internal/model/careful/tools"
	"github.com/carefuly/carefuly-admin-go-gin/pkg/ginx/query/filters"
	"gorm.io/gorm"
)

var (
	ErrDictRecordNotFound       = gorm.ErrRecordNotFound
	ErrDictNotFound             = errors.New("记录不存在")
	ErrDuplicateDict            = errors.New("字典已存在")
	ErrDuplicateDictName        = errors.New("字典名称已存在")
	ErrDuplicateDictCode        = errors.New("字典编码已存在")
	ErrDictVersionInconsistency = errors.New("数据已被修改，请刷新后重试")
)

type DictDAO interface {
	Insert(ctx context.Context, domain tools.Dict) error
	Delete(ctx context.Context, id string) (int64, error)
	BatchDelete(ctx context.Context, ids []string) error
	Update(ctx context.Context, id string, model tools.Dict) (int64, error)
	FindById(ctx context.Context, id string) (*tools.Dict, error)
	FindByName(ctx context.Context, name string) (*tools.Dict, error)
	FindListPage(ctx context.Context, filter domainTools.DictFilter) ([]*tools.Dict, int64, error)
	FindListAll(ctx context.Context, filter domainTools.DictFilter) ([]*tools.Dict, error)
	CheckExistByName(ctx context.Context, name, excludeId string) (bool, error)
	CheckExistByCode(ctx context.Context, code, excludeId string) (bool, error)
}

type GORMDictDAO struct {
	db *gorm.DB
}

func NewDictDao(db *gorm.DB) DictDAO {
	return &GORMDictDAO{db: db}
}

// Insert 创建
func (dao *GORMDictDAO) Insert(ctx context.Context, model tools.Dict) error {
	return dao.db.WithContext(ctx).Create(&model).Error
}

// Delete 根据ID删除
func (dao *GORMDictDAO) Delete(ctx context.Context, id string) (int64, error) {
	result := dao.db.WithContext(ctx).Where("id = ?", id).Delete(&tools.Dict{})
	return result.RowsAffected, result.Error
}

// BatchDelete 批量删除
func (dao *GORMDictDAO) BatchDelete(ctx context.Context, ids []string) error {
	return dao.db.WithContext(ctx).Where("id IN ?", ids).Delete(&tools.Dict{}).Error
}

// Update 更新
func (dao *GORMDictDAO) Update(ctx context.Context, id string, model tools.Dict) (int64, error) {
	result := dao.db.WithContext(ctx).Model(&model).Where("id = ? AND version = ?", id, model.Version).
		Updates(map[string]any{
			"name":      model.Name,
			"code":      model.Code,
			"type":      model.Type,
			"typeValue": model.TypeValue,
			"version":   gorm.Expr("version + 1"),
			"modifier":  model.Modifier,
			"remark":    model.Remark,
		})
	// 处理行影响数为0的情况
	if result.RowsAffected == 0 {
		// 先检查记录是否存在
		var exists bool
		dao.db.WithContext(ctx).
			Model(&tools.Dict{}).
			Select("1").
			Where("id = ?", id).
			Limit(1).
			Find(&exists)

		if !exists {
			return result.RowsAffected, ErrDictNotFound
		}
		return result.RowsAffected, ErrDictVersionInconsistency
	}
	return result.RowsAffected, result.Error
}

// FindById 根据ID查询
func (dao *GORMDictDAO) FindById(ctx context.Context, id string) (*tools.Dict, error) {
	var model tools.Dict
	err := dao.db.WithContext(ctx).Where("id = ?", id).First(&model).Error
	return &model, err
}

// FindByName 根据Name查询
func (dao *GORMDictDAO) FindByName(ctx context.Context, name string) (*tools.Dict, error) {
	var model tools.Dict
	err := dao.db.WithContext(ctx).Where("name = ?", name).First(&model).Error
	return &model, err
}

// FindListPage 分页查询
func (dao *GORMDictDAO) FindListPage(ctx context.Context, filter domainTools.DictFilter) ([]*tools.Dict, int64, error) {
	var total int64
	var models []*tools.Dict

	query := dao.buildQuery(ctx, filter)

	err := query.Count(&total).
		Offset((filter.Page - 1) * filter.PageSize).
		Limit(filter.PageSize).
		Find(&models).Error

	return models, total, err
}

// FindListAll 查询所有列表
func (dao *GORMDictDAO) FindListAll(ctx context.Context, filter domainTools.DictFilter) ([]*tools.Dict, error) {
	query := dao.buildQuery(ctx, filter)

	var models []*tools.Dict
	err := query.Find(&models).Error

	return models, err
}

func (dao *GORMDictDAO) buildQuery(ctx context.Context, filter domainTools.DictFilter) *gorm.DB {
	builder := &domainTools.DictFilter{
		Filters: filters.Filters{
			Creator:  filter.Creator,
			Modifier: filter.Modifier,
			Status:   filter.Status,
		},
		Name:      filter.Name,
		Code:      filter.Code,
		Type:      filter.Type,
		TypeValue: filter.TypeValue,
	}

	return builder.Apply(dao.db.WithContext(ctx).Model(&tools.Dict{}))
}

// CheckExistByName 检查name是否存在
func (dao *GORMDictDAO) CheckExistByName(ctx context.Context, name, excludeId string) (bool, error) {
	var count int64
	query := dao.db.WithContext(ctx).Model(&tools.Dict{}).
		Where("name = ?", name)
	if excludeId != "" {
		query = query.Where("id != ?", excludeId)
	}
	err := query.Count(&count).Error
	return count > 0, err
}

// CheckExistByCode 检查code是否存在
func (dao *GORMDictDAO) CheckExistByCode(ctx context.Context, code, excludeId string) (bool, error) {
	var count int64
	query := dao.db.WithContext(ctx).Model(&tools.Dict{}).
		Where("code = ?", code)
	if excludeId != "" {
		query = query.Where("id != ?", excludeId)
	}
	err := query.Count(&count).Error
	return count > 0, err
}
