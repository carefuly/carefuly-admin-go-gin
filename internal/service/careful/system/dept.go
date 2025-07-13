/**
 * Description：
 * FileName：dept.go
 * Author：CJiaの用心
 * Create：2025/5/15 17:03:07
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

// DeptTree 部门树形结构
type DeptTree struct {
	domainSystem.Dept             // 嵌入基础部门信息
	Children          []*DeptTree `json:"children"` // 子部门列表
}

var (
	ErrDeptNotFound             = repositorySystem.ErrDeptNotFound
	ErrDeptDuplicate            = repositorySystem.ErrDeptDuplicate
	ErrDeptVersionInconsistency = repositorySystem.ErrDeptVersionInconsistency
	ErrDeptChildNodes           = repositorySystem.ErrDeptChildNodes
)

type DeptService interface {
	Create(ctx context.Context, domain domainSystem.Dept) error
	Delete(ctx context.Context, id string) error
	BatchDelete(ctx context.Context, ids []string) error
	Update(ctx context.Context, domain domainSystem.Dept) error

	GetById(ctx context.Context, id string) (domainSystem.Dept, error)
	GetListAll(ctx context.Context, filter domainSystem.DeptFilter) ([]domainSystem.Dept, error)
	GetListTree(ctx context.Context, filter domainSystem.DeptFilter) ([]*DeptTree, error)
}

type deptService struct {
	repo repositorySystem.DeptRepository
}

func NewDeptService(repo repositorySystem.DeptRepository) DeptService {
	return &deptService{
		repo: repo,
	}
}

// Create 创建
func (svc *deptService) Create(ctx context.Context, domain domainSystem.Dept) error {
	// 检查name、code和parentId是否同时存在
	exists, err := svc.repo.CheckExistByNameAndCodeAndParentId(ctx, domain.Name, domain.Code, domain.ParentID, "")
	if err != nil {
		return err
	}
	if exists {
		return repositorySystem.ErrDeptDuplicate
	}

	// 创建
	if err := svc.repo.Create(ctx, domain); err != nil {
		if svc.IsDuplicateEntryError(err) {
			return repositorySystem.ErrDeptDuplicate
		}
		return err
	}

	return nil
}

// Delete 删除
func (svc *deptService) Delete(ctx context.Context, id string) error {
	rowsAffected, err := svc.repo.Delete(ctx, id)
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return repositorySystem.ErrDeptNotFound
	}
	return err
}

// BatchDelete 批量删除
func (svc *deptService) BatchDelete(ctx context.Context, ids []string) error {
	return svc.repo.BatchDelete(ctx, ids)
}

// Update 更新
func (svc *deptService) Update(ctx context.Context, domain domainSystem.Dept) error {
	// 检查name、code和parentId是否同时存在
	exists, err := svc.repo.CheckExistByNameAndCodeAndParentId(ctx, domain.Name, domain.Code, domain.ParentID, domain.Id)
	if err != nil {
		return err
	}
	if exists {
		return repositorySystem.ErrDeptDuplicate
	}

	// 更新
	if err := svc.repo.Update(ctx, domain); err != nil {
		switch {
		case svc.IsDuplicateEntryError(err):
			return repositorySystem.ErrDeptDuplicate
		case errors.Is(err, repositorySystem.ErrDeptVersionInconsistency):
			return repositorySystem.ErrDeptVersionInconsistency
		default:
			return err
		}
	}

	return nil
}

// GetById 获取详情
func (svc *deptService) GetById(ctx context.Context, id string) (domainSystem.Dept, error) {
	domain, err := svc.repo.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, repositorySystem.ErrDeptNotFound) {
			return domain, repositorySystem.ErrDeptNotFound
		}
		return domain, err
	}
	if domain.Id == "" {
		return domain, repositorySystem.ErrDeptNotFound
	}
	return domain, err
}

// GetListTree 获取树形结构
func (svc *deptService) GetListTree(ctx context.Context, filter domainSystem.DeptFilter) ([]*DeptTree, error) {
	list, err := svc.repo.GetListAll(ctx, filter)
	if err != nil {
		return nil, err
	}

	// 构建部门树
	deptMap := make(map[string]*DeptTree)
	var roots []*DeptTree

	if len(list) == 0 {
		return []*DeptTree{}, nil
	}

	// 第一遍遍历，创建所有节点
	for _, dept := range list {
		deptMap[dept.Id] = &DeptTree{
			Dept:     dept,
			Children: []*DeptTree{},
		}
	}

	// 第二遍遍历，构建树结构
	for _, node := range deptMap {
		parentID := node.Dept.ParentID
		// 关键修复：只通过ID存在性判断父节点
		if parentID == "" || deptMap[parentID] == nil {
			roots = append(roots, node) // 确认为根节点
		} else {
			// 安全添加到父节点
			parent := deptMap[parentID]
			parent.Children = append(parent.Children, node)
		}
	}

	return roots, nil
}

// GetListAll 查询所有列表
func (svc *deptService) GetListAll(ctx context.Context, filter domainSystem.DeptFilter) ([]domainSystem.Dept, error) {
	return svc.repo.GetListAll(ctx, filter)
}

// IsDuplicateEntryError 判断是否是唯一冲突错误
func (svc *deptService) IsDuplicateEntryError(err error) bool {
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		// MySQL 错误码 1062 表示唯一冲突
		return mysqlErr.Number == 1062
	}
	return false
}
