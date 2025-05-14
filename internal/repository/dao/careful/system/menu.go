/**
 * Description：
 * FileName：menu.go
 * Author：CJiaの用心
 * Create：2025/5/13 16:29:32
 * Remark：
 */

package system

import (
	"context"
	domainSystem "github.com/carefuly/carefuly-admin-go-gin/internal/domain/careful/system"
	"github.com/carefuly/carefuly-admin-go-gin/internal/model/careful/system"
	"github.com/carefuly/carefuly-admin-go-gin/pkg/ginx/filters"
	"gorm.io/gorm"
)

var (
// ErrMenuNotFound             = gorm.ErrRecordNotFound
// ErrMenuDuplicate            = errors.New("用户名已存在")
// ErrMenuDuplicate            = errors.New("用户信息已存在")
// ErrMenuVersionInconsistency = errors.New("数据已被修改，请刷新后重试")
)

type MenuDAO interface {
	FindListAll(ctx context.Context, filter domainSystem.MenuFilter) ([]*system.Menu, error)
}

type GORMMenuDAO struct {
	db *gorm.DB
}

func NewGORMMenuDAO(db *gorm.DB) MenuDAO {
	return &GORMMenuDAO{
		db: db,
	}
}

// FindListAll 获取所有列表
func (dao *GORMMenuDAO) FindListAll(ctx context.Context, filter domainSystem.MenuFilter) ([]*system.Menu, error) {
	var models []*system.Menu

	query := dao.db.WithContext(ctx).Model(&system.Menu{}).
		Order("create_time ASC, sort ASC")

	// if filter.Username != "" {
	// 	query = query.Where("username LIKE ?", "%"+filter.Username+"%")
	// }

	// 查询
	if err := query.Find(&models).Error; err != nil {
		return nil, err
	}

	return models, nil
}

// buildQuery 构建查询条件
func (dao *GORMMenuDAO) buildQuery(ctx context.Context, filter domainSystem.MenuFilter) *gorm.DB {
	builder := &domainSystem.UserFilter{
		Filters: filters.Filters{
			Creator:  filter.Creator,
			Modifier: filter.Modifier,
		},
	}
	return builder.Apply(dao.db.WithContext(ctx).Model(&system.User{}))
}
