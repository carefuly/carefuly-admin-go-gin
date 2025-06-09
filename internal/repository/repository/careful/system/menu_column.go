/**
 * Description：
 * FileName：menu_column.go
 * Author：CJiaの用心
 * Create：2025/6/9 14:06:12
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
	ErrMenuColumnNotFound             = daoSystem.ErrMenuColumnNotFound
	ErrMenuColumnVersionInconsistency = daoSystem.ErrMenuColumnVersionInconsistency
)

type MenuColumnRepository interface {
	Create(ctx context.Context, domain domainSystem.MenuColumn) error
	Delete(ctx context.Context, id string) (int64, error)
	BatchDelete(ctx context.Context, ids []string) error
	Update(ctx context.Context, domain domainSystem.MenuColumn) error

	GetById(ctx context.Context, id string) (domainSystem.MenuColumn, error)
	GetListPage(ctx context.Context, filters domainSystem.MenuColumnFilter) ([]domainSystem.MenuColumn, int64, error)
	GetListAll(ctx context.Context, filters domainSystem.MenuColumnFilter) ([]domainSystem.MenuColumn, error)
}

type menuColumnRepository struct {
	dao   daoSystem.MenuColumnDAO
	cache cacheSystem.MenuColumnCache
}

func NewMenuColumnRepository(dao daoSystem.MenuColumnDAO, cache cacheSystem.MenuColumnCache) MenuColumnRepository {
	return &menuColumnRepository{
		dao:   dao,
		cache: cache,
	}
}

// Create 创建
func (repo *menuColumnRepository) Create(ctx context.Context, domain domainSystem.MenuColumn) error {
	return repo.dao.Insert(ctx, repo.toEntity(domain))
}

// Delete 删除
func (repo *menuColumnRepository) Delete(ctx context.Context, id string) (int64, error) {
	rowsAffected, err := repo.dao.Delete(ctx, id)

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
func (repo *menuColumnRepository) BatchDelete(ctx context.Context, ids []string) error {
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
func (repo *menuColumnRepository) Update(ctx context.Context, domain domainSystem.MenuColumn) error {
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
func (repo *menuColumnRepository) GetById(ctx context.Context, id string) (domainSystem.MenuColumn, error) {
	domain, err := repo.cache.Get(ctx, id)
	if err == nil && domain != nil {
		return *domain, nil // 命中缓存
	}
	if err != nil && !errors.Is(err, cacheSystem.ErrMenuColumnNotExist) {
		// 缓存查询出错但不是"不存在"错误，记录日志但继续查DB
		zap.L().Error("缓存获取错误:", zap.Error(err))
	}

	entity, err := repo.dao.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, daoSystem.ErrMenuColumnNotFound) {
			// 数据库不存在，设置防穿透标记
			_ = repo.cache.SetNotFound(ctx, id)
			return domainSystem.MenuColumn{}, nil
		}
		return domainSystem.MenuColumn{}, err
	}

	toDomain := repo.toDomain(entity)
	if err := repo.cache.Set(ctx, toDomain); err != nil {
		// 网络崩了，也可能是 redis 崩了
		zap.L().Error("Redis异常", zap.Error(err))
	}

	return toDomain, nil
}

// GetListPage 分页查询列表
func (repo *menuColumnRepository) GetListPage(ctx context.Context, filters domainSystem.MenuColumnFilter) ([]domainSystem.MenuColumn, int64, error) {
	list, row, err := repo.dao.FindListPage(ctx, filters)
	if err != nil {
		return []domainSystem.MenuColumn{}, row, err
	}

	if len(list) == 0 {
		return []domainSystem.MenuColumn{}, row, nil
	}

	var domain []domainSystem.MenuColumn
	for _, v := range list {
		domain = append(domain, repo.toDomain(v))
	}

	return domain, row, nil
}

// GetListAll 查询所有列表
func (repo *menuColumnRepository) GetListAll(ctx context.Context, filters domainSystem.MenuColumnFilter) ([]domainSystem.MenuColumn, error) {
	list, err := repo.dao.FindListAll(ctx, filters)
	if err != nil {
		return []domainSystem.MenuColumn{}, err
	}

	if len(list) == 0 {
		return []domainSystem.MenuColumn{}, nil
	}

	var toDomain []domainSystem.MenuColumn
	for _, v := range list {
		toDomain = append(toDomain, repo.toDomain(v))
	}

	return toDomain, nil
}

// toEntity 转换为实体模型
func (repo *menuColumnRepository) toEntity(domain domainSystem.MenuColumn) modelSystem.MenuColumn {
	return modelSystem.MenuColumn{
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
		Title:  domain.Title,
		Field:  domain.Field,
		Width:  domain.Width,
	}
}

// toDomain 转换为领域模型
func (repo *menuColumnRepository) toDomain(entity *modelSystem.MenuColumn) domainSystem.MenuColumn {
	model := domainSystem.MenuColumn{
		MenuColumn: *entity,
	}

	if entity.CreateTime != nil {
		model.CreateTime = entity.CreateTime.Format("2006-01-02 15:04:05")
	}
	if entity.UpdateTime != nil {
		model.UpdateTime = entity.UpdateTime.Format("2006-01-02 15:04:05")
	}

	return model
}
