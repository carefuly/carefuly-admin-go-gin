/**
 * Description：
 * FileName：role.go
 * Author：CJiaの用心
 * Create：2025/6/10 17:19:20
 * Remark：
 */

package system

import (
	"github.com/carefuly/carefuly-admin-go-gin/internal/model/careful/system"
	"github.com/carefuly/carefuly-admin-go-gin/pkg/ginx/filters"
	"gorm.io/gorm"
)

type Role struct {
	system.Role
	CreateTime string `json:"createTime"` // 创建时间
	UpdateTime string `json:"updateTime"` // 更新时间
}

type RoleFilter struct {
	filters.Filters
	filters.Pagination
	Status bool   `json:"status"` // 状态
	Name   string `json:"name"`   // 角色名称
	Code   string `json:"code"`   // 角色编码
}

func (f *RoleFilter) Apply(query *gorm.DB) *gorm.DB {
	query = f.Filters.Apply(query).
		Where("status = ?", f.Status).
		Order("update_time DESC, sort ASC")

	if f.Name != "" {
		query = query.Where("name LIKE ?", "%"+f.Name+"%")
	}
	if f.Code != "" {
		query = query.Where("code LIKE ?", "%"+f.Code+"%")
	}

	return query
}
