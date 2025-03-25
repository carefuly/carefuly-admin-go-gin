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
)

var (
	ErrCaptchaSendTooMany = cache.ErrCaptchaSendTooMany
)

type CaptchaRepository interface {
	Set(ctx context.Context, id, code, bizType string) error
}

type captchaRepository struct {
	cache cache.CaptchaCache
}

func NewCaptchaRepository(cache cache.CaptchaCache) CaptchaRepository {
	return &captchaRepository{
		cache: cache,
	}
}

func (repo *captchaRepository) Set(ctx context.Context, id, code, bizType string) error{
	return repo.cache.Set(ctx, id, code, bizType)
}

func (repo *captchaRepository) Verify(ctx context.Context, biz, to, code string) (bool, error) {
	// return repo.cache.Verify(ctx, biz, to, code)
	return true, nil
}

