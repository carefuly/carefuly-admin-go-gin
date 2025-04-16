/**
 * Description：
 * FileName：jwt.go
 * Author：CJiaの用心
 * Create：2025/4/15 14:13:48
 * Remark：
 */

package middleware

import (
	config "github.com/carefuly/carefuly-admin-go-gin/config/file"
	"github.com/carefuly/carefuly-admin-go-gin/internal/service/careful/system"
	"github.com/carefuly/carefuly-admin-go-gin/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
)

var (
	UnauthorizedNotFound = "Token信息为空，请先登录"
	UnauthorizedInvalid  = "无效的Token"
	UnauthorizedNotify   = "Token认证失败"
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

func (l *LoginJWTMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// swagger文档
		if strings.Contains(ctx.Request.URL.Path, "swagger") {
			return
		}
		// 不需要登录校验的
		for _, path := range l.paths {
			if ctx.Request.URL.Path == path {
				return
			}
		}

		// 我现在用 JWT 来校验
		tokenHeader := ctx.GetHeader("Authorization")
		if tokenHeader == "" {
			// 没登录
			response.NewResponse().ErrorResponse(ctx, http.StatusUnauthorized, UnauthorizedNotFound, nil)
			ctx.Abort()
			return
		}
		seg := strings.Split(tokenHeader, " ")
		if len(seg) != 2 {
			// 没登录，有人瞎搞
			response.NewResponse().ErrorResponse(ctx, http.StatusUnauthorized, UnauthorizedInvalid, nil)
			ctx.Abort()
			return
		}

		tokenStr := seg[1]
		claims := &system.UserClaims{}
		// ParseWithClaims 里面，一定要传入指针
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(l.rely.Token.ApiKeyAuth), nil
		})
		if err != nil {
			// 没登录
			response.NewResponse().ErrorResponse(ctx, http.StatusUnauthorized, UnauthorizedNotify, nil)
			ctx.Abort()
			return
		}

		// err 为 nil，token 不为 nil
		if token == nil || !token.Valid {
			// 没登录
			response.NewResponse().ErrorResponse(ctx, http.StatusUnauthorized, UnauthorizedNotify, nil)
			ctx.Abort()
			return
		}

		// 校验用户代理
		// if claims.UserAgent != ctx.Request.UserAgent() {
		// 	// 严重的安全问题
		// 	// 你是要监控
		// 	ctx.AbortWithStatus(http.StatusUnauthorized)
		// 	return
		// }

		ctx.Set("claims", claims)
		ctx.Set("userId", claims.UId)
		ctx.Set("username", claims.Username)
	}
}
