/**
 * Description：
 * FileName：menu.go
 * Author：CJiaの用心
 * Create：2025/5/13 16:41:14
 * Remark：
 */

package system

import (
	config "github.com/carefuly/carefuly-admin-go-gin/config/file"
	domainSystem "github.com/carefuly/carefuly-admin-go-gin/internal/domain/careful/system"
	"github.com/carefuly/carefuly-admin-go-gin/internal/service/careful/system"
	"github.com/carefuly/carefuly-admin-go-gin/pkg/ginx/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type MenuHandler interface {
	RegisterRoutes(router *gin.RouterGroup)
	GetList(ctx *gin.Context)
}

type menuHandler struct {
	rely config.RelyConfig
	svc  system.MenuService
}

func NewMenuHandler(rely config.RelyConfig, svc system.MenuService) MenuHandler {
	return &menuHandler{
		rely: rely,
		svc:  svc,
	}
}

// RegisterRoutes 注册路由
func (h *menuHandler) RegisterRoutes(router *gin.RouterGroup) {
	base := router.Group("/menu")
	base.GET("/listRouter", h.GetList)
}

// GetList
// @Summary 获取所有菜单
// @Description 获取所有菜单列表
// @Tags 菜单管理
// @Accept application/json
// @Produce application/json
// @Param status query int false "状态" default(-1)
// @Success 200 {array} domainSystem.Menu
// @Failure 400 {object} response.Response
// @Router /v1/system/menu/listRouter [get]
// @Security LoginToken
func (h *menuHandler) GetList(ctx *gin.Context) {
	// status, _ := strconv.Atoi(ctx.DefaultQuery("status", "-1"))

	list, err := h.svc.GetListAll(ctx, domainSystem.MenuFilter{})

	if err != nil {
		zap.L().Error("查询列表异常", zap.Error(err))
		response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
		return
	}

	response.NewResponse().SuccessResponse(ctx, "查询成功", list)
}
