/**
 * Description：
 * FileName：captcha.go
 * Author：CJiaの用心
 * Create：2025/5/12 14:40:07
 * Remark：
 */

package third

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
	//go:embed lua/verify_captcha.lua
	luaVerifyCode string

	ErrCaptchaSendTooMany   = errors.New("生成验证码太频繁")
	ErrCaptchaNotFound      = errors.New("验证码不存在")
	ErrUserBlocked          = errors.New("用户被限制")
	ErrCaptchaIncorrect     = errors.New("验证码错误")
	ErrCaptchaVerifyTooMany = errors.New("验证次数已耗尽")
)

type CaptchaCache interface {
	Set(ctx context.Context, id, code string, bizType string) error
	Verify(ctx context.Context, id string, biz string, code string) (bool, error)
}

type captchaCache struct {
	cmd redis.Cmdable
}

func NewCaptchaCache(cmd redis.Cmdable) CaptchaCache {
	return &captchaCache{
		cmd: cmd,
	}
}

func (c *captchaCache) Set(ctx context.Context, id, code string, bizType string) error {
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

func (c *captchaCache) Verify(ctx context.Context, id string, biz string, code string) (bool, error) {
	res, err := c.cmd.Eval(ctx, luaVerifyCode, []string{c.key(id, biz)}, code).Int()
	if err != nil {
		// 调用 redis 出了问题
		return false, err
	}

	switch res {
	case -4:
		return false, ErrCaptchaNotFound
	case -3:
		return false, ErrUserBlocked
	case -2:
		return false, ErrCaptchaIncorrect
	case -1:
		return false, ErrCaptchaVerifyTooMany
	default:
		return true, nil
	}
}

func (c *captchaCache) key(id string, bizType string) string {
	// careful:id:bizType
	return fmt.Sprintf("careful:captcha:%s:%s", id, bizType)
}
