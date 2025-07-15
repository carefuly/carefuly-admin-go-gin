/**
 * Description：
 * FileName：bucket.go
 * Author：CJiaの用心
 * Create：2025/7/14 17:31:29
 * Remark：
 */

package tools

import (
	"context"
	"errors"
	domainTools "github.com/carefuly/carefuly-admin-go-gin/internal/domain/careful/tools"
	repositoryTools "github.com/carefuly/carefuly-admin-go-gin/internal/repository/repository/careful/tools"
	"github.com/go-sql-driver/mysql"
	"strings"
)

var (
	ErrBucketNotFound             = repositoryTools.ErrBucketNotFound
	ErrBucketNameDuplicate        = repositoryTools.ErrBucketNameDuplicate
	ErrBucketCodeDuplicate        = repositoryTools.ErrBucketCodeDuplicate
	ErrBucketDuplicate            = repositoryTools.ErrBucketDuplicate
	ErrBucketVersionInconsistency = repositoryTools.ErrBucketVersionInconsistency
)

type BucketService interface {
	Create(ctx context.Context, domain domainTools.Bucket) error
	Delete(ctx context.Context, id string) error
	BatchDelete(ctx context.Context, ids []string) error
	Update(ctx context.Context, domain domainTools.Bucket) error

	GetById(ctx context.Context, id string) (domainTools.Bucket, error)
	GetListPage(ctx context.Context, filters domainTools.BucketFilter) ([]domainTools.Bucket, int64, error)
	GetListAll(ctx context.Context, filters domainTools.BucketFilter) ([]domainTools.Bucket, error)
}

type bucketService struct {
	repo repositoryTools.BucketRepository
}

func NewBucketService(repo repositoryTools.BucketRepository) BucketService {
	return &bucketService{
		repo: repo,
	}
}

// Create 创建
func (svc *bucketService) Create(ctx context.Context, domain domainTools.Bucket) error {
	exists, err := svc.repo.CheckExistByName(ctx, domain.Name, "")
	if err != nil {
		return err
	}
	if exists {
		return repositoryTools.ErrBucketNameDuplicate
	}

	exists, err = svc.repo.CheckExistByCode(ctx, domain.Code, "")
	if err != nil {
		return err
	}
	if exists {
		return repositoryTools.ErrBucketCodeDuplicate
	}

	if err := svc.repo.Create(ctx, domain); err != nil {
		// 分析具体冲突字段
		if field, isDuplicate := svc.IsDuplicateEntryError(err); isDuplicate {
			switch field {
			case "name":
				return repositoryTools.ErrBucketNameDuplicate
			case "code":
				return repositoryTools.ErrBucketCodeDuplicate
			case "all":
				return repositoryTools.ErrBucketDuplicate
			default:
				return err
			}
		}
	}

	return nil
}

// Delete 删除
func (svc *bucketService) Delete(ctx context.Context, id string) error {
	return svc.repo.Delete(ctx, id)
}

// BatchDelete 批量删除
func (svc *bucketService) BatchDelete(ctx context.Context, ids []string) error {
	return svc.repo.BatchDelete(ctx, ids)
}

// Update 更新
func (svc *bucketService) Update(ctx context.Context, domain domainTools.Bucket) error {
	exists, err := svc.repo.CheckExistByName(ctx, domain.Name, domain.Id)
	if err != nil {
		return err
	}
	if exists {
		return repositoryTools.ErrBucketNameDuplicate
	}

	exists, err = svc.repo.CheckExistByCode(ctx, domain.Code, domain.Id)
	if err != nil {
		return err
	}
	if exists {
		return repositoryTools.ErrBucketCodeDuplicate
	}

	err = svc.repo.Update(ctx, domain)
	if err != nil {
		// 分析具体冲突字段
		if field, isDuplicate := svc.IsDuplicateEntryError(err); isDuplicate {
			switch field {
			case "name":
				return repositoryTools.ErrBucketNameDuplicate
			case "code":
				return repositoryTools.ErrBucketCodeDuplicate
			case "all":
				return repositoryTools.ErrBucketDuplicate
			default:
				return err
			}
		}
	}

	return err
}

// GetById 获取详情
func (svc *bucketService) GetById(ctx context.Context, id string) (domainTools.Bucket, error) {
	return svc.repo.GetById(ctx, id)
}

// GetListPage 分页查询列表
func (svc *bucketService) GetListPage(ctx context.Context, filters domainTools.BucketFilter) ([]domainTools.Bucket, int64, error) {
	return svc.repo.GetListPage(ctx, filters)
}

// GetListAll 查询所有列表
func (svc *bucketService) GetListAll(ctx context.Context, filters domainTools.BucketFilter) ([]domainTools.Bucket, error) {
	return svc.repo.GetListAll(ctx, filters)
}

// IsDuplicateEntryError 分析错误消息中的索引名
func (svc *bucketService) IsDuplicateEntryError(err error) (string, bool) {
	var mysqlErr *mysql.MySQLError
	if !errors.As(err, &mysqlErr) {
		return "", false
	}

	// MySQL 错误码 1062 表示唯一冲突
	if mysqlErr.Number != 1062 {
		return "", false
	}

	// 分析错误消息中的索引名
	switch {
	case strings.Contains(mysqlErr.Message, "uni_careful_tools_bucket_name"):
		return "name", true
	case strings.Contains(mysqlErr.Message, "uni_careful_tools_bucket_code"):
		return "code", true
	default:
		return "all", true // 未知唯一键冲突
	}
}
