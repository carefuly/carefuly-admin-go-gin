/**
 * Description：
 * FileName：choices.go
 * Author：CJiaの用心
 * Create：2025/3/27 11:38:37
 * Remark：
 */

package _const

type BizTypeCaptcha string // 验证码业务类型

const (
	BizTypeCaptchaLogin BizTypeCaptcha = "BizCaptchaLogin" // 密码登录
)

type GenderConst int

const (
	GenderMale    GenderConst = iota // 男
	GenderFemale                     // 女
	GenderUnknown                    // 未知
)
