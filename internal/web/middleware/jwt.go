/**
 * Description：
 * FileName：jwt.go
 * Author：CJiaの用心
 * Create：2025/5/13 01:04:02
 * Remark：
 */

package middleware

import (
	"errors"
	config "github.com/carefuly/carefuly-admin-go-gin/config/file"
	"github.com/carefuly/carefuly-admin-go-gin/pkg/ginx/response"
	"github.com/carefuly/carefuly-admin-go-gin/pkg/utils/jwt"
	"github.com/carefuly/carefuly-admin-go-gin/pkg/utils/requestUtils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

var (
	UnauthorizedNotFound = "请求未携带token，无权限访问"
	UnauthorizedInvalid  = "无效的Token"
)

// LoginJWTMiddlewareBuilder JWT 登录校验
type LoginJWTMiddlewareBuilder struct {
	paths []string
	rely  config.RelyConfig
}

func NewLoginJWTMiddlewareBuilder(rely config.RelyConfig) *LoginJWTMiddlewareBuilder {
	return &LoginJWTMiddlewareBuilder{
		rely: rely,
	}
}

func (l *LoginJWTMiddlewareBuilder) IgnorePaths(path string) *LoginJWTMiddlewareBuilder {
	l.paths = append(l.paths, path)
	return l
}

// FailedWithStatus 响应失败并设置HTTP状态码
func (l *LoginJWTMiddlewareBuilder) FailedWithStatus(ctx *gin.Context, httpStatus, code int, msg string) {
	ctx.JSON(httpStatus, gin.H{
		"code":    code,
		"message": msg,
		"data":    nil,
	})
}

// Unauthorized 响应未授权
func (l *LoginJWTMiddlewareBuilder) Unauthorized(ctx *gin.Context, msg string) {
	if msg == "" {
		msg = "未授权，请先登录"
	}
	l.FailedWithStatus(ctx, http.StatusUnauthorized, http.StatusUnauthorized, msg)
}

// Build JWT认证中间件
func (l *LoginJWTMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// swagger文档
		if l.containsAnySubstring(ctx.Request.URL.Path, []string{"swagger", "static"}) {
			return
		}
		// 不需要登录校验的
		for _, path := range l.paths {
			if ctx.Request.URL.Path == path {
				return
			}
		}

		// 获取Authorization头
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			response.NewResponse().ErrorResponse(ctx, http.StatusUnauthorized, UnauthorizedNotFound, nil)
			ctx.Abort()
			return
		}
		seg := strings.Split(authHeader, " ")
		if len(seg) != 2 {
			// 没登录，有人瞎搞
			response.NewResponse().ErrorResponse(ctx, http.StatusUnauthorized, UnauthorizedInvalid, nil)
			ctx.Abort()
			return
		}

		tokenStr := seg[1]

		// 检查token是否在黑名单中
		tokenBlacklist := jwt.NewTokenBlacklist(l.rely.Redis)
		blacklisted, err := tokenBlacklist.IsBlacklisted(ctx, tokenStr)
		if err != nil {
			zap.L().Error("检查token黑名单失败", zap.Error(err))
			response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器内部错误", nil)
			ctx.Abort()
			return
		}

		if blacklisted {
			response.NewResponse().ErrorResponse(ctx, http.StatusUnauthorized, "Token已失效，请重新登录", nil)
			ctx.Abort()
			return
		}

		// 解析token
		claims, err := jwt.ParseToken(tokenStr, l.rely.Token.Secret)
		if err != nil {
			switch {
			case errors.Is(err, jwt.ErrExpiredToken):
				response.NewResponse().ErrorResponse(ctx, http.StatusUnauthorized, "Token已过期", nil)
			case errors.Is(err, jwt.ErrInvalidToken):
				response.NewResponse().ErrorResponse(ctx, http.StatusUnauthorized, "Token无效", nil)
			default:
				response.NewResponse().ErrorResponse(ctx, http.StatusUnauthorized, "Token认证失败", nil)
			}
			ctx.Abort()
			return
		}

		// gin.Context.Set() 方法将数据存储到上下文，可以在后续的中间件或处理程序中访问。
		// 通过 gin.Context.Get() 方法获取存储在上下文中的数据。
		// 通过 gin.Context.Set() 方法存储数据时，需要指定一个键，以便在后续的中间件或处理程序中访问该数据。
		// 通过 gin.Context.Get() 方法获取数据时，需要指定相同的键。
		ctx.Set("requestIp", requestUtils.NormalizeIP(ctx))
		ctx.Set("request", ctx.Request)

		// 将用户信息存储到上下文

		ctx.Set("claims", claims)
		ctx.Set("userId", claims.UserId)
		ctx.Set("username", claims.Username)
		ctx.Set("userType", claims.UserType)
		ctx.Set("deptId", claims.DeptId)
	}
}

// JWTAuthMiddleware JWT认证中间件
func (l *LoginJWTMiddlewareBuilder) JWTAuthMiddleware(tokenConfig config.TokenConfig) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 获取Authorization头
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			l.Unauthorized(ctx, "请求未携带token，无权限访问")
			ctx.Abort()
			return
		}

		// 检查Bearer前缀
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			l.Unauthorized(ctx, "token格式错误")
			ctx.Abort()
			return
		}

		// 解析token
		claims, err := jwt.ParseToken(parts[1], tokenConfig.Secret)
		if err != nil {
			switch {
			case errors.Is(err, jwt.ErrExpiredToken):
				l.Unauthorized(ctx, "token已过期")
			case errors.Is(err, jwt.ErrInvalidToken):
				l.Unauthorized(ctx, "token无效")
			default:
				l.Unauthorized(ctx, "认证失败")
			}
			ctx.Abort()
			return
		}

		// 将用户信息存储到上下文
		ctx.Set("userId", claims.UserId)
		ctx.Set("username", claims.Username)

		ctx.Next()
	}
}

// containsAnySubstring 检查字符串是否包含切片中的任意一个子串
func (l *LoginJWTMiddlewareBuilder) containsAnySubstring(str string, subs []string) bool {
	for _, sub := range subs {
		if strings.Contains(str, sub) {
			return true
		}
	}
	return false
}
