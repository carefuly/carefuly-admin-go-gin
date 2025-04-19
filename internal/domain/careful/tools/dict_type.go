/**
 * Description：
 * FileName：dict_type.go
 * Author：CJiaの用心
 * Create：2025/4/17 10:34:41
 * Remark：
 */

package tools

import (
	"github.com/carefuly/carefuly-admin-go-gin/internal/model/careful/tools"
	"github.com/carefuly/carefuly-admin-go-gin/pkg/ginx/query/filters"
	"gorm.io/gorm"
)

type DictType struct {
	tools.DictType
	StrValue   string `json:"strValue"`   // 字符串-字典信息值
	IntValue   int64  `json:"intValue"`   // 整型-字典信息值
	BoolValue  bool   `json:"boolValue"`  // 布尔-字典信息值
	CreateTime string `json:"createTime"` // 创建时间
	UpdateTime string `json:"updateTime"` // 更新时间
}

type DictTypeFilter struct {
	filters.Filters
	filters.Pagination
	Name    string `json:"name"`    // 字典信息名称
	DictTag string `json:"dictTag"` // 标签类型
	DictId  string `json:"dict_id"` // 字典ID
}

func (f *DictTypeFilter) Apply(query *gorm.DB) *gorm.DB {
	query = f.Filters.Apply(query)

	if f.Name != "" {
		query = query.Where("name LIKE ?", "%"+f.Name+"%")
	}
	if f.DictTag != "" {
		query = query.Where("dictTag LIKE ?", "%"+f.DictTag+"%")
	}
	if f.DictId != "" {
		query = query.Where("dict_id = ?", f.DictId)
	}

	return query
}
