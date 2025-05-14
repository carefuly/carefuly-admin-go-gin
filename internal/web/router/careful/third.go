/**
 * Description：
 * FileName：third.go
 * Author：CJiaの用心
 * Create：2025/5/13 15:27:35
 * Remark：
 */

package careful

import (
	config "github.com/carefuly/carefuly-admin-go-gin/config/file"
	"github.com/carefuly/carefuly-admin-go-gin/internal/repository/cache/careful/third"
	repositoryThird "github.com/carefuly/carefuly-admin-go-gin/internal/repository/repository/careful/third"
	serviceThird "github.com/carefuly/carefuly-admin-go-gin/internal/service/careful/third"
	handlerThird "github.com/carefuly/carefuly-admin-go-gin/internal/web/handler/careful/third"
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

func (r *ThirdRouter) RegisterRouter(router *gin.RouterGroup) {
	baseRouter := router.Group("/third")

	captchaCache := third.NewCaptchaCache(r.rely.Redis)
	captchaRepository := repositoryThird.NewCaptchaRepository(captchaCache)
	captchaService := serviceThird.NewCaptchaService(captchaRepository)
	captchaHandler := handlerThird.NewCaptchaController(r.rely, captchaService)
	captchaHandler.RegisterRoutes(baseRouter)
}
