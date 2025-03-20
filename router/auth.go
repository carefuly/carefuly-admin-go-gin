/**
 * Description：
 * FileName：auth.go
 * Author：CJiaの用心
 * Create：2025/3/20 23:23:11
 * Remark：
 */

package router

import (
	config "github.com/carefuly/carefuly-admin-go-gin/config/file"
	controller "github.com/carefuly/carefuly-admin-go-gin/internal/web/controller/auth"
	"github.com/gin-gonic/gin"
)

type AuthRouter struct {
	rely config.RelyConfig
}

func NewAuthRouter(rely config.RelyConfig) *AuthRouter {
	return &AuthRouter{
		rely: rely,
	}
}

func (r *AuthRouter) RegisterAuthRouter(router *gin.RouterGroup) {
	authRouter := router.Group("/auth")

	registerHandler := controller.NewRegisterController(r.rely)
	registerHandler.RegisterRoutes(authRouter)
}
