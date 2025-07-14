/**
 * Description：
 * FileName：post.go
 * Author：CJiaの用心
 * Create：2025/6/13 17:34:29
 * Remark：
 */

package system

import (
	"context"
	"errors"
	domainSystem "github.com/carefuly/carefuly-admin-go-gin/internal/domain/careful/system"
	modelSystem "github.com/carefuly/carefuly-admin-go-gin/internal/model/careful/system"
	cacheSystem "github.com/carefuly/carefuly-admin-go-gin/internal/repository/cache/careful/system"
	daoSystem "github.com/carefuly/carefuly-admin-go-gin/internal/repository/dao/careful/system"
	"github.com/carefuly/carefuly-admin-go-gin/pkg/models"
	"go.uber.org/zap"
)

var (
	ErrPostNotFound             = daoSystem.ErrPostNotFound
	ErrPostDuplicate            = daoSystem.ErrPostDuplicate
	ErrPostVersionInconsistency = daoSystem.ErrPostVersionInconsistency
)

type PostRepository interface {
	Create(ctx context.Context, domain domainSystem.Post) error
	Delete(ctx context.Context, id string) (int64, error)
	BatchDelete(ctx context.Context, ids []string) error
	Update(ctx context.Context, domain domainSystem.Post) error

	GetById(ctx context.Context, id string) (domainSystem.Post, error)
	GetListPage(ctx context.Context, filters domainSystem.PostFilter) ([]domainSystem.Post, int64, error)
	GetListAll(ctx context.Context, filters domainSystem.PostFilter) ([]domainSystem.Post, error)

	CheckExistByNameAndCode(ctx context.Context, name, code, excludeId string) (bool, error)
}

type postRepository struct {
	dao   daoSystem.PostDAO
	cache cacheSystem.PostCache
}

func NewPostRepository(dao daoSystem.PostDAO, cache cacheSystem.PostCache) PostRepository {
	return &postRepository{
		dao:   dao,
		cache: cache,
	}
}

// Create 创建
func (repo *postRepository) Create(ctx context.Context, domain domainSystem.Post) error {
	return repo.dao.Insert(ctx, repo.toEntity(domain))
}

// Delete 删除
func (repo *postRepository) Delete(ctx context.Context, id string) (int64, error) {
	rowsAffected, err := repo.dao.Delete(ctx, id)
	if err != nil {
		return rowsAffected, err
	}

	// 删除缓存
	err = repo.cache.Del(ctx, id)
	if err != nil {
		// 网络崩了，也可能是 redis 崩了
		zap.L().Error("Redis异常", zap.Error(err))
		return rowsAffected, err
	}

	return rowsAffected, err
}

// BatchDelete 批量删除
func (repo *postRepository) BatchDelete(ctx context.Context, ids []string) error {
	err := repo.dao.BatchDelete(ctx, ids)
	if err != nil {
		return err
	}

	// 删除缓存
	for _, val := range ids {
		err = repo.cache.Del(ctx, val)
		if err != nil {
			// 网络崩了，也可能是 redis 崩了
			zap.L().Error("Redis异常", zap.Error(err))
			return err
		}
	}

	return err
}

// Update 更新
func (repo *postRepository) Update(ctx context.Context, domain domainSystem.Post) error {
	err := repo.dao.Update(ctx, repo.toEntity(domain))
	if err != nil {
		return err
	}

	// 删除缓存
	err = repo.cache.Del(ctx, domain.Id)
	if err != nil {
		// 网络崩了，也可能是 redis 崩了
		zap.L().Error("Redis异常", zap.Error(err))
		return err
	}

	return nil
}

// GetById 根据ID获取
func (repo *postRepository) GetById(ctx context.Context, id string) (domainSystem.Post, error) {
	domain, err := repo.cache.Get(ctx, id)
	if err == nil && domain != nil {
		return *domain, nil // 命中缓存
	}
	if err != nil && !errors.Is(err, cacheSystem.ErrPostNotExist) {
		// 缓存查询出错但不是"不存在"错误，记录日志但继续查DB
		zap.L().Error("缓存获取错误:", zap.Error(err))
	}

	entity, err := repo.dao.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, daoSystem.ErrPostNotFound) {
			// 数据库不存在，设置防穿透标记
			_ = repo.cache.SetNotFound(ctx, id)
			return domainSystem.Post{}, nil
		}
		return domainSystem.Post{}, err
	}

	toDomain := repo.toDomain(entity)
	if err := repo.cache.Set(ctx, toDomain); err != nil {
		// 网络崩了，也可能是 redis 崩了
		zap.L().Error("Redis异常", zap.Error(err))
	}

	return toDomain, nil
}

// GetListPage 分页查询列表
func (repo *postRepository) GetListPage(ctx context.Context, filters domainSystem.PostFilter) ([]domainSystem.Post, int64, error) {
	list, row, err := repo.dao.FindListPage(ctx, filters)
	if err != nil {
		return []domainSystem.Post{}, row, err
	}

	if len(list) == 0 {
		return []domainSystem.Post{}, 0, nil
	}

	var toDomain []domainSystem.Post
	for _, v := range list {
		toDomain = append(toDomain, repo.toDomain(v))
	}

	return toDomain, row, nil
}

// GetListAll 查询所有列表
func (repo *postRepository) GetListAll(ctx context.Context, filters domainSystem.PostFilter) ([]domainSystem.Post, error) {
	list, err := repo.dao.FindListAll(ctx, filters)
	if err != nil {
		return []domainSystem.Post{}, err
	}

	if len(list) == 0 {
		return []domainSystem.Post{}, nil
	}

	var toDomain []domainSystem.Post
	for _, v := range list {
		toDomain = append(toDomain, repo.toDomain(v))
	}

	return toDomain, nil
}

// CheckExistByNameAndCode 检查name、code是否同时存在
func (repo *postRepository) CheckExistByNameAndCode(ctx context.Context, name, code, excludeId string) (bool, error) {
	return repo.dao.CheckExistByNameAndCode(ctx, name, code, excludeId)
}

// toEntity 转换为实体模型
func (repo *postRepository) toEntity(domain domainSystem.Post) modelSystem.Post {
	return modelSystem.Post{
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
	}
}

// toDomain 转换为领域模型
func (repo *postRepository) toDomain(entity *modelSystem.Post) domainSystem.Post {
	domain := domainSystem.Post{
		Post: *entity,
	}

	if entity.CreateTime != nil {
		domain.CreateTime = entity.CreateTime.Format("2006-01-02 15:04:05")
	}
	if entity.UpdateTime != nil {
		domain.UpdateTime = entity.UpdateTime.Format("2006-01-02 15:04:05")
	}

	return domain
}
