/**
 * Description：
 * FileName：login.go
 * Author：CJiaの用心
 * Create：2025/3/26 22:25:21
 * Remark：
 */

package domain

import _const "github.com/carefuly/carefuly-admin-go-gin/pkg/const"

type Login struct {
	Username string                `json:"username"` // 用户账号
	Password string                `json:"password"` // 密码
	Id       string                `json:"id" `      // 验证码
	Code     string                `json:"code"`     // 验证码
	BizType  _const.BizTypeCaptcha `json:"bizType"`  // 验证码类型
}
