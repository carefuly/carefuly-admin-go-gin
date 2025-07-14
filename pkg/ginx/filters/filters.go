/**
 * Description：
 * FileName：filters.go
 * Author：CJiaの用心
 * Create：2025/5/12 17:00:50
 * Remark：
 */

package filters

import "gorm.io/gorm"

// QueryFiltersBuilder 查询构建器接口
type QueryFiltersBuilder interface {
	Apply(query *gorm.DB) *gorm.DB
}

// Pagination 分页查询构建器
type Pagination struct {
	Page     int `json:"page"`     // 当前页
	PageSize int `json:"pageSize"` // 每页显示的条数
}

// Filters 基础查询构建器
type Filters struct {
	Creator    string `json:"creator"`    // 创建人
	Modifier   string `json:"modifier"`   // 修改人
	BelongDept string `json:"belongDept"` // 数据归属部门
}

func (f *Filters) Apply(query *gorm.DB) *gorm.DB {
	if f.Creator != "" {
		query = query.Where("creator LIKE ?", "%"+f.Creator+"%")
	}
	if f.Modifier != "" {
		query = query.Where("modifier LIKE ?", "%"+f.Modifier+"%")
	}
	return query
}
