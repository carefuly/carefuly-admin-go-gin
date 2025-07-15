/**
 * Description：
 * FileName：bucket.go
 * Author：CJiaの用心
 * Create：2025/7/14 16:47:02
 * Remark：
 */

package tools

import (
	"context"
	"errors"
	domainTools "github.com/carefuly/carefuly-admin-go-gin/internal/domain/careful/tools"
	modelTools "github.com/carefuly/carefuly-admin-go-gin/internal/model/careful/tools"
	cacheTools "github.com/carefuly/carefuly-admin-go-gin/internal/repository/cache/careful/tools"
	daoTools "github.com/carefuly/carefuly-admin-go-gin/internal/repository/dao/careful/tools"
	"github.com/carefuly/carefuly-admin-go-gin/pkg/models"
	"go.uber.org/zap"
)

var (
	ErrBucketNotFound             = daoTools.ErrBucketNotFound
	ErrBucketNameDuplicate        = daoTools.ErrBucketNameDuplicate
	ErrBucketCodeDuplicate        = daoTools.ErrBucketCodeDuplicate
	ErrBucketDuplicate            = daoTools.ErrBucketDuplicate
	ErrBucketVersionInconsistency = daoTools.ErrBucketVersionInconsistency
)

type BucketRepository interface {
	Create(ctx context.Context, domain domainTools.Bucket) error
	Delete(ctx context.Context, id string) error
	BatchDelete(ctx context.Context, ids []string) error
	Update(ctx context.Context, domain domainTools.Bucket) error

	GetById(ctx context.Context, id string) (domainTools.Bucket, error)
	GetListPage(ctx context.Context, filters domainTools.BucketFilter) ([]domainTools.Bucket, int64, error)
	GetListAll(ctx context.Context, filters domainTools.BucketFilter) ([]domainTools.Bucket, error)

	CheckExistByCode(ctx context.Context, code, excludeId string) (bool, error)
	CheckExistByName(ctx context.Context, name, excludeId string) (bool, error)
}

type bucketRepository struct {
	dao   daoTools.BucketDAO
	cache cacheTools.BucketCache
}

func NewBucketRepository(dao daoTools.BucketDAO, cache cacheTools.BucketCache) BucketRepository {
	return &bucketRepository{
		dao:   dao,
		cache: cache,
	}
}

// Create 创建
func (repo *bucketRepository) Create(ctx context.Context, domain domainTools.Bucket) error {
	return repo.dao.Insert(ctx, repo.toEntity(domain))
}

// Delete 删除
func (repo *bucketRepository) Delete(ctx context.Context, id string) error {
	// 删除缓存
	if err := repo.cache.Del(ctx, id); err != nil {
		// 网络崩了，也可能是 redis 崩了
		// 缓存删除失败不影响主流程，记录日志即可
		zap.L().Error("删除缓存失败", zap.Error(err))
	}

	return repo.dao.Delete(ctx, id)
}

// BatchDelete 批量删除
func (repo *bucketRepository) BatchDelete(ctx context.Context, ids []string) error {
	// 删除缓存
	for _, id := range ids {
		if err := repo.cache.Del(ctx, id); err != nil {
			// 网络崩了，也可能是 redis 崩了
			// 缓存删除失败不影响主流程，记录日志即可
			zap.L().Error("删除缓存失败", zap.String("id", id), zap.Error(err))
			return err
		}
	}

	return repo.dao.BatchDelete(ctx, ids)
}

// Update 更新
func (repo *bucketRepository) Update(ctx context.Context, domain domainTools.Bucket) error {
	if err := repo.dao.Update(ctx, repo.toEntity(domain)); err != nil {
		return err
	}

	// 删除缓存
	if err := repo.cache.Del(ctx, domain.Id); err != nil {
		// 网络崩了，也可能是 redis 崩了
		// 缓存删除失败不影响主流程，记录日志即可
		zap.L().Error("删除缓存失败", zap.Error(err))
	}

	return nil
}

