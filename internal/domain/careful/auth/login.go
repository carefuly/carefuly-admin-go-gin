/**
 * Description：
 * FileName：login.go
 * Author：CJiaの用心
 * Create：2025/3/28 11:31:18
 * Remark：
 */

package auth

import _const "github.com/carefuly/carefuly-admin-go-gin/pkg/const"

type Login struct {
	Id       string                `json:"id" `      // 验证码
	Code     string                `json:"code"`     // 验证码
	BizType  _const.BizTypeCaptcha `json:"bizType"`  // 验证码类型
	Username string                `json:"username"` // 用户账号
	Password string                `json:"password"` // 密码
}
