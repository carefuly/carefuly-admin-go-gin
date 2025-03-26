/**
 * Description：
 * FileName：register.go
 * Author：CJiaの用心
 * Create：2025/3/26 12:54:21
 * Remark：
 */

package domain

import model "github.com/carefuly/carefuly-admin-go-gin/internal/model/system"

type Register struct {
	Username string `json:"username"` // 用户账号
	Password string `json:"password"` // 密码
}

type User struct {
	model.User
	CreateTime string `json:"createTime"` // 创建时间
	UpdateTime string `json:"updateTime"` // 更新时间
}
