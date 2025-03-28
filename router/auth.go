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
	"github.com/carefuly/carefuly-admin-go-gin/v1/internal/cache/third"
	dao2 "github.com/carefuly/carefuly-admin-go-gin/v1/internal/dao/system"
	repository2 "github.com/carefuly/carefuly-admin-go-gin/v1/internal/repository/system"
	repositoryThird "github.com/carefuly/carefuly-admin-go-gin/v1/internal/repository/third"
	"github.com/carefuly/carefuly-admin-go-gin/v1/internal/service/auth"
	serviceSystem "github.com/carefuly/carefuly-admin-go-gin/v1/internal/service/system"
	serviceThird "github.com/carefuly/carefuly-admin-go-gin/v1/internal/service/third"
	controller2 "github.com/carefuly/carefuly-admin-go-gin/v1/internal/web/controller/auth"
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

	userDAO := dao2.NewGORMUserDAO(r.rely.Db.Careful)
	userRepository := repository2.NewUserRepository(userDAO)
	userPasswordDAO := dao2.NewUserPasswordDAO(r.rely.Db.Careful)
	userPasswordRepository := repository2.NewUserPassWordRepository(userPasswordDAO)
	userService := serviceSystem.NewUserService(userRepository)

	registerService := service.NewRegisterService(userRepository, userPasswordRepository)
	registerHandler := controller2.NewRegisterController(r.rely, registerService)
	registerHandler.RegisterRoutes(authRouter)

	loginHandler := controller2.NewLoginController(r.rely, userService, captchaService)
	loginHandler.RegisterRoutes(authRouter)
}
