/**
 * Description：
 * FileName：captcha.go
 * Author：CJiaの用心
 * Create：2025/3/27 11:50:13
 * Remark：
 */

package third

import (
	"context"
	"errors"
	"github.com/carefuly/carefuly-admin-go-gin/internal/repository/careful/third"
	_const "github.com/carefuly/carefuly-admin-go-gin/pkg/const"
	"github.com/carefuly/carefuly-admin-go-gin/pkg/third/captcha"
)

var (
	ErrCaptchaSendTooMany   = third.ErrCaptchaSendTooMany
	ErrCaptchaNotFound      = third.ErrCaptchaNotFound
	ErrUserBlocked          = third.ErrUserBlocked
	ErrCaptchaIncorrect     = third.ErrCaptchaIncorrect
	ErrCaptchaVerifyTooMany = third.ErrCaptchaVerifyTooMany
)

type CaptchaService interface {
	Generate(ctx context.Context, t captcha.TypeCaptcha, bizType _const.BizTypeCaptcha) (string, string, string, error)
	Verify(ctx context.Context, id string, biz _const.BizTypeCaptcha, inputCode string) (bool, error)
}

type captchaService struct {
	repo  third.CaptchaRepository
	digit captcha.Captcha
}

func NewCaptchaService(repo third.CaptchaRepository) CaptchaService {
	return &captchaService{
		repo: repo,
	}
}

func (svc *captchaService) Generate(ctx context.Context, t captcha.TypeCaptcha, bizType _const.BizTypeCaptcha) (string, string, string, error) {
	// 根据类型生成验证码
	svc.digit = svc.NewCaptchaGenerator(t)

	id, b64s, code, err := svc.digit.Generate()
	if err != nil {
		return id, b64s, code, err
	}

	return id, b64s, code, svc.repo.Set(ctx, id, code, bizType)
}

func (svc *captchaService) Verify(ctx context.Context, id string, biz _const.BizTypeCaptcha, inputCode string) (bool, error) {
	ok, err := svc.repo.Verify(ctx, id, biz, inputCode)

	switch {
	case errors.Is(err, third.ErrUserBlocked):
		return false, third.ErrUserBlocked
	case errors.Is(err, third.ErrCaptchaNotFound):
		return false, third.ErrCaptchaNotFound
	case errors.Is(err, third.ErrCaptchaVerifyTooMany):
		return false, third.ErrCaptchaVerifyTooMany
	case errors.Is(err, third.ErrCaptchaIncorrect):
		return false, third.ErrCaptchaIncorrect
	}

	return ok, err
}

func (svc *captchaService) NewCaptchaGenerator(t captcha.TypeCaptcha) captcha.Captcha {
	switch t {
	case captcha.DigitIotaCaptcha:
		return captcha.NewDigitCaptcha(6)
	default:
		return captcha.NewDigitCaptcha(6)
	}
}
