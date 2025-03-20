/**
 * Description：
 * FileName：register.go
 * Author：CJiaの用心
 * Create：2025/3/20 23:26:20
 * Remark：
 */

package controller

import (
	config "github.com/carefuly/carefuly-admin-go-gin/config/file"
	"github.com/gin-gonic/gin"
)

type RegisterController interface {
	RegisterRoutes(router *gin.RouterGroup)
	PassWordRegisterHandler(ctx *gin.Context)
}

type registerController struct {
	rely config.RelyConfig
}

func NewRegisterController(rely config.RelyConfig) RegisterController {
	return &registerController{
		rely: rely,
	}
}

func (c *registerController) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/send-register-captcha", c.SendEmailCaptchaRegisterHandler)
	router.POST("/password-register", c.PassWordRegisterHandler)
}

func (c *registerController) SendEmailCaptchaRegisterHandler(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "success",
	})
}

func (c *registerController) PassWordRegisterHandler(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "success",
	})
}
