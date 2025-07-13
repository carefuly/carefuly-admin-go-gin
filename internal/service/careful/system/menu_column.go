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
	"github.com/carefuly/carefuly-admin-go-gin/pkg/constants/careful/system/menu"
)

type MenuAndColumn struct {
	Id       string         `json:"id"`        // 菜单列id
	Title    string         `json:"title"`     // 菜单列名称
	ParentID string         `json:"parent_id"` // 父菜单id
	Type     menu.TypeConst `json:"type"`      // 菜单列类型
	Disabled bool           `json:"disabled"`  // 是否禁用
}

type MenuAndColumnTree struct {
	MenuAndColumn                      // 嵌入基础菜单列信息
	Children      []*MenuAndColumnTree `json:"children"` // 子菜单列列表
}

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
	GetListByMenuIds(ctx context.Context, menuIds []string) ([]*MenuAndColumnTree, error)
	GetListAll(ctx context.Context, filters domainSystem.MenuColumnFilter) ([]domainSystem.MenuColumn, error)
}

type menuColumnService struct {
	repo     repositorySystem.MenuColumnRepository
	menuRepo repositorySystem.MenuRepository
}

func NewMenuColumnService(repo repositorySystem.MenuColumnRepository, menuRepo repositorySystem.MenuRepository) MenuColumnService {
	return &menuColumnService{
		repo:     repo,
		menuRepo: menuRepo,
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

// GetListByMenuIds 获取指定菜单下的所有列
func (svc *menuColumnService) GetListByMenuIds(ctx context.Context, menuIds []string) ([]*MenuAndColumnTree, error) {
	menuIdMap := make(map[string]bool)
	for _, id := range menuIds {
		menuIdMap[id] = true
	}

	// 获取全部菜单
	var menuAndColumn []MenuAndColumn

	menuList, _ := svc.menuRepo.GetListAll(ctx, domainSystem.MenuFilter{Status: true})
	for _, m := range menuList {
		if menuIdMap[m.Id] {
			menuAndColumn = append(menuAndColumn, MenuAndColumn{
				Id:       m.Id,
				Title:    m.Title,
				ParentID: m.ParentID,
				Type:     m.Type,
				Disabled: true,
			})
		}
	}
	// 获取指定菜单列
	menuColumnList, _ := svc.repo.GetListByMenuIds(ctx)
	for _, menuColumn := range menuColumnList {
		if menuIdMap[menuColumn.MenuId] {
			menuAndColumn = append(menuAndColumn, MenuAndColumn{
				Id:       menuColumn.Id,
				Title:    menuColumn.Title,
				ParentID: menuColumn.MenuId,
				Type:     4,
				Disabled: false,
			})
		}
	}

	// 构建菜单按钮树
	menuColumnMap := make(map[string]*MenuAndColumnTree)
	var roots []*MenuAndColumnTree

	if len(menuAndColumn) == 0 {
		return []*MenuAndColumnTree{}, nil
	}

	// 第一遍遍历，创建所有节点
	for _, menuColumn := range menuAndColumn {
		menuColumnMap[menuColumn.Id] = &MenuAndColumnTree{
			MenuAndColumn: menuColumn,
			Children:      []*MenuAndColumnTree{},
		}
	}

	// 第二遍遍历，构建树结构
	for _, menuColumn := range menuAndColumn {
		node := menuColumnMap[menuColumn.Id]
		if menuColumn.ParentID == "" || menuColumnMap[menuColumn.ParentID] == nil {
			roots = append(roots, node)
		} else {
			parent := menuColumnMap[menuColumn.ParentID]
			parent.Children = append(parent.Children, node)
		}
	}

	return roots, nil
}

// GetListAll 查询所有列表
func (svc *menuColumnService) GetListAll(ctx context.Context, filters domainSystem.MenuColumnFilter) ([]domainSystem.MenuColumn, error) {
	return svc.repo.GetListAll(ctx, filters)
}
