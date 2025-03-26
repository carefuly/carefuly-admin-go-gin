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
	cache "github.com/carefuly/carefuly-admin-go-gin/internal/cache/third"
	dao "github.com/carefuly/carefuly-admin-go-gin/internal/dao/system"
	repository "github.com/carefuly/carefuly-admin-go-gin/internal/repository/system"
	repositoryThird "github.com/carefuly/carefuly-admin-go-gin/internal/repository/third"
	service "github.com/carefuly/carefuly-admin-go-gin/internal/service/auth"
	serviceSystem "github.com/carefuly/carefuly-admin-go-gin/internal/service/system"
	serviceThird "github.com/carefuly/carefuly-admin-go-gin/internal/service/third"
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

	captchaCache := cache.NewCaptchaCache(r.rely.Redis)
	captchaRepository := repositoryThird.NewCaptchaRepository(captchaCache)
	captchaService := serviceThird.NewCaptchaService(captchaRepository)

	userDAO := dao.NewGORMUserDAO(r.rely.Db.Careful)
	userRepository := repository.NewUserRepository(userDAO)
	userPasswordDAO := dao.NewUserPasswordDAO(r.rely.Db.Careful)
	userPasswordRepository := repository.NewUserPassWordRepository(userPasswordDAO)
	userService := serviceSystem.NewUserService(userRepository)

	registerService := service.NewRegisterService(userRepository, userPasswordRepository)
	registerHandler := controller.NewRegisterController(r.rely, registerService)
	registerHandler.RegisterRoutes(authRouter)

	loginHandler := controller.NewLoginController(r.rely, userService, captchaService)
	loginHandler.RegisterRoutes(authRouter)
}
