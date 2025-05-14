/**
 * Description：
 * FileName：menu.go
 * Author：CJiaの用心
 * Create：2025/5/13 16:31:27
 * Remark：
 */

package system

import (
	"github.com/carefuly/carefuly-admin-go-gin/internal/model/careful/system"
	"github.com/carefuly/carefuly-admin-go-gin/pkg/ginx/filters"
)

type Menu struct {
	system.Menu
	CreateTime string `json:"createTime"` // 创建时间
	UpdateTime string `json:"updateTime"` // 更新时间
}

type MenuFilter struct {
	filters.Filters
	filters.Pagination
	// Username string `json:"username"` // 用户名
}
