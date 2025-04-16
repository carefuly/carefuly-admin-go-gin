/**
 * Description：
 * FileName：filters.go
 * Author：CJiaの用心
 * Create：2025/4/14 20:57:58
 * Remark：
 */

package filters

import (
	"gorm.io/gorm"
)

type QueryFiltersBuilder interface {
	Apply(query *gorm.DB) *gorm.DB
}

// Filters 基础查询构建器
type Filters struct {
	Creator    string `json:"creator"`    // 创建人
	Modifier   string `json:"modifier"`   // 修改人
	BelongDept string `json:"belongDept"` // 数据归属部门
	Status     bool   `json:"status"`     // 状态
}

func (f *Filters) Apply(query *gorm.DB) *gorm.DB {
	if f.Creator != "" {
		query = query.Where("creator LIKE ?", "%"+f.Creator+"%")
	}
	if f.Modifier != "" {
		query = query.Where("modifier LIKE ?", "%"+f.Modifier+"%")
	}
	query = query.Where("status = ?", f.Status).
		Order("update_time DESC, sort ASC")

	return query
}

type Pagination struct {
	Page     int `json:"page"`     // 当前页
	PageSize int `json:"pageSize"` // 每页显示的条数
}
