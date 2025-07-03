/**
 * Description：
 * FileName：menu_column.go
 * Author：CJiaの用心
 * Create：2025/6/9 11:44:59
 * Remark：
 */

package system

import (
	"github.com/carefuly/carefuly-admin-go-gin/internal/model/careful/system"
	"github.com/carefuly/carefuly-admin-go-gin/pkg/ginx/filters"
	"gorm.io/gorm"
)

type MenuColumn struct {
	system.MenuColumn
	CreateTime string `json:"createTime"` // 创建时间
	UpdateTime string `json:"updateTime"` // 更新时间
}

type MenuColumnFilter struct {
	filters.Filters
	filters.Pagination
	Status bool   `json:"status"`  // 状态
	Title  string `json:"title"`   // 标题
	Field  string `json:"field"`   // 字段名
	MenuId string `json:"menu_id"` // 菜单ID
}

func (f *MenuColumnFilter) Apply(query *gorm.DB) *gorm.DB {
	query = f.Filters.Apply(query).
		Where("status = ?", f.Status).
		Order("sort ASC, update_time DESC")

	if f.Title != "" {
		query = query.Where("title LIKE ?", "%"+f.Title+"%")
	}
	if f.Field != "" {
		query = query.Where("field LIKE ?", "%"+f.Field+"%")
	}
	if f.MenuId != "" {
		query = query.Where("menu_id = ?", f.MenuId)
	}

	return query
}
