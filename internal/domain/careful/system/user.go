/**
 * Description：
 * FileName：user.go
 * Author：CJiaの用心
 * Create：2025/3/28 11:33:55
 * Remark：
 */

package system

import "github.com/carefuly/carefuly-admin-go-gin/internal/model/careful/system"

type User struct {
	system.User
	CreateTime string `json:"createTime"` // 创建时间
	UpdateTime string `json:"updateTime"` // 更新时间
}
