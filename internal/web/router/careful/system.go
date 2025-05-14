/**
 * Description：
 * FileName：system.go
 * Author：CJiaの用心
 * Create：2025/5/13 16:46:58
 * Remark：
 */

package careful

import (
	config "github.com/carefuly/carefuly-admin-go-gin/config/file"
	"github.com/carefuly/carefuly-admin-go-gin/internal/repository/dao/careful/system"
	repositorySystem "github.com/carefuly/carefuly-admin-go-gin/internal/repository/repository/careful/system"
	serviceSystem "github.com/carefuly/carefuly-admin-go-gin/internal/service/careful/system"
	handlerSystem "github.com/carefuly/carefuly-admin-go-gin/internal/web/handler/careful/system"
	"github.com/gin-gonic/gin"
)

type SystemRouter struct {
	rely config.RelyConfig
}

func NewSystemRouter(rely config.RelyConfig) *SystemRouter {
	return &SystemRouter{
		rely: rely,
	}
}

func (r *SystemRouter) RegisterRouter(router *gin.RouterGroup) {
	baseRouter := router.Group("/system")

	menuDAO := system.NewGORMMenuDAO(r.rely.Db.Careful)
	menuRepository := repositorySystem.NewMenuRepository(menuDAO)
	menuService := serviceSystem.NewMenuService(menuRepository)
	menuHandler := handlerSystem.NewMenuHandler(r.rely, menuService)
	menuHandler.RegisterRoutes(baseRouter)
}
