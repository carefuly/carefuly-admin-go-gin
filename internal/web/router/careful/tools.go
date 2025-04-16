/**
 * Description：
 * FileName：tools.go
 * Author：CJiaの用心
 * Create：2025/4/15 15:01:42
 * Remark：
 */

package careful

import (
	config "github.com/carefuly/carefuly-admin-go-gin/config/file"
	cacheTools "github.com/carefuly/carefuly-admin-go-gin/internal/cache/careful/tools"
	daoTools "github.com/carefuly/carefuly-admin-go-gin/internal/dao/careful/tools"
	repositoryTools "github.com/carefuly/carefuly-admin-go-gin/internal/repository/careful/tools"
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

	redisDictCache := cacheTools.NewRedisDictCache(r.rely.Redis)
	dictDAO := daoTools.NewDictDao(r.rely.Db.Careful)
	dictRepository := repositoryTools.NewDictRepository(dictDAO, redisDictCache)
	dictService := serviceTools.NewDictService(dictRepository)
	dictHandler := handlerTools.NewDictHandler(r.rely, dictService)
	dictHandler.RegisterRoutes(baseRouter)
}
