/**
 * Description：
 * FileName：menu_column.go
 * Author：CJiaの用心
 * Create：2025/6/9 14:26:57
 * Remark：
 */

package system

import (
	"context"
	"errors"
	domainSystem "github.com/carefuly/carefuly-admin-go-gin/internal/domain/careful/system"
	repositorySystem "github.com/carefuly/carefuly-admin-go-gin/internal/repository/repository/careful/system"
)

var (
	ErrMenuColumnNotFound             = repositorySystem.ErrMenuColumnNotFound
	ErrMenuColumnVersionInconsistency = repositorySystem.ErrMenuColumnVersionInconsistency
)

type MenuColumnService interface {
	Create(ctx context.Context, domain domainSystem.MenuColumn) error
	Delete(ctx context.Context, id string) error
	BatchDelete(ctx context.Context, ids []string) error
	Update(ctx context.Context, domain domainSystem.MenuColumn) error

	GetById(ctx context.Context, id string) (domainSystem.MenuColumn, error)
	GetListPage(ctx context.Context, filters domainSystem.MenuColumnFilter) ([]domainSystem.MenuColumn, int64, error)
	GetListAll(ctx context.Context, filters domainSystem.MenuColumnFilter) ([]domainSystem.MenuColumn, error)
}

type menuColumnService struct {
	repo repositorySystem.MenuColumnRepository
}

func NewMenuColumnService(repo repositorySystem.MenuColumnRepository) MenuColumnService {
	return &menuColumnService{
		repo: repo,
	}
}

// Create 创建
func (svc *menuColumnService) Create(ctx context.Context, domain domainSystem.MenuColumn) error {
	return svc.repo.Create(ctx, domain)
}

// Delete 删除
func (svc *menuColumnService) Delete(ctx context.Context, id string) error {
	rowsAffected, err := svc.repo.Delete(ctx, id)
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return repositorySystem.ErrMenuColumnNotFound
	}
	return err
}

// BatchDelete 批量删除
func (svc *menuColumnService) BatchDelete(ctx context.Context, ids []string) error {
	return svc.repo.BatchDelete(ctx, ids)
}

// Update 更新
func (svc *menuColumnService) Update(ctx context.Context, domain domainSystem.MenuColumn) error {
	err := svc.repo.Update(ctx, domain)
	if err != nil {
		if errors.Is(err, repositorySystem.ErrMenuColumnVersionInconsistency) {
			return repositorySystem.ErrMenuColumnVersionInconsistency
		}
		return err
	}
	return err
}

// GetById 获取详情
func (svc *menuColumnService) GetById(ctx context.Context, id string) (domainSystem.MenuColumn, error) {
	domain, err := svc.repo.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, repositorySystem.ErrMenuColumnNotFound) {
			return domain, repositorySystem.ErrMenuColumnNotFound
		}
		return domain, err
	}
	if domain.Id == "" {
		return domain, repositorySystem.ErrMenuColumnNotFound
	}
	return domain, err
}

// GetListPage 分页查询列表
func (svc *menuColumnService) GetListPage(ctx context.Context, filters domainSystem.MenuColumnFilter) ([]domainSystem.MenuColumn, int64, error) {
	return svc.repo.GetListPage(ctx, filters)
}

// GetListAll 查询所有列表
func (svc *menuColumnService) GetListAll(ctx context.Context, filters domainSystem.MenuColumnFilter) ([]domainSystem.MenuColumn, error) {
	return svc.repo.GetListAll(ctx, filters)
}
