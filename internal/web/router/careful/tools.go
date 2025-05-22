/**
 * Description：
 * FileName：tools.go
 * Author：CJiaの用心
 * Create：2025/5/22 11:33:34
 * Remark：
 */

package careful

import (
	config "github.com/carefuly/carefuly-admin-go-gin/config/file"
	cacheSystem "github.com/carefuly/carefuly-admin-go-gin/internal/repository/cache/careful/system"
	cacheTools "github.com/carefuly/carefuly-admin-go-gin/internal/repository/cache/careful/tools"
	daoSystem "github.com/carefuly/carefuly-admin-go-gin/internal/repository/dao/careful/system"
	daoTools "github.com/carefuly/carefuly-admin-go-gin/internal/repository/dao/careful/tools"
	repositorySystem "github.com/carefuly/carefuly-admin-go-gin/internal/repository/repository/careful/system"
	repositoryTools "github.com/carefuly/carefuly-admin-go-gin/internal/repository/repository/careful/tools"
	serviceSystem "github.com/carefuly/carefuly-admin-go-gin/internal/service/careful/system"
	serviceTools "github.com/carefuly/carefuly-admin-go-gin/internal/service/careful/tools"
	handlerTools "github.com/carefuly/carefuly-admin-go-gin/internal/web/handler/careful/tools"
	"github.com/gin-gonic/gin"
)

type ToolsRouter struct {
	rely config.RelyConfig
}

func NewToolsRouter(rely config.RelyConfig) *ToolsRouter {
	return &ToolsRouter{
		rely: rely,
	}
}

func (r *ToolsRouter) RegisterRouter(router *gin.RouterGroup) {
	baseRouter := router.Group("/tools")

	userCache := cacheSystem.NewRedisUserCache(r.rely.Redis)
	userDAO := daoSystem.NewGORMUserDAO(r.rely.Db.Careful)
	userRepository := repositorySystem.NewUserRepository(userDAO, userCache)
	userService := serviceSystem.NewUserService(userRepository)

	dictCache := cacheTools.NewRedisDictCache(r.rely.Redis)
	dictDAO := daoTools.NewGORMDictDAO(r.rely.Db.Careful)
	dictRepository := repositoryTools.NewDictRepository(dictDAO, dictCache)
	dictService := serviceTools.NewDictService(dictRepository)
	dictHandler := handlerTools.NewDictHandler(r.rely, dictService, userService)
	dictHandler.RegisterRoutes(baseRouter)
}
