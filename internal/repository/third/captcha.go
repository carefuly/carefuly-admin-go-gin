/**
 * Description：
 * FileName：captcha.go
 * Author：CJiaの用心
 * Create：2025/3/25 14:06:47
 * Remark：
 */

package repository

import (
	"context"
	cache "github.com/carefuly/carefuly-admin-go-gin/internal/cache/third"
	_const "github.com/carefuly/carefuly-admin-go-gin/pkg/const"
)

var (
	ErrCaptchaSendTooMany   = cache.ErrCaptchaSendTooMany
	ErrCaptchaNotFound      = cache.ErrCaptchaNotFound
	ErrUserBlocked          = cache.ErrUserBlocked
	ErrCaptchaIncorrect     = cache.ErrCaptchaIncorrect
	ErrCaptchaVerifyTooMany = cache.ErrCaptchaVerifyTooMany
)

type CaptchaRepository interface {
	Set(ctx context.Context, id, code string, bizType _const.BizTypeCaptcha) error
	Verify(ctx context.Context, id string, biz _const.BizTypeCaptcha, code string) (bool, error)
}

type captchaRepository struct {
	cache cache.CaptchaCache
}

func NewCaptchaRepository(cache cache.CaptchaCache) CaptchaRepository {
	return &captchaRepository{
		cache: cache,
	}
}

func (repo *captchaRepository) Set(ctx context.Context, id, code string, bizType _const.BizTypeCaptcha) error {
	return repo.cache.Set(ctx, id, code, bizType)
}

func (repo *captchaRepository) Verify(ctx context.Context, id string, biz _const.BizTypeCaptcha, code string) (bool, error) {
	return repo.cache.Verify(ctx, id, biz, code)
}
