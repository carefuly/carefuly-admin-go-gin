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
	ErrDeptNotFound      = repositorySystem.ErrDeptNotFound
	ErrDeptNameDuplicate = repositorySystem.ErrDeptNameDuplicate
	ErrDeptCodeDuplicate = repositorySystem.ErrDeptCodeDuplicate
	ErrDeptDuplicate     = repositorySystem.ErrDeptDuplicate
)

type DeptService interface {
	Create(ctx context.Context, domain domainSystem.Dept) error

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
	exists, err := svc.repo.CheckExistByName(ctx, domain.Name, "")
	if err != nil {
		return err
	}
	if exists {
		return repositorySystem.ErrDeptNameDuplicate
	}

	exists, err = svc.repo.CheckExistByCode(ctx, domain.Code, "")
	if err != nil {
		return err
	}
	if exists {
		return repositorySystem.ErrDeptCodeDuplicate
	}

	if err := svc.repo.Create(ctx, domain); err != nil {
		if svc.IsDuplicateEntryError(err) {
			return repositorySystem.ErrDeptDuplicate
		}
		return err
	}

	return nil
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
	for _, dept := range list {
		node := deptMap[dept.Id]
		if dept.ParentID == "" || deptMap[dept.ParentID] == nil {
			roots = append(roots, node)
		} else {
			parent := deptMap[dept.ParentID]
			parent.Children = append(parent.Children, node)
		}
	}

	return roots, nil
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
