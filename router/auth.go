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
	dao "github.com/carefuly/carefuly-admin-go-gin/internal/dao/system"
	repository "github.com/carefuly/carefuly-admin-go-gin/internal/repository/system"
	service "github.com/carefuly/carefuly-admin-go-gin/internal/service/auth"
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

	userDAO := dao.NewGORMUserDAO(r.rely.Db.Careful)
	userRepository := repository.NewUserRepository(userDAO)
	userPasswordDAO := dao.NewUserPasswordDAO(r.rely.Db.Careful)
	userPasswordRepository := repository.NewUserPassWordRepository(userPasswordDAO)

	registerService := service.NewRegisterService(userRepository, userPasswordRepository)
	registerHandler := controller.NewRegisterController(r.rely, registerService)
	registerHandler.RegisterRoutes(authRouter)
}
