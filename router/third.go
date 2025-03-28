/**
 * Description：
 * FileName：third.go
 * Author：CJiaの用心
 * Create：2025/3/25 11:44:06
 * Remark：
 */

package router

import (
	config "github.com/carefuly/carefuly-admin-go-gin/config/file"
	"github.com/carefuly/carefuly-admin-go-gin/v1/internal/cache/third"
	"github.com/carefuly/carefuly-admin-go-gin/v1/internal/repository/third"
	"github.com/carefuly/carefuly-admin-go-gin/v1/internal/service/third"
	"github.com/carefuly/carefuly-admin-go-gin/v1/internal/web/controller/third"
	"github.com/gin-gonic/gin"
)

type ThirdRouter struct {
	rely config.RelyConfig
}

func NewThirdRouter(rely config.RelyConfig) *ThirdRouter {
	return &ThirdRouter{
		rely: rely,
	}
}

func (r *ThirdRouter) RegisterAuthRouter(router *gin.RouterGroup) {
	thirdRouter := router.Group("/third")

	captchaCache := cache.NewCaptchaCache(r.rely.Redis)
	captchaRepository := repository.NewCaptchaRepository(captchaCache)
	captchaService := service.NewCaptchaService(captchaRepository)

	captchaHandler := controller.NewCaptchaController(r.rely, captchaService)
	captchaHandler.RegisterRoutes(thirdRouter)
}
