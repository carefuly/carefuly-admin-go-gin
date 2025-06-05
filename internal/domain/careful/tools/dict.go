/**
 * Description：
 * FileName：dict.go
 * Author：CJiaの用心
 * Create：2025/5/14 11:26:32
 * Remark：
 */

package tools

import (
	"github.com/carefuly/carefuly-admin-go-gin/internal/model/careful/tools"
	"github.com/carefuly/carefuly-admin-go-gin/pkg/constants/careful/tools/dict"
	"github.com/carefuly/carefuly-admin-go-gin/pkg/ginx/filters"
	"gorm.io/gorm"
)

type Dict struct {
	tools.Dict
	CreateTime string `json:"createTime"` // 创建时间
	UpdateTime string `json:"updateTime"` // 更新时间
}

type DictFilter struct {
	filters.Pagination
	filters.Filters
	Status    bool                `json:"status"`    // 状态
	Name      string              `json:"name"`      // 字典名称
	Code      string              `json:"code"`      // 字典编码
	Type      dict.TypeConst      `json:"type"`      // 字典分类
	ValueType dict.TypeValueConst `json:"valueType"` // 字典值类型
}

func (f *DictFilter) Apply(query *gorm.DB) *gorm.DB {
	query = f.Filters.Apply(query).
		Where("status = ?", f.Status).
		Order("update_time DESC, sort ASC")

	if f.Name != "" {
		query = query.Where("name LIKE ?", "%"+f.Name+"%")
	}
	if f.Code != "" {
		query = query.Where("code LIKE ?", "%"+f.Code+"%")
	}
	if f.Type > 0 {
		query = query.Where("type = ?", f.Type)
	}
	if f.ValueType > 0 {
		query = query.Where("valueType = ?", f.ValueType)
	}

	return query
}
