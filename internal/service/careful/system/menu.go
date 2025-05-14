/**
 * Description：
 * FileName：menu.go
 * Author：CJiaの用心
 * Create：2025/5/13 16:39:24
 * Remark：
 */

package system

import (
	"context"
	domainSystem "github.com/carefuly/carefuly-admin-go-gin/internal/domain/careful/system"
	repositorySystem "github.com/carefuly/carefuly-admin-go-gin/internal/repository/repository/careful/system"
)

type MenuService interface {
	GetListAll(ctx context.Context, filter domainSystem.MenuFilter) ([]domainSystem.Menu, error)
}

type menuService struct {
	repo repositorySystem.MenuRepository
}

func NewMenuService(repo repositorySystem.MenuRepository) MenuService {
	return &menuService{
		repo: repo,
	}
}

// GetListAll 查询所有列表
func (svc *menuService) GetListAll(ctx context.Context, filter domainSystem.MenuFilter) ([]domainSystem.Menu, error) {
	return svc.repo.GetListAll(ctx, filter)
}
