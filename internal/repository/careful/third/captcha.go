/**
 * Description：
 * FileName：captcha.go
 * Author：CJiaの用心
 * Create：2025/3/27 11:42:03
 * Remark：
 */

package third

import (
	"context"
	"github.com/carefuly/carefuly-admin-go-gin/internal/cache/careful/third"
	constantsCaptcha "github.com/carefuly/carefuly-admin-go-gin/pkg/constants/third/captcha"
)

var (
	ErrCaptchaSendTooMany   = third.ErrCaptchaSendTooMany
	ErrCaptchaNotFound      = third.ErrCaptchaNotFound
	ErrUserBlocked          = third.ErrUserBlocked
	ErrCaptchaIncorrect     = third.ErrCaptchaIncorrect
	ErrCaptchaVerifyTooMany = third.ErrCaptchaVerifyTooMany
)

type CaptchaRepository interface {
	Set(ctx context.Context, id, code string, bizType constantsCaptcha.BizTypeCaptcha) error
	Verify(ctx context.Context, id string, biz constantsCaptcha.BizTypeCaptcha, code string) (bool, error)
}

type captchaRepository struct {
	cache third.CaptchaCache
}

func NewCaptchaRepository(cache third.CaptchaCache) CaptchaRepository {
	return &captchaRepository{
		cache: cache,
	}
}

func (repo *captchaRepository) Set(ctx context.Context, id, code string, bizType constantsCaptcha.BizTypeCaptcha) error {
	return repo.cache.Set(ctx, id, code, bizType)
}

func (repo *captchaRepository) Verify(ctx context.Context, id string, biz constantsCaptcha.BizTypeCaptcha, code string) (bool, error) {
	return repo.cache.Verify(ctx, id, biz, code)
}
