/**
 * Description：
 * FileName：dict_type.go
 * Author：CJiaの用心
 * Create：2025/4/17 11:33:20
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

type DictTypeDAO interface {
	Insert(ctx context.Context, model tools.DictType) error
	// Delete(ctx context.Context, id string) (int64, error)
	// BatchDelete(ctx context.Context, ids []string) error
	Update(ctx context.Context, id string, model tools.DictType) (int64, error)
	FindById(ctx context.Context, id string) (*tools.DictType, error)
	// FindListPage(ctx context.Context, filter domainTools.DictTypeFilter) ([]*tools.DictType, int64, error)
	// FindListAll(ctx context.Context, filter domainTools.DictTypeFilter) ([]*tools.DictType, error)
	// CheckExistByDictIdAndNameAndValue(ctx context.Context, name, strValue string, intValue int64, dictId string) (bool, error)
}

var (
	ErrDictTypeRecordNotFound       = gorm.ErrRecordNotFound
	ErrDictTypeNotFound             = errors.New("记录不存在")
	ErrDuplicateDictType            = tools.ErrDictTypeUniqueIndex
	ErrDictTypeVersionInconsistency = errors.New("数据已被修改，请刷新后重试")
)

type GORMDictTypeDAO struct {
	db *gorm.DB
}

func NewDictTypeDao(db *gorm.DB) DictTypeDAO {
	return &GORMDictTypeDAO{db: db}
}

// Insert 创建
func (dao *GORMDictTypeDAO) Insert(ctx context.Context, model tools.DictType) error {
	return dao.db.WithContext(ctx).Create(&model).Error
}

// Delete 根据ID删除
func (dao *GORMDictTypeDAO) Delete(ctx context.Context, id string) (int64, error) {
	result := dao.db.WithContext(ctx).Where("id = ?", id).Delete(&tools.DictType{})
	return result.RowsAffected, result.Error
}

// BatchDelete 批量删除
func (dao *GORMDictTypeDAO) BatchDelete(ctx context.Context, ids []string) error {
	return dao.db.WithContext(ctx).Where("id IN ?", ids).Delete(&tools.DictType{}).Error
}

// Update 更新
func (dao *GORMDictTypeDAO) Update(ctx context.Context, id string, model tools.DictType) (int64, error) {
	result := dao.db.WithContext(ctx).Model(&model).Where("id = ? AND version = ?", id, model.Version).
		Updates(map[string]any{
			"name":      model.Name,
			"strValue":  model.StrValue,
			"intValue":  model.IntValue,
			"boolValue": model.BoolValue,
			"dictTag":   model.DictTag,
			"dictColor": model.DictColor,
			"dictName":  model.DictName,
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
			Model(&tools.DictType{}).
			Select("1").
			Where("id = ?", id).
			Limit(1).
			Find(&exists)

		if !exists {
			return result.RowsAffected, ErrDictTypeNotFound
		}
		return result.RowsAffected, ErrDictTypeVersionInconsistency
	}
	return result.RowsAffected, result.Error
}

// FindById 根据ID查询
func (dao *GORMDictTypeDAO) FindById(ctx context.Context, id string) (*tools.DictType, error) {
	var model tools.DictType
	err := dao.db.WithContext(ctx).
		Preload("Dict").
		Where("id = ?", id).First(&model).Error
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
	query := dao.buildQuery(ctx, filter)

	var models []*tools.DictType
	err := query.Find(&models).Error

	return models, err
}

func (dao *GORMDictTypeDAO) buildQuery(ctx context.Context, filter domainTools.DictTypeFilter) *gorm.DB {
	builder := &domainTools.DictTypeFilter{
		Filters: filters.Filters{
			Creator:  filter.Creator,
			Modifier: filter.Modifier,
			Status:   filter.Status,
		},
		Name:    filter.Name,
		DictTag: filter.DictTag,
		DictId:  filter.DictId,
	}

	return builder.Apply(dao.db.WithContext(ctx).Model(&tools.Dict{}))
}

// CheckExistByDictIdAndNameAndValue 检查同一个字典下是否存在相同名称和值
func (dao *GORMDictTypeDAO) CheckExistByDictIdAndNameAndValue(ctx context.Context, name, strValue string, intValue int64, dictId string) (bool, error) {
	var count int64
	query := dao.db.WithContext(ctx).Model(&tools.DictType{}).
		Where("name = ? AND dict_id = ?", name, dictId)

	if strValue != "" {
		query = query.Where("str_value = ?", strValue)
	} else {
		query = query.Where("int_value = ?", intValue)
	}

	err := query.Count(&count).Error
	return count > 0, err
}
