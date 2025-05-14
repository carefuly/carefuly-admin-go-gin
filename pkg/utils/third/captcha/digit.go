/**
 * Description：
 * FileName：digit.go
 * Author：CJiaの用心
 * Create：2025/5/13 00:21:27
 * Remark：
 */

package captcha

import (
	"github.com/mojocn/base64Captcha"
	"sync"
)

// DigitCaptcha 数字字母图形验证码
type DigitCaptcha struct {
	store  base64Captcha.Store
	driver *base64Captcha.DriverDigit
	mu     sync.Mutex
}

// NewDigitCaptcha 创建数字字母验证码实例
func NewDigitCaptcha(length int) Captcha {
	return &DigitCaptcha{
		store: base64Captcha.DefaultMemStore,
		driver: &base64Captcha.DriverDigit{
			Height:   80,
			Width:    240,
			Length:   length,
			MaxSkew:  0.7,
			DotCount: 80,
		},
	}
}

func (d *DigitCaptcha) Generate() (string, string, string, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	c := base64Captcha.NewCaptcha(d.driver, d.store)
	id, b64s, answer, err := c.Generate()
	return id, b64s, answer, err
}

func (d *DigitCaptcha) Verify(id, answer string) bool {
	return d.store.Verify(id, answer, true)
}

func (d *DigitCaptcha) GetImage(id string) ([]byte, error) {
	// 这个实现中图片已经以base64形式返回，所以这里不需要额外实现
	return nil, nil
}
