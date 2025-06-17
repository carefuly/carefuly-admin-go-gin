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
	"gorm.io/gorm"
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
	Name     string `json:"name"`     // 姓名
	Email    string `json:"email"`    // 邮箱
	Mobile   string `json:"mobile"`   // 手机号
}

func (f *UserFilter) Apply(query *gorm.DB) *gorm.DB {
	query = f.Filters.Apply(query).
		Where("status = ?", f.Status).
		Order("update_time DESC, sort ASC")

	if f.Username != "" {
		query = query.Where("username LIKE ?", "%"+f.Username+"%")
	}
	if f.Name != "" {
		query = query.Where("name LIKE ?", "%"+f.Name+"%")
	}
	if f.Email != "" {
		query = query.Where("email LIKE ?", "%"+f.Email+"%")
	}
	if f.Mobile != "" {
		query = query.Where("mobile LIKE ?", "%"+f.Mobile+"%")
	}

	return query
}
