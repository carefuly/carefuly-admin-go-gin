/**
 * Description：
 * FileName：dept.go
 * Author：CJiaの用心
 * Create：2025/5/15 16:28:21
 * Remark：
 */

package system

import (
	"context"
	"errors"
	domainSystem "github.com/carefuly/carefuly-admin-go-gin/internal/domain/careful/system"
	modelSystem "github.com/carefuly/carefuly-admin-go-gin/internal/model/careful/system"
	cacheSystem "github.com/carefuly/carefuly-admin-go-gin/internal/repository/cache/careful/system"
	daoSystem "github.com/carefuly/carefuly-admin-go-gin/internal/repository/dao/careful/system"
	"github.com/carefuly/carefuly-admin-go-gin/pkg/models"
	"go.uber.org/zap"
)

var (
	ErrDeptNotFound             = daoSystem.ErrDeptNotFound
	ErrDeptDuplicate            = daoSystem.ErrDeptDuplicate
	ErrDeptVersionInconsistency = daoSystem.ErrDeptVersionInconsistency
	ErrDeptChildNodes           = daoSystem.ErrDeptChildNodes
)

type DeptRepository interface {
	Create(ctx context.Context, domain domainSystem.Dept) error
	Delete(ctx context.Context, id string) (int64, error)
	BatchDelete(ctx context.Context, ids []string) error
	Update(ctx context.Context, domain domainSystem.Dept) error

	GetById(ctx context.Context, id string) (domainSystem.Dept, error)
	GetListAll(ctx context.Context, filters domainSystem.DeptFilter) ([]domainSystem.Dept, error)
	CheckExistByNameAndCodeAndParentId(ctx context.Context, name, code, parentId, excludeId string) (bool, error)
}

type deptRepository struct {
	dao   daoSystem.DeptDAO
	cache cacheSystem.DeptCache
}

func NewDeptRepository(dao daoSystem.DeptDAO, cache cacheSystem.DeptCache) DeptRepository {
	return &deptRepository{
		dao:   dao,
		cache: cache,
	}
}

// Create 创建
func (repo *deptRepository) Create(ctx context.Context, domain domainSystem.Dept) error {
	return repo.dao.Insert(ctx, repo.toEntity(domain))
}

// Delete 删除
func (repo *deptRepository) Delete(ctx context.Context, id string) (int64, error) {
	exists, err := repo.dao.CheckExistByIdAndParentId(ctx, id)
	if exists {
		return 0, daoSystem.ErrDeptChildNodes
	}
	if err != nil {
		return 0, err
	}

	rowsAffected, err := repo.dao.Delete(ctx, id)
	if err != nil {
		return rowsAffected, err
	}

	// 删除缓存
	err = repo.cache.Del(ctx, id)
	if err != nil {
		// 网络崩了，也可能是 redis 崩了
		zap.L().Error("Redis异常", zap.Error(err))
		return rowsAffected, err
	}

	return rowsAffected, err
}

// BatchDelete 批量删除
func (repo *deptRepository) BatchDelete(ctx context.Context, ids []string) error {
	var DeptIds []string
	// 检查是否存在子部门
	for _, id := range ids {
		exists, err := repo.dao.CheckExistByIdAndParentId(ctx, id)
		if err != nil {
			zap.L().Error("检查部门是否存在子部门异常", zap.Error(err))
			continue
		}
		if exists {
			continue
		} else {
			DeptIds = append(DeptIds, id)
		}
	}

	err := repo.dao.BatchDelete(ctx, DeptIds)
	if err != nil {
		return err
	}

	// 删除缓存
	for _, val := range ids {
		err = repo.cache.Del(ctx, val)
		if err != nil {
			// 网络崩了，也可能是 redis 崩了
			zap.L().Error("Redis异常", zap.Error(err))
			return err
		}
	}

	return err
}

