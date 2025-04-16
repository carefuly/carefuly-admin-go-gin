/**
 * Description：
 * FileName：types.go
 * Author：CJiaの用心
 * Create：2025/3/27 13:58:06
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

	NewToolsRouter(r.rely).RegisterRouter(r.router)

	NewThirdRouter(r.rely).RegisterRouter(r.router)
}
