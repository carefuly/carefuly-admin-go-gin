/**
 * Description：
 * FileName：system.go
 * Author：CJiaの用心
 * Create：2025/3/28 11:57:01
 * Remark：
 */

package careful

import (
	config "github.com/carefuly/carefuly-admin-go-gin/config/file"
)

type SystemRouter struct {
	rely config.RelyConfig
}

func NewSystemRouter(rely config.RelyConfig) *SystemRouter {
	return &SystemRouter{
		rely: rely,
	}
}

// func (r *SystemRouter) RegisterRouter(router *gin.RouterGroup) {
//
// 	baseRouter := router.Group("/third")
//
// 	captchaCache := third.NewCaptchaCache(r.rely.Redis)
// 	captchaRepository := thirdRepository.NewCaptchaRepository(captchaCache)
// 	captchaService := thirdService.NewCaptchaService(captchaRepository)
// 	captchaHandler := thirdHandler.NewCaptchaController(r.rely, captchaService)
// 	captchaHandler.RegisterRoutes(baseRouter)
// }
