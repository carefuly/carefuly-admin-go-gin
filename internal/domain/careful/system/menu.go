/**
 * Description：
 * FileName：menu.go
 * Author：CJiaの用心
 * Create：2025/5/13 16:31:27
 * Remark：
 */

package system

import (
	"github.com/carefuly/carefuly-admin-go-gin/internal/model/careful/system"
	"github.com/carefuly/carefuly-admin-go-gin/pkg/ginx/filters"
	"gorm.io/gorm"
)

type Menu struct {
	system.Menu
	Meta       map[string]any `json:"meta"`       // 元信息
	CreateTime string         `json:"createTime"` // 创建时间
	UpdateTime string         `json:"updateTime"` // 更新时间
}

type MenuFilter struct {
	filters.Filters
	filters.Pagination
	Status bool   `json:"status"` // 状态
	Title  string `json:"title"`  // 菜单名称
}

func (f *MenuFilter) Apply(query *gorm.DB) *gorm.DB {
	query = f.Filters.Apply(query).
		Where("status = ?", f.Status).
		Order("sort ASC, update_time DESC")

	if f.Title != "" {
		query = query.Where("title LIKE ?", "%"+f.Title+"%")
	}

	return query
}
