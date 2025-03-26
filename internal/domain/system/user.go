/**
 * Description：
 * FileName：user.go
 * Author：CJiaの用心
 * Create：2025/3/26 22:26:09
 * Remark：
 */

package domain

import model "github.com/carefuly/carefuly-admin-go-gin/internal/model/system"

type User struct {
	model.User
	CreateTime string `json:"createTime"` // 创建时间
	UpdateTime string `json:"updateTime"` // 更新时间
}

