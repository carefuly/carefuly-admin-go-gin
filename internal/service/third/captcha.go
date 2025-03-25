/**
 * Description：
 * FileName：captcha.go
 * Author：CJiaの用心
 * Create：2025/3/25 11:34:47
 * Remark：
 */

package service

import (
	"context"
	repository "github.com/carefuly/carefuly-admin-go-gin/internal/repository/third"
	"github.com/carefuly/carefuly-admin-go-gin/internal/service/captcha"
)

var (
	ErrCaptchaSendTooMany = repository.ErrCaptchaSendTooMany
)

type CaptchaService interface {
	Generate(ctx context.Context, t captcha.TypeCaptcha, bizType string) (string, string, string, error)
}

type captchaService struct {
	repo repository.CaptchaRepository
	digit captcha.Captcha
}

func NewCaptchaService(repo repository.CaptchaRepository) CaptchaService {
	return &captchaService{
		repo: repo,
	}
}

func (svc *captchaService) Generate(ctx context.Context, t captcha.TypeCaptcha, bizType string) (string, string, string, error) {
	// 根据类型生成验证码
	svc.digit = captcha.NewCaptchaGenerator(t)

	id, b64s, code, err := svc.digit.Generate()
	if err != nil {
		return id, b64s, code, err
	}

	return id, b64s, code, svc.repo.Set(ctx, id, code, bizType)
}
