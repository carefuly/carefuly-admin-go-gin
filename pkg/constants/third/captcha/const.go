/**
 * Description：
 * FileName：const.go
 * Author：CJiaの用心
 * Create：2025/4/16 15:50:48
 * Remark：
 */

package captcha

type BizTypeCaptcha string // 验证码业务类型

const (
	BizTypeCaptchaLogin BizTypeCaptcha = "BizCaptchaLogin" // 验证码密码登录
)
