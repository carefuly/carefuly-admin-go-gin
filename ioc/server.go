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
	"github.com/carefuly/carefuly-admin-go-gin/router"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	entranslations "github.com/go-playground/validator/v10/translations/en"
	zhtranslations "github.com/go-playground/validator/v10/translations/zh"
	"reflect"
	"strings"
)

type Server struct {
	rely   config.RelyConfig
	locale string
}

func NewServer(rely config.RelyConfig) *Server {
	return &Server{
		rely: rely,
	}
}

func (s *Server) InitGinMiddlewares() []gin.HandlerFunc {
	return []gin.HandlerFunc{}
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

func (s *Server) InitWebServer(middle []gin.HandlerFunc) *gin.Engine {
	server := gin.Default()
	server.Use(middle...)

	ApiGroup := server.Group("/dev-api")
	v1 := ApiGroup.Group("/v1")

	router.NewRouter(s.rely, v1).RegisterRoutes()

	return server
}
