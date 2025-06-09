/**
 * Description：
 * FileName：dictType.go
 * Author：CJiaの用心
 * Create：2025/5/23 16:38:20
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
	ErrDictTypeNotFound             = gorm.ErrRecordNotFound
	ErrDictTypeDuplicate            = tools.ErrDictTypeUniqueIndex
	ErrDictTypeVersionInconsistency = errors.New("数据已被修改，请刷新后重试")
)

type DictTypeDAO interface {
	Insert(ctx context.Context, model tools.DictType) error
	Delete(ctx context.Context, id string) (int64, error)
	BatchDelete(ctx context.Context, ids []string) error
	Update(ctx context.Context, model tools.DictType) error

	FindById(ctx context.Context, id string) (*tools.DictType, error)
	FindListPage(ctx context.Context, filter domainTools.DictTypeFilter) ([]*tools.DictType, int64, error)
	FindListAll(ctx context.Context, filter domainTools.DictTypeFilter) ([]*tools.DictType, error)
}

type GORMDictTypeDAO struct {
	db *gorm.DB
}

func NewGORMDictTypeDAO(db *gorm.DB) DictTypeDAO {
	return &GORMDictTypeDAO{db: db}
}

// Insert 新增
func (dao *GORMDictTypeDAO) Insert(ctx context.Context, model tools.DictType) error {
	return dao.db.WithContext(ctx).Create(&model).Error
}

// Delete 删除
func (dao *GORMDictTypeDAO) Delete(ctx context.Context, id string) (int64, error) {
	result := dao.db.WithContext(ctx).Where("id = ?", id).Delete(&tools.DictType{})
	return result.RowsAffected, result.Error
}

// BatchDelete 批量删除
func (dao *GORMDictTypeDAO) BatchDelete(ctx context.Context, ids []string) error {
	return dao.db.WithContext(ctx).Where("id IN ?", ids).Delete(&tools.DictType{}).Error
}

// Update 更新
func (dao *GORMDictTypeDAO) Update(ctx context.Context, model tools.DictType) error {
	result := dao.db.WithContext(ctx).Model(&model).
		Where("id = ? AND version = ?", model.Id, model.Version).
		Updates(map[string]any{
			// "strValue":  model.StrValue,
			// "intValue":  model.IntValue,
			// "boolValue": model.BoolValue,
			"dictTag":   model.DictTag,
			"dictColor": model.DictColor,
			"sort":      model.Sort,
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
			Where("id = ?", model.Id).
			Limit(1).
			Find(&exists)

		if !exists {
			return ErrDictTypeNotFound
		}
		return ErrDictTypeVersionInconsistency
	}
	return result.Error
}

// FindById 根据id获取详情
func (dao *GORMDictTypeDAO) FindById(ctx context.Context, id string) (*tools.DictType, error) {
	var model tools.DictType
	err := dao.db.WithContext(ctx).
		Preload("Dict").
		Where("id = ?", id).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &model, ErrDictTypeNotFound
		}
		return &model, err
	}
	return &model, err
}

// FindListPage 分页查询
func (dao *GORMDictTypeDAO) FindListPage(ctx context.Context, filter domainTools.DictTypeFilter) ([]*tools.DictType, int64, error) {
	var total int64
	var models []*tools.DictType

	query := dao.buildQuery(ctx, filter)

	err := query.Count(&total).
		Offset((filter.Page - 1) * filter.PageSize).
		Limit(filter.PageSize).
		Find(&models).Error

	return models, total, err
}

// FindListAll 查询所有字典
func (dao *GORMDictTypeDAO) FindListAll(ctx context.Context, filter domainTools.DictTypeFilter) ([]*tools.DictType, error) {
	var models []*tools.DictType

	query := dao.buildQuery(ctx, filter)

	// 查询
	if err := query.Find(&models).Error; err != nil {
		return nil, err
	}

	return models, nil
}

// buildQuery 构建查询条件
func (dao *GORMDictTypeDAO) buildQuery(ctx context.Context, filter domainTools.DictTypeFilter) *gorm.DB {
	builder := &domainTools.DictTypeFilter{
		Filters: filters.Filters{
			Creator:    filter.Creator,
			Modifier:   filter.Modifier,
			BelongDept: filter.BelongDept,
		},
		Status:   filter.Status,
		Name:     filter.Name,
		DictTag:  filter.DictTag,
		DictName: filter.DictName,
		DictId:   filter.DictId,
	}
	return builder.Apply(dao.db.WithContext(ctx).Model(&tools.DictType{}))
}
