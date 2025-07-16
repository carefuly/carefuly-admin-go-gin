/**
 * Description：
 * FileName：auth.go
 * Author：CJiaの用心
 * Create：2025/5/13 00:29:51
 * Remark：
 */

package auth

import (
	"errors"
	config "github.com/carefuly/carefuly-admin-go-gin/config/file"
	domainSystem "github.com/carefuly/carefuly-admin-go-gin/internal/domain/careful/system"
	modelSystem "github.com/carefuly/carefuly-admin-go-gin/internal/model/careful/system"
	"github.com/carefuly/carefuly-admin-go-gin/internal/service/careful/system"
	serviceSystem "github.com/carefuly/carefuly-admin-go-gin/internal/service/careful/system"
	"github.com/carefuly/carefuly-admin-go-gin/internal/service/careful/third"
	"github.com/carefuly/carefuly-admin-go-gin/pkg/ginx/response"
	"github.com/carefuly/carefuly-admin-go-gin/pkg/utils/jwt"
	validate "github.com/carefuly/carefuly-admin-go-gin/pkg/validator"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strings"
	"time"
)

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50" example:"demo"`   // 用户名
	Password string `json:"password" binding:"required,min=6,max=20" example:"123456"` // 密码
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required" example:"demo"`   // 用户名
	Password string `json:"password" binding:"required" example:"123456"` // 密码
}

// UserTypeLoginRequest 多用户类型登录请求
type UserTypeLoginRequest struct {
	Username string `json:"username" binding:"required" example:"admin"`  // 用户名
	Password string `json:"password" binding:"required" example:"123456"` // 密码
	UserType int    `json:"userType" binding:"required" example:"1"`      // 用户类型
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token  string            `json:"token"`  // JWT令牌
	User   domainSystem.User `json:"user"`   // 用户信息
	Expire int               `json:"expire"` // 过期时间(秒)
}

// RefreshTokenRequest 刷新令牌请求
type RefreshTokenRequest struct {
	Token string `json:"token" binding:"required" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6Ikp..."` // 旧的JWT令牌
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword" binding:"required" example:"123456"`              // 旧密码
	NewPassword string `json:"newPassword" binding:"required,min=6,max=20" example:"654321"` // 新密码
}

type AuthsHandler interface {
	RegisterRoutes(router *gin.RouterGroup)
	RegisterHandler(ctx *gin.Context)
	LoginHandler(ctx *gin.Context)
	RefreshTokenHandler(ctx *gin.Context)
	LogoutHandler(ctx *gin.Context)
	GetCurrentUserHandler(ctx *gin.Context)
	ChangePasswordHandler(ctx *gin.Context)
}

type authHandler struct {
	rely       config.RelyConfig
	userSvc    serviceSystem.UserService
	captchaSvc third.CaptchaService
}

func NewRegisterHandler(rely config.RelyConfig, svc serviceSystem.UserService, captchaSvc third.CaptchaService) AuthsHandler {
	return &authHandler{
		rely:       rely,
		userSvc:    svc,
		captchaSvc: captchaSvc,
	}
}

// RegisterRoutes 注册路由
func (h *authHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/register", h.RegisterHandler)
	router.POST("/login", h.LoginHandler)
	router.POST("/refresh-token", h.RefreshTokenHandler)
	router.POST("/logout", h.LogoutHandler)
	router.GET("/userinfo", h.GetCurrentUserHandler)
	router.POST("/change-password", h.ChangePasswordHandler)
}

// RegisterHandler
// @Summary 用户注册
// @Description 用户注册
// @Tags 认证管理/用户注册
// @Accept application/json
// @Produce application/json
// @Param RegisterRequest body RegisterRequest true "请求"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /v1/auth/register [post]
func (h *authHandler) RegisterHandler(ctx *gin.Context) {
	var req RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validate.NewValidatorError(h.rely.Trans).HandleValidatorError(ctx, err)
		return
	}

	// 转换为领域模型
	user := domainSystem.User{
		User: modelSystem.User{
			Username: req.Username,
			Password: req.Password,
			DeptId:   "ADE5E818-5F3C-46F0-A366-837A8B13089E",
		},
	}

	// 调用业务逻辑
	if err := h.userSvc.Register(ctx, user); err != nil {
		switch {
		case errors.Is(err, serviceSystem.ErrUsernameDuplicate):
			response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "用户名已存在，请重新输入", nil)
			return
		default:
			ctx.Set("internal", err.Error())
			zap.L().Error("用户注册失败", zap.Error(err))
			response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
			return
		}
	}

	response.NewResponse().SuccessResponse(ctx, "注册成功", nil)
}

