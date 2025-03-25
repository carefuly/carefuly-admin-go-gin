/**
 * Description：
 * FileName：captcha.go
 * Author：CJiaの用心
 * Create：2025/3/25 14:08:32
 * Remark：
 */

package cache

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
)

var (
	//go:embed lua/set_captcha.lua
	luaSetCode string
	//go:embed lua/verify_code.lua
	luaVerifyCode string

	ErrCaptchaSendTooMany   = errors.New("发送太频繁")
	ErrUserBlocked       = errors.New("用户被限制")
	ErrCodeNotFound      = errors.New("未发送验证码")
	ErrCodeVerifyTooMany = errors.New("验证次数已耗尽")
	ErrCodeIncorrect     = errors.New("验证码错误")
)

type CaptchaCache interface {
	Set(ctx context.Context, id, code, bizType string) error
}

type captchaCache struct {
	cmd redis.Cmdable
}

func NewCaptchaCache(cmd redis.Cmdable) CaptchaCache {
	return &captchaCache{
		cmd: cmd,
	}
}

func (c *captchaCache) Set(ctx context.Context, id, code, bizType string) error {
	res, err := c.cmd.Eval(ctx, luaSetCode, []string{c.key(id, bizType)}, code).Int()
	if err != nil {
		// 调用 redis 出了问题
		return err
	}

	switch res {
	case -2:
		return errors.New("验证码存在，但是没有过期时间")
	case -1:
		return ErrCaptchaSendTooMany
	default:
		return nil
	}
}

func (c *captchaCache) Verify(ctx context.Context, biz, to, code string) (bool, error) {
	return true, nil
}

func (c *captchaCache) key(id, bizType string) string {
	// careful:id:bizType
	return fmt.Sprintf("careful:%s:%s", id, bizType)
}
