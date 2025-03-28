/**
 * Description：
 * FileName：auth.go
 * Author：CJiaの用心
 * Create：2025/3/28 11:57:48
 * Remark：
 */

package careful

import (
	config "github.com/carefuly/carefuly-admin-go-gin/config/file"
	"github.com/carefuly/carefuly-admin-go-gin/internal/cache/careful/third"
	"github.com/carefuly/carefuly-admin-go-gin/internal/dao/careful/system"
	systemRepository "github.com/carefuly/carefuly-admin-go-gin/internal/repository/careful/system"
	thirdRepository "github.com/carefuly/carefuly-admin-go-gin/internal/repository/careful/third"
	sysetmService "github.com/carefuly/carefuly-admin-go-gin/internal/service/careful/system"
	thirdService "github.com/carefuly/carefuly-admin-go-gin/internal/service/careful/third"
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
	captchaRepository := thirdRepository.NewCaptchaRepository(captchaCache)
	captchaService := thirdService.NewCaptchaService(captchaRepository)

	userDAO := system.NewGORMUserDAO(r.rely.Db.Careful)
	userRepository := systemRepository.NewUserRepository(userDAO)
	userPasswordDAO := system.NewUserPasswordDAO(r.rely.Db.Careful)
	userPasswordRepository := systemRepository.NewUserPassWordRepository(userPasswordDAO)
	userService := sysetmService.NewUserService(userRepository, userPasswordRepository)

	registerHandler := auth.NewRegisterHandler(r.rely, userService, captchaService)
	registerHandler.RegisterRoutes(baseRouter)
}
