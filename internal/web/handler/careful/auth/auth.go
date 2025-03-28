/**
 * Description：
 * FileName：auth.go
 * Author：CJiaの用心
 * Create：2025/3/28 11:50:29
 * Remark：
 */

package auth

import (
	"errors"
	config "github.com/carefuly/carefuly-admin-go-gin/config/file"
	"github.com/carefuly/carefuly-admin-go-gin/internal/domain/careful/auth"
	"github.com/carefuly/carefuly-admin-go-gin/internal/service/careful/system"
	"github.com/carefuly/carefuly-admin-go-gin/internal/service/careful/third"
	_const "github.com/carefuly/carefuly-admin-go-gin/pkg/const"
	"github.com/carefuly/carefuly-admin-go-gin/pkg/response"
	validate "github.com/carefuly/carefuly-admin-go-gin/pkg/validator"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type RegisterHandler interface {
	RegisterRoutes(router *gin.RouterGroup)
	PassWordRegisterHandler(ctx *gin.Context)
	PasswordCaptchaLoginHandler(ctx *gin.Context)
}

type registerHandler struct {
	rely       config.RelyConfig
	svc        system.UserService
	captchaSvc third.CaptchaService
}

func NewRegisterHandler(rely config.RelyConfig, svc system.UserService, captchaSvc third.CaptchaService) RegisterHandler {
	return &registerHandler{
		rely:       rely,
		svc:        svc,
		captchaSvc: captchaSvc,
	}
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"` // 用户账号
	Password string `json:"password" binding:"required,min=3,max=20"` // 密码
}

type LoginRequest struct {
	Username string                `json:"username" binding:"required,min=3,max=50"` // 用户账号
	Password string                `json:"password" binding:"required,min=3,max=20"` // 密码
	Id       string                `json:"id" binding:"required"`                    // 验证码
	Code     string                `json:"code" binding:"required"`                  // 验证码
	BizType  _const.BizTypeCaptcha `json:"bizType" binding:"required"`               // 验证码类型
}

type LoginResponse struct {
	Token string `json:"token"` // 登录令牌
}

func (c *registerHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/password-register", c.PassWordRegisterHandler)
	router.POST("/password-login", c.PasswordCaptchaLoginHandler)
}

// PassWordRegisterHandler
// @id PassWordRegisterHandler
// @Summary 账号密码注册
// @Description 账号密码注册
// @Tags 认证管理
// @Accept application/json
// @Produce application/json
// @Param RegisterRequest body RegisterRequest true "参数"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /v1/auth/password-register [post]
func (c *registerHandler) PassWordRegisterHandler(ctx *gin.Context) {
	var req RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validate.NewValidatorError(c.rely.Trans).HandleValidatorError(ctx, err)
		return
	}

	err := c.svc.UsernameAndPasswordCreate(ctx, auth.Register{
		Username: req.Username,
		Password: req.Password,
	})

	switch {
	case err == nil:
		response.NewResponse().SuccessResponse(ctx, "注册成功", nil)
	case errors.Is(err, system.ErrDuplicateUsername):
		response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "用户账号已存在，请重新输入", nil)
	default:
		ctx.Set("internal", err.Error())
		zap.L().Error("密码注册异常", zap.Error(err))
		response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
	}
}

// PasswordCaptchaLoginHandler
// @id PasswordCaptchaLoginHandler
// @Summary 图形验证码密码登录
// @Description 图形验证码密码登录
// @Tags 认证管理
// @Accept application/json
// @Produce application/json
// @Param LoginRequest body LoginRequest true "参数"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} response.Response
// @Router /v1/auth/password-login [post]
func (c *registerHandler) PasswordCaptchaLoginHandler(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validate.NewValidatorError(c.rely.Trans).HandleValidatorError(ctx, err)
		return
	}

	ok, err := c.captchaSvc.Verify(ctx, req.Id, req.BizType, req.Code)

	switch {
	case ok:
		token, err := c.svc.FindByUserName(ctx, c.rely, auth.Login{
			Username: req.Username,
			Password: req.Password,
		})

		switch {
		case err == nil:
			response.NewResponse().SuccessResponse(ctx, "登录成功", LoginResponse{Token: "Bearer " + token})
		case errors.Is(err, system.ErrInvalidUserOrPassword):
			response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "用户账号或者密码不对", nil)
		default:
			ctx.Set("internal", err.Error())
			zap.L().Error("密码登录异常", zap.Error(err))
			response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
		}
	case errors.Is(err, third.ErrUserBlocked):
		response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "操作过于频繁，请10分钟后再试", nil)
	case errors.Is(err, third.ErrCaptchaNotFound):
		response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "验证码已过期/验证码不存在", nil)
	case errors.Is(err, third.ErrCaptchaVerifyTooMany):
		response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "验证次数已耗尽，请10分钟后再试", nil)
	case errors.Is(err, third.ErrCaptchaIncorrect):
		response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "验证码错误，请重新输入", nil)
	default:
		ctx.Set("internal", err.Error())
		zap.L().Error("验证验证码异常", zap.Error(err))
		response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
	}
}
