/**
 * Description：
 * FileName：server.go
 * Author：CJiaの用心
 * Create：2025/3/20 23:17:03
 * Remark：
 */

package ioc

import (
	"fmt"
	config "github.com/carefuly/carefuly-admin-go-gin/config/file"
	"github.com/carefuly/carefuly-admin-go-gin/docs"
	"github.com/carefuly/carefuly-admin-go-gin/internal/web/middleware"
	"github.com/carefuly/carefuly-admin-go-gin/internal/web/router/careful"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	entranslations "github.com/go-playground/validator/v10/translations/en"
	zhtranslations "github.com/go-playground/validator/v10/translations/zh"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"reflect"
	"strings"
)

type Server struct {
	rely   config.RelyConfig
	locale string
}

func NewServer(rely config.RelyConfig, locale string) *Server {
	return &Server{
		rely:   rely,
		locale: locale,
	}
}

func (s *Server) InitGinMiddlewares(rely config.RelyConfig) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		middleware.CORSMiddleware(),
		middleware.NewLoginJWTMiddlewareBuilder(rely).
			IgnorePaths("/dev-api/v1/third/generateCaptcha").
			IgnorePaths("/dev-api/v1/auth/register").
			IgnorePaths("/dev-api/v1/auth/login").
			IgnorePaths("/dev-api/v1/auth/type-login").
			IgnorePaths("/dev-api/v1/auth/refresh-token").
			Build(),
		middleware.NewLogger(rely.Logger).Logger(),
		middleware.NewStorage().StorageLogger(rely.Db.Careful),
	}
}

func (s *Server) InitGinTrans() (ut.Translator, error) {
	var trans ut.Translator
	// 修改gin框架中的validator引擎属性, 实现定制
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// 注册一个获取json的tag的自定义方法
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})

		zhT := zh.New() // 中文翻译器
		enT := en.New() // 英文翻译器
		// 第一个参数是备用的语言环境，后面的参数是应该支持的语言环境
		uni := ut.New(enT, zhT, enT)
		trans, ok = uni.GetTranslator(s.locale)
		if !ok {
			return trans, fmt.Errorf("uni.GetTranslator(%s)", s.locale)
		}

		switch s.locale {
		case "en":
			err := entranslations.RegisterDefaultTranslations(v, trans)
			if err != nil {
				return nil, err
			}
		case "zh":
			err := zhtranslations.RegisterDefaultTranslations(v, trans)
			if err != nil {
				return nil, err
			}
		default:
			err := entranslations.RegisterDefaultTranslations(v, trans)
			if err != nil {
				return nil, err
			}
		}
		return trans, nil
	}
	return trans, nil
}

func (s *Server) InitWebServer(middle []gin.HandlerFunc, rely config.RelyConfig) *gin.Engine {
	server := gin.Default()
	server.Use(middle...)

	// 配置接口前缀
	docs.SwaggerInfo.BasePath = "/dev-api"
	// 配置接口文档
	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	ApiGroup := server.Group("/dev-api")
	v1 := ApiGroup.Group("/v1")

	careful.NewRouter(rely, v1).RegisterRoutes()

	return server
}
