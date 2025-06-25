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
	"errors"
	domainSystem "github.com/carefuly/carefuly-admin-go-gin/internal/domain/careful/system"
	repositorySystem "github.com/carefuly/carefuly-admin-go-gin/internal/repository/repository/careful/system"
	"github.com/go-sql-driver/mysql"
)

// MenuTree 菜单树形结构
type MenuTree struct {
	domainSystem.Menu             // 嵌入基础菜单信息
	Children          []*MenuTree `json:"children"` // 子部门列表
}

var (
	ErrMenuNotFound             = repositorySystem.ErrMenuNotFound
	ErrMenuNameDuplicate        = repositorySystem.ErrMenuNameDuplicate
	ErrMenuDuplicate            = repositorySystem.ErrMenuDuplicate
	ErrMenuVersionInconsistency = repositorySystem.ErrMenuVersionInconsistency
)

type MenuService interface {
	Create(ctx context.Context, domain domainSystem.Menu) error
	Delete(ctx context.Context, id string) error
	BatchDelete(ctx context.Context, ids []string) error
	Update(ctx context.Context, domain domainSystem.Menu) error

	GetById(ctx context.Context, id string) (domainSystem.Menu, error)
	GetListTree(ctx context.Context, filter domainSystem.MenuFilter) ([]*MenuTree, error)
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

// Create 创建
func (svc *menuService) Create(ctx context.Context, domain domainSystem.Menu) error {
	// 检查type、title和parentId是否同时存在
	exists, err := svc.repo.CheckExistByTypeAndTitleAndParentId(ctx, int(domain.Type), domain.Title, domain.ParentID, "")
	if err != nil {
		return err
	}
	if exists {
		return repositorySystem.ErrMenuNameDuplicate
	}

	// 创建用户
	if err := svc.repo.Create(ctx, domain); err != nil {
		if svc.IsDuplicateEntryError(err) {
			return repositorySystem.ErrMenuNameDuplicate
		}
		return err
	}

	return nil
}

// Delete 删除
func (svc *menuService) Delete(ctx context.Context, id string) error {
	rowsAffected, err := svc.repo.Delete(ctx, id)
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return repositorySystem.ErrMenuNotFound
	}
	return err
}

// BatchDelete 批量删除
func (svc *menuService) BatchDelete(ctx context.Context, ids []string) error {
	return svc.repo.BatchDelete(ctx, ids)
}

// Update 更新
func (svc *menuService) Update(ctx context.Context, domain domainSystem.Menu) error {
	// 检查type、title和parentId是否同时存在
	exists, err := svc.repo.CheckExistByTypeAndTitleAndParentId(ctx, int(domain.Type), domain.Title, domain.ParentID, domain.Id)
	if err != nil {
		return err
	}
	if exists {
		return repositorySystem.ErrMenuNameDuplicate
	}

	// 更新用户
	if err := svc.repo.Update(ctx, domain); err != nil {
		switch {
		case svc.IsDuplicateEntryError(err):
			return repositorySystem.ErrMenuNameDuplicate
		case errors.Is(err, repositorySystem.ErrMenuVersionInconsistency):
			return repositorySystem.ErrMenuVersionInconsistency
		default:
			return err
		}
	}

	return nil
}

// GetById 获取详情
func (svc *menuService) GetById(ctx context.Context, id string) (domainSystem.Menu, error) {
	domain, err := svc.repo.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, repositorySystem.ErrMenuNotFound) {
			return domain, repositorySystem.ErrMenuNotFound
		}
		return domain, err
	}
	if domain.Id == "" {
		return domain, repositorySystem.ErrMenuNotFound
	}
	return domain, err
}

// GetListTree 获取菜单树形结构
func (svc *menuService) GetListTree(ctx context.Context, filter domainSystem.MenuFilter) ([]*MenuTree, error) {
	list, err := svc.repo.GetListAll(ctx, filter)
	if err != nil {
		return nil, err
	}

	// 构建部门树
	menuMap := make(map[string]*MenuTree)
	var roots []*MenuTree

	if len(list) == 0 {
		return []*MenuTree{}, nil
	}

	// 第一遍遍历，创建所有节点
	for _, menu := range list {
		menu.Meta = map[string]any{
			"icon":      menu.Icon,
			"title":     menu.Title,
			"isHide":    menu.IsHide,
			"isLink":    menu.IsLink,
			"isFull":    menu.IsFull,
			"keepAlive": menu.IsKeepAlive,
			"fixedTab":  menu.IsAffix,
		}
		menuMap[menu.Id] = &MenuTree{
			Menu:     menu,
			Children: []*MenuTree{},
		}
	}

	// 第二遍遍历，构建树结构
	for _, menu := range list {
		node := menuMap[menu.Id]
		if menu.ParentID == "" || menuMap[menu.ParentID] == nil {
			roots = append(roots, node)
		} else {
			parent := menuMap[menu.ParentID]
			parent.Children = append(parent.Children, node)
		}
	}

	return roots, nil
}

// GetListAll 查询所有列表
func (svc *menuService) GetListAll(ctx context.Context, filter domainSystem.MenuFilter) ([]domainSystem.Menu, error) {
	return svc.repo.GetListAll(ctx, filter)
}

// IsDuplicateEntryError 判断是否是唯一冲突错误
func (svc *menuService) IsDuplicateEntryError(err error) bool {
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		// MySQL 错误码 1062 表示唯一冲突
		return mysqlErr.Number == 1062
	}
	return false
}
