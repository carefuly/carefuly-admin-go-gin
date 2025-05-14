/**
 * Description：
 * FileName：auth.go
 * Author：CJiaの用心
 * Create：2025/5/13 00:55:07
 * Remark：
 */

package careful

import (
	config "github.com/carefuly/carefuly-admin-go-gin/config/file"
	cacheSystem "github.com/carefuly/carefuly-admin-go-gin/internal/repository/cache/careful/system"
	"github.com/carefuly/carefuly-admin-go-gin/internal/repository/cache/careful/third"
	"github.com/carefuly/carefuly-admin-go-gin/internal/repository/dao/careful/system"
	repositorySystem "github.com/carefuly/carefuly-admin-go-gin/internal/repository/repository/careful/system"
	repositoryThird "github.com/carefuly/carefuly-admin-go-gin/internal/repository/repository/careful/third"
	serviceSystem "github.com/carefuly/carefuly-admin-go-gin/internal/service/careful/system"
	serviceThird "github.com/carefuly/carefuly-admin-go-gin/internal/service/careful/third"
	"github.com/carefuly/carefuly-admin-go-gin/internal/web/handler/careful/auth"
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

func (r *AuthRouter) RegisterRouter(router *gin.RouterGroup) {
	baseRouter := router.Group("/auth")

	captchaCache := third.NewCaptchaCache(r.rely.Redis)
	captchaRepository := repositoryThird.NewCaptchaRepository(captchaCache)
	captchaService := serviceThird.NewCaptchaService(captchaRepository)

	userCache := cacheSystem.NewRedisUserCache(r.rely.Redis)

	userDAO := system.NewGORMUserDAO(r.rely.Db.Careful)
	userRepository := repositorySystem.NewUserRepository(userDAO, userCache)
	userService := serviceSystem.NewUserService(userRepository)

	registerHandler := auth.NewRegisterHandler(r.rely, userService, captchaService)
	registerHandler.RegisterRoutes(baseRouter)
}
