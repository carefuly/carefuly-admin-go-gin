/**
 * Description：
 * FileName：captcha.go
 * Author：CJiaの用心
 * Create：2025/5/13 00:17:51
 * Remark：
 */

package third

import (
	"context"
	"github.com/carefuly/carefuly-admin-go-gin/internal/repository/cache/careful/third"
)

var (
	ErrCaptchaSendTooMany   = third.ErrCaptchaSendTooMany
	ErrCaptchaNotFound      = third.ErrCaptchaNotFound
	ErrUserBlocked          = third.ErrUserBlocked
	ErrCaptchaIncorrect     = third.ErrCaptchaIncorrect
	ErrCaptchaVerifyTooMany = third.ErrCaptchaVerifyTooMany
)

type CaptchaRepository interface {
	Set(ctx context.Context, id, code string, bizType string) error
	Verify(ctx context.Context, id string, biz string, code string) (bool, error)
}

type captchaRepository struct {
	cache third.CaptchaCache
}

func NewCaptchaRepository(cache third.CaptchaCache) CaptchaRepository {
	return &captchaRepository{
		cache: cache,
	}
}

func (repo *captchaRepository) Set(ctx context.Context, id, code string, bizType string) error {
	return repo.cache.Set(ctx, id, code, bizType)
}

func (repo *captchaRepository) Verify(ctx context.Context, id string, biz string, code string) (bool, error) {
	return repo.cache.Verify(ctx, id, biz, code)
}
