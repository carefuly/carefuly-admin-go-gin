/**
 * Description：
 * FileName：login.go
 * Author：CJiaの用心
 * Create：2025/3/25 11:03:56
 * Remark：
 */

package controller

import (
	config "github.com/carefuly/carefuly-admin-go-gin/config/file"
	"github.com/gin-gonic/gin"
)

type LoginController interface {
	RegisterRoutes(router *gin.RouterGroup)
}

type loginController struct {
	rely config.RelyConfig
}

func (c *loginController) NewLoginController(rely config.RelyConfig) LoginController {
	return &loginController{
		rely: rely,
	}
}

func (c *loginController) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/generateDigitCaptcha", c.GenerateDigitCaptcha)
}

func (c *loginController) GenerateDigitCaptcha(ctx *gin.Context) {

}
