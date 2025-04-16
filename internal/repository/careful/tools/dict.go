/**
 * Description：
 * FileName：dict.go
 * Author：CJiaの用心
 * Create：2025/4/15 00:08:05
 * Remark：
 */

package tools

import (
	"context"
	"errors"
	cacheTools "github.com/carefuly/carefuly-admin-go-gin/internal/cache/careful/tools"
	"github.com/carefuly/carefuly-admin-go-gin/internal/dao/careful/tools"
	domainTools "github.com/carefuly/carefuly-admin-go-gin/internal/domain/careful/tools"
	modelTools "github.com/carefuly/carefuly-admin-go-gin/internal/model/careful/tools"
	"github.com/carefuly/carefuly-admin-go-gin/pkg/models"
	"go.uber.org/zap"
)

var (
	ErrDictRecordNotFound   = tools.ErrDictRecordNotFound
	ErrDictNotFound         = tools.ErrDictNotFound
	ErrDuplicateDict        = tools.ErrDuplicateDict
	ErrDuplicateDictName    = tools.ErrDuplicateDictName
	ErrDuplicateDictCode    = tools.ErrDuplicateDictCode
	ErrVersionInconsistency = tools.ErrVersionInconsistency
)

type DictRepository interface {
	Create(ctx context.Context, domain domainTools.Dict) error
	Delete(ctx context.Context, id string) (int64, error)
	BatchDelete(ctx context.Context, ids []string) error
	Update(ctx context.Context, id string, d domainTools.Dict) (int64, error)
	GetById(ctx context.Context, id string) (domainTools.Dict, error)
	GetListPage(ctx context.Context, filters domainTools.DictFilter) ([]domainTools.Dict, int64, error)
	GetListAll(ctx context.Context, filters domainTools.DictFilter) ([]domainTools.Dict, error)
	CheckExistByName(ctx context.Context, name, excludeId string) (bool, error)
	CheckExistByCode(ctx context.Context, code, excludeId string) (bool, error)
}

type dictRepository struct {
	dao   tools.DictDAO
	cache cacheTools.DictCache
}

func NewDictRepository(dao tools.DictDAO, cache cacheTools.DictCache) DictRepository {
	return &dictRepository{
		dao:   dao,
		cache: cache,
	}
}

// Create 创建字典
func (repo *dictRepository) Create(ctx context.Context, domain domainTools.Dict) error {
	return repo.dao.Insert(ctx, repo.toEntity(domain))
}

// Delete 删除字典
func (repo *dictRepository) Delete(ctx context.Context, id string) (int64, error) {
	rowsAffected, err := repo.dao.Delete(ctx, id)

	// 删除缓存
	err = repo.cache.Del(ctx, id)
	if err != nil {
		// 网络崩了，也可能是 redis 崩了
		zap.L().Error("Redis异常", zap.Error(err))
		return 0, err
	}

	return rowsAffected, err
}

// BatchDelete 批量删除字典
func (repo *dictRepository) BatchDelete(ctx context.Context, ids []string) error {
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

// Update 更新字典
func (repo *dictRepository) Update(ctx context.Context, id string, d domainTools.Dict) (int64, error) {
	rowsAffected, err := repo.dao.Update(ctx, id, repo.toEntity(d))
	if err != nil {
		return 0, err
	}

	// 删除缓存
	err = repo.cache.Del(ctx, id)
	if err != nil {
		// 网络崩了，也可能是 redis 崩了
		zap.L().Error("Redis异常", zap.Error(err))
		return 0, err
	}

	return rowsAffected, err
}

// GetById 根据ID获取字典
func (repo *dictRepository) GetById(ctx context.Context, id string) (domainTools.Dict, error) {
	domain, err := repo.cache.Get(ctx, id)
	if err == nil && domain != nil {
		return *domain, nil // 命中缓存
	}
	if err != nil && !errors.Is(err, cacheTools.ErrDictNotExist) {
		// 缓存查询出错但不是"不存在"错误，记录日志但继续查DB
		zap.L().Error("缓存获取错误:", zap.Error(err))
	}

	model, err := repo.dao.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, tools.ErrDictRecordNotFound) {
			// 数据库不存在，设置防穿透标记
			_ = repo.cache.SetNotFound(ctx, id)
			return domainTools.Dict{}, nil
		}
		return domainTools.Dict{}, err
	}

	toDomain := repo.toDomain(model)
	if err := repo.cache.Set(ctx, toDomain); err != nil {
		// 网络崩了，也可能是 redis 崩了
		zap.L().Error("Redis异常", zap.Error(err))
	}

	return toDomain, nil
}

// GetListPage 分页查询字典列表
func (repo *dictRepository) GetListPage(ctx context.Context, filters domainTools.DictFilter) ([]domainTools.Dict, int64, error) {
	list, row, err := repo.dao.FindListPage(ctx, filters)
	if err != nil {
		return []domainTools.Dict{}, row, err
	}

	var domain []domainTools.Dict
	for _, v := range list {
		domain = append(domain, repo.toDomain(v))
	}

	return domain, row, nil
}

// GetListAll 查询所有字典
func (repo *dictRepository) GetListAll(ctx context.Context, filters domainTools.DictFilter) ([]domainTools.Dict, error) {
	list, err := repo.dao.FindListAll(ctx, filters)
	if err != nil {
		return []domainTools.Dict{}, err
	}

	if len(list) == 0 {
		return []domainTools.Dict{}, nil
	}

	var domain []domainTools.Dict
	for _, v := range list {
		domain = append(domain, repo.toDomain(v))
	}

	return domain, nil
}

// CheckExistByName 检查字典名称是否已存在
func (repo *dictRepository) CheckExistByName(ctx context.Context, name, excludeId string) (bool, error) {
	return repo.dao.CheckExistByName(ctx, name, excludeId)
}

// CheckExistByCode 检查字典编码是否已存在
func (repo *dictRepository) CheckExistByCode(ctx context.Context, code, excludeId string) (bool, error) {
	return repo.dao.CheckExistByCode(ctx, code, excludeId)
}

// toEntity 转换为实体模型
func (repo *dictRepository) toEntity(domain domainTools.Dict) modelTools.Dict {
	return modelTools.Dict{
		CoreModels: models.CoreModels{
			Version:    domain.Version,
			Creator:    domain.Creator,
			Modifier:   domain.Modifier,
			BelongDept: domain.BelongDept,
			Remark:     domain.Remark,
		},
		Name:      domain.Name,
		Code:      domain.Code,
		Type:      domain.Type,
		TypeValue: domain.TypeValue,
	}
}

// toDomain 转换为领域模型
func (repo *dictRepository) toDomain(model *modelTools.Dict) domainTools.Dict {
	return domainTools.Dict{
		Dict:       *model,
		CreateTime: model.CoreModels.CreateTime.Format("2006-01-02 15:04:05.000"),
		UpdateTime: model.CoreModels.UpdateTime.Format("2006-01-02 15:04:05.000"),
	}
}
