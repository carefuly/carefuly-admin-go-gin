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
	cacheSystem "github.com/carefuly/carefuly-admin-go-gin/internal/repository/cache/careful/system"
	daoSystem "github.com/carefuly/carefuly-admin-go-gin/internal/repository/dao/careful/system"
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

	// 用户
	userCache := cacheSystem.NewRedisUserCache(r.rely.Redis)
	userDAO := daoSystem.NewGORMUserDAO(r.rely.Db.Careful)
	userRepository := repositorySystem.NewUserRepository(userDAO, userCache)
	userService := serviceSystem.NewUserService(userRepository)
	userHandler := handlerSystem.NewUserHandler(r.rely, userService)
	userHandler.RegisterRoutes(baseRouter)

	// 菜单
	menuCache := cacheSystem.NewRedisMenuCache(r.rely.Redis)
	menuDAO := daoSystem.NewGORMMenuDAO(r.rely.Db.Careful)
	menuRepository := repositorySystem.NewMenuRepository(menuDAO, menuCache)
	menuService := serviceSystem.NewMenuService(menuRepository)
	menuHandler := handlerSystem.NewMenuHandler(r.rely, menuService, userService)
	menuHandler.RegisterRoutes(baseRouter)

	// 菜单权限
	menuButtonCache := cacheSystem.NewRedisMenuButtonCache(r.rely.Redis)
	menuButtonDAO := daoSystem.NewGORMMenuButtonDAO(r.rely.Db.Careful)
	menuButtonRepository := repositorySystem.NewMenuButtonRepository(menuButtonDAO, menuButtonCache)
	menuButtonService := serviceSystem.NewMenuButtonService(menuButtonRepository, menuRepository)
	menuButtonHandler := handlerSystem.NewMenuButtonHandler(r.rely, menuButtonService, userService)
	menuButtonHandler.RegisterRoutes(baseRouter)

	// 菜单数据列
	menuColumnCache := cacheSystem.NewRedisMenuColumnCache(r.rely.Redis)
	menuColumnDAO := daoSystem.NewGORMMenuColumnDAO(r.rely.Db.Careful)
	menuColumnRepository := repositorySystem.NewMenuColumnRepository(menuColumnDAO, menuColumnCache)
	menuColumnService := serviceSystem.NewMenuColumnService(menuColumnRepository, menuRepository)
	menuColumnHandler := handlerSystem.NewMenuColumnHandler(r.rely, menuColumnService, userService)
	menuColumnHandler.RegisterRoutes(baseRouter)

	// 部门
	deptCache := cacheSystem.NewRedisDeptCache(r.rely.Redis)
	deptDAO := daoSystem.NewGORMDeptDAO(r.rely.Db.Careful)
	deptRepository := repositorySystem.NewDeptRepository(deptDAO, deptCache)
	deptService := serviceSystem.NewDeptService(deptRepository)
	deptHandler := handlerSystem.NewDeptHandler(r.rely, deptService, userService)
	deptHandler.RegisterRoutes(baseRouter)

	// 角色
	roleCache := cacheSystem.NewRedisRoleCache(r.rely.Redis)
	roleDAO := daoSystem.NewGORMRoleDAO(r.rely.Db.Careful, deptDAO, menuDAO, menuButtonDAO, menuColumnDAO)
	roleRepository := repositorySystem.NewRoleRepository(roleDAO, roleCache)
	roleService := serviceSystem.NewRoleService(roleRepository)
	roleHandler := handlerSystem.NewRoleHandler(r.rely, roleService, userService)
	roleHandler.RegisterRoutes(baseRouter)

	// 岗位
	postCache := cacheSystem.NewRedisPostCache(r.rely.Redis)
	postDAO := daoSystem.NewGORMPostDAO(r.rely.Db.Careful)
	postRepository := repositorySystem.NewPostRepository(postDAO, postCache)
	postService := serviceSystem.NewPostService(postRepository)
	postHandler := handlerSystem.NewPostHandler(r.rely, postService, userService)
	postHandler.RegisterRoutes(baseRouter)
}
