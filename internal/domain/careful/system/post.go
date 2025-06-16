/**
 * Description：
 * FileName：post.go
 * Author：CJiaの用心
 * Create：2025/6/13 17:13:57
 * Remark：
 */

package system

import (
	"github.com/carefuly/carefuly-admin-go-gin/internal/model/careful/system"
	"github.com/carefuly/carefuly-admin-go-gin/pkg/ginx/filters"
	"gorm.io/gorm"
)

type Post struct {
	system.Post
	CreateTime string `json:"createTime"` // 创建时间
	UpdateTime string `json:"updateTime"` // 更新时间
}

type PostFilter struct {
	filters.Filters
	filters.Pagination
	Status bool   `json:"status"` // 状态
	Name   string `json:"name"`   // 岗位名称
	Code   string `json:"code"`   // 岗位编码
}

func (f *PostFilter) Apply(query *gorm.DB) *gorm.DB {
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