// LoginHandler
// @Summary 账号密码登录
// @Description 账号密码登录
// @Tags 认证管理/账号密码登录
// @Accept application/json
// @Produce application/json
// @Param LoginRequest body LoginRequest true "请求"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} response.Response
// @Router /v1/auth/login [post]
func (h *authHandler) LoginHandler(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validate.NewValidatorError(h.rely.Trans).HandleValidatorError(ctx, err)
		return
	}

	// 调用业务逻辑
	user, err := h.userSvc.Login(ctx, req.Username, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, system.ErrUserInvalidCredential):
			response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "用户名或密码错误", nil)
			return
		default:
			ctx.Set("internal", err.Error())
			zap.L().Error("登录失败", zap.Error(err))
			response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
			return
		}
	}

	// 生成JWT令牌
	token, err := jwt.GenerateToken(ctx, user.Id, user.Username, int(user.UserType), user.DeptId, h.rely.Token.Secret, h.rely.Token.Expire)
	if err != nil {
		zap.S().Errorf("生成令牌失败: %v", err)
		response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "生成令牌失败: "+err.Error(), nil)
		return
	}

	// 返回用户信息和令牌
	response.NewResponse().SuccessResponse(ctx, "登录成功", LoginResponse{
		Token:  token,
		User:   user,
		Expire: h.rely.Token.Expire * 3600,
	})
}

// RefreshTokenHandler
// @Summary 刷新令牌
// @Description 使用旧的JWT令牌获取新的令牌
// @Tags 认证管理/刷新令牌
// @Accept application/json
// @Produce application/json
// @Param RefreshTokenRequest body RefreshTokenRequest true "刷新令牌参数"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /v1/auth/refresh-token [post]
func (h *authHandler) RefreshTokenHandler(ctx *gin.Context) {
	var req RefreshTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validate.NewValidatorError(h.rely.Trans).HandleValidatorError(ctx, err)
		return
	}

	// 解析旧令牌
	claims, err := jwt.ParseToken(req.Token, h.rely.Token.Secret)
	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrExpiredToken):
			response.NewResponse().ErrorResponse(ctx, http.StatusUnauthorized, "令牌已过期，请重新登录", nil)
		case errors.Is(err, jwt.ErrInvalidToken):
			response.NewResponse().ErrorResponse(ctx, http.StatusUnauthorized, "无效的令牌", nil)
		default:
			response.NewResponse().ErrorResponse(ctx, http.StatusUnauthorized, "认证失败", nil)
		}
		return
	}

	// 获取用户信息
	user, err := h.userSvc.GetById(ctx, claims.UserId)
	if err != nil {
		if errors.Is(err, system.ErrUserNotFound) {
			response.NewResponse().ErrorResponse(ctx, http.StatusUnauthorized, "用户不存在", nil)
			return
		}
		zap.L().Error("刷新令牌获取用户信息异常", zap.Error(err))
		response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
		return
	}

	// 生成新的JWT令牌
	newToken, err := jwt.GenerateToken(ctx, user.Id, user.Username, int(user.UserType), user.DeptId, h.rely.Token.Secret, h.rely.Token.Expire)
	if err != nil {
		zap.L().Error("生成新令牌失败", zap.Error(err))
		response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
		return
	}

	// 返回新令牌和用户信息
	response.NewResponse().SuccessResponse(ctx, "刷新令牌成功", LoginResponse{
		Token:  newToken,
		User:   user,
		Expire: h.rely.Token.Expire * 3600,
	})
}

