/**
 * Description：
 * FileName：captcha.go
 * Author：CJiaの用心
 * Create：2025/5/13 00:25:22
 * Remark：
 */

package third

import (
	"errors"
	config "github.com/carefuly/carefuly-admin-go-gin/config/file"
	"github.com/carefuly/carefuly-admin-go-gin/internal/service/careful/third"
	"github.com/carefuly/carefuly-admin-go-gin/pkg/ginx/response"
	"github.com/carefuly/carefuly-admin-go-gin/pkg/utils/third/captcha"
	validate "github.com/carefuly/carefuly-admin-go-gin/pkg/validator"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type CaptchaHandler interface {
	RegisterRoutes(router *gin.RouterGroup)
	GenerateCaptchaHandler(ctx *gin.Context)
}

type captchaHandler struct {
	rely config.RelyConfig
	svc  third.CaptchaService
}

func NewCaptchaController(rely config.RelyConfig, svc third.CaptchaService) CaptchaHandler {
	return &captchaHandler{
		rely: rely,
		svc:  svc,
	}
}

type CaptchaRequest struct {
	Type    captcha.TypeCaptcha `form:"type" binding:"required"`                              // 验证码类型
	BizType string              `form:"bizType" binding:"required" example:"BizCaptchaLogin"` // 业务类型
}

type CaptchaResponse struct {
	Id   string `json:"id"`   // 验证码Id
	Img  string `json:"img"`  // 验证码图片
	Code string `json:"code"` // 验证码
}

func (h *captchaHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/generateCaptcha", h.GenerateCaptchaHandler)
}

// GenerateCaptchaHandler
// @Summary 生成指定业务验证码
// @Description 生成指定业务验证码
// @Tags 第三方业务管理
// @Accept application/json
// @Produce application/json
// @Param CaptchaRequest query CaptchaRequest true "请求参数"
// @Success 200 {object} CaptchaResponse
// @Failure 400 {object} response.Response
// @Router /v1/third/generateCaptcha [get]
func (h *captchaHandler) GenerateCaptchaHandler(ctx *gin.Context) {
	var req CaptchaRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		validate.NewValidatorError(h.rely.Trans).HandleValidatorError(ctx, err)
		return
	}

	id, b64s, code, err := h.svc.Generate(ctx, req.Type, req.BizType)
	// 不管成功还是失败, 控制台都要返回验证码
	zap.L().Info("当前生成的验证码", zap.String("id", id), zap.String("code", code))

	switch {
	case err == nil:
		response.NewResponse().SuccessResponse(ctx, "验证码生成成功", CaptchaResponse{
			Id:   id,
			Img:  b64s,
			Code: code,
		})
	case errors.Is(err, third.ErrCaptchaSendTooMany):
		response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "验证码发送太频繁，请稍后再试", nil)
	default:
		ctx.Set("internal", err.Error())
		zap.L().Error("验证码生成异常", zap.Error(err))
		response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
	}
}
