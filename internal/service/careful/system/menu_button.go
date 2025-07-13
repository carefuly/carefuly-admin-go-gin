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
	"github.com/carefuly/carefuly-admin-go-gin/pkg/constants/careful/system/menu"
)

type MenuAndButton struct {
	Id       string         `json:"id"`        // 菜单按钮id
	Title    string         `json:"title"`     // 菜单按钮名称
	ParentID string         `json:"parent_id"` // 父菜单id
	Type     menu.TypeConst `json:"type"`      // 菜单按钮类型
	Disabled  bool           `json:"disabled"`   // 是否禁用
}

type MenuAndButtonTree struct {
	MenuAndButton                      // 嵌入基础菜单按钮信息
	Children      []*MenuAndButtonTree `json:"children"` // 子菜单按钮列表
}

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
	GetListByMenuIds(ctx context.Context, menuIds []string) ([]*MenuAndButtonTree, error)
	GetListAll(ctx context.Context, filters domainSystem.MenuButtonFilter) ([]domainSystem.MenuButton, error)
}

type menuButtonService struct {
	repo     repositorySystem.MenuButtonRepository
	menuRepo repositorySystem.MenuRepository
}

func NewMenuButtonService(repo repositorySystem.MenuButtonRepository, menuRepo repositorySystem.MenuRepository) MenuButtonService {
	return &menuButtonService{
		repo:     repo,
		menuRepo: menuRepo,
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

// GetListByMenuIds 获取指定菜单下的所有按钮
func (svc *menuButtonService) GetListByMenuIds(ctx context.Context, menuIds []string) ([]*MenuAndButtonTree, error) {
	menuIdMap := make(map[string]bool)
	for _, id := range menuIds {
		menuIdMap[id] = true
	}

	// 获取全部菜单
	var menuAndButton []MenuAndButton

	menuList, _ := svc.menuRepo.GetListAll(ctx, domainSystem.MenuFilter{Status: true})
	for _, m := range menuList {
		if menuIdMap[m.Id] {
			menuAndButton = append(menuAndButton, MenuAndButton{
				Id:       m.Id,
				Title:    m.Title,
				ParentID: m.ParentID,
				Type:     m.Type,
				Disabled:  true,
			})
		}
	}
	// 获取指定菜单按钮
	menuButtonList, _ := svc.repo.GetListByMenuIds(ctx)
	for _, menuButton := range menuButtonList {
		if menuIdMap[menuButton.MenuId] {
			menuAndButton = append(menuAndButton, MenuAndButton{
				Id:       menuButton.Id,
				Title:    menuButton.Name,
				ParentID: menuButton.MenuId,
				Type:     3,
				Disabled:  false,
			})
		}
	}

	// 构建菜单按钮树
	menuButtonMap := make(map[string]*MenuAndButtonTree)
	var roots []*MenuAndButtonTree

	if len(menuAndButton) == 0 {
		return []*MenuAndButtonTree{}, nil
	}

	// 第一遍遍历，创建所有节点
	for _, menuButton := range menuAndButton {
		menuButtonMap[menuButton.Id] = &MenuAndButtonTree{
			MenuAndButton: menuButton,
			Children:      []*MenuAndButtonTree{},
		}
	}

	// 第二遍遍历，构建树结构
	for _, menuButton := range menuAndButton {
		node := menuButtonMap[menuButton.Id]
		if menuButton.ParentID == "" || menuButtonMap[menuButton.ParentID] == nil {
			roots = append(roots, node)
		} else {
			parent := menuButtonMap[menuButton.ParentID]
			parent.Children = append(parent.Children, node)
		}
	}

	return roots, nil
}

// GetListAll 查询所有列表
func (svc *menuButtonService) GetListAll(ctx context.Context, filters domainSystem.MenuButtonFilter) ([]domainSystem.MenuButton, error) {
	return svc.repo.GetListAll(ctx, filters)
}
