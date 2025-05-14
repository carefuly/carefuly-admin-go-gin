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
	domainSystem "github.com/carefuly/carefuly-admin-go-gin/internal/domain/careful/system"
	modelSystem "github.com/carefuly/carefuly-admin-go-gin/internal/model/careful/system"
	daoSystem "github.com/carefuly/carefuly-admin-go-gin/internal/repository/dao/careful/system"
	"github.com/carefuly/carefuly-admin-go-gin/pkg/models"
)

type MenuRepository interface {
	GetListAll(ctx context.Context, filters domainSystem.MenuFilter) ([]domainSystem.Menu, error)
}

type menuRepository struct {
	dao daoSystem.MenuDAO
	// cache cacheSystem.UserCache
}

func NewMenuRepository(dao daoSystem.MenuDAO) MenuRepository {
	return &menuRepository{
		dao: dao,
		// cache: cache,
	}
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

// toEntity 转换为实体模型
func (repo *menuRepository) toEntity(domain domainSystem.Menu) modelSystem.Menu {
	return modelSystem.Menu{
		CoreModels: models.CoreModels{
			Creator:    domain.Creator,
			Modifier:   domain.Modifier,
			BelongDept: domain.BelongDept,
			Remark:     domain.Remark,
		},
		Type:        domain.Type,
		Icon:        domain.Icon,
		Title:       domain.Title,
		Permission:  domain.Permission,
		Name:        domain.Name,
		Component:   domain.Component,
		Api:         domain.Api,
		Method:      domain.Method,
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
	user := domainSystem.Menu{
		Menu: *entity,
	}

	if entity.CreateTime != nil {
		user.CreateTime = entity.CreateTime.Format("2006-01-02 15:04:05")
	}
	if entity.UpdateTime != nil {
		user.UpdateTime = entity.UpdateTime.Format("2006-01-02 15:04:05")
	}

	return user
}
