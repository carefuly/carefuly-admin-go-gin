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
	domainSystem "github.com/carefuly/carefuly-admin-go-gin/internal/domain/careful/system"
	modelSystem "github.com/carefuly/carefuly-admin-go-gin/internal/model/careful/system"
	daoSystem "github.com/carefuly/carefuly-admin-go-gin/internal/repository/dao/careful/system"
	"github.com/carefuly/carefuly-admin-go-gin/pkg/models"
)

var (
	ErrDeptNotFound             = daoSystem.ErrDeptNotFound
	ErrDeptNameDuplicate        = daoSystem.ErrDeptNameDuplicate
	ErrDeptCodeDuplicate        = daoSystem.ErrDeptCodeDuplicate
	ErrDeptDuplicate            = daoSystem.ErrDeptDuplicate
	ErrDeptVersionInconsistency = daoSystem.ErrDeptVersionInconsistency
)

type DeptRepository interface {
	Create(ctx context.Context, domain domainSystem.Dept) error

	GetListAll(ctx context.Context, filters domainSystem.DeptFilter) ([]domainSystem.Dept, error)

	CheckExistByName(ctx context.Context, username, excludeId string) (bool, error)
	CheckExistByCode(ctx context.Context, username, excludeId string) (bool, error)
}

type deptRepository struct {
	dao daoSystem.DeptDAO
	// cache cacheSystem.
}

func NewDeptRepository(dao daoSystem.DeptDAO) DeptRepository {
	return &deptRepository{
		dao: dao,
		// cache: cache,
	}
}

// Create 创建
func (repo *deptRepository) Create(ctx context.Context, domain domainSystem.Dept) error {
	return repo.dao.Insert(ctx, repo.toEntity(domain))
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
		toDomain = append(toDomain, repo.toDomain(v))
	}

	return toDomain, nil
}

// CheckExistByName 检查name是否存在
func (repo *deptRepository) CheckExistByName(ctx context.Context, name, excludeId string) (bool, error) {
	return repo.dao.CheckExistByName(ctx, name, excludeId)
}

// CheckExistByCode 检查code是否存在
func (repo *deptRepository) CheckExistByCode(ctx context.Context, code, excludeId string) (bool, error) {
	return repo.dao.CheckExistByCode(ctx, code, excludeId)
}

// toEntity 转换为实体模型
func (repo *deptRepository) toEntity(domain domainSystem.Dept) modelSystem.Dept {
	return modelSystem.Dept{
		CoreModels: models.CoreModels{
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
