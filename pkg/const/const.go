/**
 * Description：
 * FileName：const.go
 * Author：CJiaの用心
 * Create：2025/3/26 11:44:20
 * Remark：
 */

package _const

type BizTypeCaptcha string // 验证码业务类型

const (
	BizTypeCaptchaLogin BizTypeCaptcha = "passLogin" // 密码登录
)

type GenderConst int

const (
	GenderMale    GenderConst = iota // 男
	GenderFemale                     // 女
	GenderUnknown                    // 未知
)
