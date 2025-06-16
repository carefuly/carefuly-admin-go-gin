/**
 * Description：岗位服务层
 * FileName：post.go
 * Author：CJiaの用心
 * Create：2025/6/13 17:22:26
 * Remark：
 */

package system

import (
	"context"
	"errors"
	repositorySystem "github.com/carefuly/carefuly-admin-go-gin/internal/repository/repository/careful/system"
	"github.com/go-sql-driver/mysql"

	domainSystem "github.com/carefuly/carefuly-admin-go-gin/internal/domain/careful/system"
)

var (
	ErrPostNotFound             = repositorySystem.ErrPostNotFound
	ErrPostDuplicate            = repositorySystem.ErrPostDuplicate
	ErrPostVersionInconsistency = repositorySystem.ErrPostVersionInconsistency
)

type PostService interface {
	Create(ctx context.Context, domain domainSystem.Post) error
	Delete(ctx context.Context, id string) error
	BatchDelete(ctx context.Context, ids []string) error
	Update(ctx context.Context, domain domainSystem.Post) error

	GetById(ctx context.Context, id string) (domainSystem.Post, error)
	GetListPage(ctx context.Context, filter domainSystem.PostFilter) ([]domainSystem.Post, int64, error)
	GetListAll(ctx context.Context, filter domainSystem.PostFilter) ([]domainSystem.Post, error)
}

type postService struct {
	repo repositorySystem.PostRepository
}

func NewPostService(repo repositorySystem.PostRepository) PostService {
	return &postService{
		repo: repo,
	}
}

// Create 创建
func (svc *postService) Create(ctx context.Context, domain domainSystem.Post) error {
	// 检查name、code是否同时存在
	exists, err := svc.repo.CheckExistByNameAndCode(ctx, domain.Name, domain.Code, "")
	if err != nil {
		return err
	}
	if exists {
		return repositorySystem.ErrPostDuplicate
	}

	// 创建用户
	if err := svc.repo.Create(ctx, domain); err != nil {
		if svc.IsDuplicateEntryError(err) {
			return repositorySystem.ErrPostDuplicate
		}
		return err
	}

	return nil
}

// Delete 删除
func (svc *postService) Delete(ctx context.Context, id string) error {
	rowsAffected, err := svc.repo.Delete(ctx, id)
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return repositorySystem.ErrPostNotFound
	}
	return err
}

// BatchDelete 批量删除
func (svc *postService) BatchDelete(ctx context.Context, ids []string) error {
	return svc.repo.BatchDelete(ctx, ids)
}

// Update 更新
func (svc *postService) Update(ctx context.Context, domain domainSystem.Post) error {
	// 检查name、code是否同时存在
	exists, err := svc.repo.CheckExistByNameAndCode(ctx, domain.Name, domain.Code, domain.Id)
	if err != nil {
		return err
	}
	if exists {
		return repositorySystem.ErrPostDuplicate
	}

	// 更新用户
	if err := svc.repo.Update(ctx, domain); err != nil {
		switch {
		case svc.IsDuplicateEntryError(err):
			return repositorySystem.ErrPostDuplicate
		case errors.Is(err, repositorySystem.ErrPostVersionInconsistency):
			return repositorySystem.ErrPostVersionInconsistency
		default:
			return err
		}
	}

	return nil
}

// GetById 获取详情
func (svc *postService) GetById(ctx context.Context, id string) (domainSystem.Post, error) {
	domain, err := svc.repo.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, repositorySystem.ErrPostNotFound) {
			return domain, repositorySystem.ErrPostNotFound
		}
		return domain, err
	}
	if domain.Id == "" {
		return domain, repositorySystem.ErrPostNotFound
	}
	return domain, err
}

// GetListPage 分页查询列表
func (svc *postService) GetListPage(ctx context.Context, filter domainSystem.PostFilter) ([]domainSystem.Post, int64, error) {
	return svc.repo.GetListPage(ctx, filter)
}

// GetListAll 查询所有列表
func (svc *postService) GetListAll(ctx context.Context, filter domainSystem.PostFilter) ([]domainSystem.Post, error) {
	return svc.repo.GetListAll(ctx, filter)
}

// IsDuplicateEntryError 判断是否是唯一冲突错误
func (svc *postService) IsDuplicateEntryError(err error) bool {
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		// MySQL 错误码 1062 表示唯一冲突
		return mysqlErr.Number == 1062
	}
	return false
}