// Update 更新
func (repo *deptRepository) Update(ctx context.Context, domain domainSystem.Dept) error {
	err := repo.dao.Update(ctx, repo.toEntity(domain))
	if err != nil {
		return err
	}

	// 删除缓存
	err = repo.cache.Del(ctx, domain.Id)
	if err != nil {
		// 网络崩了，也可能是 redis 崩了
		zap.L().Error("Redis异常", zap.Error(err))
		return err
	}

	return nil
}

// GetById 根据ID获取
func (repo *deptRepository) GetById(ctx context.Context, id string) (domainSystem.Dept, error) {
	domain, err := repo.cache.Get(ctx, id)
	if err == nil && domain != nil {
		return *domain, nil // 命中缓存
	}
	if err != nil && !errors.Is(err, cacheSystem.ErrDeptNotExist) {
		// 缓存查询出错但不是"不存在"错误，记录日志但继续查DB
		zap.L().Error("缓存获取错误:", zap.Error(err))
	}

	entity, err := repo.dao.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, daoSystem.ErrDeptNotFound) {
			// 数据库不存在，设置防穿透标记
			_ = repo.cache.SetNotFound(ctx, id)
			return domainSystem.Dept{}, nil
		}
		return domainSystem.Dept{}, err
	}

	toDomain := repo.toDomain(entity)
	if err := repo.cache.Set(ctx, toDomain); err != nil {
		// 网络崩了，也可能是 redis 崩了
		zap.L().Error("Redis异常", zap.Error(err))
	}

	return toDomain, nil
}

// GetListAll 查询所有列表
func (repo *deptRepository) GetListAll(ctx context.Context, filters domainSystem.DeptFilter) ([]domainSystem.Dept, error) {
	list, err := repo.dao.FindListAll(ctx, filters)
	if err != nil {
		return []domainSystem.Dept{}, err
	}

	if len(list) == 0 {
		return []domainSystem.Dept{}, nil
	}

	var toDomain []domainSystem.Dept
	for _, v := range list {
		domainDept := repo.toDomain(v)

		// 只有存在有效ParentID时才查询父部门
		if v.ParentID != "" { // 关键修复：空字符串不查询
			dept, err := repo.dao.FindById(ctx, v.ParentID)
			if err == nil { // 仅当查询成功时设置ParentName
				domainDept.ParentName = dept.Name
			} else {
				// 可选：记录错误或设置默认值
				domainDept.ParentName = "未知部门"
			}
		} else {
			// 明确标记顶级部门
			domainDept.ParentName = "" // 或"顶级部门"等业务标识
		}

		toDomain = append(toDomain, domainDept)
	}

	return toDomain, nil
}

// CheckExistByNameAndCodeAndParentId 检查name、code和parentId是否同时存在
func (repo *deptRepository) CheckExistByNameAndCodeAndParentId(ctx context.Context, name, code, parentId, excludeId string) (bool, error) {
	return repo.dao.CheckExistByNameAndCodeAndParentId(ctx, name, code, parentId, excludeId)
}

// toEntity 转换为实体模型
func (repo *deptRepository) toEntity(domain domainSystem.Dept) modelSystem.Dept {
	return modelSystem.Dept{
		CoreModels: models.CoreModels{
			Id:         domain.Id,
			Sort:       domain.Sort,
			Version:    domain.Version,
			Creator:    domain.Creator,
			Modifier:   domain.Modifier,
			BelongDept: domain.BelongDept,
			Remark:     domain.Remark,
		},
		Status:   domain.Status,
		Name:     domain.Name,
		Code:     domain.Code,
		Owner:    domain.Owner,
		Phone:    domain.Phone,
		Email:    domain.Email,
		ParentID: domain.ParentID,
	}
}

// toDomain 转换为领域模型
func (repo *deptRepository) toDomain(entity *modelSystem.Dept) domainSystem.Dept {
	domain := domainSystem.Dept{
		Dept: *entity,
	}

	if entity.CreateTime != nil {
		domain.CreateTime = entity.CreateTime.Format("2006-01-02 15:04:05")
	}
	if entity.UpdateTime != nil {
		domain.UpdateTime = entity.UpdateTime.Format("2006-01-02 15:04:05")
	}

	return domain
}
