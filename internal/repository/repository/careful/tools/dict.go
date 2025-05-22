/**
 * Description：
 * FileName：dict.go
 * Author：CJiaの用心
 * Create：2025/5/14 16:12:03
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
	ErrDictNotFound             = daoTools.ErrDictNotFound
	ErrDictNameDuplicate        = daoTools.ErrDictNameDuplicate
	ErrDictCodeDuplicate        = daoTools.ErrDictCodeDuplicate
	ErrDictDuplicate            = daoTools.ErrDictDuplicate
	ErrDictVersionInconsistency = daoTools.ErrDictVersionInconsistency
)

type DictRepository interface {
	Create(ctx context.Context, domain domainTools.Dict) error
	Delete(ctx context.Context, id string) (int64, error)
	Update(ctx context.Context, domain domainTools.Dict) error

	GetById(ctx context.Context, id string) (domainTools.Dict, error)
	GetByName(ctx context.Context, name string) (domainTools.Dict, error)
	GetListPage(ctx context.Context, filters domainTools.DictFilter) ([]domainTools.Dict, int64, error)
	GetListAll(ctx context.Context, filters domainTools.DictFilter) ([]domainTools.Dict, error)

	CheckExistByName(ctx context.Context, name, excludeId string) (bool, error)
	CheckExistByCode(ctx context.Context, code, excludeId string) (bool, error)
}

type dictRepository struct {
	dao   daoTools.DictDAO
	cache cacheTools.DictCache
}

func NewDictRepository(dao daoTools.DictDAO, cache cacheTools.DictCache) DictRepository {
	return &dictRepository{
		dao:   dao,
		cache: cache,
	}
}

// Create 创建
func (repo *dictRepository) Create(ctx context.Context, domain domainTools.Dict) error {
	return repo.dao.Insert(ctx, repo.toEntity(domain))
}

// Delete 删除
func (repo *dictRepository) Delete(ctx context.Context, id string) (int64, error) {
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

// Update 更新
func (repo *dictRepository) Update(ctx context.Context, domain domainTools.Dict) error {
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
func (repo *dictRepository) GetById(ctx context.Context, id string) (domainTools.Dict, error) {
	domain, err := repo.cache.Get(ctx, id)
	if err == nil && domain != nil {
		return *domain, nil // 命中缓存
	}
	if err != nil && !errors.Is(err, cacheTools.ErrDictNotExist) {
		// 缓存查询出错但不是"不存在"错误，记录日志但继续查DB
		zap.L().Error("缓存获取错误:", zap.Error(err))
	}

	entity, err := repo.dao.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, daoTools.ErrDictNotFound) {
			// 数据库不存在，设置防穿透标记
			_ = repo.cache.SetNotFound(ctx, id)
			return domainTools.Dict{}, nil
		}
		return domainTools.Dict{}, err
	}

	toDomain := repo.toDomain(entity)
	if err := repo.cache.Set(ctx, toDomain); err != nil {
		// 网络崩了，也可能是 redis 崩了
		zap.L().Error("Redis异常", zap.Error(err))
	}

	return toDomain, nil
}

// GetByName 根据name获取
func (repo *dictRepository) GetByName(ctx context.Context, name string) (domainTools.Dict, error) {
	model, err := repo.dao.FindByName(ctx, name)
	if err != nil {
		return domainTools.Dict{}, err
	}
	return repo.toDomain(model), nil
}

// GetListPage 分页查询列表
func (repo *dictRepository) GetListPage(ctx context.Context, filters domainTools.DictFilter) ([]domainTools.Dict, int64, error) {
	list, row, err := repo.dao.FindListPage(ctx, filters)
	if err != nil {
		return []domainTools.Dict{}, row, err
	}

	if len(list) == 0 {
		return []domainTools.Dict{}, row, nil
	}

	var domain []domainTools.Dict
	for _, v := range list {
		domain = append(domain, repo.toDomain(v))
	}

	return domain, row, nil
}

// GetListAll 查询所有列表
func (repo *dictRepository) GetListAll(ctx context.Context, filters domainTools.DictFilter) ([]domainTools.Dict, error) {
	list, err := repo.dao.FindListAll(ctx, filters)
	if err != nil {
		return []domainTools.Dict{}, err
	}

	if len(list) == 0 {
		return []domainTools.Dict{}, nil
	}

	var toDomain []domainTools.Dict
	for _, v := range list {
		toDomain = append(toDomain, repo.toDomain(v))
	}

	return toDomain, nil
}

// CheckExistByName 检查name是否存在
func (repo *dictRepository) CheckExistByName(ctx context.Context, name, excludeId string) (bool, error) {
	return repo.dao.CheckExistByName(ctx, name, excludeId)
}

// CheckExistByCode 检查code是否存在
func (repo *dictRepository) CheckExistByCode(ctx context.Context, code, excludeId string) (bool, error) {
	return repo.dao.CheckExistByCode(ctx, code, excludeId)
}

// toEntity 转换为实体模型
func (repo *dictRepository) toEntity(domain domainTools.Dict) modelTools.Dict {
	return modelTools.Dict{
		CoreModels: models.CoreModels{
			Id:         domain.Id,
			Sort:       domain.Sort,
			Creator:    domain.Creator,
			Version:    domain.Version,
			Modifier:   domain.Modifier,
			BelongDept: domain.BelongDept,
			Remark:     domain.Remark,
		},
		Status:    domain.Status,
		Name:      domain.Name,
		Code:      domain.Code,
		Type:      domain.Type,
		ValueType: domain.ValueType,
	}
}

// toDomain 转换为领域模型
func (repo *dictRepository) toDomain(entity *modelTools.Dict) domainTools.Dict {
	model := domainTools.Dict{
		Dict: *entity,
	}

	if entity.CreateTime != nil {
		model.CreateTime = entity.CreateTime.Format("2006-01-02 15:04:05")
	}
	if entity.UpdateTime != nil {
		model.UpdateTime = entity.UpdateTime.Format("2006-01-02 15:04:05")
	}

	return model
}
