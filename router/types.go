/**
 * Description：
 * FileName：types.go
 * Author：CJiaの用心
 * Create：2025/3/20 23:22:58
 * Remark：
 */

package router

import (
	config "github.com/carefuly/carefuly-admin-go-gin/config/file"
	"github.com/gin-gonic/gin"
)

type Router struct {
	rely   config.RelyConfig
	router *gin.RouterGroup
}

func NewRouter(rely config.RelyConfig, router *gin.RouterGroup) *Router {
	return &Router{
		rely:   rely,
		router: router,
	}
}

func (r *Router) RegisterRoutes() {
	NewAuthRouter(r.rely).RegisterAuthRouter(r.router)
}
