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
	"errors"
	repository "github.com/carefuly/carefuly-admin-go-gin/internal/repository/third"
	"github.com/carefuly/carefuly-admin-go-gin/internal/service/captcha"
	_const "github.com/carefuly/carefuly-admin-go-gin/pkg/const"
)

var (
	ErrCaptchaSendTooMany   = repository.ErrCaptchaSendTooMany
	ErrCaptchaNotFound      = repository.ErrCaptchaNotFound
	ErrUserBlocked          = repository.ErrUserBlocked
	ErrCaptchaIncorrect     = repository.ErrCaptchaIncorrect
	ErrCaptchaVerifyTooMany = repository.ErrCaptchaVerifyTooMany
)

type CaptchaService interface {
	Generate(ctx context.Context, t captcha.TypeCaptcha, bizType _const.BizTypeCaptcha) (string, string, string, error)
	Verify(ctx context.Context, id string, biz _const.BizTypeCaptcha, inputCode string) (bool, error)
}

type captchaService struct {
	repo  repository.CaptchaRepository
	digit captcha.Captcha
}

func NewCaptchaService(repo repository.CaptchaRepository) CaptchaService {
	return &captchaService{
		repo: repo,
	}
}

func (svc *captchaService) Generate(ctx context.Context, t captcha.TypeCaptcha, bizType _const.BizTypeCaptcha) (string, string, string, error) {
	// 根据类型生成验证码
	svc.digit = captcha.NewCaptchaGenerator(t)

	id, b64s, code, err := svc.digit.Generate()
	if err != nil {
		return id, b64s, code, err
	}

	return id, b64s, code, svc.repo.Set(ctx, id, code, bizType)
}

func (svc *captchaService) Verify(ctx context.Context, id string, biz _const.BizTypeCaptcha, inputCode string) (bool, error) {
	ok, err := svc.repo.Verify(ctx, id, biz, inputCode)

	switch {
	case errors.Is(err, repository.ErrUserBlocked):
		return false, repository.ErrUserBlocked
	case errors.Is(err, repository.ErrCaptchaNotFound):
		return false, repository.ErrCaptchaNotFound
	case errors.Is(err, repository.ErrCaptchaVerifyTooMany):
		return false, repository.ErrCaptchaVerifyTooMany
	case errors.Is(err, repository.ErrCaptchaIncorrect):
		return false, repository.ErrCaptchaIncorrect
	}

	return ok, err
}
