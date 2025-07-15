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
	serviceThird "github.com/carefuly/carefuly-admin-go-gin/internal/service/careful/third"
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

	// 用户
	userCache := cacheSystem.NewRedisUserCache(r.rely.Redis)
	userDAO := daoSystem.NewGORMUserDAO(r.rely.Db.Careful)
	userRepository := repositorySystem.NewUserRepository(userDAO, userCache)
	userService := serviceSystem.NewUserService(userRepository)

	// 数据字典
	dictCache := cacheTools.NewRedisDictCache(r.rely.Redis)
	dictDAO := daoTools.NewGORMDictDAO(r.rely.Db.Careful)
	dictRepository := repositoryTools.NewDictRepository(dictDAO, dictCache)
	dictService := serviceTools.NewDictService(dictRepository)
	dictHandler := handlerTools.NewDictHandler(r.rely, dictService, userService)
	dictHandler.RegisterRoutes(baseRouter)

	// 字典项
	dictTypeCache := cacheTools.NewRedisDictTypeCache(r.rely.Redis)
	dictTypeDAO := daoTools.NewGORMDictTypeDAO(r.rely.Db.Careful)
	dictTypeRepository := repositoryTools.NewDictTypeRepository(dictTypeDAO, dictTypeCache)
	dictTypeService := serviceTools.NewDictTypeService(dictTypeRepository, dictRepository)
	dictTypeHandler := handlerTools.NewDictTypeHandler(r.rely, dictTypeService, userService)
	dictTypeHandler.RegisterRoutes(baseRouter)

	// 文件
	fileService := serviceThird.NewBucketFileService()

	// 存储桶
	bucketCache := cacheTools.NewRedisBucketCache(r.rely.Redis)
	bucketDAO := daoTools.NewGORMBucketDAO(r.rely.Db.Careful)
	bucketRepository := repositoryTools.NewBucketRepository(bucketDAO, bucketCache)
	bucketService := serviceTools.NewBucketService(bucketRepository)
	bucketHandler := handlerTools.NewBucketHandler(r.rely, bucketService, userService, fileService)
	bucketHandler.RegisterRoutes(baseRouter)
}