// GetCurrentUserHandler
// @Summary 获取当前登录用户信息
// @Description 获取当前登录用户的详细信息
// @Tags 认证管理/获取用户信息
// @Accept application/json
// @Produce application/json
// @Security BearerAuth
// @Success 200 {object} domainSystem.User
// @Failure 401 {object} response.Response
// @Router /v1/auth/userinfo [get]
// @Security LoginToken
func (h *authHandler) GetCurrentUserHandler(ctx *gin.Context) {
	// 从上下文中获取用户ID
	userId, exists := ctx.MustGet("userId").(string)
	if !exists {
		zap.S().Error("未找到用户认证信息", userId)
		response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
		return
	}

	// 根据id获取用户信息
	user, err := h.userSvc.GetById(ctx, userId)
	if err != nil {
		if errors.Is(err, serviceSystem.ErrUserNotFound) {
			response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "用户不存在", nil)
			return
		}
		zap.L().Error("获取用户信息异常", zap.String("userId", userId), zap.Error(err))
		response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
		return
	}

	// 返回用户信息
	response.NewResponse().SuccessResponse(ctx, "获取成功", user)
}

// LogoutHandler
// @Summary 退出登录
// @Description 用户退出登录
// @Tags 认证管理/退出登录
// @Accept application/json
// @Produce application/json
// @Security BearerAuth
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /v1/auth/logout [post]
// @Security LoginToken
func (h *authHandler) LogoutHandler(ctx *gin.Context) {
	// 从上下文中获取用户ID
	userId, exists := ctx.MustGet("userId").(string)
	if !exists {
		zap.S().Error("未找到用户认证信息", userId)
		response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
		return
	}

	// 获取Authorization头信息中的token
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		response.NewResponse().ErrorResponse(ctx, http.StatusUnauthorized, "退出登录失败：未携带令牌", nil)
		return
	}

	// 提取token
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		response.NewResponse().ErrorResponse(ctx, http.StatusUnauthorized, "退出登录失败：无效的令牌格式", nil)
		return
	}
	tokenStr := parts[1]

	// 解析token以获取过期时间
	claims, err := jwt.ParseToken(tokenStr, h.rely.Token.Secret)
	if err != nil {
		zap.S().Errorf("解析token失败: %v", err)
		response.NewResponse().ErrorResponse(ctx, http.StatusUnauthorized, "退出登录失败：无效的令牌", nil)
		return
	}

	// 计算token的剩余有效期
	expirationTime := claims.ExpiresAt.Time
	remainingTime := time.Until(expirationTime)
	if remainingTime <= 0 {
		// 如果token已过期，直接返回成功
		response.NewResponse().SuccessResponse(ctx, "令牌已过期，登出成功", nil)
		return
	}

	// 将token加入黑名单
	tokenBlacklist := jwt.NewTokenBlacklist(h.rely.Redis)
	if err := tokenBlacklist.Add(ctx, tokenStr, remainingTime); err != nil {
		zap.S().Errorf("将token加入黑名单失败: %v", err)
		response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "退出登录失败：服务器内部错误", nil)
		return
	}

	zap.S().Infof("用户登出成功, userId: %s", userId)

	// 返回成功信息
	response.NewResponse().SuccessResponse(ctx, "退出登录成功", nil)
}

// ChangePasswordHandler
// @Summary 修改密码
// @Description 修改当前登录用户的密码
// @Tags 认证管理/修改密码
// @Accept application/json
// @Produce application/json
// @Security BearerAuth
// @Param ChangePasswordRequest body ChangePasswordRequest true "修改密码参数"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /v1/auth/change-password [post]
// @Security LoginToken
func (h *authHandler) ChangePasswordHandler(ctx *gin.Context) {
	// 从上下文中获取用户ID
	userId, exists := ctx.MustGet("userId").(string)
	if !exists {
		zap.S().Error("未找到用户认证信息", userId)
		response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
		return
	}

	var req ChangePasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validate.NewValidatorError(h.rely.Trans).HandleValidatorError(ctx, err)
		return
	}

	// 调用业务逻辑修改密码
	err := h.userSvc.ChangePassword(ctx, userId, req.OldPassword, req.NewPassword)
	if err != nil {
		switch {
		case errors.Is(err, system.ErrUserInvalidCredential):
			response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "旧密码错误", nil)
		case errors.Is(err, system.ErrUserNotFound):
			response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "用户不存在", nil)
		default:
			zap.L().Error("修改密码失败", zap.Error(err))
			response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "修改密码失败", nil)
		}
		return
	}

	response.NewResponse().SuccessResponse(ctx, "密码修改成功", nil)
}
