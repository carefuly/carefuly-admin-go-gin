/**
 * Description：
 * FileName：types.go
 * Author：CJiaの用心
 * Create：2025/3/25 10:55:06
 * Remark：
 */

package captcha

// TypeCaptcha 验证码类型枚举
type TypeCaptcha int

const (
	DigitIotaCaptcha TypeCaptcha = 1 // 数字字母验证码
	MathIotaCaptcha  TypeCaptcha = 2 // 数学验证码
	SlideIotaCaptcha TypeCaptcha = 3 // 滑动验证码
)

// Captcha 验证码生成器接口
type Captcha interface {
	// Generate 生成验证码，返回验证码ID、问题和答案
	Generate() (id string, question string, answer string, err error)
	// Verify 验证验证码是否正确
	Verify(id, answer string) bool
	// GetImage 获取验证码图片(如果有的话)
	GetImage(id string) ([]byte, error)
}

func NewCaptchaGenerator(t TypeCaptcha) Captcha {
	switch t {
	case DigitIotaCaptcha:
		return NewDigitCaptcha()
	default:
		return NewDigitCaptcha()
	}
}
