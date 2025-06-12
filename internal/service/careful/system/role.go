/**
 * Description：
 * FileName：role.go
 * Author：CJiaの用心
 * Create：2025/6/12 13:46:53
 * Remark：
 */

package system

import (
	"context"
	"errors"
	domainSystem "github.com/carefuly/carefuly-admin-go-gin/internal/domain/careful/system"
	repositorySystem "github.com/carefuly/carefuly-admin-go-gin/internal/repository/repository/careful/system"
	"github.com/go-sql-driver/mysql"
)

var (
	ErrRoleNotFound             = repositorySystem.ErrRoleNotFound
	ErrRoleCodeDuplicate        = repositorySystem.ErrRoleCodeDuplicate
	ErrRoleDuplicate            = repositorySystem.ErrRoleDuplicate
	ErrRoleVersionInconsistency = repositorySystem.ErrRoleVersionInconsistency
)

type RoleService interface {
	Create(ctx context.Context, domain domainSystem.Role) error
	Delete(ctx context.Context, id string) error
	BatchDelete(ctx context.Context, ids []string) error
	Update(ctx context.Context, domain domainSystem.Role) error

	GetById(ctx context.Context, id string) (domainSystem.Role, error)
	GetListPage(ctx context.Context, filter domainSystem.RoleFilter) ([]domainSystem.Role, int64, error)
	GetListAll(ctx context.Context, filter domainSystem.RoleFilter) ([]domainSystem.Role, error)
}

type roleService struct {
	repo repositorySystem.RoleRepository
}

func NewRoleService(repo repositorySystem.RoleRepository) RoleService {
	return &roleService{
		repo: repo,
	}
}

// Create 创建
func (svc *roleService) Create(ctx context.Context, domain domainSystem.Role) error {
	// 检查code是否存在
	exists, err := svc.repo.CheckExistByCode(ctx, domain.Code, "")
	if err != nil {
		return err
	}
	if exists {
		return repositorySystem.ErrRoleCodeDuplicate
	}

	// 创建用户
	if err := svc.repo.Create(ctx, domain); err != nil {
		if svc.IsDuplicateEntryError(err) {
			return repositorySystem.ErrRoleCodeDuplicate
		}
		return err
	}

	return nil
}

// Delete 删除
func (svc *roleService) Delete(ctx context.Context, id string) error {
	rowsAffected, err := svc.repo.Delete(ctx, id)
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return repositorySystem.ErrRoleNotFound
	}
	return err
}

// BatchDelete 批量删除
func (svc *roleService) BatchDelete(ctx context.Context, ids []string) error {
	return svc.repo.BatchDelete(ctx, ids)
}

// Update 更新
func (svc *roleService) Update(ctx context.Context, domain domainSystem.Role) error {
	// 检查code是否存在
	exists, err := svc.repo.CheckExistByCode(ctx, domain.Code, domain.Id)
	if err != nil {
		return err
	}
	if exists {
		return repositorySystem.ErrRoleCodeDuplicate
	}

	// 更新用户
	if err := svc.repo.Update(ctx, domain); err != nil {
		switch {
		case svc.IsDuplicateEntryError(err):
			return repositorySystem.ErrRoleCodeDuplicate
		case errors.Is(err, repositorySystem.ErrRoleVersionInconsistency):
			return repositorySystem.ErrRoleVersionInconsistency
		default:
			return err
		}
	}

	return nil
}

// GetById 获取详情
func (svc *roleService) GetById(ctx context.Context, id string) (domainSystem.Role, error) {
	domain, err := svc.repo.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, repositorySystem.ErrRoleNotFound) {
			return domain, repositorySystem.ErrRoleNotFound
		}
		return domain, err
	}
	if domain.Id == "" {
		return domain, repositorySystem.ErrRoleNotFound
	}
	return domain, err
}

// GetListPage 分页查询列表
func (svc *roleService) GetListPage(ctx context.Context, filter domainSystem.RoleFilter) ([]domainSystem.Role, int64, error) {
	return svc.repo.GetListPage(ctx, filter)
}

// GetListAll 查询所有列表
func (svc *roleService) GetListAll(ctx context.Context, filter domainSystem.RoleFilter) ([]domainSystem.Role, error) {
	return svc.repo.GetListAll(ctx, filter)
}

// IsDuplicateEntryError 判断是否是唯一冲突错误
func (svc *roleService) IsDuplicateEntryError(err error) bool {
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		// MySQL 错误码 1062 表示唯一冲突
		return mysqlErr.Number == 1062
	}
	return false
}
