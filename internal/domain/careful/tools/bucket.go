/**
 * Description：
 * FileName：bucket.go
 * Author：CJiaの用心
 * Create：2025/7/14 16:35:52
 * Remark：
 */

package tools

import (
	"github.com/carefuly/carefuly-admin-go-gin/internal/model/careful/tools"
	"github.com/carefuly/carefuly-admin-go-gin/pkg/ginx/filters"
	"gorm.io/gorm"
)

type Bucket struct {
	tools.Bucket
	CreateTime string `json:"createTime"` // 创建时间
	UpdateTime string `json:"updateTime"` // 更新时间
}

type BucketFilter struct {
	filters.Pagination
	filters.Filters
	Status bool   `json:"status"` // 状态
	Name   string `json:"name"`   // 存储桶名称
	Code   string `json:"code"`   // 存储桶编码
}

func (f *BucketFilter) Apply(query *gorm.DB) *gorm.DB {
	query = f.Filters.Apply(query).
		Where("status = ?", f.Status).
		Order("sort ASC, update_time DESC")

	if f.Name != "" {
		query = query.Where("name LIKE ?", "%"+f.Name+"%")
	}
	if f.Code != "" {
		query = query.Where("code LIKE ?", "%"+f.Code+"%")
	}

	return query
}
