/**
 * Description：
 * FileName：menu_button.go
 * Author：CJiaの用心
 * Create：2025/6/9 14:26:50
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
	ErrMenuButtonNotFound             = repositorySystem.ErrMenuButtonNotFound
	ErrMenuButtonVersionInconsistency = repositorySystem.ErrMenuButtonVersionInconsistency
)

type MenuButtonService interface {
	Create(ctx context.Context, domain domainSystem.MenuButton) error
	Delete(ctx context.Context, id string) error
	BatchDelete(ctx context.Context, ids []string) error
	Update(ctx context.Context, domain domainSystem.MenuButton) error
	GetById(ctx context.Context, id string) (domainSystem.MenuButton, error)
	GetListPage(ctx context.Context, filters domainSystem.MenuButtonFilter) ([]domainSystem.MenuButton, int64, error)
	GetListAll(ctx context.Context, filters domainSystem.MenuButtonFilter) ([]domainSystem.MenuButton, error)
}

type menuButtonService struct {
	repo repositorySystem.MenuButtonRepository
}

func NewMenuButtonService(repo repositorySystem.MenuButtonRepository) MenuButtonService {
	return &menuButtonService{
		repo: repo,
	}
}

// Create 创建
func (svc *menuButtonService) Create(ctx context.Context, domain domainSystem.MenuButton) error {
	return svc.repo.Create(ctx, domain)
}

// Delete 删除
func (svc *menuButtonService) Delete(ctx context.Context, id string) error {
	rowsAffected, err := svc.repo.Delete(ctx, id)
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return repositorySystem.ErrMenuButtonNotFound
	}
	return err
}

// BatchDelete 批量删除
func (svc *menuButtonService) BatchDelete(ctx context.Context, ids []string) error {
	return svc.repo.BatchDelete(ctx, ids)
}

// Update 更新
func (svc *menuButtonService) Update(ctx context.Context, domain domainSystem.MenuButton) error {
	err := svc.repo.Update(ctx, domain)
	if err != nil {
		if errors.Is(err, repositorySystem.ErrMenuButtonVersionInconsistency) {
			return repositorySystem.ErrMenuButtonVersionInconsistency
		}
		return err
	}
	return err
}

// GetById 获取详情
func (svc *menuButtonService) GetById(ctx context.Context, id string) (domainSystem.MenuButton, error) {
	domain, err := svc.repo.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, repositorySystem.ErrMenuButtonNotFound) {
			return domain, repositorySystem.ErrMenuButtonNotFound
		}
		return domain, err
	}
	if domain.Id == "" {
		return domain, repositorySystem.ErrMenuButtonNotFound
	}
	return domain, err
}

// GetListPage 分页查询列表
func (svc *menuButtonService) GetListPage(ctx context.Context, filters domainSystem.MenuButtonFilter) ([]domainSystem.MenuButton, int64, error) {
	return svc.repo.GetListPage(ctx, filters)
}

// GetListAll 查询所有列表
func (svc *menuButtonService) GetListAll(ctx context.Context, filters domainSystem.MenuButtonFilter) ([]domainSystem.MenuButton, error) {
	return svc.repo.GetListAll(ctx, filters)
}
