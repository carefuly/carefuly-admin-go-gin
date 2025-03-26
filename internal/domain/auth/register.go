/**
 * Description：
 * FileName：register.go
 * Author：CJiaの用心
 * Create：2025/3/26 12:54:21
 * Remark：
 */

package domain

type Register struct {
	Username string `json:"username"` // 用户账号
	Password string `json:"password"` // 密码
}

