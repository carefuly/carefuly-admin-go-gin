/**
 * Description：
 * FileName：types.go
 * Author：CJiaの用心
 * Create：2025/5/13 01:00:06
 * Remark：
 */

package careful

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
	NewAuthRouter(r.rely).RegisterRouter(r.router)
	NewSystemRouter(r.rely).RegisterRouter(r.router)
	NewToolsRouter(r.rely).RegisterRouter(r.router)
	NewThirdRouter(r.rely).RegisterRouter(r.router)
}
