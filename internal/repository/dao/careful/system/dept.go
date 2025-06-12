/**
 * Description：
 * FileName：dept.go
 * Author：CJiaの用心
 * Create：2025/5/15 16:11:08
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
	ErrDeptNotFound             = gorm.ErrRecordNotFound
	ErrDeptDuplicate            = errors.New("部门信息已存在")
	ErrDeptVersionInconsistency = errors.New("数据已被修改，请刷新后重试")
	ErrDeptChildNodes           = errors.New("请先删除当前部门下的子部门")
)

type DeptDAO interface {
	Insert(ctx context.Context, model system.Dept) error
	Delete(ctx context.Context, id string) (int64, error)
	BatchDelete(ctx context.Context, ids []string) error
	Update(ctx context.Context, model system.Dept) error

	FindById(ctx context.Context, id string) (*system.Dept, error)
	FindListAll(ctx context.Context, filter domainSystem.DeptFilter) ([]*system.Dept, error)

	CheckExistByIdAndParentId(ctx context.Context, id string) (bool, error)
	CheckExistByNameAndCodeAndParentId(ctx context.Context, name, code, parentId, excludeId string) (bool, error)
}

type GORMDeptDAO struct {
	db *gorm.DB
}

func NewGORMDeptDAO(db *gorm.DB) DeptDAO {
	return &GORMDeptDAO{
		db: db,
	}
}

// Insert 新增
func (dao *GORMDeptDAO) Insert(ctx context.Context, model system.Dept) error {
	return dao.db.WithContext(ctx).Create(&model).Error
}

// Delete 删除
func (dao *GORMDeptDAO) Delete(ctx context.Context, id string) (int64, error) {
	result := dao.db.WithContext(ctx).Where("id = ?", id).Delete(&system.Dept{})
	return result.RowsAffected, result.Error
}

// BatchDelete 批量删除
func (dao *GORMDeptDAO) BatchDelete(ctx context.Context, ids []string) error {
	return dao.db.WithContext(ctx).Where("id IN ?", ids).Delete(&system.Dept{}).Error
}

// Update 更新
func (dao *GORMDeptDAO) Update(ctx context.Context, model system.Dept) error {
	result := dao.db.WithContext(ctx).Model(&model).
		Where("id = ? AND version = ?", model.Id, model.Version).
		Updates(map[string]any{
			"name":      model.Name,
			"code":      model.Code,
			"owner":     model.Owner,
			"phone":     model.Phone,
			"email":     model.Email,
			"parent_id": model.ParentID,
			"sort":      model.Sort,
			"status":    model.Status,
			"version":   gorm.Expr("version + 1"),
			"modifier":  model.Modifier,
			"remark":    model.Remark,
		})
	// 处理行影响数为0的情况
	if result.RowsAffected == 0 {
		// 先检查记录是否存在
		var exists bool
		dao.db.WithContext(ctx).
			Model(&system.Dept{}).
			Select("1").
			Where("id = ?", model.Id).
			Limit(1).
			Find(&exists)

		if !exists {
			return ErrDeptNotFound
		}
		return ErrDeptVersionInconsistency
	}

	return result.Error
}

// FindById 根据id获取详情
func (dao *GORMDeptDAO) FindById(ctx context.Context, id string) (*system.Dept, error) {
	var model system.Dept
	err := dao.db.WithContext(ctx).Where("id = ?", id).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &model, ErrDeptNotFound
		}
		return &model, err
	}
	return &model, err
}

// FindListAll 获取所有列表
func (dao *GORMDeptDAO) FindListAll(ctx context.Context, filter domainSystem.DeptFilter) ([]*system.Dept, error) {
	var models []*system.Dept

	query := dao.buildQuery(ctx, filter)

	// 查询
	if err := query.Find(&models).Error; err != nil {
		return nil, err
	}

	return models, nil
}

// buildQuery 构建查询条件
func (dao *GORMDeptDAO) buildQuery(ctx context.Context, filter domainSystem.DeptFilter) *gorm.DB {
	builder := &domainSystem.DeptFilter{
		Filters: filters.Filters{
			Creator:    filter.Creator,
			Modifier:   filter.Modifier,
			BelongDept: filter.BelongDept,
		},
		Status: filter.Status,
		Name:   filter.Name,
		Code:   filter.Code,
	}
	return builder.Apply(dao.db.WithContext(ctx).Model(&system.Dept{}))
}

// CheckExistByIdAndParentId 检查当前id是否存在子节点
func (dao *GORMDeptDAO) CheckExistByIdAndParentId(ctx context.Context, id string) (bool, error) {
	var model system.Dept
	query := dao.db.WithContext(ctx).Model(&system.Dept{}).
		Select("id"). // 只查询必要的字段
		Where("parent_id = ?", id)

	// 使用 LIMIT 1 快速判断是否存在
	err := query.Limit(1).First(&model).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil // 不存在
	}
	return err == nil, err // 存在或查询出错
}

// CheckExistByNameAndCodeAndParentId 检查name、code和parentId是否同时存在
func (dao *GORMDeptDAO) CheckExistByNameAndCodeAndParentId(ctx context.Context, name, code, parentId, excludeId string) (bool, error) {
	var model system.Dept
	query := dao.db.WithContext(ctx).Model(&system.Dept{}).
		Select("id"). // 只查询必要的字段
		Where("name = ? AND code = ? AND parent_id = ?", name, code, parentId)

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
