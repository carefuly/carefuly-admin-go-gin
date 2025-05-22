/**
 * Description：
 * FileName：user.go
 * Author：CJiaの用心
 * Create：2025/5/12 14:34:23
 * Remark：
 */

package system

import (
	"github.com/carefuly/carefuly-admin-go-gin/internal/model/careful/system"
	"github.com/carefuly/carefuly-admin-go-gin/pkg/ginx/filters"
)

type User struct {
	system.User
	CreateTime string `json:"createTime"` // 创建时间
	UpdateTime string `json:"updateTime"` // 更新时间
}

type UserFilter struct {
	filters.Filters
	filters.Pagination
	Status   bool   `json:"status"`   // 状态
	Username string `json:"username"` // 用户名
}
