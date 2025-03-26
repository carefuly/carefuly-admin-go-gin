/**
 * Description：
 * FileName：register.go
 * Author：CJiaの用心
 * Create：2025/3/20 23:26:20
 * Remark：
 */

package controller

import (
	"errors"
	config "github.com/carefuly/carefuly-admin-go-gin/config/file"
	domain "github.com/carefuly/carefuly-admin-go-gin/internal/domain/auth"
	service "github.com/carefuly/carefuly-admin-go-gin/internal/service/auth"
	"github.com/carefuly/carefuly-admin-go-gin/pkg/response"
	validate "github.com/carefuly/carefuly-admin-go-gin/pkg/validator"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type RegisterController interface {
	RegisterRoutes(router *gin.RouterGroup)
	PassWordRegisterHandler(ctx *gin.Context)
}

type registerController struct {
	rely config.RelyConfig
	svc  service.RegisterService
}

func NewRegisterController(rely config.RelyConfig, svc service.RegisterService) RegisterController {
	return &registerController{
		rely: rely,
		svc:  svc,
	}
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"` // 用户账号
	Password string `json:"password" binding:"required,min=3,max=20"` // 密码
}

func (c *registerController) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/password-register", c.PassWordRegisterHandler)
}

// PassWordRegisterHandler
// @id PassWordRegisterHandler
// @Summary 密码注册
// @Description 密码注册
// @Tags 认证管理
// @Accept application/json
// @Produce application/json
// @Param RegisterRequest body RegisterRequest true "参数"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /v1/auth/password-register [post]
func (c *registerController) PassWordRegisterHandler(ctx *gin.Context) {
	var req RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validate.NewValidatorError(c.rely.Trans).HandleValidatorError(ctx, err)
		return
	}

	err := c.svc.Register(ctx, domain.Register{
		Username: req.Username,
		Password: req.Password,
	})

	switch {
	case err == nil:
		response.NewResponse().SuccessResponse(ctx, "注册成功", nil)
	case errors.Is(err, service.ErrDuplicateUsername):
		response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "用户账号已存在，请重新输入", nil)
	default:
		ctx.Set("internal", err.Error())
		zap.L().Error("密码注册异常", zap.Error(err))
		response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
	}
}
