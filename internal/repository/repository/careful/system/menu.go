/**
 * Description：
 * FileName：menu.go
 * Author：CJiaの用心
 * Create：2025/5/13 16:33:06
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
	ErrMenuNotFound             = daoSystem.ErrMenuNotFound
	ErrMenuNameDuplicate        = daoSystem.ErrMenuNameDuplicate
	ErrMenuDuplicate            = daoSystem.ErrMenuDuplicate
	ErrMenuVersionInconsistency = daoSystem.ErrMenuVersionInconsistency
)

type MenuRepository interface {
	Create(ctx context.Context, domain domainSystem.Menu) error
	Delete(ctx context.Context, id string) (int64, error)
	BatchDelete(ctx context.Context, ids []string) error
	Update(ctx context.Context, domain domainSystem.Menu) error

	GetById(ctx context.Context, id string) (domainSystem.Menu, error)
	GetListAll(ctx context.Context, filters domainSystem.MenuFilter) ([]domainSystem.Menu, error)

	CheckExistByTypeAndTitleAndParentId(ctx context.Context, menuType int, title, parentId, excludeId string) (bool, error)
}

type menuRepository struct {
	dao   daoSystem.MenuDAO
	cache cacheSystem.MenuCache
}

func NewMenuRepository(dao daoSystem.MenuDAO, cache cacheSystem.MenuCache) MenuRepository {
	return &menuRepository{
		dao:   dao,
		cache: cache,
	}
}

// Create 创建
func (repo *menuRepository) Create(ctx context.Context, domain domainSystem.Menu) error {
	return repo.dao.Insert(ctx, repo.toEntity(domain))
}

// Delete 删除
func (repo *menuRepository) Delete(ctx context.Context, id string) (int64, error) {
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
func (repo *menuRepository) BatchDelete(ctx context.Context, ids []string) error {
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
func (repo *menuRepository) Update(ctx context.Context, domain domainSystem.Menu) error {
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
func (repo *menuRepository) GetById(ctx context.Context, id string) (domainSystem.Menu, error) {
	domain, err := repo.cache.Get(ctx, id)
	if err == nil && domain != nil {
		return *domain, nil // 命中缓存
	}
	if err != nil && !errors.Is(err, cacheSystem.ErrMenuNotExist) {
		// 缓存查询出错但不是"不存在"错误，记录日志但继续查DB
		zap.L().Error("缓存获取错误:", zap.Error(err))
	}

	entity, err := repo.dao.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, daoSystem.ErrMenuNotFound) {
			// 数据库不存在，设置防穿透标记
			_ = repo.cache.SetNotFound(ctx, id)
			return domainSystem.Menu{}, nil
		}
		return domainSystem.Menu{}, err
	}

	toDomain := repo.toDomain(entity)
	if err := repo.cache.Set(ctx, toDomain); err != nil {
		// 网络崩了，也可能是 redis 崩了
		zap.L().Error("Redis异常", zap.Error(err))
	}

	return toDomain, nil
}

// GetListAll 查询所有列表
func (repo *menuRepository) GetListAll(ctx context.Context, filters domainSystem.MenuFilter) ([]domainSystem.Menu, error) {
	list, err := repo.dao.FindListAll(ctx, filters)
	if err != nil {
		return []domainSystem.Menu{}, err
	}

	if len(list) == 0 {
		return []domainSystem.Menu{}, nil
	}

	var toDomain []domainSystem.Menu
	for _, v := range list {
		toDomain = append(toDomain, repo.toDomain(v))
	}

	return toDomain, nil
}

// CheckExistByTypeAndTitleAndParentId 检查type、title和parentId是否同时存在
func (repo *menuRepository) CheckExistByTypeAndTitleAndParentId(ctx context.Context, menuType int, title, parentId, excludeId string) (bool, error) {
	return repo.dao.CheckExistByTypeAndTitleAndParentId(ctx, menuType, title, parentId, excludeId)
}

// toEntity 转换为实体模型
func (repo *menuRepository) toEntity(domain domainSystem.Menu) modelSystem.Menu {
	return modelSystem.Menu{
		CoreModels: models.CoreModels{
			Id:         domain.Id,
			Sort:       domain.Sort,
			Version:    domain.Version,
			Creator:    domain.Creator,
			Modifier:   domain.Modifier,
			BelongDept: domain.BelongDept,
			Remark:     domain.Remark,
		},
		Type:        domain.Type,
		Icon:        domain.Icon,
		Title:       domain.Title,
		Name:        domain.Name,
		Component:   domain.Component,
		Path:        domain.Path,
		Redirect:    domain.Redirect,
		IsHide:      domain.IsHide,
		IsLink:      domain.IsLink,
		IsKeepAlive: domain.IsKeepAlive,
		IsFull:      domain.IsFull,
		IsAffix:     domain.IsAffix,
		ParentID:    domain.ParentID,
	}
}

// toDomain 转换为领域模型
func (repo *menuRepository) toDomain(entity *modelSystem.Menu) domainSystem.Menu {
	domain := domainSystem.Menu{
		Menu: *entity,
	}

	if entity.CreateTime != nil {
		domain.CreateTime = entity.CreateTime.Format("2006-01-02 15:04:05")
	}
	if entity.UpdateTime != nil {
		domain.UpdateTime = entity.UpdateTime.Format("2006-01-02 15:04:05")
	}

	return domain
}
