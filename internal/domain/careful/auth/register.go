/**
 * Description：
 * FileName：register.go
 * Author：CJiaの用心
 * Create：2025/3/28 11:31:10
 * Remark：
 */

package auth

type Register struct {
	Username string `json:"username"` // 用户账号
	Password string `json:"password"` // 密码
}
