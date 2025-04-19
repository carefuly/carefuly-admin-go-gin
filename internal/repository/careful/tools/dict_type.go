/**
 * Description：
 * FileName：dict_type.go
 * Author：CJiaの用心
 * Create：2025/4/17 14:37:40
 * Remark：
 */

package tools

import (
	"context"
	"database/sql"
	"errors"
	cacheTools "github.com/carefuly/carefuly-admin-go-gin/internal/cache/careful/tools"
	"github.com/carefuly/carefuly-admin-go-gin/internal/dao/careful/tools"
	domainTools "github.com/carefuly/carefuly-admin-go-gin/internal/domain/careful/tools"
	modelTools "github.com/carefuly/carefuly-admin-go-gin/internal/model/careful/tools"
	"github.com/carefuly/carefuly-admin-go-gin/pkg/models"
	"go.uber.org/zap"
)

var (
	ErrNotSupportedTypeValue        = errors.New("不支持的字典类型")
	ErrDictTypeRecordNotFound       = tools.ErrDictTypeRecordNotFound
	ErrDictTypeNotFound             = tools.ErrDictTypeNotFound
	ErrDuplicateDictType            = tools.ErrDuplicateDictType
	ErrDictTypeVersionInconsistency = tools.ErrDictTypeVersionInconsistency
)

type DictTypeRepository interface {
	Create(ctx context.Context, domain domainTools.DictType) error

	Update(ctx context.Context, id string, domain domainTools.DictType) (int64, error)
	GetById(ctx context.Context, id string) (domainTools.DictType, error)
}

type dictTypeRepository struct {
	dao   tools.DictTypeDAO
	cache cacheTools.DictTypeCache
}

func NewDictTypeRepository(dao tools.DictTypeDAO, cache cacheTools.DictTypeCache) DictTypeRepository {
	return &dictTypeRepository{
		dao:   dao,
		cache: cache,
	}
}

// Create 创建
func (repo *dictTypeRepository) Create(ctx context.Context, domain domainTools.DictType) error {
	model, err := repo.toEntity(domain)
	if err != nil {
		return err
	}
	return repo.dao.Insert(ctx, model)
}

// Update 更新
func (repo *dictTypeRepository) Update(ctx context.Context, id string, domain domainTools.DictType) (int64, error) {
	model, err := repo.toEntity(domain)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := repo.dao.Update(ctx, id, model)
	if err != nil {
		return rowsAffected, err
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

// GetById 根据ID获取
func (repo *dictTypeRepository) GetById(ctx context.Context, id string) (domainTools.DictType, error) {
	domain, err := repo.cache.Get(ctx, id)
	if err == nil && domain != nil {
		return *domain, nil // 命中缓存
	}
	if err != nil && !errors.Is(err, cacheTools.ErrDictTypeNotExist) {
		// 缓存查询出错但不是"不存在"错误，记录日志但继续查DB
		zap.L().Error("缓存获取错误:", zap.Error(err))
	}

	model, err := repo.dao.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, tools.ErrDictTypeRecordNotFound) {
			// 数据库不存在，设置防穿透标记
			_ = repo.cache.SetNotFound(ctx, id)
			return domainTools.DictType{}, err
		}
		return domainTools.DictType{}, err
	}

	toDomain := repo.toDomain(model)
	if err := repo.cache.Set(ctx, toDomain); err != nil {
		// 网络崩了，也可能是 redis 崩了
		zap.L().Error("Redis异常", zap.Error(err))
	}

	return toDomain, nil
}

// toEntity 转换为实体模型
func (repo *dictTypeRepository) toEntity(domain domainTools.DictType) (modelTools.DictType, error) {
	// 公共字段
	model := modelTools.DictType{
		CoreModels: models.CoreModels{
			Version:    domain.Version,
			Creator:    domain.Creator,
			Modifier:   domain.Modifier,
			BelongDept: domain.BelongDept,
			Remark:     domain.Remark,
		},
		Name:      domain.Name,
		DictTag:   domain.DictTag,
		DictColor: domain.DictColor,
		DictName:  domain.DictName,
		TypeValue: domain.TypeValue,
		DictId:    domain.DictId,
	}

	// 根据类型设置值
	switch domain.TypeValue {
	case 0: // 字符串
		model.StrValue = sql.NullString{
			Valid:  true,
			String: domain.StrValue,
		}
		model.IntValue = sql.NullInt64{Valid: false}
		model.BoolValue = sql.NullBool{Valid: false}
	case 1: // 整型
		model.StrValue = sql.NullString{Valid: false}
		model.IntValue = sql.NullInt64{
			Valid: true,
			Int64: domain.IntValue,
		}
		model.BoolValue = sql.NullBool{Valid: false}
	case 2: // 布尔
		model.StrValue = sql.NullString{Valid: false}
		model.IntValue = sql.NullInt64{Valid: false}
		model.BoolValue = sql.NullBool{
			Valid: true,
			Bool:  domain.BoolValue,
		}
	default:
		return modelTools.DictType{}, ErrNotSupportedTypeValue
	}

	return model, nil
}

// toDomain 转换为领域模型
func (repo *dictTypeRepository) toDomain(model *modelTools.DictType) domainTools.DictType {
	return domainTools.DictType{
		DictType:   *model,
		CreateTime: model.CoreModels.CreateTime.Format("2006-01-02 15:04:05.000"),
		UpdateTime: model.CoreModels.UpdateTime.Format("2006-01-02 15:04:05.000"),
	}
}
