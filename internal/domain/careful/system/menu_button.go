/**
 * Description：
 * FileName：menu_button.go
 * Author：CJiaの用心
 * Create：2025/6/9 11:44:48
 * Remark：
 */

package system

import (
	"github.com/carefuly/carefuly-admin-go-gin/internal/model/careful/system"
	"github.com/carefuly/carefuly-admin-go-gin/pkg/ginx/filters"
	"gorm.io/gorm"
)

type MenuButton struct {
	system.MenuButton
	CreateTime string `json:"createTime"` // 创建时间
	UpdateTime string `json:"updateTime"` // 更新时间
}

type MenuButtonFilter struct {
	filters.Filters
	filters.Pagination
	Status bool   `json:"status"`  // 状态
	Name   string `json:"name"`    // 名称
	Code   string `json:"code"`    // 权限值
	MenuId string `json:"menu_id"` // 菜单ID
}

func (f *MenuButtonFilter) Apply(query *gorm.DB) *gorm.DB {
	query = f.Filters.Apply(query).
		Where("status = ?", f.Status).
		Order("sort ASC, update_time DESC")

	if f.Name != "" {
		query = query.Where("name LIKE ?", "%"+f.Name+"%")
	}
	if f.Code != "" {
		query = query.Where("code LIKE ?", "%"+f.Code+"%")
	}
	if f.MenuId != "" {
		query = query.Where("menu_id = ?", f.MenuId)
	}

	return query
}
