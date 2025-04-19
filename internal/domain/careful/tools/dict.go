/**
 * Description：
 * FileName：dict.go
 * Author：CJiaの用心
 * Create：2025/4/14 20:40:41
 * Remark：
 */

package tools

import (
	"github.com/carefuly/carefuly-admin-go-gin/internal/model/careful/tools"
	"github.com/carefuly/carefuly-admin-go-gin/pkg/ginx/query/filters"
	"gorm.io/gorm"
)

type Dict struct {
	tools.Dict
	CreateTime string `json:"createTime"` // 创建时间
	UpdateTime string `json:"updateTime"` // 更新时间
}

type DictFilter struct {
	filters.Filters
	filters.Pagination
	Name      string `json:"name"`      // 字典名称
	Code      string `json:"code"`      // 字典编码
	Type      int    `json:"type"`      // 字典类型
	TypeValue int    `json:"typeValue"` // 字典类型值
}

func (f *DictFilter) Apply(query *gorm.DB) *gorm.DB {
	query = f.Filters.Apply(query)

	if f.Name != "" {
		query = query.Where("name LIKE ?", "%"+f.Name+"%")
	}
	if f.Code != "" {
		query = query.Where("code LIKE ?", "%"+f.Code+"%")
	}
	if f.Type >= 0 {
		query = query.Where("type = ?", f.Type)
	}
	if f.TypeValue >= 0 {
		query = query.Where("typeValue = ?", f.TypeValue)
	}

	return query
}
