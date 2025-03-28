/**
 * Description：
 * FileName：third.go
 * Author：CJiaの用心
 * Create：2025/3/27 11:54:38
 * Remark：
 */

package careful

import (
	config "github.com/carefuly/carefuly-admin-go-gin/config/file"
	"github.com/carefuly/carefuly-admin-go-gin/internal/cache/careful/third"
	thirdRepository "github.com/carefuly/carefuly-admin-go-gin/internal/repository/careful/third"
	thirdService "github.com/carefuly/carefuly-admin-go-gin/internal/service/careful/third"
	thirdHandler "github.com/carefuly/carefuly-admin-go-gin/internal/web/handler/careful/third"
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
	captchaRepository := thirdRepository.NewCaptchaRepository(captchaCache)
	captchaService := thirdService.NewCaptchaService(captchaRepository)
	captchaHandler := thirdHandler.NewCaptchaController(r.rely, captchaService)
	captchaHandler.RegisterRoutes(baseRouter)
}