// GetById 根据ID获取
func (repo *bucketRepository) GetById(ctx context.Context, id string) (domainTools.Bucket, error) {
	domain, err := repo.cache.Get(ctx, id)
	if err == nil && domain != nil {
		return *domain, nil // 命中缓存
	}

	if err != nil && !errors.Is(err, cacheTools.ErrBucketNotExist) {
		// 缓存查询出错但不是"不存在"错误，记录日志但继续查DB
		zap.L().Error("缓存获取错误:", zap.Error(err))
	}

	entity, err := repo.dao.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, daoTools.ErrBucketNotFound) {
			// 数据库不存在，设置防穿透标记
			_ = repo.cache.SetNotFound(ctx, id)
			return domainTools.Bucket{}, daoTools.ErrBucketNotFound
		}
		return domainTools.Bucket{}, err
	}

	toDomain := repo.toDomain(entity)

	if err := repo.cache.Set(ctx, toDomain); err != nil {
		// 网络崩了，也可能是 redis 崩了
		// 缓存删除失败不影响主流程，记录日志即可
		zap.L().Error("设置缓存失败异常", zap.Error(err))
	}

	return toDomain, nil
}

// GetListPage 分页查询列表
func (repo *bucketRepository) GetListPage(ctx context.Context, filters domainTools.BucketFilter) ([]domainTools.Bucket, int64, error) {
	list, row, err := repo.dao.FindListPage(ctx, filters)
	if err != nil {
		return []domainTools.Bucket{}, 0, err
	}

	if len(list) == 0 {
		return []domainTools.Bucket{}, 0, nil
	}

	var domain []domainTools.Bucket
	for _, v := range list {
		domain = append(domain, repo.toDomain(v))
	}

	return domain, row, nil
}

// GetListAll 查询所有列表
func (repo *bucketRepository) GetListAll(ctx context.Context, filters domainTools.BucketFilter) ([]domainTools.Bucket, error) {
	list, err := repo.dao.FindListAll(ctx, filters)
	if err != nil {
		return []domainTools.Bucket{}, err
	}

	if len(list) == 0 {
		return []domainTools.Bucket{}, nil
	}

	var toDomain []domainTools.Bucket
	for _, v := range list {
		toDomain = append(toDomain, repo.toDomain(v))
	}

	return toDomain, nil
}

// CheckExistByName 检查name是否存在
func (repo *bucketRepository) CheckExistByName(ctx context.Context, name, excludeId string) (bool, error) {
	return repo.dao.CheckExistByName(ctx, name, excludeId)
}

// CheckExistByCode 检查code是否存在
func (repo *bucketRepository) CheckExistByCode(ctx context.Context, code, excludeId string) (bool, error) {
	return repo.dao.CheckExistByCode(ctx, code, excludeId)
}

// toEntity 转换为实体模型
func (repo *bucketRepository) toEntity(domain domainTools.Bucket) modelTools.Bucket {
	return modelTools.Bucket{
		CoreModels: models.CoreModels{
			Id:         domain.Id,
			Sort:       domain.Sort,
			Version:    domain.Version,
			Creator:    domain.Creator,
			Modifier:   domain.Modifier,
			BelongDept: domain.BelongDept,
			Remark:     domain.Remark,
		},
		Status: domain.Status,
		Name:   domain.Name,
		Code:   domain.Code,
		Size:   domain.Size,
	}
}

// toDomain 转换为领域模型
func (repo *bucketRepository) toDomain(entity *modelTools.Bucket) domainTools.Bucket {
	model := domainTools.Bucket{
		Bucket: *entity,
	}

	if entity.CreateTime != nil {
		model.CreateTime = entity.CreateTime.Format("2006-01-02 15:04:05")
	}
	if entity.UpdateTime != nil {
		model.UpdateTime = entity.UpdateTime.Format("2006-01-02 15:04:05")
	}

	return model
}
